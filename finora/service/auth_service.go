package service

import (
	"fmt"
	"log"
	"time"

	"github.com/ayyoob-k-a/finora/configs"
	"github.com/ayyoob-k-a/finora/domain"
	"github.com/ayyoob-k-a/finora/model/dto"
	"github.com/ayyoob-k-a/finora/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	db         *gorm.DB
	config     configs.Config
	emailService *utils.EmailService
	smsService   *utils.SMSService
}

func NewAuthService(db *gorm.DB, config configs.Config) *AuthService {
	return &AuthService{
		db:         db,
		config:     config,
		emailService: utils.NewEmailService(config.SMTP),
		smsService:   utils.NewSMSService(config.Twilio),
	}
}

// SendOTP handles sending OTP to email or phone
func (s *AuthService) SendOTP(req dto.SendOTPRequest) error {
	// Validate input
	if req.Phone == "" && req.Email == "" {
		return fmt.Errorf("either phone or email is required")
	}

	if req.Phone != "" {
		if err := utils.ValidatePhoneNumber(req.Phone); err != nil {
			return fmt.Errorf("invalid phone number: %w", err)
		}
	}

	if req.Email != "" {
		if err := utils.ValidateEmail(req.Email); err != nil {
			return fmt.Errorf("invalid email: %w", err)
		}
	}

	// Generate OTP
	otpCode, err := utils.GenerateOTP()
	if err != nil {
		return fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Create OTP record
	otp := domain.OTP{
		Phone:     req.Phone,
		Email:     req.Email,
		OTPCode:   otpCode,
		ExpiresAt: utils.GenerateOTPExpiry(),
		IsUsed:    false,
		CreatedAt: time.Now(),
	}

	if err := s.db.Create(&otp).Error; err != nil {
		return fmt.Errorf("failed to save OTP: %w", err)
	}

	// Send OTP via SMS or Email
	if req.Phone != "" {
		if err := s.smsService.SendOTPSMS(req.Phone, otpCode); err != nil {
			log.Printf("Failed to send SMS OTP: %v", err)
			// Don't return error here, try email as fallback or continue
		} else {
			log.Printf("OTP sent successfully via SMS to %s", req.Phone)
		}
	}

	if req.Email != "" {
		if err := s.emailService.SendOTPEmail(req.Email, otpCode); err != nil {
			log.Printf("Failed to send email OTP: %v", err)
			return fmt.Errorf("failed to send OTP email: %w", err)
		} else {
			log.Printf("OTP sent successfully via email to %s", req.Email)
		}
	}

	return nil
}

// VerifyOTP handles OTP verification and user authentication
func (s *AuthService) VerifyOTP(req dto.VerifyOTPRequest) (*dto.AuthResponse, error) {
	// Validate input
	if err := utils.ValidatePhoneNumber(req.Phone); err != nil {
		return nil, fmt.Errorf("invalid phone number: %w", err)
	}

	if len(req.OTP) != 6 {
		return nil, fmt.Errorf("OTP must be 6 digits")
	}

	// Find OTP record
	var otp domain.OTP
	err := s.db.Where("phone = ? AND otp_code = ? AND is_used = ?", req.Phone, req.OTP, false).
		Order("created_at DESC").First(&otp).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid or expired OTP")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Check if OTP is expired
	if utils.IsOTPExpired(otp.ExpiresAt) {
		return nil, fmt.Errorf("OTP has expired")
	}

	// Mark OTP as used
	otp.IsUsed = true
	if err := s.db.Save(&otp).Error; err != nil {
		return nil, fmt.Errorf("failed to update OTP: %w", err)
	}

	// Find or create user
	var user domain.User
	isNewUser := false
	
	err = s.db.Where("phone = ?", req.Phone).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new user
			user = domain.User{
				Phone:           req.Phone,
				Name:            fmt.Sprintf("User %s", req.Phone[len(req.Phone)-4:]), // Default name
				DefaultCurrency: "USD",
				IsVerified:      true,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			}
			
			if err := s.db.Create(&user).Error; err != nil {
				return nil, fmt.Errorf("failed to create user: %w", err)
			}
			
			isNewUser = true
			log.Printf("Created new user with phone: %s", req.Phone)
		} else {
			return nil, fmt.Errorf("database error: %w", err)
		}
	} else {
		// Update existing user verification status
		user.IsVerified = true
		if err := s.db.Save(&user).Error; err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Phone, s.config.JWT.Secret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Send welcome messages for new users
	if isNewUser {
		go func() {
			if user.Email != "" {
				if err := s.emailService.SendWelcomeEmail(user.Email, user.Name); err != nil {
					log.Printf("Failed to send welcome email: %v", err)
				}
			}
			
			if err := s.smsService.SendWelcomeSMS(user.Phone, user.Name); err != nil {
				log.Printf("Failed to send welcome SMS: %v", err)
			}
		}()
	}

	// Prepare response
	userResponse := &dto.UserResponse{
		ID:              user.ID,
		Name:            user.Name,
		Phone:           user.Phone,
		Email:           user.Email,
		PhotoURL:        user.PhotoURL,
		DefaultCurrency: user.DefaultCurrency,
		MonthlyIncome:   user.MonthlyIncome,
		IsNewUser:       isNewUser,
		CreatedAt:       user.CreatedAt,
	}

	return &dto.AuthResponse{
		Success: true,
		Token:   token,
		User:    userResponse,
	}, nil
}

// RefreshToken handles JWT token refresh (placeholder for now)
func (s *AuthService) RefreshToken(req dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	// For now, we'll implement a simple approach
	// In a production app, you'd want to implement proper refresh token handling
	// with separate refresh tokens stored in database
	
	// Validate the existing token
	claims, err := utils.ValidateJWT(req.RefreshToken, s.config.JWT.Secret)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Find user
	var user domain.User
	if err := s.db.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Generate new token
	newToken, err := utils.GenerateJWT(user.ID, user.Phone, s.config.JWT.Secret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new token: %w", err)
	}

	return &dto.AuthResponse{
		Success: true,
		Token:   newToken,
	}, nil
}

// CleanupExpiredOTPs removes expired OTP records
func (s *AuthService) CleanupExpiredOTPs() error {
	result := s.db.Where("expires_at < ? OR is_used = ?", time.Now(), true).Delete(&domain.OTP{})
	if result.Error != nil {
		return fmt.Errorf("failed to cleanup expired OTPs: %w", result.Error)
	}
	
	if result.RowsAffected > 0 {
		log.Printf("Cleaned up %d expired/used OTP records", result.RowsAffected)
	}
	
	return nil
}
