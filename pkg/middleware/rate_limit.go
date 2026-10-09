package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/yoyojr404/pkg/apperror"
	"github.com/yoyojr404/pkg/response"
	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple in-memory sliding-window rate limiter per client IP.
type RateLimiter struct {
	mu       sync.Mutex
	limits   map[string][]time.Time
	maxReqs  int
	duration time.Duration
}

// NewRateLimiter creates a RateLimiter that allows maxReqs within the duration window.
func NewRateLimiter(maxReqs int, duration time.Duration) *RateLimiter {
	return &RateLimiter{
		limits:   make(map[string][]time.Time),
		maxReqs:  maxReqs,
		duration: duration,
	}
}

// Limit returns a Gin middleware that throttles excessive requests.
func (rl *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		rl.mu.Lock()
		defer rl.mu.Unlock()

		windowStart := now.Add(-rl.duration)
		var valid []time.Time
		for _, t := range rl.limits[ip] {
			if t.After(windowStart) {
				valid = append(valid, t)
			}
		}

		if len(valid) >= rl.maxReqs {
			err := apperror.New(http.StatusTooManyRequests, "TOO_MANY_REQUESTS", "Rate limit exceeded. Please try again later.", nil)
			response.Error(c, err)
			c.Abort()
			return
		}

		rl.limits[ip] = append(valid, now)
		c.Next()
	}
}
