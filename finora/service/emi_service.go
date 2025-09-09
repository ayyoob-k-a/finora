package service

import (
	"fmt"
	"time"

	"github.com/ayyoob-k-a/finora/domain"
	"github.com/ayyoob-k-a/finora/model/dto"
	"gorm.io/gorm"
)

type EMIService struct {
	db *gorm.DB
}

func NewEMIService(db *gorm.DB) *EMIService {
	return &EMIService{
		db: db,
	}
}

// CreateEMI creates a new EMI plan
func (s *EMIService) CreateEMI(userID string, req dto.CreateEMIRequest) (*domain.EMI, error) {
	emi := domain.EMI{
		UserID:        userID,
		Title:         req.Title,
		TotalAmount:   req.TotalAmount,
		MonthlyAmount: req.MonthlyAmount,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		DueDate:       req.DueDate,
		Description:   req.Description,
		IsActive:      true,
	}

	err := s.db.Create(&emi).Error
	if err != nil {
		return nil, err
	}

	return &emi, nil
}

// GetUserEMIs retrieves all EMIs for a user
func (s *EMIService) GetUserEMIs(userID string) ([]domain.EMI, error) {
	var emis []domain.EMI

	err := s.db.Where("user_id = ?", userID).
		Order("next_due_date ASC").
		Find(&emis).Error
	if err != nil {
		return nil, err
	}

	// Calculate next due dates and remaining months for each EMI
	for i := range emis {
		s.calculateEMIDetails(&emis[i])
	}

	return emis, nil
}

// GetEMIByID retrieves an EMI by ID
func (s *EMIService) GetEMIByID(userID, emiID string) (*domain.EMI, error) {
	var emi domain.EMI

	err := s.db.Where("id = ? AND user_id = ?", emiID, userID).First(&emi).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("EMI not found")
		}
		return nil, err
	}

	s.calculateEMIDetails(&emi)
	return &emi, nil
}

// RecordEMIPayment records a payment for an EMI
func (s *EMIService) RecordEMIPayment(userID, emiID string, req dto.CreateEMIPaymentRequest) (*domain.EMIPayment, error) {
	// First, verify the EMI exists and belongs to the user
	var emi domain.EMI
	err := s.db.Where("id = ? AND user_id = ?", emiID, userID).First(&emi).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("EMI not found")
		}
		return nil, err
	}

	// Check if payment for this month already exists
	var existingPayment domain.EMIPayment
	err = s.db.Where("emi_id = ? AND due_month = ?", emiID, req.DueMonth).First(&existingPayment).Error
	if err == nil {
		return nil, fmt.Errorf("payment for this month already exists")
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Create the payment record
	payment := domain.EMIPayment{
		EMIID:       emiID,
		Amount:      req.Amount,
		PaymentDate: req.PaymentDate,
		DueMonth:    req.DueMonth,
		Notes:       req.Notes,
		IsPaid:      true,
	}

	err = s.db.Create(&payment).Error
	if err != nil {
		return nil, err
	}

	// Update EMI's last payment date
	s.db.Model(&emi).Update("last_payment_date", req.PaymentDate)

	return &payment, nil
}

// GetEMIPayments retrieves payment history for an EMI
func (s *EMIService) GetEMIPayments(userID, emiID string) ([]domain.EMIPayment, error) {
	// First, verify the EMI exists and belongs to the user
	var emi domain.EMI
	err := s.db.Where("id = ? AND user_id = ?", emiID, userID).First(&emi).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("EMI not found")
		}
		return nil, err
	}

	var payments []domain.EMIPayment
	err = s.db.Where("emi_id = ?", emiID).
		Order("due_month DESC").
		Find(&payments).Error
	if err != nil {
		return nil, err
	}

	return payments, nil
}

// GetUpcomingEMIs retrieves EMIs due within the next 30 days
func (s *EMIService) GetUpcomingEMIs(userID string, days int) ([]domain.EMI, error) {
	if days <= 0 {
		days = 30 // Default to 30 days
	}

	var emis []domain.EMI
	cutoffDate := time.Now().AddDate(0, 0, days)

	err := s.db.Where("user_id = ? AND is_active = ? AND next_due_date <= ?",
		userID, true, cutoffDate).
		Order("next_due_date ASC").
		Find(&emis).Error
	if err != nil {
		return nil, err
	}

	// Calculate details for each EMI
	for i := range emis {
		s.calculateEMIDetails(&emis[i])
	}

	return emis, nil
}

// calculateEMIDetails calculates next due date and remaining months for an EMI
func (s *EMIService) calculateEMIDetails(emi *domain.EMI) {
	now := time.Now()

	// Calculate next due date
	year := now.Year()
	month := now.Month()

	// If today's date is past the due date for this month, move to next month
	if now.Day() > emi.DueDate {
		if month == 12 {
			month = 1
			year++
		} else {
			month++
		}
	}

	emi.NextDueDate = time.Date(year, month, emi.DueDate, 0, 0, 0, 0, now.Location())

	// Calculate remaining months
	totalMonths := int(emi.EndDate.Sub(emi.StartDate).Hours() / (24 * 30)) // Approximate
	passedMonths := int(now.Sub(emi.StartDate).Hours() / (24 * 30))        // Approximate
	emi.RemainingMonths = totalMonths - passedMonths

	if emi.RemainingMonths < 0 {
		emi.RemainingMonths = 0
	}

	// Check if EMI should be marked as inactive
	if now.After(emi.EndDate) {
		emi.IsActive = false
		s.db.Model(emi).Update("is_active", false)
	}
}
