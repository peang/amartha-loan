package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/peang/gojek-taxi/dto"
	"github.com/peang/gojek-taxi/repositories"
	"github.com/peang/gojek-taxi/utils"
)

type tripHandler struct {
	taxiTripRepository repositories.TaxiTripRepositoryInterface
}

func NewTripHandler(e *echo.Echo, taxiTripRepository repositories.TaxiTripRepositoryInterface) *tripHandler {
	handler := &tripHandler{
		taxiTripRepository: taxiTripRepository,
	}

	e.GET("/total_trips", handler.getTotalTrips)

	return handler
}

func (h *tripHandler) getTotalTrips(ctx echo.Context) error {
	startDate := ctx.QueryParam("start")
	endDate := ctx.QueryParam("end")

	dto := dto.GetTotalTrip{
		StartDate: startDate,
		EndDate:   endDate,
	}

	if err := ctx.Validate(dto); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return utils.NewValidationError(ctx, err.(validator.ValidationErrors))
		}
		return err
	}

	parsedStartDate, err := utils.DateParser(startDate)
	if err != nil {
		return utils.NewInternalServerError(ctx.Echo().AcquireContext(), err)
	}

	parsedEndDate, err := utils.DateParser(endDate)
	if err != nil {
		return utils.NewInternalServerError(ctx.Echo().AcquireContext(), err)
	}

	result, err := h.taxiTripRepository.TotalTrips(ctx.Request().Context(), *parsedStartDate, *parsedEndDate)
	if err != nil {
		return utils.NewInternalServerError(ctx.Echo().AcquireContext(), err)
	}

	return ctx.JSON(http.StatusOK, utils.Response{
		Data: result,
	})
}
