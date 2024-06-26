package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaxiTrip struct {
	UniqueKey            string    `bson:"unique_key"`
	TaxiID               string    `bson:"taxi_id"`
	TripStartTimestamp   time.Time `bson:"trip_start_timestamp"`
	TripEndTimestamp     time.Time `bson:"trip_end_timestamp"`
	TripSeconds          float64   `bson:"trip_seconds"`
	TripMiles            float64   `bson:"trip_miles"`
	PickupCensusTract    float64   `bson:"pickup_census_tract"`
	DropoffCensusTract   float64   `bson:"dropoff_census_tract"`
	PickupCommunityArea  float64   `bson:"pickup_community_area"`
	DropoffCommunityArea float64   `bson:"dropoff_community_area"`
	Fare                 float64   `bson:"fare"`
	Tips                 float64   `bson:"tips"`
	Tolls                float64   `bson:"tolls"`
	Extras               float64   `bson:"extras"`
	TripTotal            float64   `bson:"trip_total"`
	PaymentType          string    `bson:"payment_type"`
	Company              string    `bson:"company"`
	PickupLatitude       float64   `bson:"pickup_latitude"`
	PickupLongitude      float64   `bson:"pickup_longitude"`
	PickupLocation       Point     `bson:"pickup_location"`
	DropoffLatitude      float64   `bson:"dropoff_latitude"`
	DropoffLongitude     float64   `bson:"dropoff_longitude"`
	DropoffLocation      Point     `bson:"dropoff_location"`
}

type Point struct {
	Type        string      `bson:"type"`
	Coordinates primitive.A `bson:"coordinates"`
}
