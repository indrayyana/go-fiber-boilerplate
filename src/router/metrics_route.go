package router

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of response latency (seconds) for HTTP requests",
			Buckets: prometheus.DefBuckets, // Default buckets for duration tracking
		},
		[]string{"method", "path", "status_code"}, // Labels to track request metrics
	)
)

func MetricsRoutes(app *fiber.App) {
	app.Use(trackRequestDuration())
	reg := prometheus.NewRegistry()
	reg.MustRegister(httpRequestDuration)
	reg.MustRegister(collectors.NewGoCollector())
	promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
}

func trackRequestDuration() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Measure duration
		duration := time.Since(start).Seconds()

		// Record duration in the histogram
		httpRequestDuration.WithLabelValues(c.Method(), c.Path(), strconv.Itoa(c.Response().StatusCode())).Observe(duration)

		return err
	}
}
