package main

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/labstack/echo-contrib/prometheus"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Create new server
	e := echo.New()

	// Logging Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: metricsSkipper,
	}))

	// Prometheus Metrics Middleware
	// Depends on github.com/carlosedp/echo-contrib@prometheus
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	// Application routes
	e.GET("/", rootHandler)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// metricsSkipper ignores metrics route on some middleware
func metricsSkipper(c echo.Context) bool {
	if strings.HasPrefix(c.Path(), "/metrics") {
		return true
	}
	return false
}

// Handlers
func rootHandler(c echo.Context) error {
	os := runtime.GOOS
	arch := runtime.GOARCH
	name := c.QueryParam("name")

	if name == "" {
		name = "World"
	}
	output := fmt.Sprintf("Hello, %s! I'm running on %s/%s inside a container!", name, os, arch)
	return c.String(http.StatusOK, output)
}
