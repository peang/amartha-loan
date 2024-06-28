package dto

import "time"

type GetTotalTripDTO struct {
	StartDate string `validate:"required" json:"start"`
	EndDate   string `validate:"required,enddate" json:"end"`
}

type GetFareHeatmapDTO struct {
	Date    time.Time `validate:"required" json:"date"`
	Page    float64
	PerPage float64
}

type GetAverageSpeedDTO struct {
	Date string `validate:"required" json:"date"`
}
