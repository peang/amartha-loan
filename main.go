package main

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/labstack/echo/v4"
)

func main() {
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
