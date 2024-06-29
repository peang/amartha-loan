package middleware

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/peang/amartha-loan-service/utils"
)

func (m *Middleware) RBACMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			payload := c.Get("payload").(utils.Payload)
			role := strconv.FormatInt(int64(payload.Role), 10)

			obj := c.Path()
			act := c.Request().Method

			if ok, err := m.enforcer.Enforce(role, obj, act); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to enforce scope policy")
			} else if !ok {
				return echo.NewHTTPError(http.StatusForbidden, "Unauthorized")
			}

			return next(c)
		}
	}

}
