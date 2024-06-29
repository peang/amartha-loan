package utils

import (
	"fmt"
	"math"
	"strconv"
)

type Meta struct {
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}

func GenerateSort(str string) string {
	len := len(str)
	flag := string(str[0])
	if flag == "-" {
		return fmt.Sprintf("%v %v", str[1:len], "desc")
	} else {
		return fmt.Sprintf("%v %v", str[0:len], "asc")
	}
}

func GenerateOffsetLimit(page, perPage int) (offset, limit int) {
	if page < 1 {
		page = 1
	}

	if perPage < 1 {
		perPage = 10
	}

	offset = (page - 1) * perPage
	limit = perPage

	return offset, limit
}

func GenerateMeta(pageString string, perPageString string, count int) *Meta {
	page, err := strconv.Atoi(pageString)
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(perPageString)
	if err != nil || perPage < 1 || perPage > 100 {
		perPage = 10
	}

	return &Meta{
		Page:      page,
		PerPage:   perPage,
		Total:     count,
		TotalPage: int(math.Ceil(float64(count) / float64(perPage))),
	}
}
