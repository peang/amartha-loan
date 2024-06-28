package dto

type GetTotalTripDTO struct {
	StartDate string `validate:"required" json:"start"`
	EndDate   string `validate:"required,enddate" json:"end"`
}

type GetFareHeatmapDTO struct {
	Date string `validate:"required" json:"date"`
}

type GetAverageSpeedDTO struct {
	Date string `validate:"required" json:"date"`
}
