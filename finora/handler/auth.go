package handler

import (
	"net/http"

	"github.com/ayyoob-k-a/finora/domain"
	"github.com/ayyoob-k-a/finora/model/inbound"
	"github.com/ayyoob-k-a/finora/model/response"
	"github.com/ayyoob-k-a/finora/usecase"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase *usecase.Usecase
}

func NewHandler(usecase *usecase.Usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

//if user already signed up return string

func (h *Handler) Signup(ctx *gin.Context) {
	var signupData domain.User
	if err := ctx.ShouldBindJSON(&signupData); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	err := h.usecase.Signup(signupData)
	if err != nil {
		if err.Error() == "user already exists" {
			ctx.JSON(409, gin.H{"error": "User already signed up"})
			return
		}
		ctx.JSON(500, gin.H{"error": "Failed to sign up"})
		return
	}
	ctx.JSON(200, gin.H{"message": "User signed up successfully"})
}

func (h *Handler) Login(c *gin.Context) {
	var req inbound.Login
	if err := c.ShouldBindJSON(&req); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		response.NewCommonResponse(c, "Invalid input", "error", err, http.StatusBadRequest, nil)
		return
	}

	res, err := h.usecase.Login(req)
	if err != nil {
		// c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		response.NewCommonResponse(c, "Login failed", "error", err, http.StatusUnauthorized, nil)
		return
	}

	// c.JSON(http.StatusOK, res)
	response.NewCommonResponse(c, "Login successful", "success", nil, http.StatusOK, res)
}
