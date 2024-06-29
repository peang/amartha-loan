package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/peang/gojek-taxi/dto"
	"github.com/peang/gojek-taxi/repositories"
	"github.com/peang/gojek-taxi/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTaxiTripRepository struct {
	mock.Mock
}

func (m *MockTaxiTripRepository) GetTotalTrips(ctx context.Context, startTime time.Time, endTime time.Time) (*[]repositories.TotalTripResponse, error) {
	args := m.Called(ctx, startTime, endTime)
	return args.Get(0).(*[]repositories.TotalTripResponse), args.Error(1)
}

func (m *MockTaxiTripRepository) GetFareHeatmap(ctx context.Context, dto *dto.GetFareHeatmapDTO) (*repositories.FareHeatmapResponse, error) {
	args := m.Called(ctx, dto)
	return args.Get(0).(*repositories.FareHeatmapResponse), args.Error(1)
}

func (m *MockTaxiTripRepository) GetAverageSpeed(ctx context.Context, time time.Time) (*repositories.AverageSpeedResponse, error) {
	args := m.Called(ctx, time)
	return args.Get(0).(*repositories.AverageSpeedResponse), args.Error(1)
}

func TestGetTotalTrips(t *testing.T) {
	e := echo.New()

	validators := validator.New()
	validators.RegisterValidation("enddate", utils.ValidateEndDate)
	e.Validator = &utils.Validator{Validator: validators}

	req := httptest.NewRequest(http.MethodGet, "/total_trips?start=2020-01-01&end=2020-01-31", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := new(MockTaxiTripRepository)
	handler := NewTripHandler(e, mockRepo)

	startDate, _ := utils.DateParser("2020-01-01")
	endDate, _ := utils.DateParser("2020-01-31")

	expectedResponse := &[]repositories.TotalTripResponse{
		{Date: "2020-01-01", TotalTrip: 10},
	}

	mockRepo.On("GetTotalTrips", mock.Anything, *startDate, *endDate).Return(expectedResponse, nil)

	if assert.NoError(t, handler.getTotalTrips(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"data":[{"date":"2020-01-01","total_trips":10}]}`, strings.TrimSpace(rec.Body.String()))
	}
}

func TestGetFareHeatmap(t *testing.T) {
	e := echo.New()

	validators := validator.New()
	validators.RegisterValidation("enddate", utils.ValidateEndDate)
	e.Validator = &utils.Validator{Validator: validators}

	req := httptest.NewRequest(http.MethodGet, "/average_fare_heatmap?date=2020-01-01&page=1&perPage=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := new(MockTaxiTripRepository)
	handler := NewTripHandler(e, mockRepo)

	parsedDate, _ := utils.DateParser("2020-01-01")
	dto := &dto.GetFareHeatmapDTO{
		Date:    *parsedDate,
		Page:    1,
		PerPage: 10,
	}

	expectedResponse := &repositories.FareHeatmapResponse{
		Data: []interface{}{"dummy data"},
		Meta: repositories.MetaResponse{
			Page:    1,
			PerPage: 10},
	}

	mockRepo.On("GetFareHeatmap", mock.Anything, dto).Return(expectedResponse, nil)

	if assert.NoError(t, handler.getFareHeatmap(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"data":["dummy data"],"meta":{"page":1,"perPage":10}}`, strings.TrimSpace(rec.Body.String()))
	}
}

func TestGetAverageSpeed(t *testing.T) {
	e := echo.New()
	validators := validator.New()
	validators.RegisterValidation("enddate", utils.ValidateEndDate)
	e.Validator = &utils.Validator{Validator: validators}

	req := httptest.NewRequest(http.MethodGet, "/average_speed_24hrs?date=2020-01-01", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := new(MockTaxiTripRepository)
	handler := NewTripHandler(e, mockRepo)

	parsedDate, _ := utils.DateParser("2020-01-01")

	expectedResponse := &repositories.AverageSpeedResponse{
		AverageSpeed: 50.5,
	}

	mockRepo.On("GetAverageSpeed", mock.Anything, *parsedDate).Return(expectedResponse, nil)

	if assert.NoError(t, handler.getAverageSpeed(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"data":{"average_speed":50.5}}`, strings.TrimSpace(rec.Body.String()))
	}
}
