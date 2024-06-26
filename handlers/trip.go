package handlers

import "github.com/labstack/echo/v4"

type tripHandler struct{}

func NewTripHandler(e *echo.Echo) *tripHandler {
	handler := &tripHandler{}

	e.GET("/register", handler.getTotalTrips)

	return handler
}

func (h *tripHandler) getTotalTrips(ctx echo.Context) error {
	return nil
}
