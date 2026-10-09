package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/yoyojr404/pkg/apperror"
	"github.com/yoyojr404/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// FieldError represents a single validation failure on a field.
type FieldError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

// BindAndValidate binds JSON/query into target struct and handles validation errors.
// Returns false if validation or binding failed and error response was already written.
func BindAndValidate(c *gin.Context, obj any) bool {
	if err := c.ShouldBind(obj); err != nil {
		if valErrs, ok := errors.AsType[validator.ValidationErrors](err); ok {
			details := make([]FieldError, 0, len(valErrs))
			for _, fe := range valErrs {
				details = append(details, FieldError{
					Field:   strings.ToLower(fe.Field()),
					Tag:     fe.Tag(),
					Message: formatValidationMessage(fe),
				})
			}
			response.ValidationError(c, details)
			return false
		}

		// Malformed JSON or type mismatch
		appErr := apperror.NewBadRequest(fmt.Sprintf("Invalid request payload: %v", err.Error()), err)
		response.Error(c, appErr)
		return false
	}

	return true
}

func formatValidationMessage(fe validator.FieldError) string {
	field := strings.ToLower(fe.Field())
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters/items", field, fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters/items", field, fe.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	default:
		return fmt.Sprintf("%s failed validation on '%s'", field, fe.Tag())
	}
}
