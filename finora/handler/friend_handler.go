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

type FriendHandler struct {
	friendService *service.FriendService
}

func NewFriendHandler(friendService *service.FriendService) *FriendHandler {
	return &FriendHandler{
		friendService: friendService,
	}
}

// SendFriendRequest handles POST /api/friends/request
func (h *FriendHandler) SendFriendRequest(c *gin.Context) {
	// Check if service is available
	if h.friendService == nil {
		log.Println("⚠️  FriendService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	var req dto.SendFriendRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid send friend request: %v", err)
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid input format"))
		return
	}

	// Validate phone number
	if err := utils.ValidatePhoneNumber(req.Phone); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid phone number: "+err.Error()))
		return
	}

	friend, err := h.friendService.SendFriendRequest(userID, req)
	if err != nil {
		log.Printf("Failed to send friend request from user %s: %v", userID, err)
		if err.Error() == "user with phone number not found" {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("User with this phone number not found"))
			return
		}
		if err.Error() == "cannot send friend request to yourself" ||
			err.Error() == "already friends with this user" ||
			err.Error() == "friend request already sent or received" {
			c.JSON(http.StatusConflict, utils.ErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to send friend request"))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("Friend request sent successfully", friend))
}

// GetFriendsList handles GET /api/friends
func (h *FriendHandler) GetFriendsList(c *gin.Context) {
	// Check if service is available
	if h.friendService == nil {
		log.Println("⚠️  FriendService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	friendsList, err := h.friendService.GetFriendsList(userID)
	if err != nil {
		log.Printf("Failed to get friends list for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve friends list"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Friends list retrieved successfully", friendsList))
}

// HandleFriendRequest handles PUT /api/friends/request/:id
func (h *FriendHandler) HandleFriendRequest(c *gin.Context) {
	// Check if service is available
	if h.friendService == nil {
		log.Println("⚠️  FriendService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	requestID := c.Param("id")
	if requestID == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Request ID is required"))
		return
	}

	var req dto.HandleFriendRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid handle friend request: %v", err)
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid input format"))
		return
	}

	// Validate action
	if req.Action != "accept" && req.Action != "reject" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Action must be 'accept' or 'reject'"))
		return
	}

	err := h.friendService.HandleFriendRequest(userID, requestID, req)
	if err != nil {
		log.Printf("Failed to handle friend request %s for user %s: %v", requestID, userID, err)
		if err.Error() == "friend request not found" {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Friend request not found"))
			return
		}
		if err.Error() == "invalid action: must be 'accept' or 'reject'" {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to handle friend request"))
		return
	}

	action := "accepted"
	if req.Action == "reject" {
		action = "rejected"
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Friend request "+action+" successfully", nil))
}

// RemoveFriend handles DELETE /api/friends/:id
func (h *FriendHandler) RemoveFriend(c *gin.Context) {
	// Check if service is available
	if h.friendService == nil {
		log.Println("⚠️  FriendService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	friendID := c.Param("id")
	if friendID == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Friend ID is required"))
		return
	}

	err := h.friendService.RemoveFriend(userID, friendID)
	if err != nil {
		log.Printf("Failed to remove friend %s for user %s: %v", friendID, userID, err)
		if err.Error() == "friend relationship not found" {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Friend not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to remove friend"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Friend removed successfully", nil))
}
