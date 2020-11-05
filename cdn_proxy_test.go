package cdnproxy_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/eiri/echo-cdn-proxy"
	"github.com/labstack/echo/v4"
	"github.com/steinfletcher/apitest"
	"github.com/steinfletcher/apitest-jsonpath"
)

func TestProxyHandler(t *testing.T) {
	e := echo.New()
	e.Use(cdnproxy.Proxy)
	e.GET("/time", func(c echo.Context) error {
		now := map[string]string{
			"time": time.Now().Format("15:04:05"),
		}
		return c.JSON(http.StatusOK, now)
	})

	apitest.New().
		Handler(e).
		Get("/time").
		Expect(t).
		Status(http.StatusOK).
		Header("X-Proxy", "CDN").
		Assert(jsonpath.Present(`$.time`)).
		End()
}
