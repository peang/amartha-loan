package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Env  string
	Port string

	SQLHost     string
	SQLPort     int
	SQLUsername string
	SQLPassword string
	SQLDatabase string
	SQLSSL      string
}

func LoadConfig() (c *Config) {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	sqlPort, _ := strconv.Atoi(os.Getenv("SQL_PORT"))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	return &Config{
		Env:  os.Getenv("ENV"),
		Port: port,

		SQLHost:     os.Getenv("SQL_HOST"),
		SQLPort:     sqlPort,
		SQLUsername: os.Getenv("SQL_USERNAME"),
		SQLPassword: os.Getenv("SQL_PASSWORD"),
		SQLDatabase: os.Getenv("SQL_DATABASE"),
		SQLSSL:      os.Getenv("SQL_SSL"),
	}
}
