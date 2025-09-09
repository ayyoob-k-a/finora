package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ayyoob-k-a/finora/middleware"
	"github.com/ayyoob-k-a/finora/model/dto"
	"github.com/ayyoob-k-a/finora/service"
	"github.com/ayyoob-k-a/finora/utils"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
}

func NewTransactionHandler(transactionService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// CreateTransaction handles POST /api/transactions
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	// Check if service is available
	if h.transactionService == nil {
		log.Println("⚠️  TransactionService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid create transaction request: %v", err)
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid input format"))
		return
	}

	// Validate the request
	if err := utils.ValidateCreateTransactionRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	transaction, err := h.transactionService.CreateTransaction(userID, req)
	if err != nil {
		log.Printf("Failed to create transaction for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create transaction"))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("Transaction created successfully", transaction))
}

// GetTransactions handles GET /api/transactions
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	// Check if service is available
	if h.transactionService == nil {
		log.Println("⚠️  TransactionService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	// Parse query parameters
	filters := dto.TransactionFilters{
		Page:  1,
		Limit: 20,
	}

	// Parse page
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			filters.Page = page
		}
	}

	// Parse limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			filters.Limit = limit
		}
	}

	// Parse type filter
	filters.Type = c.Query("type")

	// Parse category filter
	filters.CategoryID = c.Query("category_id")

	// Parse search filter
	filters.Search = c.Query("search")

	// Parse date filters
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filters.StartDate = startDate
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			filters.EndDate = endDate
		}
	}

	result, err := h.transactionService.GetTransactions(userID, filters)
	if err != nil {
		log.Printf("Failed to get transactions for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve transactions"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Transactions retrieved successfully", result))
}

// GetTransactionByID handles GET /api/transactions/:id
func (h *TransactionHandler) GetTransactionByID(c *gin.Context) {
	// Check if service is available
	if h.transactionService == nil {
		log.Println("⚠️  TransactionService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	transactionID := c.Param("id")
	if transactionID == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Transaction ID is required"))
		return
	}

	transaction, err := h.transactionService.GetTransactionByID(userID, transactionID)
	if err != nil {
		log.Printf("Failed to get transaction %s for user %s: %v", transactionID, userID, err)
		if err.Error() == "transaction not found" {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Transaction not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve transaction"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Transaction retrieved successfully", transaction))
}

// UpdateTransaction handles PUT /api/transactions/:id
func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	// Check if service is available
	if h.transactionService == nil {
		log.Println("⚠️  TransactionService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	transactionID := c.Param("id")
	if transactionID == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Transaction ID is required"))
		return
	}

	var req dto.UpdateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid update transaction request: %v", err)
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid input format"))
		return
	}

	transaction, err := h.transactionService.UpdateTransaction(userID, transactionID, req)
	if err != nil {
		log.Printf("Failed to update transaction %s for user %s: %v", transactionID, userID, err)
		if err.Error() == "transaction not found" {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Transaction not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update transaction"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Transaction updated successfully", transaction))
}

// DeleteTransaction handles DELETE /api/transactions/:id
func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	// Check if service is available
	if h.transactionService == nil {
		log.Println("⚠️  TransactionService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	transactionID := c.Param("id")
	if transactionID == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Transaction ID is required"))
		return
	}

	err := h.transactionService.DeleteTransaction(userID, transactionID)
	if err != nil {
		log.Printf("Failed to delete transaction %s for user %s: %v", transactionID, userID, err)
		if err.Error() == "transaction not found" {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Transaction not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to delete transaction"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Transaction deleted successfully", nil))
}
