package service

import (
	"fmt"

	"github.com/ayyoob-k-a/finora/domain"
	"github.com/ayyoob-k-a/finora/model/dto"
	"gorm.io/gorm"
)

type NotificationService struct {
	db *gorm.DB
}

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{
		db: db,
	}
}

// GetUserNotifications retrieves notifications for a user with pagination and filtering
func (s *NotificationService) GetUserNotifications(userID string, page, limit int, unreadOnly bool) (*dto.NotificationsResponse, error) {
	var notifications []domain.Notification
	var total int64
	var unreadCount int64

	query := s.db.Where("user_id = ?", userID)

	// Apply unread filter if requested
	if unreadOnly {
		query = query.Where("is_read = ?", false)
	}

	// Count total notifications
	err := query.Model(&domain.Notification{}).Count(&total).Error
	if err != nil {
		return nil, err
	}

	// Count unread notifications (always needed for response)
	err = s.db.Model(&domain.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&unreadCount).Error
	if err != nil {
		return nil, err
	}

	// Apply pagination and ordering
	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&notifications).Error
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var notificationResponses []dto.NotificationResponse
	for _, notification := range notifications {
		notificationResponses = append(notificationResponses, dto.NotificationResponse{
			ID:        notification.ID,
			Title:     notification.Title,
			Message:   notification.Message,
			Type:      notification.Type,
			IsRead:    notification.IsRead,
			Metadata:  notification.Metadata,
			CreatedAt: notification.CreatedAt,
		})
	}

	return &dto.NotificationsResponse{
		Notifications: notificationResponses,
		UnreadCount:   int(unreadCount),
	}, nil
}

// MarkNotificationAsRead marks a specific notification as read
func (s *NotificationService) MarkNotificationAsRead(userID, notificationID string) error {
	result := s.db.Model(&domain.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("notification not found")
	}

	return nil
}

// MarkAllNotificationsAsRead marks all user notifications as read
func (s *NotificationService) MarkAllNotificationsAsRead(userID string) error {
	err := s.db.Model(&domain.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error

	return err
}

// DeleteNotification deletes a specific notification
func (s *NotificationService) DeleteNotification(userID, notificationID string) error {
	result := s.db.Where("id = ? AND user_id = ?", notificationID, userID).
		Delete(&domain.Notification{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("notification not found")
	}

	return nil
}

// CreateNotification creates a new notification (for internal use)
func (s *NotificationService) CreateNotification(userID, title, message, notificationType string, metadata map[string]interface{}) (*domain.Notification, error) {
	notification := domain.Notification{
		UserID:   userID,
		Title:    title,
		Message:  message,
		Type:     notificationType,
		IsRead:   false,
		Metadata: metadata,
	}

	err := s.db.Create(&notification).Error
	if err != nil {
		return nil, err
	}

	return &notification, nil
}

// GetUnreadNotificationCount returns the count of unread notifications for a user
func (s *NotificationService) GetUnreadNotificationCount(userID string) (int, error) {
	var count int64
	err := s.db.Model(&domain.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// DeleteAllReadNotifications deletes all read notifications for a user (cleanup)
func (s *NotificationService) DeleteAllReadNotifications(userID string) error {
	err := s.db.Where("user_id = ? AND is_read = ?", userID, true).
		Delete(&domain.Notification{}).Error

	return err
}

// CreateEMIReminderNotifications creates reminders for upcoming EMI payments
func (s *NotificationService) CreateEMIReminderNotifications(userID string) error {
	// Get user's active EMIs that are due within 7 days
	var emis []domain.EMI
	err := s.db.Where("user_id = ? AND is_active = ?", userID, true).
		Find(&emis).Error
	if err != nil {
		return err
	}

	for _, emi := range emis {
		// Check if reminder already exists for this month
		var existingNotification domain.Notification
		err = s.db.Where("user_id = ? AND type = ? AND metadata->>'emi_id' = ?",
			userID, "emi_reminder", emi.ID).First(&existingNotification).Error

		if err == gorm.ErrRecordNotFound {
			// Create new EMI reminder
			metadata := map[string]interface{}{
				"emi_id":    emi.ID,
				"amount":    emi.MonthlyAmount,
				"due_date":  emi.DueDate,
				"emi_title": emi.Title,
			}

			_, err = s.CreateNotification(
				userID,
				"EMI Payment Reminder",
				fmt.Sprintf("Your EMI payment for %s (₹%.2f) is due soon", emi.Title, emi.MonthlyAmount),
				"emi_reminder",
				metadata,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CreateFriendRequestNotification creates a notification for friend requests
func (s *NotificationService) CreateFriendRequestNotification(userID, fromUserID, fromUserName string) error {
	metadata := map[string]interface{}{
		"from_user_id":   fromUserID,
		"from_user_name": fromUserName,
		"type":           "friend_request",
	}

	_, err := s.CreateNotification(
		userID,
		"New Friend Request",
		fmt.Sprintf("%s sent you a friend request", fromUserName),
		"friend_request",
		metadata,
	)

	return err
}

// CreateGroupInviteNotification creates a notification for group invitations
func (s *NotificationService) CreateGroupInviteNotification(userID, groupID, groupName, inviterName string) error {
	metadata := map[string]interface{}{
		"group_id":     groupID,
		"group_name":   groupName,
		"inviter_name": inviterName,
		"type":         "group_invite",
	}

	_, err := s.CreateNotification(
		userID,
		"Group Invitation",
		fmt.Sprintf("%s invited you to join group '%s'", inviterName, groupName),
		"group_invite",
		metadata,
	)

	return err
}
