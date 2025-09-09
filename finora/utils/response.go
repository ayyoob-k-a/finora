package utils

import (
	"fmt"
	"math"
	"time"
)

// Response utilities
func SuccessResponse(message string, data interface{}) map[string]interface{} {
	response := map[string]interface{}{
		"success": true,
	}
	
	if message != "" {
		response["message"] = message
	}
	
	if data != nil {
		response["data"] = data
	}
	
	return response
}

func ErrorResponse(message string) map[string]interface{} {
	return map[string]interface{}{
		"success": false,
		"error":   message,
	}
}

func ErrorResponseWithCode(message string, code int) map[string]interface{} {
	return map[string]interface{}{
		"success": false,
		"error":   message,
		"code":    code,
	}
}

func PaginatedResponse(data interface{}, total, page, limit int) map[string]interface{} {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	
	return map[string]interface{}{
		"success":     true,
		"data":        data,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	}
}

// Time utilities
func ParseDateString(dateStr string) (time.Time, error) {
	// Try different date formats
	formats := []string{
		"2006-01-02",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05-07:00",
		time.RFC3339,
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}
	
	return time.Time{}, fmt.Errorf("invalid date format: %s", dateStr)
}

func GetMonthStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func GetMonthEnd(t time.Time) time.Time {
	return GetMonthStart(t).AddDate(0, 1, 0).Add(-time.Second)
}

// Pagination utilities
func CalculatePagination(page, limit, total int) (int, int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100 // Max limit
	}
	
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	
	return offset, limit, totalPages
}

// Currency utilities
func FormatCurrency(amount float64, currency string) string {
	return fmt.Sprintf("%.2f %s", amount, currency)
}
