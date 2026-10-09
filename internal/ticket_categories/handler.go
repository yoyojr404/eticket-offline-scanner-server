package ticket_categories

import (
	"net/http"
	"strconv"

	"github.com/yoyojr404/pkg/apperror"
	"github.com/yoyojr404/pkg/pagination"
	"github.com/yoyojr404/pkg/response"
	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for ticket categories.
// It is strictly focused on request processing, invoking services, and rendering responses.
type Handler struct {
	service Service
}

// NewHandler creates a new Handler with injected Service (DI).
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetAll handles listing ticket categories with pagination.
func (h *Handler) GetAll(c *gin.Context) {
	p := pagination.GetParams(c)
	result, err := h.service.GetAll(c.Request.Context(), p.Page, p.PerPage)
	if err != nil {
		c.Error(err)
		return
	}

	meta := response.NewPaginationMeta(result.CurrentPage, result.TotalPerPage, int(result.Total))
	resp := ToResponseList(result.Entities)

	response.SuccessWithMeta(c, http.StatusOK, "Categories retrieved successfully", resp, meta)
}

// GetByID handles retrieving a single ticket category by ID.
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperror.NewBadRequest("Invalid category ID format", err))
		return
	}

	item, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, "Category retrieved successfully", item)
}
