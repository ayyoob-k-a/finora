package dto

import (
	"time"
)

// Authentication DTOs
type SendOTPRequest struct {
	Phone string `json:"phone" binding:"omitempty" validate:"e164"`
	Email string `json:"email" binding:"omitempty" validate:"email"`
}

type VerifyOTPRequest struct {
	Phone string `json:"phone" binding:"required" validate:"e164"`
	OTP   string `json:"otp" binding:"required" validate:"len=6"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AuthResponse struct {
	Success   bool          `json:"success"`
	Token     string        `json:"token,omitempty"`
	User      *UserResponse `json:"user,omitempty"`
	Message   string        `json:"message,omitempty"`
	ExpiresIn int           `json:"expires_in,omitempty"`
}

// User DTOs
type UserResponse struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Phone           string    `json:"phone"`
	Email           string    `json:"email"`
	PhotoURL        string    `json:"photo_url"`
	DefaultCurrency string    `json:"default_currency"`
	MonthlyIncome   float64   `json:"monthly_income"`
	IsNewUser       bool      `json:"is_new_user,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type UpdateProfileRequest struct {
	Name            string  `json:"name" binding:"omitempty" validate:"max=255"`
	PhotoURL        string  `json:"photo_url" binding:"omitempty" validate:"url"`
	DefaultCurrency string  `json:"default_currency" binding:"omitempty" validate:"len=3"`
	MonthlyIncome   float64 `json:"monthly_income" binding:"omitempty" validate:"gte=0"`
}

type DashboardResponse struct {
	Success bool          `json:"success"`
	Data    DashboardData `json:"data"`
}

type DashboardData struct {
	TotalBalance       float64               `json:"total_balance"`
	MonthlyIncome      float64               `json:"monthly_income"`
	MonthlyExpense     float64               `json:"monthly_expense"`
	UpcomingEMIs       []UpcomingEMI         `json:"upcoming_emis"`
	RecentTransactions []TransactionResponse `json:"recent_transactions"`
	GroupSummary       GroupSummary          `json:"group_summary"`
}

type UpcomingEMI struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Amount  float64   `json:"amount"`
	DueDate time.Time `json:"due_date"`
}

type GroupSummary struct {
	TotalOwed    float64 `json:"total_owed"`
	TotalLending float64 `json:"total_lending"`
}

// Transaction DTOs
type CreateTransactionRequest struct {
	Type               string    `json:"type" binding:"required" validate:"oneof=income expense lend borrow"`
	Amount             float64   `json:"amount" binding:"required" validate:"gt=0"`
	CategoryID         *string   `json:"category_id" binding:"omitempty" validate:"uuid"`
	Description        string    `json:"description" binding:"omitempty" validate:"max=500"`
	TransactionDate    time.Time `json:"transaction_date" binding:"required"`
	IsRecurring        bool      `json:"is_recurring"`
	RecurringFrequency *string   `json:"recurring_frequency" binding:"omitempty" validate:"omitempty,oneof=daily weekly monthly yearly"`
}

type UpdateTransactionRequest struct {
	Type               string    `json:"type" binding:"omitempty" validate:"omitempty,oneof=income expense lend borrow"`
	Amount             float64   `json:"amount" binding:"omitempty" validate:"omitempty,gt=0"`
	CategoryID         string    `json:"category_id" binding:"omitempty" validate:"omitempty,uuid"`
	Description        string    `json:"description" binding:"omitempty" validate:"max=500"`
	TransactionDate    time.Time `json:"transaction_date" binding:"omitempty"`
	IsRecurring        *bool     `json:"is_recurring" binding:"omitempty"`
	RecurringFrequency *string   `json:"recurring_frequency" binding:"omitempty" validate:"omitempty,oneof=daily weekly monthly yearly"`
}

type TransactionResponse struct {
	ID                 string            `json:"id"`
	Type               string            `json:"type"`
	Amount             float64           `json:"amount"`
	Category           *CategoryResponse `json:"category,omitempty"`
	Description        string            `json:"description"`
	TransactionDate    time.Time         `json:"transaction_date"`
	IsRecurring        bool              `json:"is_recurring"`
	RecurringFrequency *string           `json:"recurring_frequency,omitempty"`
	CreatedAt          time.Time         `json:"created_at"`
}

type TransactionFilters struct {
	Page       int       `json:"page" form:"page"`
	Limit      int       `json:"limit" form:"limit"`
	Type       string    `json:"type" form:"type"`
	CategoryID string    `json:"category_id" form:"category_id"`
	StartDate  time.Time `json:"start_date" form:"start_date"`
	EndDate    time.Time `json:"end_date" form:"end_date"`
	Search     string    `json:"search" form:"search"`
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type PaginatedTransactions struct {
	Transactions []TransactionResponse `json:"transactions"`
	Pagination   Pagination           `json:"pagination"`
}

// Category DTOs
type CategoryResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
	Type  string `json:"type"`
}

// EMI DTOs
type CreateEMIRequest struct {
	Title         string    `json:"title" binding:"required" validate:"max=255"`
	TotalAmount   float64   `json:"total_amount" binding:"required" validate:"gt=0"`
	MonthlyAmount float64   `json:"monthly_amount" binding:"required" validate:"gt=0"`
	StartDate     time.Time `json:"start_date" binding:"required"`
	EndDate       time.Time `json:"end_date" binding:"required"`
	DueDate       int       `json:"due_date" binding:"required" validate:"gte=1,lte=31"`
	Description   string    `json:"description" binding:"omitempty" validate:"max=500"`
}

type CreateEMIPaymentRequest struct {
	Amount      float64   `json:"amount" binding:"required" validate:"gt=0"`
	PaymentDate time.Time `json:"payment_date" binding:"required"`
	DueMonth    time.Time `json:"due_month" binding:"required"`
	Notes       string    `json:"notes" binding:"omitempty" validate:"max=500"`
}

type EMIResponse struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	TotalAmount     float64   `json:"total_amount"`
	MonthlyAmount   float64   `json:"monthly_amount"`
	NextDueDate     time.Time `json:"next_due_date"`
	RemainingMonths int       `json:"remaining_months"`
	IsActive        bool      `json:"is_active"`
}

type RecordEMIPaymentRequest struct {
	Amount      float64   `json:"amount" binding:"required" validate:"gt=0"`
	PaymentDate time.Time `json:"payment_date" binding:"required"`
	DueMonth    time.Time `json:"due_month" binding:"required"`
}

type EMIPaymentResponse struct {
	ID          string    `json:"id"`
	Amount      float64   `json:"amount"`
	PaymentDate time.Time `json:"payment_date"`
	DueMonth    time.Time `json:"due_month"`
	IsPaid      bool      `json:"is_paid"`
}

// Friend DTOs
type SendFriendRequestRequest struct {
	Phone string `json:"phone" binding:"required" validate:"e164"`
}

type FriendResponse struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Phone        string  `json:"phone"`
	PhotoURL     string  `json:"photo_url"`
	TotalOwed    float64 `json:"total_owed"`
	TotalLending float64 `json:"total_lending"`
}

type FriendRequestResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Phone       string    `json:"phone"`
	RequestedBy string    `json:"requested_by"`
	CreatedAt   time.Time `json:"created_at"`
}

type FriendsListResponse struct {
	Friends         []FriendResponse        `json:"friends"`
	PendingRequests []FriendRequestResponse `json:"pending_requests"`
}

type HandleFriendRequestRequest struct {
	Action string `json:"action" binding:"required" validate:"oneof=accept reject"`
}

// Group DTOs
type CreateGroupRequest struct {
	Name        string   `json:"name" binding:"required" validate:"max=255"`
	Description string   `json:"description" binding:"omitempty" validate:"max=1000"`
	MemberIDs   []string `json:"member_ids" binding:"required" validate:"min=1,dive,uuid"`
}

type GroupResponse struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	MemberCount   int     `json:"member_count"`
	TotalExpenses float64 `json:"total_expenses"`
	YourShare     float64 `json:"your_share"`
	YourBalance   float64 `json:"your_balance"`
}

type GroupDetailsResponse struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Members       []GroupMemberResponse  `json:"members"`
	Expenses      []GroupExpenseResponse `json:"expenses"`
	TotalExpenses float64                `json:"total_expenses"`
	Balances      map[string]float64     `json:"balances"`
}

type GroupMemberResponse struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	PhotoURL string  `json:"photo_url"`
	Balance  float64 `json:"balance"`
}

type AddGroupExpenseRequest struct {
	Amount      float64           `json:"amount" binding:"required" validate:"gt=0"`
	Description string            `json:"description" binding:"required" validate:"max=255"`
	ExpenseDate time.Time         `json:"expense_date" binding:"required"`
	SplitType   string            `json:"split_type" binding:"required" validate:"oneof=equal custom"`
	Splits      []ExpenseSplitDTO `json:"splits" binding:"omitempty"`
}

type ExpenseSplitDTO struct {
	UserID string  `json:"user_id" binding:"required" validate:"uuid"`
	Amount float64 `json:"amount" binding:"required" validate:"gt=0"`
}

type GroupExpenseResponse struct {
	ID          string    `json:"id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	PaidBy      string    `json:"paid_by"`
	PayerName   string    `json:"payer_name"`
	ExpenseDate time.Time `json:"expense_date"`
	CreatedAt   time.Time `json:"created_at"`
}

type SettleGroupRequest struct {
	Settlements []SettlementDTO `json:"settlements" binding:"required" validate:"min=1,dive"`
}

type SettlementDTO struct {
	FromUserID  string    `json:"from_user_id" binding:"required" validate:"uuid"`
	ToUserID    string    `json:"to_user_id" binding:"required" validate:"uuid"`
	Amount      float64   `json:"amount" binding:"required" validate:"gt=0"`
	Description string    `json:"description" binding:"omitempty"`
	ExpenseDate time.Time `json:"expense_date" binding:"omitempty"`
}

// Report DTOs
type MonthlyReportResponse struct {
	Month             string              `json:"month"`
	TotalIncome       float64             `json:"total_income"`
	TotalExpense      float64             `json:"total_expense"`
	NetBalance        float64             `json:"net_balance"`
	CategoryBreakdown []CategoryBreakdown `json:"category_breakdown"`
	DailyTrend        []DailyTrendData    `json:"daily_trend"`
}

type CategoryBreakdown struct {
	CategoryName string  `json:"category_name"`
	Amount       float64 `json:"amount"`
	Percentage   float64 `json:"percentage"`
}

type DailyTrendData struct {
	Date    string  `json:"date"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

type CategoryReportResponse struct {
	Category              CategoryResponse      `json:"category"`
	TotalSpent            float64               `json:"total_spent"`
	TransactionCount      int                   `json:"transaction_count"`
	AveragePerTransaction float64               `json:"average_per_transaction"`
	MonthlyTrend          []MonthlyTrendData    `json:"monthly_trend"`
	RecentTransactions    []TransactionResponse `json:"recent_transactions"`
}

type MonthlyTrendData struct {
	Month  string  `json:"month"`
	Amount float64 `json:"amount"`
}

type YearlyReportResponse struct {
	Year             int                    `json:"year"`
	TotalIncome      float64                `json:"total_income"`
	TotalExpense     float64                `json:"total_expense"`
	NetBalance       float64                `json:"net_balance"`
	MonthlyBreakdown []MonthlyBreakdownData `json:"monthly_breakdown"`
	TopCategories    []CategoryBreakdown    `json:"top_categories"`
	Savings          float64                `json:"savings"`
	SavingsRate      float64                `json:"savings_rate"`
}

type MonthlyBreakdownData struct {
	Month   string  `json:"month"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Balance float64 `json:"balance"`
}

// Notification DTOs
type NotificationResponse struct {
	ID        string                 `json:"id"`
	Title     string                 `json:"title"`
	Message   string                 `json:"message"`
	Type      string                 `json:"type"`
	IsRead    bool                   `json:"is_read"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
}

type NotificationsResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	UnreadCount   int                    `json:"unread_count"`
}

// Generic Response DTOs
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code,omitempty"`
}

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}
