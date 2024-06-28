package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func NewInternalServerError(c echo.Context, e error) error {
	return c.JSON(http.StatusInternalServerError, Error{
		Code:  http.StatusInternalServerError,
		Error: e.Error(),
	})
}

func NewValidationError(c echo.Context, e validator.ValidationErrors) error {
	errorMessages := make(map[string]string)

	for _, err := range e {
		fieldName := strings.ToLower(err.StructField())
		errorMessages[fieldName] = fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", err.StructField(), err.Tag())
	}

	return c.JSON(http.StatusBadRequest, Error{
		Code:  http.StatusBadRequest,
		Error: errorMessages,
	})
}
