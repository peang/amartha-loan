package repositories

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/golang/geo/s2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaxiTripRepositoryInterface interface {
	GetTotalTrips(ctx context.Context, startTime time.Time, endTime time.Time) (*[]TotalTripResponse, error)
	GetFareHeatmap(ctx context.Context, time time.Time) (*[]FareHeatmapResponse, error)
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
	S2ID string  `json:"s2id"`
	Fare float64 `json:"fare"`
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

func (r *taxiTripRepository) GetFareHeatmap(ctx context.Context, time time.Time) (*[]FareHeatmapResponse, error) {
	collection := r.client.Database("gojek").Collection("taxi_trips")

	nextDay := time.AddDate(0, 0, 1)

	options := options.Find().SetBatchSize(100)

	cursor, err := collection.Find(ctx, bson.M{
		"trip_start_timestamp": bson.M{
			"$gte": time,
			"$lt":  nextDay,
		},
	}, options)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []FareHeatmapResponse
	chTrip := make(chan primitive.M, 100)
	var wg sync.WaitGroup
	var mt sync.Mutex

	for i := 0; i <= 5; i++ {
		wg.Add(1)
		go fareHeatmapWorker(&wg, &mt, chTrip, &results)
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

	return &results, nil
}

func (r *taxiTripRepository) GetAverageSpeed(ctx context.Context, time time.Time) (*AverageSpeedResponse, error) {
	collection := r.client.Database("gojek").Collection("taxi_trips")

	nextDay := time.AddDate(0, 0, 1)
	options := options.Find().SetBatchSize(100)

	cursor, err := collection.Find(ctx, bson.M{
		"trip_start_timestamp": bson.M{
			"$gte": time,
			"$lt":  nextDay,
		},
		"trip_seconds": bson.M{
			"$ne": 0, // This is Importan because division by Zero could lead to NaN value
		},
	}, options)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var totalSpeed float64
	var totalData float64

	chTrip := make(chan primitive.M, 10)
	var wg sync.WaitGroup
	var mt sync.Mutex

	for i := 0; i <= 10; i++ {
		wg.Add(1)
		go averageSpeedWorker(&wg, &mt, chTrip, &totalSpeed, &totalData)
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

	results := totalSpeed / totalData
	resultsString := fmt.Sprintf("%.2f", results)
	resultsFloat, err := strconv.ParseFloat(resultsString, 64)
	if err != nil {
		return nil, err
	}

	return &AverageSpeedResponse{
		AverageSpeed: resultsFloat,
	}, nil
}

func fareHeatmapWorker(wg *sync.WaitGroup, mutex *sync.Mutex, cursorChannel <-chan primitive.M, result *[]FareHeatmapResponse) error {
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
		*result = append(*result, FareHeatmapResponse{
			S2ID: cellID.ToToken(),
			Fare: parsedFloat,
		})
		mutex.Unlock()
	}

	return nil
}

func averageSpeedWorker(wg *sync.WaitGroup, mutex *sync.Mutex, cursorChannel <-chan primitive.M, totalSpeed *float64, totalData *float64) error {
	defer wg.Done()

	for trip := range cursorChannel {
		tripTime, _ := trip["trip_seconds"].(float64)
		tripRange, _ := trip["trip_miles"].(float64)

		tripRangeKilometer := tripRange * 1.60934
		tripSpeed := tripRangeKilometer / (tripTime / 3600)

		mutex.Lock()
		*totalSpeed += tripSpeed
		*totalData += 1
		mutex.Unlock()

		fmt.Printf("Time %f, Range %f, Speed %f\n", tripTime, tripRange, tripSpeed)
	}

	return nil
}
