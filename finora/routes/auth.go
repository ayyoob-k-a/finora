package routes

import (
	"github.com/ayyoob-k-a/finora/handler"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine, handler *handler.Handler) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/signup", handler.Signup)
		authGroup.POST("/login", handler.Login)
		// router.POST("/resend-otp", handler.ResendOTP)
		// router.POST("/verify-otp", handler.VerifyOTP)


	}
}
