package main

import (
	mw "echo-tracing-metrics/middleware"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	hitMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "webserver_get_root",
		Help: "The total number GET calls on each path",
	}, []string{"code", "path"})

	name string
)

func init() {
	prometheus.Register(hitMetric)
}

// PromHitMetric is the middleware function to increase hit metric.
func PromHitMetric(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			c.Error(err)
		}
		if c.Path() != "/metrics" {
			hitMetric.WithLabelValues(fmt.Sprintf("%d", http.StatusOK), c.Path()).Inc()
		}
		return nil
	}
}

// MetricsSkipper ignores metrics route on middleware
func MetricsSkipper(c echo.Context) bool {
	if strings.HasPrefix(c.Path(), "/metrics") {
		return true
	}
	return false
}

func main() {
	// Add Opentracing instrumentation
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}
	tracer, closer, _ := cfg.New(
		"echo-tracing",
		config.Logger(jaeger.StdLogger),
	)
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	// Create new server
	e := echo.New()

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Skipper: MetricsSkipper}))
	e.Use(PromHitMetric)
	e.Use(mw.TraceWithConfig(mw.TraceConfig{
		Tracer:  tracer,
		Skipper: MetricsSkipper}))

	// Application routes
	e.GET("/", rootHandler)
	e.GET("/test", testHandler)

	// Route to Prometheus Metrics
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// Handlers
func rootHandler(c echo.Context) error {
	os := runtime.GOOS
	arch := runtime.GOARCH
	time.Sleep(50 * time.Millisecond)
	name = c.QueryParam("name")
	if name == "" {
		fmt.Printf("%v", c)
		if span := opentracing.SpanFromContext(c.Request().Context()); span != nil {
			fmt.Printf("---> %v\n", span)
			span.SetBaggageItem("Got OS/ARCH", fmt.Sprintf("%s/%s", os, arch))
		} else {
			fmt.Println("---> none\n")
		}
		time.Sleep(100 * time.Millisecond)
		name = "World"
	}
	output := fmt.Sprintf("Hello, %s! I'm running on %s/%s inside a container!", name, os, arch)
	return c.String(http.StatusOK, output)
}

func testHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Test path")
}
