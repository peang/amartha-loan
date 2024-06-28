db.getSiblingDB("gojek");
db.createCollection("taxi_trips");

db.taxi_trips.createIndex({ "trip_start_timestamp": 1 }, { name: "trip_start_timestamp_index" })
db.taxi_trips.createIndex({ "trip_seconds": 1 }, { name: "trip_seconds_index" })

db.taxi_trips.createIndex({ pickup_location: "2dsphere" });
db.taxi_trips.createIndex({ dropoff_location: "2dsphere" });
