package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ayyoob-k-a/finora/configs"
	"github.com/ayyoob-k-a/finora/db"
	"github.com/ayyoob-k-a/finora/handler"
	"github.com/ayyoob-k-a/finora/middleware"
	"github.com/ayyoob-k-a/finora/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("=== FINORA API STARTING ===")
	log.Println("Step 1: Application started")

	// Load environment variables
	log.Println("Step 2: Loading environment variables...")
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using environment variables")
	}
	log.Println("Step 3: Environment variables loaded")

	// Get configuration
	config := configs.GetConfig()
	log.Printf("Starting Finora API server...")
	log.Printf("Server configuration: Port=%s, Mode=%s", config.Server.Port, config.Server.Mode)

	// Initialize database
	log.Printf("Attempting to connect to database with config: Host=%s, Port=%s, User=%s, Name=%s",
		config.Database.Host, config.Database.Port, config.Database.User, config.Database.Name)

	database, err := db.InitDB(config)
	if err != nil {
		log.Printf("⚠️  Database connection failed: %v", err)
		log.Printf("🔄 Starting in API-only mode without database features...")
		log.Printf("💡 To use full features, please set up PostgreSQL or use Docker Compose")

		// Continue without database for basic API testing
		database = nil
	} else {
		log.Println("✅ Database connection successful!")
	}

	// Initialize services (handle nil database gracefully)
	var authService *service.AuthService
	var userService *service.UserService
	var categoryService *service.CategoryService
	var transactionService *service.TransactionService
	var emiService *service.EMIService
	var friendService *service.FriendService
	var groupService *service.GroupService
	var reportService *service.ReportService
	var notificationService *service.NotificationService

	if database != nil {
		authService = service.NewAuthService(database, config)
		userService = service.NewUserService(database)
		categoryService = service.NewCategoryService(database)
		transactionService = service.NewTransactionService(database)
		emiService = service.NewEMIService(database)
		friendService = service.NewFriendService(database)
		groupService = service.NewGroupService(database)
		reportService = service.NewReportService(database)
		notificationService = service.NewNotificationService(database)
		log.Println("✅ All services initialized with database connection")
	} else {
		log.Println("⚠️  Services initialized without database - limited functionality")
	}

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	emiHandler := handler.NewEMIHandler(emiService)
	friendHandler := handler.NewFriendHandler(friendService)
	groupHandler := handler.NewGroupHandler(groupService)
	reportHandler := handler.NewReportHandler(reportService)
	notificationHandler := handler.NewNotificationHandler(notificationService)

	// Initialize Gin router
	if config.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.SecurityHeadersMiddleware())
	router.Use(middleware.ValidateJSONMiddleware())
	router.Use(middleware.RequestSizeMiddleware(10 << 20)) // 10MB limit

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		status := "healthy"
		dbStatus := "disconnected"
		if database != nil {
			dbStatus = "connected"
		}

		c.JSON(200, gin.H{
			"status":    status,
			"version":   "1.0.0",
			"database":  dbStatus,
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		})
	})

	// Setup routes
	setupRoutes(router, config, authHandler, userHandler, categoryHandler, transactionHandler, emiHandler, friendHandler, groupHandler, reportHandler, notificationHandler)

	// Start server
	port := ":" + config.Server.Port
	log.Printf("Starting server on port %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRoutes(router *gin.Engine, config configs.Config,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	categoryHandler *handler.CategoryHandler,
	transactionHandler *handler.TransactionHandler,
	emiHandler *handler.EMIHandler,
	friendHandler *handler.FriendHandler,
	groupHandler *handler.GroupHandler,
	reportHandler *handler.ReportHandler,
	notificationHandler *handler.NotificationHandler) {
	// API v1 routes
	api := router.Group("/api")

	// Authentication routes (no auth required)
	auth := api.Group("/auth")
	{
		// Rate limiting for OTP requests
		auth.POST("/send-otp", middleware.GlobalRateLimiter.OTPRateLimitMiddleware(3), authHandler.SendOTP)
		auth.POST("/verify-otp", authHandler.VerifyOTP)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	// User routes (authentication required)
	user := api.Group("/user")
	user.Use(middleware.AuthMiddleware(config))
	{
		user.GET("/profile", userHandler.GetProfile)
		user.PUT("/profile", userHandler.UpdateProfile)
		user.GET("/dashboard", userHandler.GetDashboard)
	}

	// Categories routes (authentication required)
	categories := api.Group("/categories")
	categories.Use(middleware.AuthMiddleware(config))
	{
		categories.GET("", categoryHandler.GetAllCategories)
	}

	// Transactions routes (authentication required)
	transactions := api.Group("/transactions")
	transactions.Use(middleware.AuthMiddleware(config))
	{
		transactions.POST("", transactionHandler.CreateTransaction)
		transactions.GET("", transactionHandler.GetTransactions)
		transactions.GET("/:id", transactionHandler.GetTransactionByID)
		transactions.PUT("/:id", transactionHandler.UpdateTransaction)
		transactions.DELETE("/:id", transactionHandler.DeleteTransaction)
	}

	// EMI routes (authentication required)
	emis := api.Group("/emis")
	emis.Use(middleware.AuthMiddleware(config))
	{
		emis.POST("", emiHandler.CreateEMI)
		emis.GET("", emiHandler.GetUserEMIs)
		emis.POST("/:id/payment", emiHandler.RecordEMIPayment)
		emis.GET("/:id/payments", emiHandler.GetEMIPayments)
	}

	// Friends routes (authentication required)
	friends := api.Group("/friends")
	friends.Use(middleware.AuthMiddleware(config))
	{
		friends.POST("/request", friendHandler.SendFriendRequest)
		friends.GET("", friendHandler.GetFriendsList)
		friends.PUT("/request/:id", friendHandler.HandleFriendRequest)
		friends.DELETE("/:id", friendHandler.RemoveFriend)
	}

	// Groups routes (authentication required)
	groups := api.Group("/groups")
	groups.Use(middleware.AuthMiddleware(config))
	{
		groups.POST("", groupHandler.CreateGroup)
		groups.GET("", groupHandler.GetUserGroups)
		groups.GET("/:id", groupHandler.GetGroupDetails)
		groups.POST("/:id/expenses", groupHandler.AddGroupExpense)
		groups.POST("/:id/settle", groupHandler.SettleGroupBalances)
	}

	// Reports routes (authentication required)
	reports := api.Group("/reports")
	reports.Use(middleware.AuthMiddleware(config))
	{
		reports.GET("/monthly", reportHandler.GetMonthlyReport)
		reports.GET("/category/:id", reportHandler.GetCategoryReport)
		reports.GET("/yearly", reportHandler.GetYearlyReport)
	}

	// Notifications routes (authentication required)
	notifications := api.Group("/notifications")
	notifications.Use(middleware.AuthMiddleware(config))
	{
		notifications.GET("", notificationHandler.GetNotifications)
		notifications.PUT("/:id/read", notificationHandler.MarkNotificationAsRead)
		notifications.PUT("/mark-all-read", notificationHandler.MarkAllNotificationsAsRead)
		notifications.DELETE("/:id", notificationHandler.DeleteNotification)
		notifications.GET("/unread-count", notificationHandler.GetUnreadCount) // Bonus endpoint
	}
}
