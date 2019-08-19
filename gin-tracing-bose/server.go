package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	ginopentracing "github.com/Bose/go-gin-opentracing"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

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
	// initialize the global singleton for tracing...
	tracer, reporter, closer, err := ginopentracing.InitTracing("gin-tracing", "localhost:5775", ginopentracing.WithEnableInfoLog(true))
	if err != nil {
		panic("unable to init tracing")
	}
	defer closer.Close()
	defer reporter.Close()
	opentracing.SetGlobalTracer(tracer)

	// create the middleware
	p := ginopentracing.OpenTracer([]byte("api-request-"))

	// tell gin to use the middleware
	r.Use(p)

	// Prometheus metrics
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Application routes
	r.GET("/", func(c *gin.Context) {
		os := runtime.GOOS
		arch := runtime.GOARCH
		start := time.Now()
		if span, ok := c.Get("tracing-context"); ok {
			fmt.Printf("%v", span)
			// span.SetBaggageItem("Got OS/ARCH", fmt.Sprintf("%s/%s", os, arch))
		}
		output := fmt.Sprintf("Hello, world! I'm running on %s/%s using Gin Framework!", os, arch)
		c.String(http.StatusOK, output)
		duration := time.Since(start)
		d := float64(duration) / 1000
		rootHits.WithLabelValues(fmt.Sprintf("%d", http.StatusOK)).Inc()
		histogram.WithLabelValues(fmt.Sprintf("%d", http.StatusOK)).Observe(d)
	})
	r.Run(":8082") // listen and serve on 0.0.0.0:8082
}
