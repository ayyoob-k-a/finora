package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ayyoob-k-a/finora/model/dto"
)

// ValidatePhoneNumber validates phone number in E.164 format
func ValidatePhoneNumber(phone string) error {
	// Basic E.164 format validation
	phoneRegex := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	if !phoneRegex.MatchString(phone) {
		return fmt.Errorf("invalid phone number format, must be in E.164 format (e.g., +1234567890)")
	}
	return nil
}

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// SanitizeInput removes potentially harmful characters
func SanitizeInput(input string) string {
	// Remove potentially harmful characters
	input = strings.TrimSpace(input)

	// Remove SQL injection patterns (basic protection)
	sqlPatterns := []string{"'", "\"", ";", "--", "/*", "*/", "xp_", "sp_"}
	for _, pattern := range sqlPatterns {
		input = strings.ReplaceAll(input, pattern, "")
	}

	return input
}

// ValidateTransactionType validates transaction types
func ValidateTransactionType(transactionType string) error {
	validTypes := map[string]bool{
		"income":  true,
		"expense": true,
		"lend":    true,
		"borrow":  true,
	}

	if !validTypes[transactionType] {
		return fmt.Errorf("invalid transaction type: %s. Must be one of: income, expense, lend, borrow", transactionType)
	}

	return nil
}

// ValidateRecurringFrequency validates recurring frequency
func ValidateRecurringFrequency(frequency string) error {
	if frequency == "" {
		return nil // Optional field
	}

	validFrequencies := map[string]bool{
		"daily":   true,
		"weekly":  true,
		"monthly": true,
		"yearly":  true,
	}

	if !validFrequencies[frequency] {
		return fmt.Errorf("invalid recurring frequency: %s. Must be one of: daily, weekly, monthly, yearly", frequency)
	}

	return nil
}

// ValidateCurrency validates currency code (3-letter ISO code)
func ValidateCurrency(currency string) error {
	if len(currency) != 3 {
		return fmt.Errorf("currency code must be exactly 3 characters")
	}

	// Add more comprehensive currency validation here if needed
	commonCurrencies := map[string]bool{
		"USD": true, "EUR": true, "GBP": true, "JPY": true,
		"CAD": true, "AUD": true, "CHF": true, "CNY": true,
		"INR": true, "BRL": true, "MXN": true, "KRW": true,
	}

	upperCurrency := strings.ToUpper(currency)
	if !commonCurrencies[upperCurrency] {
		// For now, just log a warning but don't reject
		// In production, you might want to validate against a comprehensive list
	}

	return nil
}

// ValidateAmount validates that amount is positive
func ValidateAmount(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}
	return nil
}

// ValidateDueDate validates EMI due date (1-31)
func ValidateDueDate(dueDate int) error {
	if dueDate < 1 || dueDate > 31 {
		return fmt.Errorf("due date must be between 1 and 31")
	}
	return nil
}

// ValidateCreateTransactionRequest validates the transaction creation request
func ValidateCreateTransactionRequest(req dto.CreateTransactionRequest) error {
	// Validate transaction type
	if err := ValidateTransactionType(req.Type); err != nil {
		return err
	}

	// Validate amount
	if err := ValidateAmount(req.Amount); err != nil {
		return err
	}

	// Validate description
	if strings.TrimSpace(req.Description) == "" {
		return fmt.Errorf("description is required")
	}

	// Validate category ID if provided
	if req.CategoryID != nil && *req.CategoryID != "" {
		if !isValidUUID(*req.CategoryID) {
			return fmt.Errorf("invalid category ID format")
		}
	}

	// Validate recurring frequency if recurring is enabled
	if req.IsRecurring && req.RecurringFrequency != nil {
		if err := ValidateRecurringFrequency(*req.RecurringFrequency); err != nil {
			return err
		}
	}

	return nil
}

// ValidateUpdateTransactionRequest validates transaction update request
func ValidateUpdateTransactionRequest(req dto.UpdateTransactionRequest) error {
	// Validate transaction type if provided
	if req.Type != "" {
		if err := ValidateTransactionType(req.Type); err != nil {
			return err
		}
	}

	// Validate amount if provided
	if req.Amount != 0 {
		if err := ValidateAmount(req.Amount); err != nil {
			return err
		}
	}

	// Validate category ID if provided
	if req.CategoryID != "" && !isValidUUID(req.CategoryID) {
		return fmt.Errorf("invalid category ID format")
	}

	// Validate recurring frequency if provided
	if req.RecurringFrequency != nil && *req.RecurringFrequency != "" {
		if err := ValidateRecurringFrequency(*req.RecurringFrequency); err != nil {
			return err
		}
	}

	return nil
}

// isValidUUID checks if a string is a valid UUID
func isValidUUID(uuid string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	return uuidRegex.MatchString(uuid)
}
