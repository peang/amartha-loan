db.getSiblingDB("gojek");
db.createCollection("taxi_trips");

db.taxi_trips.createIndex({ pickup_location: "2dsphere" });
db.taxi_trips.createIndex({ dropoff_location: "2dsphere" });
