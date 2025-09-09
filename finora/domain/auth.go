package domain

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents the users table
type User struct {
	ID              string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Phone           string    `json:"phone" gorm:"uniqueIndex"`
	Email           string    `json:"email" gorm:"uniqueIndex"`
	Name            string    `json:"name" gorm:"not null"`
	PhotoURL        string    `json:"photo_url"`
	DefaultCurrency string    `json:"default_currency" gorm:"default:'USD'"`
	MonthlyIncome   float64   `json:"monthly_income"`
	IsVerified      bool      `json:"is_verified" gorm:"default:false"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// OTP represents the otps table for authentication
type OTP struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	OTPCode   string    `json:"otp_code" gorm:"not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	IsUsed    bool      `json:"is_used" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
}

func (o *OTP) BeforeCreate(tx *gorm.DB) error {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	return nil
}

// Category represents expense/income categories
type Category struct {
	ID        string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string `json:"name" gorm:"not null"`
	Icon      string `json:"icon"`
	Color     string `json:"color"`
	Type      string `json:"type" gorm:"check:type IN ('expense', 'income')"`
	IsDefault bool   `json:"is_default" gorm:"default:false"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

// Transaction represents the transactions table
type Transaction struct {
	ID                 string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID             string    `json:"user_id" gorm:"not null"`
	Type               string    `json:"type" gorm:"not null;check:type IN ('income', 'expense', 'lend', 'borrow')"`
	Amount             float64   `json:"amount" gorm:"not null"`
	CategoryID         *string   `json:"category_id"`
	Description        string    `json:"description"`
	TransactionDate    time.Time `json:"transaction_date" gorm:"not null"`
	IsRecurring        bool      `json:"is_recurring" gorm:"default:false"`
	RecurringFrequency *string   `json:"recurring_frequency" gorm:"check:recurring_frequency IN ('daily', 'weekly', 'monthly', 'yearly')"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	
	// Associations
	User     User      `json:"user" gorm:"foreignKey:UserID"`
	Category *Category `json:"category" gorm:"foreignKey:CategoryID"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

// EMI represents the emis table
type EMI struct {
	ID            string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID        string    `json:"user_id" gorm:"not null"`
	Title         string    `json:"title" gorm:"not null"`
	Description   string    `json:"description"`
	TotalAmount   float64   `json:"total_amount" gorm:"not null"`
	MonthlyAmount float64   `json:"monthly_amount" gorm:"not null"`
	StartDate     time.Time `json:"start_date" gorm:"not null"`
	EndDate       time.Time `json:"end_date" gorm:"not null"`
	DueDate       int       `json:"due_date" gorm:"not null"` // Day of month (1-31)
	IsActive      bool      `json:"is_active" gorm:"default:true"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	
	// Computed fields (not stored in database)
	NextDueDate      time.Time `json:"next_due_date" gorm:"-"`
	RemainingMonths  int       `json:"remaining_months" gorm:"-"`
	
	// Associations
	User     User         `json:"user" gorm:"foreignKey:UserID"`
	Payments []EMIPayment `json:"payments" gorm:"foreignKey:EMIID"`
}

func (e *EMI) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return nil
}

// EMIPayment represents the emi_payments table
type EMIPayment struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	EMIID       string    `json:"emi_id" gorm:"not null"`
	Amount      float64   `json:"amount" gorm:"not null"`
	PaymentDate time.Time `json:"payment_date" gorm:"not null"`
	DueMonth    time.Time `json:"due_month" gorm:"not null"` // First day of the month this payment is for
	Notes       string    `json:"notes"`
	IsPaid      bool      `json:"is_paid" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	
	// Associations
	EMI EMI `json:"emi" gorm:"foreignKey:EMIID"`
}

func (ep *EMIPayment) BeforeCreate(tx *gorm.DB) error {
	if ep.ID == "" {
		ep.ID = uuid.New().String()
	}
	return nil
}

// Friend represents the friends table
type Friend struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID      string    `json:"user_id" gorm:"not null"`
	FriendID    string    `json:"friend_id" gorm:"not null"`
	Status      string    `json:"status" gorm:"default:'pending';check:status IN ('pending', 'accepted', 'rejected', 'blocked')"`
	RequestedBy string    `json:"requested_by" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Associations
	User         User `json:"user" gorm:"foreignKey:UserID"`
	FriendUser   User `json:"friend_user" gorm:"foreignKey:FriendID"`
	RequestedByUser User `json:"requested_by_user" gorm:"foreignKey:RequestedBy"`
}

func (f *Friend) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = uuid.New().String()
	}
	return nil
}

// Group represents the groups table
type Group struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"created_by" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Associations
	Creator User          `json:"creator" gorm:"foreignKey:CreatedBy"`
	Members []GroupMember `json:"members" gorm:"foreignKey:GroupID"`
	Expenses []GroupExpense `json:"expenses" gorm:"foreignKey:GroupID"`
}

func (g *Group) BeforeCreate(tx *gorm.DB) error {
	if g.ID == "" {
		g.ID = uuid.New().String()
	}
	return nil
}

// GroupMember represents the group_members table
type GroupMember struct {
	ID       string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	GroupID  string    `json:"group_id" gorm:"not null"`
	UserID   string    `json:"user_id" gorm:"not null"`
	JoinedAt time.Time `json:"joined_at"`
	
	// Associations
	Group Group `json:"group" gorm:"foreignKey:GroupID"`
	User  User  `json:"user" gorm:"foreignKey:UserID"`
}

func (gm *GroupMember) BeforeCreate(tx *gorm.DB) error {
	if gm.ID == "" {
		gm.ID = uuid.New().String()
	}
	return nil
}

// GroupExpense represents the group_expenses table
type GroupExpense struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	GroupID     string    `json:"group_id" gorm:"not null"`
	PaidBy      string    `json:"paid_by" gorm:"not null"`
	Amount      float64   `json:"amount" gorm:"not null"`
	Description string    `json:"description" gorm:"not null"`
	ExpenseDate time.Time `json:"expense_date" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Associations
	Group  Group         `json:"group" gorm:"foreignKey:GroupID"`
	Payer  User          `json:"payer" gorm:"foreignKey:PaidBy"`
	Splits []ExpenseSplit `json:"splits" gorm:"foreignKey:ExpenseID"`
}

func (ge *GroupExpense) BeforeCreate(tx *gorm.DB) error {
	if ge.ID == "" {
		ge.ID = uuid.New().String()
	}
	return nil
}

// ExpenseSplit represents the expense_splits table
type ExpenseSplit struct {
	ID        string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ExpenseID string     `json:"expense_id" gorm:"not null"`
	UserID    string     `json:"user_id" gorm:"not null"`
	Amount    float64    `json:"amount" gorm:"not null"`
	IsSettled bool       `json:"is_settled" gorm:"default:false"`
	SettledAt *time.Time `json:"settled_at"`
	CreatedAt time.Time  `json:"created_at"`
	
	// Associations
	Expense GroupExpense `json:"expense" gorm:"foreignKey:ExpenseID"`
	User    User         `json:"user" gorm:"foreignKey:UserID"`
}

func (es *ExpenseSplit) BeforeCreate(tx *gorm.DB) error {
	if es.ID == "" {
		es.ID = uuid.New().String()
	}
	return nil
}

// Notification represents the notifications table
type Notification struct {
	ID        string                 `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    string                 `json:"user_id" gorm:"not null"`
	Title     string                 `json:"title" gorm:"not null"`
	Message   string                 `json:"message" gorm:"not null"`
	Type      string                 `json:"type" gorm:"not null"`
	IsRead    bool                   `json:"is_read" gorm:"default:false"`
	Metadata  map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	CreatedAt time.Time              `json:"created_at"`
	
	// Associations
	User User `json:"user" gorm:"foreignKey:UserID"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	return nil
}
