package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

// GenerateOTP generates a 6-digit OTP
func GenerateOTP() (string, error) {
	max := big.NewInt(999999)
	min := big.NewInt(100000)

	// Generate a random number between min and max (inclusive)
	n, err := rand.Int(rand.Reader, new(big.Int).Sub(max, min))
	if err != nil {
		return "", fmt.Errorf("failed to generate random number: %w", err)
	}

	// Add min to the result to get the desired range
	otp := new(big.Int).Add(n, min).Int64()

	return fmt.Sprintf("%06d", otp), nil
}

// GenerateOTPExpiry returns expiration time for OTP (5 minutes from now)
func GenerateOTPExpiry() time.Time {
	return time.Now().Add(5 * time.Minute)
}

// IsOTPExpired checks if OTP has expired
func IsOTPExpired(expiresAt time.Time) bool {
	return time.Now().After(expiresAt)
}
