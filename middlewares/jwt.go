package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/peang/amartha-loan-service/utils"
)

func (m *Middleware) JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenStr := c.Request().Header.Get("Authorization")

			if tokenStr == "" {
				return c.JSON(http.StatusUnauthorized, "unauthorized")
			}

			tokenInfo, err := utils.ParseToken(tokenStr)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, "unauthorized")
			}

			c.Set("payload", tokenInfo.Payload)

			return next(c)
		}
	}
}
