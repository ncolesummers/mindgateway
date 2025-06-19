package server

import (
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/company/mindgateway/internal/shared/logging"
)

// Request metrics
var (
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 15),
		},
		[]string{"method", "endpoint", "status"},
	)
	
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)
)

func init() {
	prometheus.MustRegister(requestDuration, requestCounter)
}

// MetricsMiddleware records request metrics
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		c.Next()
		
		duration := time.Since(start).Seconds()
		status := c.Writer.Status()
		
		requestDuration.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			string(rune(status)),
		).Observe(duration)
		
		requestCounter.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			string(rune(status)),
		).Inc()
	}
}

// LoggingMiddleware logs all requests
func LoggingMiddleware(logger *logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		c.Next()
		
		duration := time.Since(start)
		status := c.Writer.Status()
		
		logger.WithFields(map[string]interface{}{
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"status":   status,
			"duration": duration,
			"ip":       c.ClientIP(),
		}).Info("API request")
	}
}

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}

// AuthMiddleware validates authentication tokens
func (s *Server) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization token required"})
			return
		}
		
		// Strip 'Bearer ' prefix if present
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
		
		// Validate token with auth service
		valid, err := s.authClient.ValidateToken(c.Request.Context(), token)
		if err != nil {
			s.logger.WithError(err).Error("Failed to validate token")
			c.AbortWithStatusJSON(500, gin.H{"error": "Internal authentication error"})
			return
		}
		
		if !valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid authorization token"})
			return
		}
		
		c.Next()
	}
}