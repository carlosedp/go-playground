package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/opentracing-contrib/go-gin/ginhttp"
	"github.com/opentracing/opentracing-go"

	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	rootHits = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "webserver_get_root",
		Help: "The total number GET calls on /",
	}, []string{"code"})

	histogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "request_duration_u_seconds",
		Help:    "Time taken to answer request in microseconds",
		Buckets: []float64{10, 20, 50, 100, 200, 2000},
	}, []string{"code"})
)

func init() {
	prometheus.Register(histogram)
}

func main() {

	r := gin.Default()

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
		"gin-app",
		config.Logger(jaeger.StdLogger),
	)
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	// Tell gin to use the middleware
	r.Use(ginhttp.Middleware(tracer))

	// Prometheus metrics route
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Application routes
	r.GET("/", func(c *gin.Context) {
		os := runtime.GOOS
		arch := runtime.GOARCH
		fmt.Printf("%v", c.Request.Context())
		if span := opentracing.SpanFromContext(c.Request.Context()); span != nil {
			span.SetBaggageItem("Got OS/ARCH", fmt.Sprintf("%s/%s", os, arch))
		}
		start := time.Now()
		output := fmt.Sprintf("Hello, world! I'm running on %s/%s using Gin Framework!", os, arch)
		c.String(http.StatusOK, output)
		duration := time.Since(start)
		d := float64(duration) / 1000
		rootHits.WithLabelValues(fmt.Sprintf("%d", http.StatusOK)).Inc()
		histogram.WithLabelValues(fmt.Sprintf("%d", http.StatusOK)).Observe(d)
	})
	r.Run(":8082") // listen and serve on 0.0.0.0:8082
}
