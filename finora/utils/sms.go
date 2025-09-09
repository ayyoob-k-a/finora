package utils

import (
	"fmt"
	"log"

	"github.com/ayyoob-k-a/finora/configs"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type SMSService struct {
	client *twilio.RestClient
	config configs.TwilioConfig
}

func NewSMSService(config configs.TwilioConfig) *SMSService {
	if config.AccountSID == "" || config.AuthToken == "" {
		log.Println("Warning: Twilio credentials not configured. SMS functionality will be disabled.")
		return &SMSService{config: config}
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: config.AccountSID,
		Password: config.AuthToken,
	})

	return &SMSService{
		client: client,
		config: config,
	}
}

func (s *SMSService) SendOTPSMS(phoneNumber, otp string) error {
	if s.client == nil {
		return fmt.Errorf("SMS service not configured - Twilio credentials missing")
	}

	if s.config.PhoneNumber == "" {
		return fmt.Errorf("Twilio phone number not configured")
	}

	message := fmt.Sprintf("Your Finora OTP code is: %s\n\nThis code will expire in 5 minutes. Never share this code with anyone.", otp)

	params := &openapi.CreateMessageParams{}
	params.SetTo(phoneNumber)
	params.SetFrom(s.config.PhoneNumber)
	params.SetBody(message)

	resp, err := s.client.Api.CreateMessage(params)
	if err != nil {
		log.Printf("Failed to send SMS to %s: %v", phoneNumber, err)
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	if resp.Sid != nil {
		log.Printf("SMS sent successfully to %s, Message SID: %s", phoneNumber, *resp.Sid)
	}

	return nil
}

func (s *SMSService) SendWelcomeSMS(phoneNumber, userName string) error {
	if s.client == nil {
		return fmt.Errorf("SMS service not configured - Twilio credentials missing")
	}

	if s.config.PhoneNumber == "" {
		return fmt.Errorf("Twilio phone number not configured")
	}

	message := fmt.Sprintf("Welcome to Finora, %s! 🎉\n\nStart managing your expenses, EMIs, and group expenses today.\n\nDownload our app and take control of your finances!", userName)

	params := &openapi.CreateMessageParams{}
	params.SetTo(phoneNumber)
	params.SetFrom(s.config.PhoneNumber)
	params.SetBody(message)

	resp, err := s.client.Api.CreateMessage(params)
	if err != nil {
		log.Printf("Failed to send welcome SMS to %s: %v", phoneNumber, err)
		return fmt.Errorf("failed to send welcome SMS: %w", err)
	}

	if resp.Sid != nil {
		log.Printf("Welcome SMS sent successfully to %s, Message SID: %s", phoneNumber, *resp.Sid)
	}

	return nil
}

func (s *SMSService) SendEMIReminderSMS(phoneNumber, emiTitle string, amount float64, dueDate string) error {
	if s.client == nil {
		return fmt.Errorf("SMS service not configured - Twilio credentials missing")
	}

	if s.config.PhoneNumber == "" {
		return fmt.Errorf("Twilio phone number not configured")
	}

	message := fmt.Sprintf("📅 EMI Reminder\n\n%s payment of $%.2f is due on %s.\n\nDon't forget to make your payment on time!\n\n- Finora", emiTitle, amount, dueDate)

	params := &openapi.CreateMessageParams{}
	params.SetTo(phoneNumber)
	params.SetFrom(s.config.PhoneNumber)
	params.SetBody(message)

	resp, err := s.client.Api.CreateMessage(params)
	if err != nil {
		log.Printf("Failed to send EMI reminder SMS to %s: %v", phoneNumber, err)
		return fmt.Errorf("failed to send EMI reminder SMS: %w", err)
	}

	if resp.Sid != nil {
		log.Printf("EMI reminder SMS sent successfully to %s, Message SID: %s", phoneNumber, *resp.Sid)
	}

	return nil
}

func (s *SMSService) IsConfigured() bool {
	return s.client != nil && s.config.PhoneNumber != ""
}
