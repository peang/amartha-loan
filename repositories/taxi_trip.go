package repositories

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/golang/geo/s2"
	"github.com/peang/gojek-taxi/dto"
	"github.com/peang/gojek-taxi/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaxiTripRepositoryInterface interface {
	GetTotalTrips(ctx context.Context, startTime time.Time, endTime time.Time) (*[]TotalTripResponse, error)
	GetFareHeatmap(ctx context.Context, dto *dto.GetFareHeatmapDTO) (*FareHeatmapResponse, error)
	GetAverageSpeed(ctx context.Context, time time.Time) (*AverageSpeedResponse, error)
}

type TaxiTripRepositoryFilter struct {
	PickupCommunityArea int `json:"pickup_community_area"`
}

type taxiTripRepository struct {
	client *mongo.Client
}

type TotalTripResponse struct {
	Date      string `json:"date"`
	TotalTrip int32  `json:"total_trips"`
}

type FareHeatmapResponse struct {
	Data interface{} `json:"data"`
	Meta struct {
		Page    int64
		PerPage int64
	} `json:"meta"`
}

type AverageSpeedResponse struct {
	AverageSpeed float64 `json:"average_speed"`
}

func NewTaxiTripRepository(c *mongo.Client) TaxiTripRepositoryInterface {
	return &taxiTripRepository{
		client: c,
	}
}

func (r *taxiTripRepository) GetTotalTrips(ctx context.Context, startTime time.Time, endTime time.Time) (*[]TotalTripResponse, error) {
	collection := r.client.Database("gojek").Collection("taxi_trips")

	pipeline := mongo.Pipeline{
		bson.D{
			{"$match", bson.D{
				{"trip_start_timestamp", bson.D{
					{"$gte", startTime},
					{"$lte", endTime},
				}},
			}},
		},
		bson.D{
			{"$group", bson.D{
				{"_id", bson.D{
					{"$dateToString", bson.D{
						{"format", "%Y-%m-%d"},
						{"date", "$trip_start_timestamp"},
					}},
				}},
				{"total_trips", bson.D{
					{"$sum", 1},
				}},
			}},
		},
		bson.D{
			{"$sort", bson.D{
				{"_id", 1},
			}},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var results []TotalTripResponse
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}

		transformedResult := TotalTripResponse{
			Date:      result["_id"].(string),
			TotalTrip: result["total_trips"].(int32),
		}

		results = append(results, transformedResult)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	return &results, nil
}

func (r *taxiTripRepository) GetFareHeatmap(ctx context.Context, dto *dto.GetFareHeatmapDTO) (*FareHeatmapResponse, error) {
	collection := r.client.Database("gojek").Collection("taxi_trips")

	nextDay := dto.Date.AddDate(0, 0, 1)
	skip, limit := utils.GeneratePagination(dto.Page, dto.PerPage)
	findOptions := options.Find()
	findOptions.SetSkip(skip)
	findOptions.SetLimit(limit)

	cursor, err := collection.Find(ctx, bson.M{
		"trip_start_timestamp": bson.M{
			"$gte": dto.Date,
			"$lt":  nextDay,
		},
	}, findOptions)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var data []interface{}
	chTrip := make(chan primitive.M, 100)
	var wg sync.WaitGroup
	var mt sync.Mutex

	for i := 0; i <= 5; i++ {
		wg.Add(1)
		go fareHeatmapWorker(&wg, &mt, chTrip, &data)
	}

	for cursor.Next(ctx) {
		var trip bson.M
		if err := cursor.Decode(&trip); err != nil {
			return nil, err
		}
		chTrip <- trip
	}

	close(chTrip)
	wg.Wait()

	return &FareHeatmapResponse{
		Data: data,
		Meta: struct {
			Page    int64
			PerPage int64
		}{
			Page:    int64(dto.Page),
			PerPage: int64(dto.PerPage),
		},
	}, nil
}

func (r *taxiTripRepository) GetAverageSpeed(ctx context.Context, time time.Time) (*AverageSpeedResponse, error) {
	collection := r.client.Database("gojek").Collection("taxi_trips")
	nextDay := time.AddDate(0, 0, 1)

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"trip_start_timestamp": bson.M{
					"$gte": time,
					"$lt":  nextDay,
				},
				"trip_seconds": bson.M{
					"$ne": 0,
				},
			},
		},
		{
			"$group": bson.M{
				"_id": nil,
				"total_seconds": bson.M{
					"$sum": "$trip_seconds",
				},
				"total_miles": bson.M{
					"$sum": "$trip_miles",
				},
			},
		},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var totalHour, totalKilometer float64
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}

		totalSeconds := result["total_seconds"]
		totalMiles := result["total_miles"]

		totalHour = totalSeconds.(float64) / 3600
		totalKilometer = totalMiles.(float64) * 1.60934
	}

	speedAverage := totalKilometer / totalHour
	resultsString := fmt.Sprintf("%.2f", speedAverage)
	resultsFloat, err := strconv.ParseFloat(resultsString, 64)
	if err != nil {
		return nil, err
	}

	return &AverageSpeedResponse{
		AverageSpeed: resultsFloat,
	}, nil
}

func fareHeatmapWorker(wg *sync.WaitGroup, mutex *sync.Mutex, cursorChannel <-chan primitive.M, result *[]interface{}) error {
	defer wg.Done()
	s2Map := make(map[s2.CellID]struct {
		totalFare float64
		count     int
	})

	for trip := range cursorChannel {
		lat, _ := trip["pickup_latitude"].(float64)
		lng, _ := trip["pickup_longitude"].(float64)
		cellID := s2.CellIDFromLatLng(s2.LatLngFromDegrees(lat, lng)).Parent(16)

		fare, _ := trip["fare"].(float64)

		if val, ok := s2Map[cellID]; ok {
			s2Map[cellID] = struct {
				totalFare float64
				count     int
			}{
				totalFare: val.totalFare + fare,
				count:     val.count + 1,
			}
		} else {
			s2Map[cellID] = struct {
				totalFare float64
				count     int
			}{
				totalFare: fare,
				count:     1,
			}
		}
	}

	for cellID, data := range s2Map {
		parsedFloat, err := strconv.ParseFloat(fmt.Sprintf("%.2f", data.totalFare/float64(data.count)), 64)
		if err != nil {
			parsedFloat = 0
		}

		mutex.Lock()
		*result = append(*result, struct {
			S2ID string
			Fare float64
		}{
			S2ID: cellID.ToToken(),
			Fare: parsedFloat,
		})
		mutex.Unlock()
	}

	return nil
}
