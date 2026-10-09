package response

import (
	"math"
	"net/http"

	"github.com/yoyojr404/pkg/apperror"
	"github.com/gin-gonic/gin"
)

// Response represents standard API envelope.
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
	Error   any    `json:"error,omitempty"`
}

// PaginationMeta represents pagination metadata.
type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	Total       int `json:"total"`
	TotalPages  int `json:"total_pages"`
}

// NewPaginationMeta calculates total pages and returns PaginationMeta.
func NewPaginationMeta(currentPage, perPage, total int) *PaginationMeta {
	if perPage <= 0 {
		perPage = 10
	}
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	if totalPages == 0 {
		totalPages = 1
	}
	return &PaginationMeta{
		CurrentPage: currentPage,
		PerPage:     perPage,
		Total:       total,
		TotalPages:  totalPages,
	}
}

// Success sends an HTTP success response.
func Success(c *gin.Context, statusCode int, message string, data any) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// OK sends a 200 OK response.
func OK(c *gin.Context, message string, data any) {
	Success(c, http.StatusOK, message, data)
}

// Created sends a 201 Created response.
func Created(c *gin.Context, message string, data any) {
	Success(c, http.StatusCreated, message, data)
}

// SuccessWithMeta sends an HTTP success response with metadata (e.g., pagination).
func SuccessWithMeta(c *gin.Context, statusCode int, message string, data any, meta any) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// Error sends an HTTP error response using AppError.
func Error(c *gin.Context, err *apperror.AppError) {
	if err == nil {
		err = apperror.NewInternal("An unexpected error occurred", nil)
	}

	payload := Response{
		Success: false,
		Message: err.Message,
		Error: gin.H{
			"code":    err.Code,
			"details": err.Details,
		},
	}

	c.JSON(err.StatusCode, payload)
}

// ValidationError sends a 422 Unprocessable Entity error response.
func ValidationError(c *gin.Context, details any) {
	err := apperror.NewValidationError("Validation failed", details)
	Error(c, err)
}
