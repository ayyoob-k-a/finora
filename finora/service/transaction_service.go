package service

import (
	"fmt"
	"time"

	"github.com/ayyoob-k-a/finora/domain"
	"github.com/ayyoob-k-a/finora/model/dto"
	"gorm.io/gorm"
)

type TransactionService struct {
	db *gorm.DB
}

func NewTransactionService(db *gorm.DB) *TransactionService {
	return &TransactionService{
		db: db,
	}
}

// CreateTransaction creates a new transaction
func (s *TransactionService) CreateTransaction(userID string, req dto.CreateTransactionRequest) (*domain.Transaction, error) {
	transaction := domain.Transaction{
		UserID:          userID,
		Type:            req.Type,
		Amount:          req.Amount,
		CategoryID:      req.CategoryID,
		Description:     req.Description,
		TransactionDate: req.TransactionDate,
		IsRecurring:     req.IsRecurring,
	}

	if req.RecurringFrequency != nil {
		transaction.RecurringFrequency = req.RecurringFrequency
	}

	err := s.db.Create(&transaction).Error
	if err != nil {
		return nil, err
	}

	// Load the category information
	err = s.db.Preload("Category").First(&transaction, transaction.ID).Error
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

// GetTransactions retrieves transactions with filtering and pagination
func (s *TransactionService) GetTransactions(userID string, filters dto.TransactionFilters) (*dto.PaginatedTransactions, error) {
	var transactions []domain.Transaction
	var total int64

	query := s.db.Where("user_id = ?", userID)

	// Apply filters
	if filters.Type != "" {
		query = query.Where("type = ?", filters.Type)
	}
	if filters.CategoryID != "" {
		query = query.Where("category_id = ?", filters.CategoryID)
	}
	if !filters.StartDate.IsZero() {
		query = query.Where("transaction_date >= ?", filters.StartDate)
	}
	if !filters.EndDate.IsZero() {
		query = query.Where("transaction_date <= ?", filters.EndDate)
	}
	if filters.Search != "" {
		query = query.Where("description ILIKE ?", "%"+filters.Search+"%")
	}

	// Count total records
	countQuery := query
	err := countQuery.Model(&domain.Transaction{}).Count(&total).Error
	if err != nil {
		return nil, err
	}

	// Apply pagination
	offset := (filters.Page - 1) * filters.Limit
	query = query.Offset(offset).Limit(filters.Limit)

	// Apply ordering
	query = query.Order("transaction_date DESC, created_at DESC")

	// Preload category information
	err = query.Preload("Category").Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((total + int64(filters.Limit) - 1) / int64(filters.Limit))

	// Convert domain transactions to DTO responses
	transactionResponses := make([]dto.TransactionResponse, len(transactions))
	for i, tx := range transactions {
		var categoryResponse *dto.CategoryResponse
		if tx.Category != nil {
			categoryResponse = &dto.CategoryResponse{
				ID:    tx.Category.ID,
				Name:  tx.Category.Name,
				Icon:  tx.Category.Icon,
				Color: tx.Category.Color,
				Type:  tx.Category.Type,
			}
		}

		transactionResponses[i] = dto.TransactionResponse{
			ID:                 tx.ID,
			Type:               tx.Type,
			Amount:             tx.Amount,
			Category:           categoryResponse,
			Description:        tx.Description,
			TransactionDate:    tx.TransactionDate,
			IsRecurring:        tx.IsRecurring,
			RecurringFrequency: tx.RecurringFrequency,
			CreatedAt:          tx.CreatedAt,
		}
	}

	return &dto.PaginatedTransactions{
		Transactions: transactionResponses,
		Pagination: dto.Pagination{
			Page:       filters.Page,
			Limit:      filters.Limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}

// GetTransactionByID retrieves a transaction by ID
func (s *TransactionService) GetTransactionByID(userID, transactionID string) (*domain.Transaction, error) {
	var transaction domain.Transaction

	err := s.db.Preload("Category").
		Where("id = ? AND user_id = ?", transactionID, userID).
		First(&transaction).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, err
	}

	return &transaction, nil
}

// UpdateTransaction updates an existing transaction
func (s *TransactionService) UpdateTransaction(userID, transactionID string, req dto.UpdateTransactionRequest) (*domain.Transaction, error) {
	var transaction domain.Transaction

	// First, find the transaction
	err := s.db.Where("id = ? AND user_id = ?", transactionID, userID).First(&transaction).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, err
	}

	// Update fields
	updateData := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if req.Type != "" {
		updateData["type"] = req.Type
	}
	if req.Amount != 0 {
		updateData["amount"] = req.Amount
	}
	if req.CategoryID != "" {
		updateData["category_id"] = req.CategoryID
	}
	if req.Description != "" {
		updateData["description"] = req.Description
	}
	if !req.TransactionDate.IsZero() {
		updateData["transaction_date"] = req.TransactionDate
	}
	if req.IsRecurring != nil {
		updateData["is_recurring"] = *req.IsRecurring
	}
	if req.RecurringFrequency != nil {
		updateData["recurring_frequency"] = *req.RecurringFrequency
	}

	err = s.db.Model(&transaction).Updates(updateData).Error
	if err != nil {
		return nil, err
	}

	// Reload with category
	err = s.db.Preload("Category").First(&transaction, transaction.ID).Error
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

// DeleteTransaction deletes a transaction
func (s *TransactionService) DeleteTransaction(userID, transactionID string) error {
	result := s.db.Where("id = ? AND user_id = ?", transactionID, userID).Delete(&domain.Transaction{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("transaction not found")
	}

	return nil
}

// GetUserBalance calculates user's current balance
func (s *TransactionService) GetUserBalance(userID string) (float64, error) {
	var income, expense float64

	// Calculate total income
	err := s.db.Model(&domain.Transaction{}).
		Where("user_id = ? AND type = ?", userID, "income").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&income).Error
	if err != nil {
		return 0, err
	}

	// Calculate total expense
	err = s.db.Model(&domain.Transaction{}).
		Where("user_id = ? AND type = ?", userID, "expense").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&expense).Error
	if err != nil {
		return 0, err
	}

	return income - expense, nil
}

// GetRecentTransactions retrieves recent transactions for dashboard
func (s *TransactionService) GetRecentTransactions(userID string, limit int) ([]domain.Transaction, error) {
	var transactions []domain.Transaction

	err := s.db.Preload("Category").
		Where("user_id = ?", userID).
		Order("transaction_date DESC, created_at DESC").
		Limit(limit).
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
