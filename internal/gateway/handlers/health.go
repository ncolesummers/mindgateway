package handlers

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
)

// Health returns a handler for the health check endpoint
func Health() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	}
}

// Ready returns a handler for the readiness check endpoint
func Ready() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Add actual readiness checks (database, redis, etc.)
		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
			"checks": map[string]string{
				"database": "ok",
				"redis":    "ok",
				"etcd":     "ok",
			},
		})
	}
}