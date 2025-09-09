package service

import (
	"fmt"

	"github.com/ayyoob-k-a/finora/domain"
	"github.com/ayyoob-k-a/finora/model/dto"
	"gorm.io/gorm"
)

type FriendService struct {
	db *gorm.DB
}

func NewFriendService(db *gorm.DB) *FriendService {
	return &FriendService{
		db: db,
	}
}

// SendFriendRequest sends a friend request to another user
func (s *FriendService) SendFriendRequest(fromUserID string, req dto.SendFriendRequestRequest) (*domain.Friend, error) {
	// Find the target user by phone
	var targetUser domain.User
	err := s.db.Where("phone = ?", req.Phone).First(&targetUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with phone number not found")
		}
		return nil, err
	}

	// Check if they're trying to add themselves
	if targetUser.ID == fromUserID {
		return nil, fmt.Errorf("cannot send friend request to yourself")
	}

	// Check if friend request or friendship already exists
	var existingFriend domain.Friend
	err = s.db.Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
		fromUserID, targetUser.ID, targetUser.ID, fromUserID).First(&existingFriend).Error

	if err == nil {
		if existingFriend.Status == "accepted" {
			return nil, fmt.Errorf("already friends with this user")
		} else if existingFriend.Status == "pending" {
			return nil, fmt.Errorf("friend request already sent or received")
		}
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Create new friend request
	friend := domain.Friend{
		UserID:      fromUserID,
		FriendID:    targetUser.ID,
		Status:      "pending",
		RequestedBy: fromUserID,
	}

	err = s.db.Create(&friend).Error
	if err != nil {
		return nil, err
	}

	// Load the friend user information
	err = s.db.Preload("FriendUser").First(&friend, friend.ID).Error
	if err != nil {
		return nil, err
	}

	return &friend, nil
}

// GetFriendsList retrieves user's friends and pending requests
func (s *FriendService) GetFriendsList(userID string) (*dto.FriendsListResponse, error) {
	var friends []domain.Friend
	var pendingRequests []domain.Friend

	// Get accepted friends
	err := s.db.Preload("FriendUser").
		Where("user_id = ? AND status = ?", userID, "accepted").
		Find(&friends).Error
	if err != nil {
		return nil, err
	}

	// Also get friends where current user is the friend_id
	var reverseFriends []domain.Friend
	err = s.db.Preload("User").
		Where("friend_id = ? AND status = ?", userID, "accepted").
		Find(&reverseFriends).Error
	if err != nil {
		return nil, err
	}

	// Get pending requests received by current user
	err = s.db.Preload("User").
		Where("friend_id = ? AND status = ?", userID, "pending").
		Find(&pendingRequests).Error
	if err != nil {
		return nil, err
	}

	// Convert to response format
	friendResponses := make([]dto.FriendResponse, 0)
	for _, friend := range friends {
		friendResponses = append(friendResponses, dto.FriendResponse{
			ID:           friend.FriendUser.ID,
			Name:         friend.FriendUser.Name,
			Phone:        friend.FriendUser.Phone,
			PhotoURL:     friend.FriendUser.PhotoURL,
			TotalOwed:    0, // TODO: Calculate from transactions
			TotalLending: 0, // TODO: Calculate from transactions
		})
	}

	for _, friend := range reverseFriends {
		friendResponses = append(friendResponses, dto.FriendResponse{
			ID:           friend.User.ID,
			Name:         friend.User.Name,
			Phone:        friend.User.Phone,
			PhotoURL:     friend.User.PhotoURL,
			TotalOwed:    0, // TODO: Calculate from transactions
			TotalLending: 0, // TODO: Calculate from transactions
		})
	}

	pendingRequestResponses := make([]dto.FriendRequestResponse, 0)
	for _, request := range pendingRequests {
		pendingRequestResponses = append(pendingRequestResponses, dto.FriendRequestResponse{
			ID:          request.ID,
			Name:        request.User.Name,
			Phone:       request.User.Phone,
			RequestedBy: request.UserID,
			CreatedAt:   request.CreatedAt,
		})
	}

	return &dto.FriendsListResponse{
		Friends:         friendResponses,
		PendingRequests: pendingRequestResponses,
	}, nil
}

// HandleFriendRequest accepts or rejects a friend request
func (s *FriendService) HandleFriendRequest(userID, requestID string, req dto.HandleFriendRequestRequest) error {
	// Find the friend request
	var friendRequest domain.Friend
	err := s.db.Where("id = ? AND friend_id = ? AND status = ?", requestID, userID, "pending").First(&friendRequest).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("friend request not found")
		}
		return err
	}

	// Update the request status
	if req.Action == "accept" {
		err = s.db.Model(&friendRequest).Update("status", "accepted").Error
	} else if req.Action == "reject" {
		err = s.db.Delete(&friendRequest).Error
	} else {
		return fmt.Errorf("invalid action: must be 'accept' or 'reject'")
	}

	if err != nil {
		return err
	}

	return nil
}

// RemoveFriend removes a friend relationship
func (s *FriendService) RemoveFriend(userID, friendID string) error {
	// Find and delete the friendship (could be in either direction)
	result := s.db.Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
		userID, friendID, friendID, userID).Delete(&domain.Friend{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("friend relationship not found")
	}

	return nil
}

// GetFriendByID gets a specific friend's details
func (s *FriendService) GetFriendByID(userID, friendID string) (*dto.FriendResponse, error) {
	// Check if friendship exists
	var friend domain.Friend
	err := s.db.Where("((user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)) AND status = ?",
		userID, friendID, friendID, userID, "accepted").First(&friend).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("friend not found")
		}
		return nil, err
	}

	// Get the friend's user details
	var friendUser domain.User
	if friend.UserID == userID {
		err = s.db.Where("id = ?", friend.FriendID).First(&friendUser).Error
	} else {
		err = s.db.Where("id = ?", friend.UserID).First(&friendUser).Error
	}

	if err != nil {
		return nil, err
	}

	return &dto.FriendResponse{
		ID:           friendUser.ID,
		Name:         friendUser.Name,
		Phone:        friendUser.Phone,
		PhotoURL:     friendUser.PhotoURL,
		TotalOwed:    0, // TODO: Calculate from transactions
		TotalLending: 0, // TODO: Calculate from transactions
	}, nil
}
