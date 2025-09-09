package middleware

import (
	"net/http"
	"strings"

	"github.com/ayyoob-k-a/finora/configs"
	"github.com/ayyoob-k-a/finora/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware(config configs.Config) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Authorization header required"))
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid authorization header format"))
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Validate JWT token
		claims, err := utils.ValidateJWT(tokenString, config.JWT.Secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid or expired token"))
			c.Abort()
			return
		}

		// Store user information in context
		c.Set("user_id", claims.UserID)
		c.Set("phone", claims.Phone)
		
		c.Next()
	})
}

// OptionalAuthMiddleware validates JWT tokens but doesn't block if missing
func OptionalAuthMiddleware(config configs.Config) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) == 2 && tokenParts[0] == "Bearer" {
			tokenString := tokenParts[1]

			// Validate JWT token
			if claims, err := utils.ValidateJWT(tokenString, config.JWT.Secret); err == nil {
				// Store user information in context if valid
				c.Set("user_id", claims.UserID)
				c.Set("phone", claims.Phone)
			}
		}
		
		c.Next()
	})
}

// GetUserID extracts user ID from context
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	
	userIDStr, ok := userID.(string)
	return userIDStr, ok
}

// GetPhone extracts phone from context
func GetPhone(c *gin.Context) (string, bool) {
	phone, exists := c.Get("phone")
	if !exists {
		return "", false
	}
	
	phoneStr, ok := phone.(string)
	return phoneStr, ok
}
