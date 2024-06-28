package utils

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	Validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	if err := v.Validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func ValidateEndDate(fl validator.FieldLevel) bool {
	endDateStr := fl.Field().String()
	startDateStr := fl.Parent().FieldByName("StartDate").String()

	// Parse dates
	endDate, err1 := time.Parse("2006-01-02", endDateStr)
	startDate, err2 := time.Parse("2006-01-02", startDateStr)

	// Check if parsing succeeded
	if err1 != nil || err2 != nil {
		return false
	}

	// Compare dates
	return !endDate.Before(startDate)
}
