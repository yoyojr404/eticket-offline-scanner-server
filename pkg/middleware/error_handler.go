package middleware

import (
	"github.com/yoyojr404/pkg/apperror"
	"github.com/yoyojr404/pkg/response"
	"github.com/gin-gonic/gin"
)

// ErrorHandler is a Gin middleware that formats any unhandled errors accumulated via c.Error(err).
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// If a response was already sent by the handler, skip.
		if c.Writer.Written() {
			return
		}

		if len(c.Errors) > 0 {
			// Pick the last error recorded
			lastErr := c.Errors.Last().Err

			appErr := apperror.FromError(lastErr)
			response.Error(c, appErr)
		}
	}
}
