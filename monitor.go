package monitor

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"strings"
	"time"
)

var (
	_ResponseTime    *prometheus.HistogramVec
	_RequestsCounter *prometheus.CounterVec
)

const minus = "-"
const underLine = "_"

// Init  目前只用到了_RequestsCounter 和 _ResponseTime
func Init(namespace, subsystem string) {
	namespace = strings.Replace(namespace, minus, underLine, -1)
	subsystem = strings.Replace(subsystem, minus, underLine, -1)

	_RequestsCounter = func() *prometheus.CounterVec {
		vec := prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "requests_counter",
				Help:      "number of requests",
			},
			[]string{"api", "code"},
		)
		prometheus.MustRegister(vec)
		return vec
	}()
	_ResponseTime = func() *prometheus.HistogramVec {
		vec := prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "response_time",
				Help:      "Response time of requests",
				Buckets: []float64{
					0, 0.1, 0.2, 0.3, 0.5, 1, 5, 10, 60, 300, 1800, 3600,
				}, // time.second
			},
			[]string{"api", "code"},
		)
		prometheus.MustRegister(vec)
		return vec
	}()
}
func PrometheusMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		elapsed := time.Since(start).Seconds()

		path := c.FullPath()
		status := strconv.Itoa(c.Writer.Status())

		_RequestsCounter.WithLabelValues(path, status).Inc()
		_ResponseTime.WithLabelValues(path, status).Observe(elapsed)

	}
}
