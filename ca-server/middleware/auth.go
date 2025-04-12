package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthRequired checks for valid authentication
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization header
		authHeader := c.GetHeader("Authorization")

		// Check if authorization header exists and has valid format
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: missing or invalid authorization token",
			})
			return
		}

		// Extract token (in production, validate JWT/OAuth token here)
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: empty token",
			})
			return
		}

		// For demonstration purposes, accept any non-empty token
		// In production, validate against a real auth system

		// Set user ID in context (extracted from token in production)
		c.Set("userID", "sample-user-id")

		c.Next()
	}
}
