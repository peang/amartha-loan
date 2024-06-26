package main

import (
	"context"
	"log"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/peang/gojek-taxi/configs"
	"github.com/peang/gojek-taxi/models"
	"github.com/uptrace/bun"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
)

type taxiTrip struct {
	UniqueKey            *string  `parquet:"name=unique_key, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	TaxiID               *string  `parquet:"name=taxi_id, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	TripStartTimestamp   *int64   `parquet:"name=trip_start_timestamp, type=DOUBLE"`
	TripEndTimestamp     *int64   `parquet:"name=trip_end_timestamp, type=DOUBLE"`
	TripSeconds          *float64 `parquet:"name=trip_seconds, type=DOUBLE"`
	TripMiles            *float64 `parquet:"name=trip_miles, type=DOUBLE"`
	PickupCensusTract    *float64 `parquet:"name=pickup_census_tract, type=DOUBLE"`
	DropoffCensusTract   *float64 `parquet:"name=dropoff_census_tract, type=DOUBLE"`
	PickupCommunityArea  *float64 `parquet:"name=pickup_community_area, type=DOUBLE"`
	DropoffCommunityArea *float64 `parquet:"name=dropoff_community_area, type=DOUBLE"`
	Fare                 *float64 `parquet:"name=fare, type=DOUBLE"`
	Tips                 *float64 `parquet:"name=tips, type=DOUBLE"`
	Tolls                *float64 `parquet:"name=tolls, type=DOUBLE"`
	Extras               *float64 `parquet:"name=extras, type=DOUBLE"`
	TripTotal            *float64 `parquet:"name=trip_total, type=DOUBLE"`
	PaymentType          *string  `parquet:"name=payment_type, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	Company              *string  `parquet:"name=company, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	PickupLatitude       *float64 `parquet:"name=pickup_latitude, type=DOUBLE"`
	PickupLongitude      *float64 `parquet:"name=pickup_longitude, type=DOUBLE"`
	PickupLocation       *string  `parquet:"name=pickup_location, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
	DropoffLatitude      *float64 `parquet:"name=dropoff_latitude, type=DOUBLE"`
	DropoffLongitude     *float64 `parquet:"name=dropoff_longitude, type=DOUBLE"`
	DropoffLocation      *string  `parquet:"name=dropoff_location, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY"`
}

func main() {
	conf := configs.LoadConfig()
	db := configs.LoadDatabase(conf)

	fr, err := local.NewLocalFileReader("datasets.parquet")
	if err != nil {
		log.Println("Can't open file")
		return
	}
	defer fr.Close()

	pr, err := reader.NewParquetReader(fr, new(taxiTrip), 4)
	if err != nil {
		log.Println("Can't create parquet reader", err)
		return
	}
	defer pr.ReadStop()

	batchSize := runtime.NumCPU() * 50
	numRows := int(pr.GetNumRows())

	tripChannel := make(chan taxiTrip, 20)

	var wg sync.WaitGroup
	wg.Add(batchSize)

	for i := 0; i < batchSize; i++ {
		go persistData(db, &wg, tripChannel)
	}

	for i := 0; i < numRows; i += batchSize {
		end := i + batchSize
		if end > numRows {
			end = numRows
		}
		taxiTrips := make([]taxiTrip, end-i)
		if err = pr.Read(&taxiTrips); err != nil {
			log.Println("Can't read:", err)
			return
		}

		for _, trip := range taxiTrips {
			tripChannel <- trip
		}
	}

	close(tripChannel)
	wg.Wait()
}

func persistData(db *bun.DB, wg *sync.WaitGroup, c <-chan taxiTrip) {
	defer wg.Done()
	ctx := context.TODO()

	for trip := range c {
		tripStartTime := ConvertTime(*trip.TripStartTimestamp)
		tripEndTime := ConvertTime(*trip.TripEndTimestamp)
		model := models.TaxiTrip{
			UniqueKey:            *trip.UniqueKey,
			TaxiID:               *trip.TaxiID,
			TripStartTimestamp:   tripStartTime,
			TripEndTimestamp:     tripEndTime,
			TripSeconds:          NilCheck(trip.TripSeconds).(float64),
			TripMiles:            NilCheck(trip.TripMiles).(float64),
			PickupCensusTract:    NilCheck(trip.PickupCensusTract).(float64),
			DropoffCensusTract:   NilCheck(trip.DropoffCensusTract).(float64),
			PickupCommunityArea:  NilCheck(trip.PickupCommunityArea).(float64),
			DropoffCommunityArea: NilCheck(trip.DropoffCommunityArea).(float64),
			Fare:                 NilCheck(trip.Fare).(float64),
			Tips:                 NilCheck(trip.Tips).(float64),
			Tolls:                NilCheck(trip.Tolls).(float64),
			Extras:               NilCheck(trip.Extras).(float64),
			TripTotal:            NilCheck(trip.TripTotal).(float64),
			PaymentType:          *trip.PaymentType,
			Company:              *trip.Company,
			PickupLatitude:       NilCheck(trip.PickupLatitude).(float64),
			PickupLocation:       NilCheck(trip.PickupLocation).(string),
			DropoffLatitude:      NilCheck(trip.DropoffLatitude).(float64),
			DropoffLongitude:     NilCheck(trip.DropoffLongitude).(float64),
			DropoffLocation:      NilCheck(trip.DropoffLocation).(string),
		}

		db.NewInsert().Model(&model).Exec(ctx)
	}
}

func NilCheck(ptr interface{}) interface{} {
	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return reflect.Zero(v.Type().Elem()).Interface()
	}
	return v.Elem().Interface()
}

func ConvertTime(baseTime int64) time.Time {
	timestampMicroseconds := int64(baseTime)
	timestampSeconds := timestampMicroseconds / 1e6
	timestampNanoseconds := (timestampMicroseconds % 1e6) * 1e3

	return time.Unix(timestampSeconds, timestampNanoseconds).UTC()
}
