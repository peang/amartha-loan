package main

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/peang/gojek-taxi/configs"
	"github.com/peang/gojek-taxi/models"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

// 1m33.023855875s with 5 concurrencies
// 1m26.783613042s with 100 concurrencies
func main() {
	startExecution := time.Now()

	conf := configs.LoadConfig()
	db := configs.MongoDatabaseConnector{}

	client := db.Connect(conf)
	defer client.Disconnect(context.Background())

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

	multiplier := 5
	batchSize := runtime.NumCPU() * multiplier
	numRows := int(pr.GetNumRows())

	tripChannel := make(chan taxiTrip, 20)
	modelEscrow := make([]interface{}, 0, 50)
	modelCountEscrow := 0

	var wg sync.WaitGroup
	var mt sync.Mutex

	fmt.Println("Please Wait Converting Datasets to MongoDB ...")

	for i := 0; i < batchSize; i++ {
		wg.Add(1)
		go persistData(client, &wg, &mt, tripChannel, &modelEscrow, &modelCountEscrow)
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

	fmt.Printf("Converted %d of Data at %v", modelCountEscrow, time.Since(startExecution))
}

func persistData(
	mongo *mongo.Client,
	wg *sync.WaitGroup,
	mt *sync.Mutex,
	c <-chan taxiTrip,
	modelEscrow *[]interface{},
	modelCountEscrow *int,
) {
	defer wg.Done()
	collection := mongo.Database("gojek").Collection("taxi_trips")

	for {
		trip, open := <-c
		if !open {
			break
		}

		tripStartTime := ConvertTime(*trip.TripStartTimestamp)
		tripEndTime := ConvertTime(*trip.TripEndTimestamp)
		pickupLat, pickupLong, _ := ExtractLongitudeLatitude(trip.PickupLocation)
		dropoffLat, dropoffLong, _ := ExtractLongitudeLatitude(trip.DropoffLocation)

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
			PickupLocation: models.Point{
				Type:        "Point",
				Coordinates: primitive.A{pickupLat, pickupLong},
			},
			DropoffLatitude:  NilCheck(trip.DropoffLatitude).(float64),
			DropoffLongitude: NilCheck(trip.DropoffLongitude).(float64),
			DropoffLocation: models.Point{
				Type:        "Point",
				Coordinates: primitive.A{dropoffLat, dropoffLong},
			},
		}

		mt.Lock()
		*modelEscrow = append(*modelEscrow, model)

		if len(*modelEscrow) == 100 {
			*modelCountEscrow += len(*modelEscrow)
			fmt.Printf("\rConverting %d of 3888425", *modelCountEscrow)
			if _, err := collection.InsertMany(context.Background(), *modelEscrow); err != nil {
				panic(err)
			}

			*modelEscrow = make([]interface{}, 0, 100)
		}
		mt.Unlock()
	}

	mt.Lock()
	if len(*modelEscrow) > 0 {
		*modelCountEscrow += len(*modelEscrow)
		fmt.Printf("\rConverting %d of 3888425", *modelCountEscrow)
		if _, err := collection.InsertMany(context.Background(), *modelEscrow); err != nil {
			panic(err)
		}

		*modelEscrow = make([]interface{}, 0, 100)
	}
	mt.Unlock()
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

func ExtractLongitudeLatitude(point *string) (float64, float64, error) {
	if point == nil {
		return 0, 0, nil
	}

	pointStr := *point
	pointStr = strings.TrimPrefix(pointStr, "POINT (")
	pointStr = strings.TrimSuffix(pointStr, ")")

	parts := strings.Split(pointStr, " ")

	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid point format: %s", pointStr)
	}

	longitude, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("error parsing longitude: %v", err)
	}

	latitude, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("error parsing latitude: %v", err)
	}

	return longitude, latitude, nil
}
