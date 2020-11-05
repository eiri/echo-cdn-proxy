package cdnproxy

import (
	"github.com/labstack/echo/v4"
)

// Proxy adds a `X-Proxy` header to the response.
func Proxy(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("X-Proxy", "CDN")
		return next(c)
	}
}
