package service

import (
	"fmt"
	"time"

	"github.com/ayyoob-k-a/finora/domain"
	"github.com/ayyoob-k-a/finora/model/dto"
	"gorm.io/gorm"
)

type GroupService struct {
	db *gorm.DB
}

func NewGroupService(db *gorm.DB) *GroupService {
	return &GroupService{
		db: db,
	}
}

// CreateGroup creates a new expense group
func (s *GroupService) CreateGroup(userID string, req dto.CreateGroupRequest) (*domain.Group, error) {
	// Verify all member IDs exist and are friends with creator
	for _, memberID := range req.MemberIDs {
		var user domain.User
		err := s.db.Where("id = ?", memberID).First(&user).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, fmt.Errorf("member with ID %s not found", memberID)
			}
			return nil, err
		}

		// Check if they're friends (optional - you might want to allow any user)
		var friendship domain.Friend
		err = s.db.Where("((user_id = ? AND friend_user_id = ?) OR (user_id = ? AND friend_user_id = ?)) AND status = ?",
			userID, memberID, memberID, userID, "accepted").First(&friendship).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		// Note: We're allowing non-friends to be added to groups for flexibility
	}

	// Create the group
	group := domain.Group{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   userID,
	}

	err := s.db.Create(&group).Error
	if err != nil {
		return nil, err
	}

	// Add creator as first member
	creatorMember := domain.GroupMember{
		GroupID: group.ID,
		UserID:  userID,
	}
	err = s.db.Create(&creatorMember).Error
	if err != nil {
		return nil, err
	}

	// Add other members
	for _, memberID := range req.MemberIDs {
		if memberID != userID { // Don't add creator twice
			member := domain.GroupMember{
				GroupID: group.ID,
				UserID:  memberID,
			}
			err = s.db.Create(&member).Error
			if err != nil {
				return nil, err
			}
		}
	}

	// Load the group with members
	err = s.db.Preload("Members").Preload("Members.User").First(&group, group.ID).Error
	if err != nil {
		return nil, err
	}

	return &group, nil
}

// GetUserGroups retrieves all groups for a user
func (s *GroupService) GetUserGroups(userID string) ([]dto.GroupResponse, error) {
	var groupMembers []domain.GroupMember

	err := s.db.Preload("Group").Where("user_id = ?", userID).Find(&groupMembers).Error
	if err != nil {
		return nil, err
	}

	var groupResponses []dto.GroupResponse
	for _, member := range groupMembers {
		group := member.Group

		// Calculate totals for this group
		var totalExpenses float64
		s.db.Model(&domain.GroupExpense{}).
			Where("group_id = ?", group.ID).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&totalExpenses)

		// Calculate user's share (simplified - equal split for now)
		var memberCount int64
		s.db.Model(&domain.GroupMember{}).Where("group_id = ?", group.ID).Count(&memberCount)

		yourShare := float64(0)
		yourBalance := float64(0) // TODO: Implement proper balance calculation
		if memberCount > 0 {
			yourShare = totalExpenses / float64(memberCount)
		}

		groupResponses = append(groupResponses, dto.GroupResponse{
			ID:            group.ID,
			Name:          group.Name,
			Description:   group.Description,
			MemberCount:   int(memberCount),
			TotalExpenses: totalExpenses,
			YourShare:     yourShare,
			YourBalance:   yourBalance,
		})
	}

	return groupResponses, nil
}

// GetGroupDetails retrieves detailed information about a group
func (s *GroupService) GetGroupDetails(userID, groupID string) (*dto.GroupDetailsResponse, error) {
	// Check if user is a member of the group
	var member domain.GroupMember
	err := s.db.Where("group_id = ? AND user_id = ?", groupID, userID).First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("access denied: you are not a member of this group")
		}
		return nil, err
	}

	// Get group information
	var group domain.Group
	err = s.db.Preload("Members").Preload("Members.User").Where("id = ?", groupID).First(&group).Error
	if err != nil {
		return nil, err
	}

	// Get group expenses
	var expenses []domain.GroupExpense
	err = s.db.Preload("Payer").Where("group_id = ?", groupID).
		Order("expense_date DESC").Find(&expenses).Error
	if err != nil {
		return nil, err
	}

	// Calculate total expenses
	var totalExpenses float64
	s.db.Model(&domain.GroupExpense{}).
		Where("group_id = ?", groupID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalExpenses)

	// Convert to response format
	var memberResponses []dto.GroupMemberResponse
	balances := make(map[string]float64)

	for _, member := range group.Members {
		// TODO: Calculate actual balance for each member
		balance := float64(0) // Placeholder
		balances[member.UserID] = balance

		memberResponses = append(memberResponses, dto.GroupMemberResponse{
			ID:       member.UserID,
			Name:     member.User.Name,
			PhotoURL: member.User.PhotoURL,
			Balance:  balance,
		})
	}

	var expenseResponses []dto.GroupExpenseResponse
	for _, expense := range expenses {
		expenseResponses = append(expenseResponses, dto.GroupExpenseResponse{
			ID:          expense.ID,
			Amount:      expense.Amount,
			Description: expense.Description,
			PaidBy:      expense.PaidBy,
			PayerName:   expense.Payer.Name,
			ExpenseDate: expense.ExpenseDate,
			CreatedAt:   expense.CreatedAt,
		})
	}

	return &dto.GroupDetailsResponse{
		ID:            group.ID,
		Name:          group.Name,
		Description:   group.Description,
		Members:       memberResponses,
		Expenses:      expenseResponses,
		TotalExpenses: totalExpenses,
		Balances:      balances,
	}, nil
}

// AddGroupExpense adds an expense to a group
func (s *GroupService) AddGroupExpense(userID, groupID string, req dto.AddGroupExpenseRequest) (*domain.GroupExpense, error) {
	// Check if user is a member of the group
	var member domain.GroupMember
	err := s.db.Where("group_id = ? AND user_id = ?", groupID, userID).First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("access denied: you are not a member of this group")
		}
		return nil, err
	}

	// Create the expense
	expense := domain.GroupExpense{
		GroupID:     groupID,
		Amount:      req.Amount,
		Description: req.Description,
		PaidBy:      userID,
		ExpenseDate: req.ExpenseDate,
	}

	err = s.db.Create(&expense).Error
	if err != nil {
		return nil, err
	}

	// Handle splits
	if req.SplitType == "equal" {
		// Get all group members
		var members []domain.GroupMember
		err = s.db.Where("group_id = ?", groupID).Find(&members).Error
		if err != nil {
			return nil, err
		}

		// Create equal splits
		splitAmount := req.Amount / float64(len(members))
		for _, member := range members {
			split := domain.ExpenseSplit{
				ExpenseID: expense.ID,
				UserID:    member.UserID,
				Amount:    splitAmount,
			}
			err = s.db.Create(&split).Error
			if err != nil {
				return nil, err
			}
		}
	} else if req.SplitType == "custom" {
		// Create custom splits
		totalSplitAmount := float64(0)
		for _, split := range req.Splits {
			totalSplitAmount += split.Amount
		}

		// Validate that splits add up to total amount
		if totalSplitAmount != req.Amount {
			return nil, fmt.Errorf("split amounts must equal total expense amount")
		}

		for _, splitReq := range req.Splits {
			split := domain.ExpenseSplit{
				ExpenseID: expense.ID,
				UserID:    splitReq.UserID,
				Amount:    splitReq.Amount,
			}
			err = s.db.Create(&split).Error
			if err != nil {
				return nil, err
			}
		}
	}

	// Load the expense with relations
	err = s.db.Preload("PaidByUser").First(&expense, expense.ID).Error
	if err != nil {
		return nil, err
	}

	return &expense, nil
}

// SettleGroupBalances settles balances between group members
func (s *GroupService) SettleGroupBalances(userID, groupID string, req dto.SettleGroupRequest) error {
	// Check if user is a member of the group
	var member domain.GroupMember
	err := s.db.Where("group_id = ? AND user_id = ?", groupID, userID).First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("access denied: you are not a member of this group")
		}
		return err
	}

	// Process each settlement
	for _, settlement := range req.Settlements {
		// Verify both users are in the group
		var fromMember, toMember domain.GroupMember
		err = s.db.Where("group_id = ? AND user_id = ?", groupID, settlement.FromUserID).First(&fromMember).Error
		if err != nil {
			return fmt.Errorf("from_user_id %s is not a member of this group", settlement.FromUserID)
		}

		err = s.db.Where("group_id = ? AND user_id = ?", groupID, settlement.ToUserID).First(&toMember).Error
		if err != nil {
			return fmt.Errorf("to_user_id %s is not a member of this group", settlement.ToUserID)
		}

		// Create a settlement expense (negative amount to balance things out)
		expenseDate := settlement.ExpenseDate
		if expenseDate.IsZero() {
			expenseDate = time.Now()
		}

		description := settlement.Description
		if description == "" {
			description = "Settlement payment"
		}

		settlementExpense := domain.GroupExpense{
			GroupID:     groupID,
			Amount:      settlement.Amount,
			Description: description,
			PaidBy:      settlement.FromUserID,
			ExpenseDate: expenseDate,
		}

		err = s.db.Create(&settlementExpense).Error
		if err != nil {
			return err
		}

		// Create splits for the settlement
		fromSplit := domain.ExpenseSplit{
			ExpenseID: settlementExpense.ID,
			UserID:    settlement.FromUserID,
			Amount:    -settlement.Amount, // Negative because they're paying
		}

		toSplit := domain.ExpenseSplit{
			ExpenseID: settlementExpense.ID,
			UserID:    settlement.ToUserID,
			Amount:    settlement.Amount, // Positive because they're receiving
		}

		err = s.db.Create(&fromSplit).Error
		if err != nil {
			return err
		}

		err = s.db.Create(&toSplit).Error
		if err != nil {
			return err
		}
	}

	return nil
}
