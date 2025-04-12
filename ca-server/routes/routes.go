package routes

import (
	"ca-server/models"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(r *gin.Engine, store models.Store) {
	// Public routes
	r.GET("/", HomeHandler)
	r.GET("/health", HealthCheckHandler)

	// API group
	api := r.Group("/api")
	{
		api.GET("/ping", PingHandler)
	}

	// Setup feature-specific routes
	SetupUserRoutes(r, store)
}

// HomeHandler returns welcome message
func HomeHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome to Gin HTTP Server",
	})
}

// HealthCheckHandler returns server status
func HealthCheckHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

// PingHandler returns pong
func PingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
