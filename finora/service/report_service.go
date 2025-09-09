package service

import (
	"fmt"
	"time"

	"github.com/ayyoob-k-a/finora/domain"
	"github.com/ayyoob-k-a/finora/model/dto"
	"gorm.io/gorm"
)

type ReportService struct {
	db *gorm.DB
}

func NewReportService(db *gorm.DB) *ReportService {
	return &ReportService{
		db: db,
	}
}

// GetMonthlyReport generates a monthly financial report
func (s *ReportService) GetMonthlyReport(userID string, month string) (*dto.MonthlyReportResponse, error) {
	// Parse month (expected format: YYYY-MM)
	startDate, err := time.Parse("2006-01", month)
	if err != nil {
		return nil, fmt.Errorf("invalid month format, expected YYYY-MM")
	}

	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second) // Last second of the month

	// Calculate total income
	var totalIncome float64
	err = s.db.Model(&domain.Transaction{}).
		Where("user_id = ? AND type = ? AND transaction_date >= ? AND transaction_date <= ?",
			userID, "income", startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalIncome).Error
	if err != nil {
		return nil, err
	}

	// Calculate total expense
	var totalExpense float64
	err = s.db.Model(&domain.Transaction{}).
		Where("user_id = ? AND type = ? AND transaction_date >= ? AND transaction_date <= ?",
			userID, "expense", startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalExpense).Error
	if err != nil {
		return nil, err
	}

	// Get category breakdown for expenses
	var categoryBreakdown []dto.CategoryBreakdown
	rows, err := s.db.Table("transactions").
		Select("categories.name as category_name, SUM(transactions.amount) as amount").
		Joins("LEFT JOIN categories ON transactions.category_id = categories.id").
		Where("transactions.user_id = ? AND transactions.type = ? AND transactions.transaction_date >= ? AND transactions.transaction_date <= ?",
			userID, "expense", startDate, endDate).
		Group("categories.id, categories.name").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var breakdown dto.CategoryBreakdown
		err = rows.Scan(&breakdown.CategoryName, &breakdown.Amount)
		if err != nil {
			return nil, err
		}
		// Calculate percentage
		if totalExpense > 0 {
			breakdown.Percentage = (breakdown.Amount / totalExpense) * 100
		}
		categoryBreakdown = append(categoryBreakdown, breakdown)
	}

	// Get daily trend data
	var dailyTrend []dto.DailyTrendData

	// Generate daily trend for the month
	current := startDate
	for current.Before(endDate) || current.Equal(endDate.Truncate(24*time.Hour)) {
		dayStart := current
		dayEnd := current.Add(24*time.Hour - time.Second)

		var dailyIncome, dailyExpense float64

		// Get daily income
		s.db.Model(&domain.Transaction{}).
			Where("user_id = ? AND type = ? AND transaction_date >= ? AND transaction_date <= ?",
				userID, "income", dayStart, dayEnd).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&dailyIncome)

		// Get daily expense
		s.db.Model(&domain.Transaction{}).
			Where("user_id = ? AND type = ? AND transaction_date >= ? AND transaction_date <= ?",
				userID, "expense", dayStart, dayEnd).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&dailyExpense)

		dailyTrend = append(dailyTrend, dto.DailyTrendData{
			Date:    current.Format("2006-01-02"),
			Income:  dailyIncome,
			Expense: dailyExpense,
		})

		current = current.AddDate(0, 0, 1)
	}

	return &dto.MonthlyReportResponse{
		Month:             month,
		TotalIncome:       totalIncome,
		TotalExpense:      totalExpense,
		NetBalance:        totalIncome - totalExpense,
		CategoryBreakdown: categoryBreakdown,
		DailyTrend:        dailyTrend,
	}, nil
}

// GetCategoryReport generates a detailed report for a specific category
func (s *ReportService) GetCategoryReport(userID, categoryID string) (*dto.CategoryReportResponse, error) {
	// Get category information
	var category domain.Category
	err := s.db.Where("id = ?", categoryID).First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("category not found")
		}
		return nil, err
	}

	// Calculate total spent in this category
	var totalSpent float64
	var transactionCount int64

	err = s.db.Model(&domain.Transaction{}).
		Where("user_id = ? AND category_id = ?", userID, categoryID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalSpent).Error
	if err != nil {
		return nil, err
	}

	err = s.db.Model(&domain.Transaction{}).
		Where("user_id = ? AND category_id = ?", userID, categoryID).
		Count(&transactionCount).Error
	if err != nil {
		return nil, err
	}

	// Calculate average per transaction
	averagePerTransaction := float64(0)
	if transactionCount > 0 {
		averagePerTransaction = totalSpent / float64(transactionCount)
	}

	// Get monthly trend data for the last 12 months
	var monthlyTrend []dto.MonthlyTrendData
	now := time.Now()

	for i := 11; i >= 0; i-- {
		monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, -i, 0)
		monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Second)

		var monthlyAmount float64
		s.db.Model(&domain.Transaction{}).
			Where("user_id = ? AND category_id = ? AND transaction_date >= ? AND transaction_date <= ?",
				userID, categoryID, monthStart, monthEnd).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&monthlyAmount)

		monthlyTrend = append(monthlyTrend, dto.MonthlyTrendData{
			Month:  monthStart.Format("2006-01"),
			Amount: monthlyAmount,
		})
	}

	// Get recent transactions in this category
	var recentTransactions []domain.Transaction
	err = s.db.Preload("Category").
		Where("user_id = ? AND category_id = ?", userID, categoryID).
		Order("transaction_date DESC, created_at DESC").
		Limit(10).
		Find(&recentTransactions).Error
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var transactionResponses []dto.TransactionResponse
	for _, tx := range recentTransactions {
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

		transactionResponses = append(transactionResponses, dto.TransactionResponse{
			ID:                 tx.ID,
			Type:               tx.Type,
			Amount:             tx.Amount,
			Category:           categoryResponse,
			Description:        tx.Description,
			TransactionDate:    tx.TransactionDate,
			IsRecurring:        tx.IsRecurring,
			RecurringFrequency: tx.RecurringFrequency,
			CreatedAt:          tx.CreatedAt,
		})
	}

	return &dto.CategoryReportResponse{
		Category: dto.CategoryResponse{
			ID:    category.ID,
			Name:  category.Name,
			Icon:  category.Icon,
			Color: category.Color,
			Type:  category.Type,
		},
		TotalSpent:            totalSpent,
		TransactionCount:      int(transactionCount),
		AveragePerTransaction: averagePerTransaction,
		MonthlyTrend:          monthlyTrend,
		RecentTransactions:    transactionResponses,
	}, nil
}

// GetYearlyReport generates a yearly financial summary
func (s *ReportService) GetYearlyReport(userID string, year int) (*dto.YearlyReportResponse, error) {
	// Calculate year boundaries
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC).Add(-time.Second)

	// Calculate total income for the year
	var totalIncome float64
	err := s.db.Model(&domain.Transaction{}).
		Where("user_id = ? AND type = ? AND transaction_date >= ? AND transaction_date <= ?",
			userID, "income", yearStart, yearEnd).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalIncome).Error
	if err != nil {
		return nil, err
	}

	// Calculate total expense for the year
	var totalExpense float64
	err = s.db.Model(&domain.Transaction{}).
		Where("user_id = ? AND type = ? AND transaction_date >= ? AND transaction_date <= ?",
			userID, "expense", yearStart, yearEnd).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalExpense).Error
	if err != nil {
		return nil, err
	}

	// Get monthly breakdown
	var monthlyBreakdown []dto.MonthlyBreakdownData
	for month := 1; month <= 12; month++ {
		monthStart := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Second)

		var monthlyIncome, monthlyExpense float64

		s.db.Model(&domain.Transaction{}).
			Where("user_id = ? AND type = ? AND transaction_date >= ? AND transaction_date <= ?",
				userID, "income", monthStart, monthEnd).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&monthlyIncome)

		s.db.Model(&domain.Transaction{}).
			Where("user_id = ? AND type = ? AND transaction_date >= ? AND transaction_date <= ?",
				userID, "expense", monthStart, monthEnd).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&monthlyExpense)

		monthlyBreakdown = append(monthlyBreakdown, dto.MonthlyBreakdownData{
			Month:   monthStart.Format("2006-01"),
			Income:  monthlyIncome,
			Expense: monthlyExpense,
			Balance: monthlyIncome - monthlyExpense,
		})
	}

	// Get top categories by spending
	var topCategories []dto.CategoryBreakdown
	rows, err := s.db.Table("transactions").
		Select("categories.name as category_name, SUM(transactions.amount) as amount").
		Joins("LEFT JOIN categories ON transactions.category_id = categories.id").
		Where("transactions.user_id = ? AND transactions.type = ? AND transactions.transaction_date >= ? AND transactions.transaction_date <= ?",
			userID, "expense", yearStart, yearEnd).
		Group("categories.id, categories.name").
		Order("amount DESC").
		Limit(10).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var breakdown dto.CategoryBreakdown
		err = rows.Scan(&breakdown.CategoryName, &breakdown.Amount)
		if err != nil {
			return nil, err
		}
		// Calculate percentage
		if totalExpense > 0 {
			breakdown.Percentage = (breakdown.Amount / totalExpense) * 100
		}
		topCategories = append(topCategories, breakdown)
	}

	return &dto.YearlyReportResponse{
		Year:             year,
		TotalIncome:      totalIncome,
		TotalExpense:     totalExpense,
		NetBalance:       totalIncome - totalExpense,
		MonthlyBreakdown: monthlyBreakdown,
		TopCategories:    topCategories,
		Savings:          totalIncome - totalExpense,
		SavingsRate:      0, // Calculate if needed
	}, nil
}
