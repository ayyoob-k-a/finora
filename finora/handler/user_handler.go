package handler

import (
	"log"
	"net/http"

	"github.com/ayyoob-k-a/finora/middleware"
	"github.com/ayyoob-k-a/finora/model/dto"
	"github.com/ayyoob-k-a/finora/service"
	"github.com/ayyoob-k-a/finora/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetProfile handles GET /api/user/profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	// Check if service is available
	if h.userService == nil {
		log.Println("⚠️  UserService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	profile, err := h.userService.GetProfile(userID)
	if err != nil {
		log.Printf("Failed to get profile for user %s: %v", userID, err)
		
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("User not found"))
			return
		}
		
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve profile"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Profile retrieved successfully", profile))
}

// UpdateProfile handles PUT /api/user/profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// Check if service is available
	if h.userService == nil {
		log.Println("⚠️  UserService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid update profile request: %v", err)
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid input format"))
		return
	}

	if err := h.userService.UpdateProfile(userID, req); err != nil {
		log.Printf("Failed to update profile for user %s: %v", userID, err)
		
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("User not found"))
			return
		}
		
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Profile updated successfully", nil))
}

// GetDashboard handles GET /api/user/dashboard
func (h *UserHandler) GetDashboard(c *gin.Context) {
	// Check if service is available
	if h.userService == nil {
		log.Println("⚠️  UserService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	dashboard, err := h.userService.GetDashboard(userID)
	if err != nil {
		log.Printf("Failed to get dashboard for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve dashboard data"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Dashboard data retrieved successfully", dashboard))
}
