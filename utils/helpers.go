package utils

import "time"

func DateParser(date string) (*time.Time, error) {
	parsed, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}
