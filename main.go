package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/peang/gojek-taxi/configs"
	"github.com/peang/gojek-taxi/handlers"
	"github.com/peang/gojek-taxi/repositories"
	"github.com/peang/gojek-taxi/utils"
)

func main() {
	conf := configs.LoadConfig()
	db := configs.MongoDatabaseConnector{}

	client := db.Connect(conf)
	defer client.Disconnect(context.Background())

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	validators := validator.New()
	validators.RegisterValidation("enddate", utils.ValidateEndDate)
	e.Validator = &utils.Validator{Validator: validators}

	taxiTripRepository := repositories.NewTaxiTripRepository(client)

	handlers.NewTripHandler(e, taxiTripRepository)

	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
