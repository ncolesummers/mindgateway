package handlers

import (
	"net/http"
	"strconv"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// Prometheus metrics
var (
	requestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "mindgateway_request_latency_seconds",
			Help:    "Request latency in seconds",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 15),
		},
		[]string{"model", "endpoint", "status"},
	)
	
	requestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mindgateway_requests_total",
			Help: "Total number of requests",
		},
		[]string{"model", "endpoint", "status"},
	)
	
	workersActive = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "mindgateway_workers_active",
		Help: "Number of active workers",
	})
	
	queueDepth = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "mindgateway_queue_depth",
		Help: "Number of requests in queue",
	})
	
	tokenCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mindgateway_tokens_total",
			Help: "Total number of tokens processed",
		},
		[]string{"model", "type"},
	)
)

func init() {
	// Register metrics with Prometheus
	prometheus.MustRegister(
		requestLatency,
		requestTotal,
		workersActive,
		queueDepth,
		tokenCounter,
	)
}

// RecordRequestMetrics records request metrics
func RecordRequestMetrics(model, endpoint string, status int, startTime time.Time, inputTokens, outputTokens int) {
	duration := time.Since(startTime).Seconds()
	statusStr := strconv.Itoa(status)
	
	requestLatency.WithLabelValues(model, endpoint, statusStr).Observe(duration)
	requestTotal.WithLabelValues(model, endpoint, statusStr).Inc()
	
	tokenCounter.WithLabelValues(model, "input").Add(float64(inputTokens))
	tokenCounter.WithLabelValues(model, "output").Add(float64(outputTokens))
}

// UpdateQueueMetrics updates queue-related metrics
func UpdateQueueMetrics(queueSize int) {
	queueDepth.Set(float64(queueSize))
}

// UpdateWorkerMetrics updates worker-related metrics
func UpdateWorkerMetrics(activeWorkers int) {
	workersActive.Set(float64(activeWorkers))
}

// MetricsHandler returns a handler for custom metrics
func MetricsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get metrics from various subsystems
		c.JSON(http.StatusOK, gin.H{
			"workers": gin.H{
				"active": 0, // TODO: Get actual worker count
			},
			"queue": gin.H{
				"depth": 0, // TODO: Get actual queue depth
			},
			"requests": gin.H{
				"total":   0, // TODO: Get actual request count
				"success": 0,
				"error":   0,
			},
		})
	}
}