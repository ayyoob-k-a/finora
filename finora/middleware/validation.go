package middleware

import (
	"net/http"
	"strings"

	"github.com/ayyoob-k-a/finora/utils"
	"github.com/gin-gonic/gin"
)

// ValidateJSONMiddleware ensures content type is JSON for POST/PUT requests
func ValidateJSONMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			contentType := c.GetHeader("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				c.JSON(http.StatusBadRequest, utils.ErrorResponse("Content-Type must be application/json"))
				c.Abort()
				return
			}
		}
		c.Next()
	})
}

// RequestSizeMiddleware limits request body size
func RequestSizeMiddleware(maxSize int64) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			c.JSON(http.StatusRequestEntityTooLarge, utils.ErrorResponse("Request body too large"))
			c.Abort()
			return
		}
		c.Next()
	})
}

// SecurityHeadersMiddleware adds security headers
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Add security headers
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		c.Next()
	})
}
