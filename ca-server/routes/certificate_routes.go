package routes

import (
	"ca-server/controllers"
	"ca-server/middleware"
	"ca-server/models"

	"github.com/gin-gonic/gin"
)

// SetupCertRoutes registers all cert-related routes
func SetupCertRoutes(router *gin.Engine, store models.Store) {
	certController := controllers.NewUserController(store)

	// Public user API endpoints
	userGroup := router.Group("/api/users")
	{
		userGroup.GET("", certController.ListUsers)
		userGroup.GET("/:id", certController.GetUser)
	}

	// Protected user API endpoints - require authentication
	protectedGroup := router.Group("/api/users")
	protectedGroup.Use(middleware.AuthRequired())
	{
		protectedGroup.POST("", certController.CreateUser)
		protectedGroup.PUT("/:id", certController.UpdateUser)
		protectedGroup.DELETE("/:id", certController.DeleteUser)
	}
}
