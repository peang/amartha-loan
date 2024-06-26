package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/labstack/echo/v4"
	"github.com/peang/gojek-taxi/configs"
)

func main() {
	conf := configs.LoadConfig()
	mongo := configs.LoadDatabase(conf)
	defer mongo.Disconnect(context.Background())

	e := echo.New()

	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
