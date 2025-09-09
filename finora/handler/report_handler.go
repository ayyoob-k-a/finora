package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ayyoob-k-a/finora/middleware"
	"github.com/ayyoob-k-a/finora/service"
	"github.com/ayyoob-k-a/finora/utils"
	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	reportService *service.ReportService
}

func NewReportHandler(reportService *service.ReportService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
	}
}

// GetMonthlyReport handles GET /api/reports/monthly
func (h *ReportHandler) GetMonthlyReport(c *gin.Context) {
	// Check if service is available
	if h.reportService == nil {
		log.Println("⚠️  ReportService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	month := c.Query("month")
	if month == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Month parameter is required (format: YYYY-MM)"))
		return
	}

	report, err := h.reportService.GetMonthlyReport(userID, month)
	if err != nil {
		log.Printf("Failed to get monthly report for user %s, month %s: %v", userID, month, err)
		if err.Error() == "invalid month format, expected YYYY-MM" {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to generate monthly report"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Monthly report generated successfully", report))
}

// GetCategoryReport handles GET /api/reports/category/:id
func (h *ReportHandler) GetCategoryReport(c *gin.Context) {
	// Check if service is available
	if h.reportService == nil {
		log.Println("⚠️  ReportService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	categoryID := c.Param("id")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Category ID is required"))
		return
	}

	report, err := h.reportService.GetCategoryReport(userID, categoryID)
	if err != nil {
		log.Printf("Failed to get category report for user %s, category %s: %v", userID, categoryID, err)
		if err.Error() == "category not found" {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Category not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to generate category report"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Category report generated successfully", report))
}

// GetYearlyReport handles GET /api/reports/yearly
func (h *ReportHandler) GetYearlyReport(c *gin.Context) {
	// Check if service is available
	if h.reportService == nil {
		log.Println("⚠️  ReportService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	yearStr := c.Query("year")
	if yearStr == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Year parameter is required (format: YYYY)"))
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid year format, expected YYYY"))
		return
	}

	// Validate year range
	if year < 2000 || year > 2100 {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Year must be between 2000 and 2100"))
		return
	}

	report, err := h.reportService.GetYearlyReport(userID, year)
	if err != nil {
		log.Printf("Failed to get yearly report for user %s, year %d: %v", userID, year, err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to generate yearly report"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Yearly report generated successfully", report))
}
