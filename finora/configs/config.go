package configs

import (
	"os"
	"strconv"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

type JWTConfig struct {
	Secret        string
	Expiry        string
	RefreshExpiry string
}

type TwilioConfig struct {
	AccountSID  string
	AuthToken   string
	PhoneNumber string
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type ServerConfig struct {
	Port string
	Mode string
}

type FileConfig struct {
	MaxFileSize string
	UploadPath  string
}

type Config struct {
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Twilio   TwilioConfig
	SMTP     SMTPConfig
	Server   ServerConfig
	File     FileConfig
}

func GetConfig() Config {
	smtpPort, _ := strconv.Atoi(getEnvWithDefault("SMTP_PORT", "587"))

	return Config{
		Database: DatabaseConfig{
			Host:     getEnvWithDefault("DB_HOST", "localhost"),
			Port:     getEnvWithDefault("DB_PORT", "5432"),
			User:     getEnvWithDefault("DB_USER", "finora_user"),
			Password: getEnvWithDefault("DB_PASSWORD", "finora123"),
			Name:     getEnvWithDefault("DB_NAME", "finora_db"),
			SSLMode:  getEnvWithDefault("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnvWithDefault("REDIS_HOST", "localhost"),
			Port:     getEnvWithDefault("REDIS_PORT", "6379"),
			Password: getEnvWithDefault("REDIS_PASSWORD", ""),
		},
		JWT: JWTConfig{
			Secret:        getEnvWithDefault("JWT_SECRET", "default-secret-change-this"),
			Expiry:        getEnvWithDefault("JWT_EXPIRY", "24h"),
			RefreshExpiry: getEnvWithDefault("REFRESH_TOKEN_EXPIRY", "168h"),
		},
		Twilio: TwilioConfig{
			AccountSID:  os.Getenv("TWILIO_ACCOUNT_SID"),
			AuthToken:   os.Getenv("TWILIO_AUTH_TOKEN"),
			PhoneNumber: os.Getenv("TWILIO_PHONE_NUMBER"),
		},
		SMTP: SMTPConfig{
			Host:     getEnvWithDefault("SMTP_HOST", "smtp.gmail.com"),
			Port:     smtpPort,
			Username: os.Getenv("SMTP_USERNAME"),
			Password: os.Getenv("SMTP_PASSWORD"),
		},
		Server: ServerConfig{
			Port: getEnvWithDefault("PORT", "8081"),
			Mode: getEnvWithDefault("GIN_MODE", "debug"),
		},
		File: FileConfig{
			MaxFileSize: getEnvWithDefault("MAX_FILE_SIZE", "10MB"),
			UploadPath:  getEnvWithDefault("UPLOAD_PATH", "./uploads"),
		},
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
