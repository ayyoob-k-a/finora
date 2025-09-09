package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/ayyoob-k-a/finora/utils"
	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	visitors map[string]*Visitor
	mutex    sync.RWMutex
}

type Visitor struct {
	lastSeen time.Time
	count    int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
	}
	
	// Clean up expired visitors every minute
	go rl.cleanupExpiredVisitors()
	
	return rl
}

// RateLimitMiddleware creates rate limiting middleware
func (rl *RateLimiter) RateLimitMiddleware(requestsPerMinute int) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Use IP address as identifier
		// In production, you might want to use user ID for authenticated requests
		ip := c.ClientIP()
		
		if rl.isRateLimited(ip, requestsPerMinute) {
			c.JSON(http.StatusTooManyRequests, utils.ErrorResponse("Rate limit exceeded. Please try again later."))
			c.Abort()
			return
		}
		
		c.Next()
	})
}

// OTPRateLimitMiddleware creates special rate limiting for OTP requests
func (rl *RateLimiter) OTPRateLimitMiddleware(requestsPerHour int) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Use phone number or IP address as identifier
		identifier := c.ClientIP()
		if phone := c.PostForm("phone"); phone != "" {
			identifier = phone
		}
		
		if rl.isOTPRateLimited(identifier, requestsPerHour) {
			c.JSON(http.StatusTooManyRequests, utils.ErrorResponse("Too many OTP requests. Please try again in an hour."))
			c.Abort()
			return
		}
		
		c.Next()
	})
}

func (rl *RateLimiter) isRateLimited(identifier string, limit int) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	
	now := time.Now()
	windowStart := now.Add(-time.Minute)
	
	visitor, exists := rl.visitors[identifier]
	if !exists {
		rl.visitors[identifier] = &Visitor{
			lastSeen: now,
			count:    1,
		}
		return false
	}
	
	// Reset count if outside the time window
	if visitor.lastSeen.Before(windowStart) {
		visitor.count = 1
		visitor.lastSeen = now
		return false
	}
	
	// Increment count
	visitor.count++
	visitor.lastSeen = now
	
	return visitor.count > limit
}

func (rl *RateLimiter) isOTPRateLimited(identifier string, limit int) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	
	now := time.Now()
	windowStart := now.Add(-time.Hour) // OTP rate limiting is per hour
	
	visitor, exists := rl.visitors[identifier+"_otp"]
	if !exists {
		rl.visitors[identifier+"_otp"] = &Visitor{
			lastSeen: now,
			count:    1,
		}
		return false
	}
	
	// Reset count if outside the time window
	if visitor.lastSeen.Before(windowStart) {
		visitor.count = 1
		visitor.lastSeen = now
		return false
	}
	
	// Increment count
	visitor.count++
	visitor.lastSeen = now
	
	return visitor.count > limit
}

func (rl *RateLimiter) cleanupExpiredVisitors() {
	for {
		time.Sleep(time.Minute)
		
		rl.mutex.Lock()
		now := time.Now()
		expiry := now.Add(-time.Hour) // Remove visitors older than 1 hour
		
		for identifier, visitor := range rl.visitors {
			if visitor.lastSeen.Before(expiry) {
				delete(rl.visitors, identifier)
			}
		}
		rl.mutex.Unlock()
	}
}

// Global rate limiter instance
var GlobalRateLimiter = NewRateLimiter()
