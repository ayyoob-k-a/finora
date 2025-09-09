package handler

import (
	"log"
	"net/http"

	"github.com/ayyoob-k-a/finora/middleware"
	"github.com/ayyoob-k-a/finora/model/dto"
	"github.com/ayyoob-k-a/finora/service"
	"github.com/ayyoob-k-a/finora/utils"
	"github.com/gin-gonic/gin"
)

type EMIHandler struct {
	emiService *service.EMIService
}

func NewEMIHandler(emiService *service.EMIService) *EMIHandler {
	return &EMIHandler{
		emiService: emiService,
	}
}

// CreateEMI handles POST /api/emis
func (h *EMIHandler) CreateEMI(c *gin.Context) {
	// Check if service is available
	if h.emiService == nil {
		log.Println("⚠️  EMIService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	var req dto.CreateEMIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid create EMI request: %v", err)
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid input format"))
		return
	}

	// Validate the request
	if err := utils.ValidateAmount(req.TotalAmount); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid total amount: "+err.Error()))
		return
	}

	if err := utils.ValidateAmount(req.MonthlyAmount); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid monthly amount: "+err.Error()))
		return
	}

	if err := utils.ValidateDueDate(req.DueDate); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid due date: "+err.Error()))
		return
	}

	if req.EndDate.Before(req.StartDate) {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("End date must be after start date"))
		return
	}

	emi, err := h.emiService.CreateEMI(userID, req)
	if err != nil {
		log.Printf("Failed to create EMI for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create EMI"))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("EMI created successfully", emi))
}

// GetUserEMIs handles GET /api/emis
func (h *EMIHandler) GetUserEMIs(c *gin.Context) {
	// Check if service is available
	if h.emiService == nil {
		log.Println("⚠️  EMIService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	emis, err := h.emiService.GetUserEMIs(userID)
	if err != nil {
		log.Printf("Failed to get EMIs for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve EMIs"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("EMIs retrieved successfully", emis))
}

// RecordEMIPayment handles POST /api/emis/:id/payment
func (h *EMIHandler) RecordEMIPayment(c *gin.Context) {
	// Check if service is available
	if h.emiService == nil {
		log.Println("⚠️  EMIService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	emiID := c.Param("id")
	if emiID == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("EMI ID is required"))
		return
	}

	var req dto.CreateEMIPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid EMI payment request: %v", err)
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid input format"))
		return
	}

	// Validate the request
	if err := utils.ValidateAmount(req.Amount); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid payment amount: "+err.Error()))
		return
	}

	payment, err := h.emiService.RecordEMIPayment(userID, emiID, req)
	if err != nil {
		log.Printf("Failed to record EMI payment for user %s, EMI %s: %v", userID, emiID, err)
		if err.Error() == "EMI not found" {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("EMI not found"))
			return
		}
		if err.Error() == "payment for this month already exists" {
			c.JSON(http.StatusConflict, utils.ErrorResponse("Payment for this month already exists"))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to record payment"))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("Payment recorded successfully", payment))
}

// GetEMIPayments handles GET /api/emis/:id/payments
func (h *EMIHandler) GetEMIPayments(c *gin.Context) {
	// Check if service is available
	if h.emiService == nil {
		log.Println("⚠️  EMIService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	emiID := c.Param("id")
	if emiID == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("EMI ID is required"))
		return
	}

	payments, err := h.emiService.GetEMIPayments(userID, emiID)
	if err != nil {
		log.Printf("Failed to get EMI payments for user %s, EMI %s: %v", userID, emiID, err)
		if err.Error() == "EMI not found" {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("EMI not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve payments"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Payments retrieved successfully", payments))
}
