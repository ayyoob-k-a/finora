package db

import (
	"fmt"
	"log"

	"github.com/ayyoob-k-a/finora/configs"
	"github.com/ayyoob-k-a/finora/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg configs.Config) (*gorm.DB, error) {
	// Build connection string using the new config structure with timeouts
	psqlInfo := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=%s connect_timeout=10", 
		cfg.Database.User, 
		cfg.Database.Name, 
		cfg.Database.Password, 
		cfg.Database.Host, 
		cfg.Database.Port,
		cfg.Database.SSLMode)

	// Configure GORM logger
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(postgres.Open(psqlInfo), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database")

	// Auto-migrate all models
	if err := autoMigrateModels(db); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate models: %w", err)
	}

	// Seed default categories
	if err := seedDefaultCategories(db); err != nil {
		log.Printf("Warning: Failed to seed default categories: %v", err)
	}

	return db, nil
}

func autoMigrateModels(db *gorm.DB) error {
	models := []interface{}{
		&domain.User{},
		&domain.OTP{},
		&domain.Category{},
		&domain.Transaction{},
		&domain.EMI{},
		&domain.EMIPayment{},
		&domain.Friend{},
		&domain.Group{},
		&domain.GroupMember{},
		&domain.GroupExpense{},
		&domain.ExpenseSplit{},
		&domain.Notification{},
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate model %T: %w", model, err)
		}
	}

	log.Println("Successfully auto-migrated all models")
	return nil
}

func seedDefaultCategories(db *gorm.DB) error {
	// Check if categories already exist
	var count int64
	db.Model(&domain.Category{}).Where("is_default = ?", true).Count(&count)
	if count > 0 {
		log.Println("Default categories already exist, skipping seed")
		return nil
	}

	defaultCategories := []domain.Category{
		{Name: "Food & Dining", Icon: "🍽️", Color: "#FF6B6B", Type: "expense", IsDefault: true},
		{Name: "Transportation", Icon: "🚗", Color: "#4ECDC4", Type: "expense", IsDefault: true},
		{Name: "Shopping", Icon: "🛍️", Color: "#45B7D1", Type: "expense", IsDefault: true},
		{Name: "Entertainment", Icon: "🎬", Color: "#96CEB4", Type: "expense", IsDefault: true},
		{Name: "Bills & Utilities", Icon: "💡", Color: "#FFEAA7", Type: "expense", IsDefault: true},
		{Name: "Healthcare", Icon: "🏥", Color: "#DDA0DD", Type: "expense", IsDefault: true},
		{Name: "Education", Icon: "📚", Color: "#98D8C8", Type: "expense", IsDefault: true},
		{Name: "Travel", Icon: "✈️", Color: "#F7DC6F", Type: "expense", IsDefault: true},
		{Name: "Personal Care", Icon: "💆", Color: "#BB8FCE", Type: "expense", IsDefault: true},
		{Name: "Gifts & Donations", Icon: "🎁", Color: "#F8C471", Type: "expense", IsDefault: true},
		{Name: "Salary", Icon: "💰", Color: "#58D68D", Type: "income", IsDefault: true},
		{Name: "Business", Icon: "💼", Color: "#5DADE2", Type: "income", IsDefault: true},
		{Name: "Investments", Icon: "📈", Color: "#F4D03F", Type: "income", IsDefault: true},
		{Name: "Others", Icon: "📝", Color: "#85929E", Type: "income", IsDefault: true},
	}

	for _, category := range defaultCategories {
		if err := db.Create(&category).Error; err != nil {
			return fmt.Errorf("failed to create default category %s: %w", category.Name, err)
		}
	}

	log.Println("Successfully seeded default categories")
	return nil
}
