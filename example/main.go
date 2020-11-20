package main

import (
	"net/http"
	"time"

	"github.com/eiri/echo-cdn-proxy"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = false
	e.Debug = true
	e.Logger.SetLevel(1)

	logFmt := "${time_rfc3339} ${status} ${method} ${uri} - ${latency_human}\n"
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: logFmt}))
	e.Use(middleware.Recover())

	// Configure the middleware
	cfg := cdnproxy.NewConfig("https://cdn.jsdelivr.net", "/npm")
	e.Use(cfg.Proxy)

	// Serve static content
	e.Static("/", "example/frontend")

	// Serve dynamic content
	e.GET("/time", func(c echo.Context) error {
		return c.JSON(http.StatusOK, time.Now().Format("15:04:05"))
	})

	// Start the server
	e.Logger.Fatal(e.Start(":8000"))
}
