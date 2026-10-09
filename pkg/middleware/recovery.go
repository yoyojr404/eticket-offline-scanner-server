package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/yoyojr404/pkg/apperror"
	"github.com/yoyojr404/pkg/response"
	"github.com/gin-gonic/gin"
)

// Recovery returns a middleware that recovers from any panics and writes a standard 500 JSON response.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC RECOVERED] %v\nStack trace:\n%s", r, string(debug.Stack()))

				var err error
				if e, ok := r.(error); ok {
					err = e
				} else {
					err = fmt.Errorf("%v", r)
				}

				appErr := apperror.NewInternal("Internal server error", err)
				response.Error(c, appErr)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()

		c.Next()
	}
}
