package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ayyoob-k-a/finora/middleware"
	"github.com/ayyoob-k-a/finora/service"
	"github.com/ayyoob-k-a/finora/utils"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationService *service.NotificationService
}

func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// GetNotifications handles GET /api/notifications
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	// Check if service is available
	if h.notificationService == nil {
		log.Println("⚠️  NotificationService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	// Parse query parameters
	page := 1
	limit := 20
	unreadOnly := false

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if unreadOnlyStr := c.Query("unread_only"); unreadOnlyStr == "true" {
		unreadOnly = true
	}

	notifications, err := h.notificationService.GetUserNotifications(userID, page, limit, unreadOnly)
	if err != nil {
		log.Printf("Failed to get notifications for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve notifications"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Notifications retrieved successfully", notifications))
}

// MarkNotificationAsRead handles PUT /api/notifications/:id/read
func (h *NotificationHandler) MarkNotificationAsRead(c *gin.Context) {
	// Check if service is available
	if h.notificationService == nil {
		log.Println("⚠️  NotificationService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	notificationID := c.Param("id")
	if notificationID == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Notification ID is required"))
		return
	}

	err := h.notificationService.MarkNotificationAsRead(userID, notificationID)
	if err != nil {
		log.Printf("Failed to mark notification %s as read for user %s: %v", notificationID, userID, err)
		if err.Error() == "notification not found" {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Notification not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to mark notification as read"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Notification marked as read successfully", nil))
}

// MarkAllNotificationsAsRead handles PUT /api/notifications/mark-all-read
func (h *NotificationHandler) MarkAllNotificationsAsRead(c *gin.Context) {
	// Check if service is available
	if h.notificationService == nil {
		log.Println("⚠️  NotificationService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	err := h.notificationService.MarkAllNotificationsAsRead(userID)
	if err != nil {
		log.Printf("Failed to mark all notifications as read for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to mark all notifications as read"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("All notifications marked as read successfully", nil))
}

// DeleteNotification handles DELETE /api/notifications/:id
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	// Check if service is available
	if h.notificationService == nil {
		log.Println("⚠️  NotificationService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	notificationID := c.Param("id")
	if notificationID == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Notification ID is required"))
		return
	}

	err := h.notificationService.DeleteNotification(userID, notificationID)
	if err != nil {
		log.Printf("Failed to delete notification %s for user %s: %v", notificationID, userID, err)
		if err.Error() == "notification not found" {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Notification not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to delete notification"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Notification deleted successfully", nil))
}

// GetUnreadCount handles GET /api/notifications/unread-count (bonus endpoint)
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	// Check if service is available
	if h.notificationService == nil {
		log.Println("⚠️  NotificationService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	count, err := h.notificationService.GetUnreadNotificationCount(userID)
	if err != nil {
		log.Printf("Failed to get unread notification count for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to get unread notification count"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Unread notification count retrieved successfully", gin.H{
		"unread_count": count,
	}))
}
