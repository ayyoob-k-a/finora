package server

import "github.com/gin-gonic/gin"

func InitRouter() *gin.Engine {
	router := gin.Default()


	return router
}
func StartServer(router *gin.Engine) {
	// Start the server on port 8080
	if err := router.Run(":8080"); err != nil {
		panic("failed to start server: " + err.Error())
	}
}