package response

import "github.com/gin-gonic/gin"

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type AuthResponse struct {
	UserID       int    `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type CommonResponse struct {
	Message string      `json:"message"`
	Error   string      `json:"error,omitempty"`
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
}

func NewCommonResponse(ctx *gin.Context, message string, status string, err error, statusCode int, data interface{}) {
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	} else {
		errMsg = ""
	}
	ctx.JSON(statusCode, &CommonResponse{
		Message: message,
		Status:  status,
		Error:   errMsg,
		Data:    data,
	},
	)
}
