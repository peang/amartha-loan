CREATE TABLE taxi_trips (
    unique_key VARCHAR(255) PRIMARY KEY,
    taxi_id VARCHAR(255),
    trip_start_timestamp TIMESTAMP,
    trip_end_timestamp TIMESTAMP,
    trip_seconds FLOAT,
    trip_miles FLOAT,
    pickup_census_tract FLOAT,
    dropoff_census_tract FLOAT,
    pickup_community_area FLOAT,
    dropoff_community_area FLOAT,
    fare FLOAT,
    tips FLOAT,
    tolls FLOAT,
    extras FLOAT,
    trip_total FLOAT,
    payment_type VARCHAR(255),
    company VARCHAR(255),
    pickup_latitude FLOAT,
    pickup_longitude FLOAT,
    pickup_location VARCHAR(255),
    dropoff_latitude FLOAT,
    dropoff_longitude FLOAT,
    dropoff_location VARCHAR(255)
);

CREATE INDEX idx_trip_start_timestamp ON taxi_trips (trip_start_timestamp);
CREATE INDEX idx_trip_total ON taxi_trips (trip_total);
