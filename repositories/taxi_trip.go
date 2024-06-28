package repositories

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaxiTripRepositoryInterface interface {
	TotalTrips(ctx context.Context, startTime time.Time, endTime time.Time) (interface{}, error)
}

type TaxiTripRepositoryFilter struct {
	PickupCommunityArea int `json:"pickup_community_area"`
}

type taxiTripRepository struct {
	client *mongo.Client
}

func NewTaxiTripRepository(c *mongo.Client) TaxiTripRepositoryInterface {
	return &taxiTripRepository{
		client: c,
	}
}

func (r *taxiTripRepository) TotalTrips(ctx context.Context, startTime time.Time, endTime time.Time) (interface{}, error) {
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

	fmt.Printf("Aggregation Pipeline: %s\n", pipeline)

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}

		transformedResult := bson.M{
			"date":        result["_id"],
			"total_trips": result["total_trips"],
		}

		results = append(results, transformedResult)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	return results, nil
}
