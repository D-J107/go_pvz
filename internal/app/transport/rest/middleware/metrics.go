package middleware

import (
	prometheusMetrics "my_pvz/internal/app/transport/rest/prometheus_metrics"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func PrometheusMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		prometheusMetrics.RequestCount.WithLabelValues(c.Request.Method, path, strconv.Itoa(c.Writer.Status())).Inc()
		prometheusMetrics.RequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)
	}
}
