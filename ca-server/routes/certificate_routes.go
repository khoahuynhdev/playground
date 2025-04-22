package routes

import (
	"ca-server/controllers"
	"ca-server/models"

	"github.com/gin-gonic/gin"
)

// SetupCertRoutes registers all cert-related routes
func SetupCertRoutes(router *gin.Engine, store models.Store) {
	certController := controllers.NewCertController()

	// Public user API endpoints
	certGroup := router.Group("/api/certs")
	{
		certGroup.POST("", certController.CreateKey)
		certGroup.POST("/ca", certController.CreateCA)
		certGroup.POST("/server", certController.CreateServerCert)
		certGroup.POST("/client", certController.CreateClientCert)
	}
}
