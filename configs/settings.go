package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env  string
	Port string

	MongoHost     string
	MongoUsername string
	MongoPassword string
	MongoDatabase string
}

func LoadConfig() (c *Config) {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	return &Config{
		Env:  os.Getenv("ENV"),
		Port: port,

		MongoHost:     os.Getenv("MONGODB_HOST"),
		MongoUsername: os.Getenv("MONGODB_USERNAME"),
		MongoPassword: os.Getenv("MONGODB_PASSWORD"),
		MongoDatabase: os.Getenv("MONGODB_DATABASE"),
	}
}
