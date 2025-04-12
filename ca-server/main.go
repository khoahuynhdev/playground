package main

import (
	"fmt"
	"log"

	"ca-server/config"
	"ca-server/middleware"
	"ca-server/models"
	"ca-server/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.New()

	// Set gin mode
	gin.SetMode(cfg.Mode)

	// Create gin router with default middleware
	r := gin.New()

	// Add custom middleware
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	// Initialize store
	store := models.NewMemoryStore()

	// Setup routes
	routes.SetupRoutes(r, store)

	// Start server
	serverAddr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Starting server on %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
