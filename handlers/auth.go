package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/peang/amartha-loan-service/repositories"
	"github.com/peang/amartha-loan-service/utils"
)

type authHandler struct {
	userRepository repositories.UserRepositoryInterface
}

func NewAuthHandler(e *echo.Echo, userRepository repositories.UserRepositoryInterface) *authHandler {
	handler := &authHandler{
		userRepository: userRepository,
	}

	authGroup := e.Group("/auths")

	authGroup.POST("/login", handler.login)

	return handler
}

func (h *authHandler) login(ctx echo.Context) error {
	// For Demo Purpose, gonna use query params to determine user role
	role := ctx.QueryParam("role")
	if role == "" {
		return ctx.JSON(http.StatusBadRequest, utils.Error{
			Code:  http.StatusBadRequest,
			Error: "Query 'role' is Required",
		})
	}

	roleInt, err := strconv.Atoi(role)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.Error{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
	}

	user, err := h.userRepository.Detail(ctx.Request().Context(), uint(roleInt))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, utils.Error{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
	}

	token, err := utils.CreateJWTToken(user, true)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, utils.Error{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		})
	}

	//!TODO Implement Data Transfer Object

	//!TODO Implement Validator

	return ctx.JSON(http.StatusOK, utils.Response{
		Message: "Login Success",
		Data:    token,
	})
}
