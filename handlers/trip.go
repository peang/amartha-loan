package handlers

import (
	"net/http"
	"strconv"

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
	e.GET("/average_fare_heatmap", handler.getFareHeatmap)
	e.GET("/average_speed_24hrs", handler.getAverageSpeed)

	return handler
}

func (h *tripHandler) getTotalTrips(ctx echo.Context) error {
	startDate := ctx.QueryParam("start")
	endDate := ctx.QueryParam("end")

	dto := dto.GetTotalTripDTO{
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

	result, err := h.taxiTripRepository.GetTotalTrips(ctx.Request().Context(), *parsedStartDate, *parsedEndDate)
	if err != nil {
		return utils.NewInternalServerError(ctx.Echo().AcquireContext(), err)
	}

	return ctx.JSON(http.StatusOK, utils.Response{
		Data: result,
	})
}

func (h *tripHandler) getFareHeatmap(ctx echo.Context) error {
	date := ctx.QueryParam("date")
	page, _ := strconv.ParseFloat(ctx.QueryParam("page"), 64)
	if page == 0 {
		page = 1
	}

	perPage, _ := strconv.ParseFloat(ctx.QueryParam("perPage"), 64)
	if perPage == 0 {
		perPage = 10
	}

	if perPage > 100 {
		perPage = 100 // fallback to protect memory
	}

	parsedDate, err := utils.DateParser(date)
	if err != nil {
		return utils.NewInternalServerError(ctx.Echo().AcquireContext(), err)
	}

	dto := dto.GetFareHeatmapDTO{
		Date:    *parsedDate,
		Page:    page,
		PerPage: perPage,
	}

	if err := ctx.Validate(dto); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return utils.NewValidationError(ctx, err.(validator.ValidationErrors))
		}
		return err
	}

	result, err := h.taxiTripRepository.GetFareHeatmap(ctx.Request().Context(), &dto)
	if err != nil {
		return utils.NewInternalServerError(ctx.Echo().AcquireContext(), err)
	}

	return ctx.JSON(http.StatusOK, result)
}

func (h *tripHandler) getAverageSpeed(ctx echo.Context) error {
	date := ctx.QueryParam("date")

	dto := dto.GetAverageSpeedDTO{
		Date: date,
	}

	if err := ctx.Validate(dto); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			return utils.NewValidationError(ctx, err.(validator.ValidationErrors))
		}
		return err
	}

	parsedDate, err := utils.DateParser(date)
	if err != nil {
		return utils.NewInternalServerError(ctx.Echo().AcquireContext(), err)
	}

	result, err := h.taxiTripRepository.GetAverageSpeed(ctx.Request().Context(), *parsedDate)
	if err != nil {
		return utils.NewInternalServerError(ctx.Echo().AcquireContext(), err)
	}

	return ctx.JSON(http.StatusOK, utils.Response{
		Data: result,
	})
}
