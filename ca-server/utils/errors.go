package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse standardizes API error responses
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// NewErrorResponse creates a standard error response
func NewErrorResponse(status int, message, err string) *ErrorResponse {
	return &ErrorResponse{
		Status:  status,
		Message: message,
		Error:   err,
	}
}

// RespondWithError sends a standardized error response
func RespondWithError(c *gin.Context, status int, message, err string) {
	c.JSON(status, NewErrorResponse(status, message, err))
}

// BadRequest sends a 400 Bad Request error
func BadRequest(c *gin.Context, message string, err string) {
	RespondWithError(c, http.StatusBadRequest, message, err)
}

// Unauthorized sends a 401 Unauthorized error
func Unauthorized(c *gin.Context, message string) {
	RespondWithError(c, http.StatusUnauthorized, message, "")
}

// Forbidden sends a 403 Forbidden error
func Forbidden(c *gin.Context, message string) {
	RespondWithError(c, http.StatusForbidden, message, "")
}

// NotFound sends a 404 Not Found error
func NotFound(c *gin.Context, message string) {
	RespondWithError(c, http.StatusNotFound, message, "")
}

// InternalServerError sends a 500 Internal Server Error
func InternalServerError(c *gin.Context, err string) {
	RespondWithError(c, http.StatusInternalServerError, "Internal server error", err)
}
