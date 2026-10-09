package middleware

import (
	"time"

	"github.com/yoyojr404/pkg/config"
	"github.com/gin-gonic/gin"
)

// Provider holds initialized and configured middlewares ready to be distributed.
// Being in pkg/middleware prevents any circular dependencies with internal feature packages.
type Provider struct {
	BasicAuth gin.HandlerFunc
	RateLimit gin.HandlerFunc
}

// NewProvider constructs the middleware provider with application config.
func NewProvider(cfg *config.Config) *Provider {
	limiter := NewRateLimiter(60, time.Minute)

	return &Provider{
		BasicAuth: BasicAuthStatic(cfg.BasicAuthUser, cfg.BasicAuthPass, "Kartjis Scanner API"),
		RateLimit: limiter.Limit(),
	}
}
