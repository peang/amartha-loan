package models

import (
	"time"

	"github.com/uptrace/bun"
)

type TaxiTrip struct {
	bun.BaseModel `bun:"table:taxi_trips"`

	UniqueKey            string    `bun:"type:VARCHAR(255),pk"`
	TaxiID               string    `bun:"type:VARCHAR(255)"`
	TripStartTimestamp   time.Time `bun:"type:TIMESTAMP,index:idx_trip_start_timestamp"`
	TripEndTimestamp     time.Time `bun:"type:TIMESTAMP"`
	TripSeconds          float64   `bun:"type:INTEGER"`
	TripMiles            float64   `bun:"type:FLOAT"`
	PickupCensusTract    float64   `bun:"type:INTEGER"`
	DropoffCensusTract   float64   `bun:"type:INTEGER"`
	PickupCommunityArea  float64   `bun:"type:INTEGER"`
	DropoffCommunityArea float64   `bun:"type:INTEGER"`
	Fare                 float64   `bun:"type:FLOAT"`
	Tips                 float64   `bun:"type:FLOAT"`
	Tolls                float64   `bun:"type:FLOAT"`
	Extras               float64   `bun:"type:FLOAT"`
	TripTotal            float64   `bun:"type:FLOAT,index:idx_trip_total"`
	PaymentType          string    `bun:"type:VARCHAR(255)"`
	Company              string    `bun:"type:VARCHAR(255)"`
	PickupLatitude       float64   `bun:"type:FLOAT"`
	PickupLongitude      float64   `bun:"type:FLOAT"`
	PickupLocation       string    `bun:"type:STRING"`
	DropoffLatitude      float64   `bun:"type:FLOAT"`
	DropoffLongitude     float64   `bun:"type:FLOAT"`
	DropoffLocation      string    `bun:"type:STRING"`
}
