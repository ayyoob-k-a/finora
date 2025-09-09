package service

import (
	"github.com/ayyoob-k-a/finora/domain"
	"gorm.io/gorm"
)

type CategoryService struct {
	db *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{
		db: db,
	}
}

// GetAllCategories retrieves all categories
func (s *CategoryService) GetAllCategories() ([]domain.Category, error) {
	var categories []domain.Category

	err := s.db.Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

// GetCategoriesByType retrieves categories by type (income/expense)
func (s *CategoryService) GetCategoriesByType(categoryType string) ([]domain.Category, error) {
	var categories []domain.Category

	err := s.db.Where("type = ?", categoryType).Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

// GetCategoryByID retrieves a category by ID
func (s *CategoryService) GetCategoryByID(id string) (*domain.Category, error) {
	var category domain.Category

	err := s.db.Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}

	return &category, nil
}
