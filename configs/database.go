package configs

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func LoadDatabase(c *Config) *mongo.Client {
	clientOptions := options.Client().ApplyURI(
		fmt.Sprintf("mongodb://%s:27017/%s", c.MongoHost, c.MongoDatabase),
	).
		SetAuth(options.Credential{
			Username:   c.MongoUsername,
			Password:   c.MongoPassword,
			AuthSource: "admin",
		})
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Panicf("Faild Mongo Connection : %s", err)
	}

	return client
}
