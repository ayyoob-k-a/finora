package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"

	"github.com/ayyoob-k-a/finora/configs"
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	config configs.SMTPConfig
}

func NewEmailService(config configs.SMTPConfig) *EmailService {
	return &EmailService{
		config: config,
	}
}

func (e *EmailService) SendOTPEmail(recipientEmail, otp string) error {
	// Create HTML template for OTP email
	htmlTemplate := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Your OTP Code</title>
    <style>
        .container {
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            font-family: Arial, sans-serif;
        }
        .header {
            text-align: center;
            color: #333;
            margin-bottom: 30px;
        }
        .otp-box {
            background-color: #f8f9fa;
            border: 2px solid #007bff;
            border-radius: 8px;
            padding: 20px;
            text-align: center;
            margin: 20px 0;
        }
        .otp-code {
            font-size: 32px;
            font-weight: bold;
            color: #007bff;
            letter-spacing: 5px;
            margin: 10px 0;
        }
        .warning {
            color: #dc3545;
            font-size: 14px;
            margin-top: 20px;
        }
        .footer {
            text-align: center;
            color: #666;
            font-size: 12px;
            margin-top: 30px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Finora - Your OTP Code</h1>
        </div>
        
        <p>Hello,</p>
        
        <p>You've requested to verify your account. Please use the following One-Time Password (OTP) to complete your verification:</p>
        
        <div class="otp-box">
            <div class="otp-code">{{.OTP}}</div>
            <p style="margin: 0; color: #666;">This code will expire in 5 minutes</p>
        </div>
        
        <p>If you didn't request this code, please ignore this email.</p>
        
        <div class="warning">
            <strong>Security Notice:</strong> Never share this code with anyone. Finora will never ask for your OTP over phone or email.
        </div>
        
        <div class="footer">
            <p>This is an automated email from Finora. Please do not reply.</p>
            <p>&copy; 2024 Finora. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`

	// Parse template
	tmpl, err := template.New("otp").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	// Execute template with OTP data
	var body bytes.Buffer
	data := struct {
		OTP string
	}{
		OTP: otp,
	}

	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	// Create email message
	m := gomail.NewMessage()
	m.SetHeader("From", e.config.Username)
	m.SetHeader("To", recipientEmail)
	m.SetHeader("Subject", "Your Finora OTP Code")
	m.SetBody("text/html", body.String())

	// Create dialer
	d := gomail.NewDialer(e.config.Host, e.config.Port, e.config.Username, e.config.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send email
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email to %s: %v", recipientEmail, err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("OTP email sent successfully to %s", recipientEmail)
	return nil
}

func (e *EmailService) SendWelcomeEmail(recipientEmail, userName string) error {
	htmlTemplate := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Welcome to Finora</title>
    <style>
        .container {
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            font-family: Arial, sans-serif;
        }
        .header {
            text-align: center;
            color: #333;
            margin-bottom: 30px;
        }
        .welcome-box {
            background-color: #f8f9fa;
            border-radius: 8px;
            padding: 30px;
            text-align: center;
            margin: 20px 0;
        }
        .cta-button {
            background-color: #007bff;
            color: white;
            padding: 12px 30px;
            text-decoration: none;
            border-radius: 5px;
            display: inline-block;
            margin: 20px 0;
        }
        .footer {
            text-align: center;
            color: #666;
            font-size: 12px;
            margin-top: 30px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to Finora! 🎉</h1>
        </div>
        
        <div class="welcome-box">
            <h2>Hello {{.UserName}}!</h2>
            <p>Welcome to Finora - your personal expense management companion.</p>
            <p>Start tracking your expenses, managing EMIs, and splitting bills with friends!</p>
            <a href="#" class="cta-button">Get Started</a>
        </div>
        
        <div class="footer">
            <p>&copy; 2024 Finora. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`

	tmpl, err := template.New("welcome").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse welcome email template: %w", err)
	}

	var body bytes.Buffer
	data := struct {
		UserName string
	}{
		UserName: userName,
	}

	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute welcome email template: %w", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", e.config.Username)
	m.SetHeader("To", recipientEmail)
	m.SetHeader("Subject", "Welcome to Finora!")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(e.config.Host, e.config.Port, e.config.Username, e.config.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send welcome email: %w", err)
	}

	log.Printf("Welcome email sent successfully to %s", recipientEmail)
	return nil
}
