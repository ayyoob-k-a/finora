package handler

import (
	"log"
	"net/http"

	"github.com/ayyoob-k-a/finora/model/dto"
	"github.com/ayyoob-k-a/finora/service"
	"github.com/ayyoob-k-a/finora/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// SendOTP handles POST /api/auth/send-otp
func (h *AuthHandler) SendOTP(c *gin.Context) {
	// Check if service is available
	if h.authService == nil {
		log.Println("⚠️  AuthService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	var req dto.SendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid send OTP request: %v", err)
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid input format"))
		return
	}

	// Validate that at least one contact method is provided
	if req.Phone == "" && req.Email == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Either phone or email is required"))
		return
	}

	if err := h.authService.SendOTP(req); err != nil {
		log.Printf("Failed to send OTP: %v", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to send OTP. Please try again."))
		return
	}

	response := dto.AuthResponse{
		Success:   true,
		Message:   "OTP sent successfully",
		ExpiresIn: 300, // 5 minutes
	}

	c.JSON(http.StatusOK, response)
}

// VerifyOTP handles POST /api/auth/verify-otp
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	// Check if service is available
	if h.authService == nil {
		log.Println("⚠️  AuthService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	var req dto.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid verify OTP request: %v", err)
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid input format"))
		return
	}

	response, err := h.authService.VerifyOTP(req)
	if err != nil {
		log.Printf("OTP verification failed: %v", err)
		
		// Determine appropriate status code based on error
		statusCode := http.StatusBadRequest
		if err.Error() == "invalid or expired OTP" || err.Error() == "OTP has expired" {
			statusCode = http.StatusUnauthorized
		}
		
		c.JSON(statusCode, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken handles POST /api/auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Check if service is available
	if h.authService == nil {
		log.Println("⚠️  AuthService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid refresh token request: %v", err)
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid input format"))
		return
	}

	response, err := h.authService.RefreshToken(req)
	if err != nil {
		log.Printf("Token refresh failed: %v", err)
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid or expired refresh token"))
		return
	}

	c.JSON(http.StatusOK, response)
}
