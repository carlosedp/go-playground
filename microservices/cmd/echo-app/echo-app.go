package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strings"
	"time"

	// t "microservices/lib/tracing"

	"github.com/labstack/echo-contrib/jaegertracing"

	// "github.com/labstack/echo-contrib/prometheus"

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
	// p := prometheus.NewPrometheus("echo", nil)
	// p.Use(e)

	c := jaegertracing.New(e, metricsSkipper)
	defer c.Close()

	// Application routes
	e.GET("/", rootHandler)
	e.GET("/test/:name", testHandler)

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
	time.Sleep(50 * time.Millisecond)
	name := c.QueryParam("name")

	if name == "" {
		time.Sleep(100 * time.Millisecond)
		name = "World"
	}
	output := fmt.Sprintf("Hello, %s! I'm running on %s/%s inside a container!", name, os, arch)
	return c.String(http.StatusOK, output)
}

func testHandler(c echo.Context) error {
	// Demonstrate Span creation for custom log/tag/baggage append
	sp := jaegertracing.CreateChildSpan(c, "test handler")
	defer sp.Finish()
	var name = ""
	name = c.Param("name")
	sp.LogEvent("Called testHandler function, HTTP name param is: " + name)
	sp.SetBaggageItem("name", name)
	sp.SetTag("name_tag", "name")
	time.Sleep(10 * time.Millisecond)

	// Call slow function 5 times, it will create it's own span
	ch := make(chan string)
	for index := 0; index < 5; index++ {
		// Do in parallel
		go jaegertracing.TraceFunction(c, slowFunc, "Test String", ch)
	}
	for index := 0; index < 5; index++ {
		fmt.Println(<-ch)
	}

	ret := fmt.Sprintf("Test path, name: %s", name)
	return c.String(http.StatusOK, ret)
}

// A function to be wrapped
func slowFunc(s string, c chan string) {
	time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
	c <- "received " + s
}
