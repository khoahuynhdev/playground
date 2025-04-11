package controllers

import (
	"ca-server/models"
	"ca-server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserController handles user-related API endpoints
type UserController struct {
	store models.Store
}

// NewUserController creates a new user controller with the given store
func NewUserController(store models.Store) *UserController {
	return &UserController{
		store: store,
	}
}

// GetUser returns a user by ID
func (c *UserController) GetUser(ctx *gin.Context) {
	userID := ctx.Param("id")

	user, err := c.store.GetUser(userID)
	if err != nil {
		utils.NotFound(ctx, "User not found")
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// ListUsers returns a list of users
func (c *UserController) ListUsers(ctx *gin.Context) {
	users, err := c.store.ListUsers()
	if err != nil {
		utils.InternalServerError(ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, users)
}

// CreateUser creates a new user
func (c *UserController) CreateUser(ctx *gin.Context) {
	var user models.User

	// Bind JSON request body to user model
	if err := ctx.ShouldBindJSON(&user); err != nil {
		utils.BadRequest(ctx, "Invalid user data", err.Error())
		return
	}

	// Create the user
	if err := c.store.CreateUser(&user); err != nil {
		utils.InternalServerError(ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

// UpdateUser updates an existing user
func (c *UserController) UpdateUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	var user models.User

	// Bind JSON request to user model
	if err := ctx.ShouldBindJSON(&user); err != nil {
		utils.BadRequest(ctx, "Invalid user data", err.Error())
		return
	}

	// Ensure ID in path matches ID in body
	user.ID = userID

	// Update the user
	if err := c.store.UpdateUser(&user); err != nil {
		utils.NotFound(ctx, "User not found")
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// DeleteUser deletes a user
func (c *UserController) DeleteUser(ctx *gin.Context) {
	userID := ctx.Param("id")

	if err := c.store.DeleteUser(userID); err != nil {
		utils.NotFound(ctx, "User not found")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
		"id":      userID,
	})
}
