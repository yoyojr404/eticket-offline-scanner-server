package middleware

import (
	"crypto/subtle"

	"github.com/yoyojr404/pkg/apperror"
	"github.com/yoyojr404/pkg/response"
	"github.com/gin-gonic/gin"
)

const AuthUserKey = "auth_user"

// AuthValidator is a function type for custom credential validation.
type AuthValidator func(username, password string) bool

// BasicAuth creates a Gin middleware for HTTP Basic Authentication.
func BasicAuth(validator AuthValidator, realm string) gin.HandlerFunc {
	if realm == "" {
		realm = "Restricted Area"
	}

	return func(c *gin.Context) {
		username, password, ok := c.Request.BasicAuth()
		if !ok || !validator(username, password) {
			c.Header("WWW-Authenticate", `Basic realm="`+realm+`"`)
			err := apperror.NewUnauthorized("Invalid or missing authentication credentials", nil)
			response.Error(c, err)
			c.Abort()
			return
		}

		c.Set(AuthUserKey, username)
		c.Next()
	}
}

// BasicAuthStatic creates a Gin middleware using fixed username and password with constant-time comparison.
func BasicAuthStatic(expectedUser, expectedPass, realm string) gin.HandlerFunc {
	return BasicAuth(func(user, pass string) bool {
		userMatch := subtle.ConstantTimeCompare([]byte(user), []byte(expectedUser)) == 1
		passMatch := subtle.ConstantTimeCompare([]byte(pass), []byte(expectedPass)) == 1
		return userMatch && passMatch
	}, realm)
}

// GetAuthUser retrieves the authenticated username from the Gin context.
func GetAuthUser(c *gin.Context) (string, bool) {
	val, exists := c.Get(AuthUserKey)
	if !exists {
		return "", false
	}
	user, ok := val.(string)
	return user, ok
}
