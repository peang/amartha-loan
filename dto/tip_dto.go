package dto

type GetTotalTrip struct {
	StartDate string `validate:"required" json:"start"`
	EndDate   string `validate:"required,enddate" json:"end"`
}
