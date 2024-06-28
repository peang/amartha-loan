package utils

import "time"

func DateParser(date string) (*time.Time, error) {
	parsed, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

func GeneratePagination(page float64, perPage float64) (skip int64, limit int64) {
	skipValue := (page - 1) * perPage
	limitValue := perPage

	return int64(skipValue), int64(limitValue)
}
