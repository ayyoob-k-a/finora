package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	go_mail "gopkg.in/mail.v2"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// EmailMetrics for monitoring
type EmailMetrics struct {
	TotalSent     int64
	TotalFailed   int64
	LastSentTime  time.Time
	LastError     error
	ResponseTimes []time.Duration
	mutex         sync.RWMutex
}

type singleton struct {
	smtpHost    string
	smtpAddress string
	from        string
	password    string
}

var emailMetrics = &EmailMetrics{}

func GetEmailMetrics() EmailMetrics {
	emailMetrics.mutex.RLock()
	defer emailMetrics.mutex.RUnlock()
	return *emailMetrics
}

func GenerateToken(userID int) (string, string, error) {
	now := time.Now()

	atClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     now.Add(365 * 24 * time.Hour).Unix(),
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessToken, err := at.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	rtClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     now.Add(30 * 24 * time.Hour).Unix(),
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := rt.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

var instance *singleton
var once sync.Once

func GetInstance() *singleton {
	once.Do(func() {
		instance = newClient()
	})
	return instance
}

func newClient() *singleton {
	host := "smtp.gmail.com"
	from := os.Getenv("SMTP_FROM") // Use environment variable
	if from == "" {
		from = "finoraisme@gmail.com" // fallback
	}

	password := os.Getenv("SMTP_PASSWORD") // Use environment variable
	if password == "" {
		password = "gpyufuakpvjyddua" // fallback - REMOVE IN PRODUCTION
	}

	return &singleton{
		smtpHost:    host,
		smtpAddress: fmt.Sprintf("%s:587", host),
		from:        from,
		password:    password,
	}
}

func SendVerificationEmail(recipientEmail string, otp int, credentials string) error {

	data := struct {
		VerificationURL string
		Otp             int
	}{
		Otp: otp,
		// VerificationURL: fmt.Sprintf("http://%s/verify?email=%s&token=%s", credentials.URL, recipientEmail, verificationToken),
	}

	t, err := template.ParseFiles("otp_helper.html")

	if err != nil {

		return err

	}

	var tpl bytes.Buffer

	if err = t.Execute(&tpl, data); err != nil {

		fmt.Println("--", tpl.String())

		return err

	}

	result := tpl.String()

	m := go_mail.NewMessage()

	m.SetHeader("From", "finoraisme@gmail.com")

	m.SetHeader("To", recipientEmail)

	m.SetHeader("Subject", "Verification Mail From HiperHive!")

	m.SetBody("text/html", result)

	d := go_mail.NewDialer("smtp.gmail.com", 587, "finoraisme@gmail.com", credentials)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail

	if err := d.DialAndSend(m); err != nil {

		fmt.Println(err)

		panic(err)

	}

	return nil

}

// // Enhanced Send method with comprehensive monitoring and debugging
// func (s *singleton) Send(to string, htmlTemplate string, subject string, data any) error {
// 	startTime := time.Now()

// 	log.Printf("🚀 Starting email send to: %s", to)
// 	log.Printf("📧 Subject: %s", subject)
// 	log.Printf("📮 From: %s", s.from)
// 	log.Printf("🌐 SMTP Server: %s", s.smtpAddress)

// 	// Step 1: Test network connectivity
// 	if err := s.testConnectivity(); err != nil {
// 		s.recordFailure(err, startTime)
// 		return fmt.Errorf("connectivity test failed: %w", err)
// 	}
// 	log.Println("✅ Network connectivity: OK")

// 	// Step 2: Create email body
// 	body, err := s.createBodyWithValidation(htmlTemplate, subject, data)
// 	if err != nil {
// 		s.recordFailure(err, startTime)
// 		return fmt.Errorf("email body creation failed: %w", err)
// 	}
// 	log.Printf("✅ Email body created successfully (%d bytes)", len(body))

// 	// Step 3: Send with enhanced error handling
// 	err = s.sendWithRetry(to, body, 3)

// 	duration := time.Since(startTime)
// 	if err != nil {
// 		s.recordFailure(err, startTime)
// 		log.Printf("❌ Email sending failed after %v: %v", duration, err)
// 		return err
// 	}

// 	s.recordSuccess(startTime)
// 	log.Printf("✅ Email sent successfully to: %s (took %v)", to, duration)
// 	return nil
// }

// // Test network connectivity to SMTP server
// func (s *singleton) testConnectivity() error {
// 	log.Println("🔍 Testing network connectivity...")

// 	conn, err := net.DialTimeout("tcp", s.smtpAddress, 10*time.Second)
// 	if err != nil {
// 		return fmt.Errorf("failed to connect to SMTP server: %w", err)
// 	}
// 	defer conn.Close()

// 	log.Println("✅ Successfully connected to SMTP server")
// 	return nil
// }

// // Enhanced body creation with validation
// func (s *singleton) createBodyWithValidation(htmlTemplate string, subject string, data any) ([]byte, error) {
// 	log.Println("📝 Creating email body...")

// 	if subject == "" {
// 		return nil, fmt.Errorf("subject cannot be empty")
// 	}

// 	if htmlTemplate == "" {
// 		return nil, fmt.Errorf("HTML template cannot be empty")
// 	}

// 	var body bytes.Buffer
// 	const mime = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
// 	body.Write([]byte(fmt.Sprintf("Subject: %s\n%s\n", subject, mime)))

// 	t, err := template.New("email").Parse(htmlTemplate)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse template: %w", err)
// 	}

// 	err = t.Execute(&body, data)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to execute template: %w", err)
// 	}

// 	return body.Bytes(), nil
// }

// // Send with retry mechanism
// func (s *singleton) sendWithRetry(to string, body []byte, maxRetries int) error {
// 	var lastErr error

// 	for attempt := 1; attempt <= maxRetries; attempt++ {
// 		log.Printf("📤 Sending attempt %d/%d...", attempt, maxRetries)

// 		// Try different sending methods
// 		if attempt == 1 {
// 			lastErr = s.sendWithTLS(to, body)
// 		} else {
// 			lastErr = s.sendWithPlainAuth(to, body)
// 		}

// 		if lastErr == nil {
// 			return nil
// 		}

// 		log.Printf("⚠️ Attempt %d failed: %v", attempt, lastErr)

// 		if attempt < maxRetries {
// 			waitTime := time.Duration(attempt) * 2 * time.Second
// 			log.Printf("⏳ Waiting %v before retry...", waitTime)
// 			time.Sleep(waitTime)
// 		}
// 	}

// 	return fmt.Errorf("all %d attempts failed, last error: %w", maxRetries, lastErr)
// }

// // Send with explicit TLS
// func (s *singleton) sendWithTLS(to string, body []byte) error {
// 	log.Println("🔐 Attempting to send with TLS...")

// 	// Connect to the SMTP server
// 	log.Println("📡 Connecting to SMTP server...")
// 	conn, err := net.Dial("tcp", s.smtpAddress)
// 	if err != nil {
// 		return fmt.Errorf("failed to connect: %w", err)
// 	}
// 	defer conn.Close()
// 	log.Println("✅ TCP connection established")

// 	client, err := smtp.NewClient(conn, s.smtpHost)
// 	if err != nil {
// 		return fmt.Errorf("failed to create SMTP client: %w", err)
// 	}
// 	defer client.Quit()
// 	log.Println("✅ SMTP client created")

// 	// Start TLS
// 	log.Println("🔒 Starting TLS...")
// 	tlsConfig := &tls.Config{
// 		ServerName: s.smtpHost,
// 	}

// 	if err = client.StartTLS(tlsConfig); err != nil {
// 		return fmt.Errorf("failed to start TLS: %w", err)
// 	}
// 	log.Println("✅ TLS started successfully")

// 	// Authenticate
// 	log.Println("🔑 Authenticating...")
// 	auth := smtp.PlainAuth("", s.from, s.password, s.smtpHost)
// 	if err = client.Auth(auth); err != nil {
// 		log.Printf("❌ Authentication failed with error: %v", err)
// 		log.Printf("📧 From address: %s", s.from)
// 		log.Printf("🔐 Password length: %d", len(s.password))
// 		log.Printf("🌐 Host: %s", s.smtpHost)
// 		return fmt.Errorf("authentication failed: %w", err)
// 	}
// 	log.Println("✅ Authentication successful")

// 	// Set sender and recipient
// 	log.Printf("📤 Setting sender: %s", s.from)
// 	if err = client.Mail(s.from); err != nil {
// 		return fmt.Errorf("failed to set sender: %w", err)
// 	}
// 	log.Println("✅ Sender set successfully")

// 	log.Printf("📨 Setting recipient: %s", to)
// 	if err = client.Rcpt(to); err != nil {
// 		return fmt.Errorf("failed to set recipient: %w", err)
// 	}
// 	log.Println("✅ Recipient set successfully")

// 	// Send the email body
// 	log.Println("📝 Sending email data...")
// 	w, err := client.Data()
// 	if err != nil {
// 		return fmt.Errorf("failed to open data writer: %w", err)
// 	}

// 	_, err = w.Write(body)
// 	if err != nil {
// 		return fmt.Errorf("failed to write email body: %w", err)
// 	}

// 	err = w.Close()
// 	if err != nil {
// 		return fmt.Errorf("failed to close data writer: %w", err)
// 	}
// 	log.Println("✅ Email data sent successfully")

// 	return nil
// }

// // Fallback to plain auth method
// func (s *singleton) sendWithPlainAuth(to string, body []byte) error {
// 	log.Println("🔓 Attempting to send with plain auth...")

// 	auth := smtp.PlainAuth("", s.from, s.password, s.smtpHost)
// 	err := smtp.SendMail(s.smtpAddress, auth, s.from, []string{to}, body)
// 	if err != nil {
// 		return fmt.Errorf("plain auth send failed: %w", err)
// 	}

// 	return nil
// }

// // Record successful send
// func (s *singleton) recordSuccess(startTime time.Time) {
// 	emailMetrics.mutex.Lock()
// 	defer emailMetrics.mutex.Unlock()

// 	emailMetrics.TotalSent++
// 	emailMetrics.LastSentTime = time.Now()
// 	emailMetrics.ResponseTimes = append(emailMetrics.ResponseTimes, time.Since(startTime))

// 	// Keep only last 100 response times
// 	if len(emailMetrics.ResponseTimes) > 100 {
// 		emailMetrics.ResponseTimes = emailMetrics.ResponseTimes[1:]
// 	}
// }

// // Record failed send
// func (s *singleton) recordFailure(err error, startTime time.Time) {
// 	emailMetrics.mutex.Lock()
// 	defer emailMetrics.mutex.Unlock()

// 	emailMetrics.TotalFailed++
// 	emailMetrics.LastError = err
// 	emailMetrics.ResponseTimes = append(emailMetrics.ResponseTimes, time.Since(startTime))

// 	// Keep only last 100 response times
// 	if len(emailMetrics.ResponseTimes) > 100 {
// 		emailMetrics.ResponseTimes = emailMetrics.ResponseTimes[1:]
// 	}
// }

// // Diagnostic function
// func (s *singleton) RunDiagnostics() {
// 	log.Println("🔍 Running email system diagnostics...")

// 	// Check environment variables
// 	log.Printf("Environment check:")
// 	log.Printf("  SMTP_FROM: %s", os.Getenv("SMTP_FROM"))
// 	log.Printf("  SMTP_PASSWORD: %s", maskPassword(os.Getenv("SMTP_PASSWORD")))
// 	log.Printf("  JWT_SECRET: %s", maskPassword(os.Getenv("JWT_SECRET")))

// 	// Check configuration
// 	log.Printf("Configuration:")
// 	log.Printf("  SMTP Host: %s", s.smtpHost)
// 	log.Printf("  SMTP Address: %s", s.smtpAddress)
// 	log.Printf("  From: %s", s.from)
// 	log.Printf("  Password: %s", maskPassword(s.password))

// 	// Gmail-specific checks
// 	log.Println("Gmail Authentication Checks:")
// 	if len(s.password) == 16 && !containsSpecialChars(s.password) {
// 		log.Println("✅ Password format looks like Gmail App Password (16 chars, alphanumeric)")
// 	} else {
// 		log.Printf("⚠️  Password doesn't look like Gmail App Password. Length: %d", len(s.password))
// 		log.Println("   Expected: 16 characters, alphanumeric only")
// 		log.Println("   💡 Generate one at: https://myaccount.google.com/apppasswords")
// 	}

// 	// Test connectivity
// 	if err := s.testConnectivity(); err != nil {
// 		log.Printf("❌ Connectivity test failed: %v", err)
// 	} else {
// 		log.Println("✅ Connectivity test passed")
// 	}

// 	// Test basic SMTP capabilities
// 	log.Println("Testing SMTP capabilities...")
// 	if err := s.testSMTPCapabilities(); err != nil {
// 		log.Printf("⚠️  SMTP capabilities test: %v", err)
// 	}

// 	// Display metrics
// 	metrics := GetEmailMetrics()
// 	log.Printf("Metrics:")
// 	log.Printf("  Total sent: %d", metrics.TotalSent)
// 	log.Printf("  Total failed: %d", metrics.TotalFailed)
// 	log.Printf("  Last sent: %v", metrics.LastSentTime)
// 	if metrics.LastError != nil {
// 		log.Printf("  Last error: %v", metrics.LastError)
// 	}
// }

// // Test SMTP server capabilities
// func (s *singleton) testSMTPCapabilities() error {
// 	conn, err := net.Dial("tcp", s.smtpAddress)
// 	if err != nil {
// 		return fmt.Errorf("failed to connect: %w", err)
// 	}
// 	defer conn.Close()

// 	client, err := smtp.NewClient(conn, s.smtpHost)
// 	if err != nil {
// 		return fmt.Errorf("failed to create SMTP client: %w", err)
// 	}
// 	defer client.Quit()

// 	// Check if STARTTLS is supported
// 	ok, _ := client.Extension("STARTTLS")
// 	if ok {
// 		log.Println("✅ STARTTLS supported")
// 	} else {
// 		log.Println("⚠️  STARTTLS not supported")
// 	}

// 	// Check authentication methods
// 	ok, _ = client.Extension("AUTH")
// 	if ok {
// 		log.Println("✅ AUTH supported")
// 	} else {
// 		log.Println("⚠️  AUTH not supported")
// 	}

// 	return nil
// }

// // Check if password contains special characters
// func containsSpecialChars(s string) bool {
// 	for _, char := range s {
// 		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
// 			return true
// 		}
// 	}
// 	return false
// }

// func maskPassword(password string) string {
// 	if password == "" {
// 		return "NOT SET"
// 	}
// 	if len(password) <= 4 {
// 		return "****"
// 	}
// 	return password[:2] + "****" + password[len(password)-2:]
// }

// // Legacy function for backward compatibility
// func createBody(htmlTemplate string, subject string, data any) ([]byte, error) {
// 	instance := GetInstance()
// 	return instance.createBodyWithValidation(htmlTemplate, subject, data)
// }
