package middleware

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple in-memory sliding-window rate limiter per client IP.
type BlacklistMap struct {
	mu          sync.Mutex
	whitelistIP map[string]bool
}

// NewRateLimiter creates a RateLimiter that allows maxReqs within the duration window.
func NewBlacklistMap(whitelistIP map[string]bool) *BlacklistMap {
	return &BlacklistMap{
		whitelistIP: whitelistIP,
		mu:          sync.Mutex{},
	}
}

// Limit returns a Gin middleware that throttles excessive requests.
func (rl *BlacklistMap) Blacklist() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		rl.mu.Lock()
		defer rl.mu.Unlock()

		ip, port, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil {
			log.Println(err)
		}

		addr := fmt.Sprintf("%s:%s", ip, port)
		if _, ok := rl.whitelistIP[addr]; !ok {
			log.Println("abort IP", addr)
			c.Abort()
			return
		}

		log.Println("allow IP", addr)
		c.Next()
	}
}
