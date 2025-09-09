package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ayyoob-k-a/finora/middleware"
	"github.com/ayyoob-k-a/finora/model/dto"
	"github.com/ayyoob-k-a/finora/service"
	"github.com/ayyoob-k-a/finora/utils"
	"github.com/gin-gonic/gin"
)

type GroupHandler struct {
	groupService *service.GroupService
}

func NewGroupHandler(groupService *service.GroupService) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
	}
}

// CreateGroup handles POST /api/groups
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	// Check if service is available
	if h.groupService == nil {
		log.Println("⚠️  GroupService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	var req dto.CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid create group request: %v", err)
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid input format"))
		return
	}

	// Validate group name
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Group name is required"))
		return
	}

	// Validate member IDs
	if len(req.MemberIDs) == 0 {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("At least one member is required"))
		return
	}

	group, err := h.groupService.CreateGroup(userID, req)
	if err != nil {
		log.Printf("Failed to create group for user %s: %v", userID, err)
		if err.Error()[:10] == "member with" { // "member with ID ... not found"
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create group"))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("Group created successfully", group))
}

// GetUserGroups handles GET /api/groups
func (h *GroupHandler) GetUserGroups(c *gin.Context) {
	// Check if service is available
	if h.groupService == nil {
		log.Println("⚠️  GroupService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	groups, err := h.groupService.GetUserGroups(userID)
	if err != nil {
		log.Printf("Failed to get groups for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve groups"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Groups retrieved successfully", groups))
}

// GetGroupDetails handles GET /api/groups/:id
func (h *GroupHandler) GetGroupDetails(c *gin.Context) {
	// Check if service is available
	if h.groupService == nil {
		log.Println("⚠️  GroupService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	groupID := c.Param("id")
	if groupID == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Group ID is required"))
		return
	}

	groupDetails, err := h.groupService.GetGroupDetails(userID, groupID)
	if err != nil {
		log.Printf("Failed to get group details %s for user %s: %v", groupID, userID, err)
		if err.Error() == "access denied: you are not a member of this group" {
			c.JSON(http.StatusForbidden, utils.ErrorResponse("Access denied: you are not a member of this group"))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve group details"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Group details retrieved successfully", groupDetails))
}

// AddGroupExpense handles POST /api/groups/:id/expenses
func (h *GroupHandler) AddGroupExpense(c *gin.Context) {
	// Check if service is available
	if h.groupService == nil {
		log.Println("⚠️  GroupService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	groupID := c.Param("id")
	if groupID == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Group ID is required"))
		return
	}

	var req dto.AddGroupExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid add group expense request: %v", err)
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid input format"))
		return
	}

	// Validate the request
	if err := utils.ValidateAmount(req.Amount); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid amount: "+err.Error()))
		return
	}

	if req.Description == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Description is required"))
		return
	}

	if req.SplitType != "equal" && req.SplitType != "custom" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Split type must be 'equal' or 'custom'"))
		return
	}

	if req.SplitType == "custom" && len(req.Splits) == 0 {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Custom splits are required when split_type is 'custom'"))
		return
	}

	expense, err := h.groupService.AddGroupExpense(userID, groupID, req)
	if err != nil {
		log.Printf("Failed to add group expense to group %s by user %s: %v", groupID, userID, err)
		if err.Error() == "access denied: you are not a member of this group" {
			c.JSON(http.StatusForbidden, utils.ErrorResponse("Access denied: you are not a member of this group"))
			return
		}
		if err.Error() == "split amounts must equal total expense amount" {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to add group expense"))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("Group expense added successfully", expense))
}

// SettleGroupBalances handles POST /api/groups/:id/settle
func (h *GroupHandler) SettleGroupBalances(c *gin.Context) {
	// Check if service is available
	if h.groupService == nil {
		log.Println("⚠️  GroupService not available - database not connected")
		c.JSON(http.StatusServiceUnavailable, utils.ErrorResponse("Database not available. Please try again later or set up database connection."))
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
		return
	}

	groupID := c.Param("id")
	if groupID == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Group ID is required"))
		return
	}

	var req dto.SettleGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid settle group request: %v", err)
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid input format"))
		return
	}

	// Validate settlements
	if len(req.Settlements) == 0 {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("At least one settlement is required"))
		return
	}

	for i, settlement := range req.Settlements {
		if settlement.FromUserID == settlement.ToUserID {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(fmt.Sprintf("Settlement %d: from_user_id and to_user_id cannot be the same", i+1)))
			return
		}
		if err := utils.ValidateAmount(settlement.Amount); err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(fmt.Sprintf("Settlement %d: invalid amount - %s", i+1, err.Error())))
			return
		}
	}

	err := h.groupService.SettleGroupBalances(userID, groupID, req)
	if err != nil {
		log.Printf("Failed to settle group balances for group %s by user %s: %v", groupID, userID, err)
		if err.Error() == "access denied: you are not a member of this group" {
			c.JSON(http.StatusForbidden, utils.ErrorResponse("Access denied: you are not a member of this group"))
			return
		}
		if err.Error()[:10] == "from_user_" || err.Error()[:8] == "to_user_" { // User not in group errors
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to settle group balances"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Group balances settled successfully", nil))
}
