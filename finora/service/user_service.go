package service

import (
	"fmt"
	"log"
	"time"

	"github.com/ayyoob-k-a/finora/domain"
	"github.com/ayyoob-k-a/finora/model/dto"
	"github.com/ayyoob-k-a/finora/utils"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

// GetProfile retrieves user profile information
func (s *UserService) GetProfile(userID string) (*dto.UserResponse, error) {
	var user domain.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &dto.UserResponse{
		ID:              user.ID,
		Name:            user.Name,
		Phone:           user.Phone,
		Email:           user.Email,
		PhotoURL:        user.PhotoURL,
		DefaultCurrency: user.DefaultCurrency,
		MonthlyIncome:   user.MonthlyIncome,
		CreatedAt:       user.CreatedAt,
	}, nil
}

// UpdateProfile updates user profile information
func (s *UserService) UpdateProfile(userID string, req dto.UpdateProfileRequest) error {
	var user domain.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("database error: %w", err)
	}

	// Update fields if provided
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.PhotoURL != "" {
		user.PhotoURL = req.PhotoURL
	}
	if req.DefaultCurrency != "" {
		if err := utils.ValidateCurrency(req.DefaultCurrency); err != nil {
			return fmt.Errorf("invalid currency: %w", err)
		}
		user.DefaultCurrency = req.DefaultCurrency
	}
	if req.MonthlyIncome > 0 {
		user.MonthlyIncome = req.MonthlyIncome
	}

	user.UpdatedAt = time.Now()

	if err := s.db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	log.Printf("Profile updated for user %s", userID)
	return nil
}

// GetDashboard retrieves dashboard data for the user
func (s *UserService) GetDashboard(userID string) (*dto.DashboardData, error) {
	// Get current month boundaries
	now := time.Now()
	monthStart := utils.GetMonthStart(now)
	monthEnd := utils.GetMonthEnd(now)

	// Calculate total balance (income - expenses for this month)
	var monthlyIncome, monthlyExpense float64

	// Get monthly income
	s.db.Model(&domain.Transaction{}).
		Where("user_id = ? AND type = 'income' AND transaction_date >= ? AND transaction_date <= ?", 
			userID, monthStart, monthEnd).
		Select("COALESCE(SUM(amount), 0)").Scan(&monthlyIncome)

	// Get monthly expenses
	s.db.Model(&domain.Transaction{}).
		Where("user_id = ? AND type = 'expense' AND transaction_date >= ? AND transaction_date <= ?", 
			userID, monthStart, monthEnd).
		Select("COALESCE(SUM(amount), 0)").Scan(&monthlyExpense)

	totalBalance := monthlyIncome - monthlyExpense

	// Get upcoming EMIs (due in next 7 days)
	upcomingEMIs := []dto.UpcomingEMI{}
	nextWeek := now.AddDate(0, 0, 7)
	
	var emis []domain.EMI
	s.db.Where("user_id = ? AND is_active = ? AND due_date <= ?", 
		userID, true, nextWeek.Day()).
		Find(&emis)

	for _, emi := range emis {
		// Calculate next due date
		nextDueDate := time.Date(now.Year(), now.Month(), emi.DueDate, 0, 0, 0, 0, now.Location())
		if nextDueDate.Before(now) {
			nextDueDate = nextDueDate.AddDate(0, 1, 0) // Next month
		}

		if nextDueDate.Before(nextWeek) {
			upcomingEMIs = append(upcomingEMIs, dto.UpcomingEMI{
				ID:      emi.ID,
				Title:   emi.Title,
				Amount:  emi.MonthlyAmount,
				DueDate: nextDueDate,
			})
		}
	}

	// Get recent transactions (last 10)
	var transactions []domain.Transaction
	s.db.Where("user_id = ?", userID).
		Preload("Category").
		Order("transaction_date DESC, created_at DESC").
		Limit(10).
		Find(&transactions)

	recentTransactions := []dto.TransactionResponse{}
	for _, t := range transactions {
		var category *dto.CategoryResponse
		if t.Category != nil {
			category = &dto.CategoryResponse{
				ID:    t.Category.ID,
				Name:  t.Category.Name,
				Icon:  t.Category.Icon,
				Color: t.Category.Color,
				Type:  t.Category.Type,
			}
		}

		recentTransactions = append(recentTransactions, dto.TransactionResponse{
			ID:              t.ID,
			Type:            t.Type,
			Amount:          t.Amount,
			Category:        category,
			Description:     t.Description,
			TransactionDate: t.TransactionDate,
			IsRecurring:     t.IsRecurring,
			RecurringFrequency: t.RecurringFrequency,
			CreatedAt:       t.CreatedAt,
		})
	}

	// Get group summary (placeholder for now)
	groupSummary := dto.GroupSummary{
		TotalOwed:    0,
		TotalLending: 0,
	}

	// TODO: Implement group expense calculations when group features are ready

	return &dto.DashboardData{
		TotalBalance:       totalBalance,
		MonthlyIncome:      monthlyIncome,
		MonthlyExpense:     monthlyExpense,
		UpcomingEMIs:       upcomingEMIs,
		RecentTransactions: recentTransactions,
		GroupSummary:       groupSummary,
	}, nil
}

// GetUserByID retrieves user by ID (internal use)
func (s *UserService) GetUserByID(userID string) (*domain.User, error) {
	var user domain.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByPhone retrieves user by phone number (internal use)
func (s *UserService) GetUserByPhone(phone string) (*domain.User, error) {
	var user domain.User
	if err := s.db.Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUserEmail updates user email
func (s *UserService) UpdateUserEmail(userID, email string) error {
	if err := utils.ValidateEmail(email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}

	// Check if email is already in use
	var count int64
	s.db.Model(&domain.User{}).Where("email = ? AND id != ?", email, userID).Count(&count)
	if count > 0 {
		return fmt.Errorf("email already in use")
	}

	if err := s.db.Model(&domain.User{}).Where("id = ?", userID).
		Update("email", email).Error; err != nil {
		return fmt.Errorf("failed to update email: %w", err)
	}

	return nil
}
