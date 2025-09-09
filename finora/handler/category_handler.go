package handler

import (
	"log"
	"net/http"

	"github.com/ayyoob-k-a/finora/service"
	"github.com/ayyoob-k-a/finora/utils"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// GetAllCategories handles GET /api/categories
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	// Check if service is available
	if h.categoryService == nil {
		log.Println("⚠️  CategoryService not available - database not connected")

		// Return placeholder data for API-only mode
		placeholderCategories := []map[string]interface{}{
			{"id": "cat-1", "name": "Food & Dining", "type": "expense", "icon": "🍽️"},
			{"id": "cat-2", "name": "Transportation", "type": "expense", "icon": "🚗"},
			{"id": "cat-3", "name": "Shopping", "type": "expense", "icon": "🛒"},
			{"id": "cat-4", "name": "Entertainment", "type": "expense", "icon": "🎬"},
			{"id": "cat-5", "name": "Bills & Utilities", "type": "expense", "icon": "💡"},
			{"id": "cat-6", "name": "Healthcare", "type": "expense", "icon": "🏥"},
			{"id": "cat-7", "name": "Education", "type": "expense", "icon": "📚"},
			{"id": "cat-8", "name": "Salary", "type": "income", "icon": "💼"},
			{"id": "cat-9", "name": "Business", "type": "income", "icon": "🏢"},
			{"id": "cat-10", "name": "Investments", "type": "income", "icon": "📈"},
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Categories retrieved (placeholder data - database not connected)",
			"data":    placeholderCategories,
		})
		return
	}

	categories, err := h.categoryService.GetAllCategories()
	if err != nil {
		log.Printf("Failed to get categories: %v", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve categories"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Categories retrieved successfully", categories))
}
