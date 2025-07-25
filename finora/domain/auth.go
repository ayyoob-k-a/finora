package domain

type User struct {
	ID         int    `json:"id" gorm:"primaryKey"`
	Email      string `json:"email" gorm:"unique"`
	Password   string `json:"password"`
	Phone      string `json:"phone"`
	Username   string `json:"username"`
	IsVerified bool   `json:"is_verified"`
	CreatedAt  string `json:"created_at"`
	Otp        int    `json:"otp"`
}
