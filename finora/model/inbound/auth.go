package inbound

type Signup struct {
	Email           string `json:"email" `
	Password        string `json:"password"`
	Phone           string `json:"phone" `
	ConfirmPassword string `json:"confirm_password"`
	Username        string `json:"username" `
	IsVerified      bool   `json:"is_verified"`
}

type Login struct {
	Identifier string `json:"identifier"` // email or phone
	Password   string `json:"password"`
}
