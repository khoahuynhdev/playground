package routes

import (
	"ca-server/controllers"
	"ca-server/middleware"
	"ca-server/models"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes registers all user-related routes
func SetupUserRoutes(router *gin.Engine, store models.Store) {
	userController := controllers.NewUserController(store)

	// Public user API endpoints
	userGroup := router.Group("/api/users")
	{
		userGroup.GET("", userController.ListUsers)
		userGroup.GET("/:id", userController.GetUser)
	}

	// Protected user API endpoints - require authentication
	protectedGroup := router.Group("/api/users")
	protectedGroup.Use(middleware.AuthRequired())
	{
		protectedGroup.POST("", userController.CreateUser)
		protectedGroup.PUT("/:id", userController.UpdateUser)
		protectedGroup.DELETE("/:id", userController.DeleteUser)
	}
}
