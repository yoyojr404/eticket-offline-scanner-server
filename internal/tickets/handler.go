package tickets

import (
	"net/http"
	"strconv"

	"github.com/yoyojr404/pkg/apperror"
	"github.com/yoyojr404/pkg/pagination"
	"github.com/yoyojr404/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAll(c *gin.Context) {
	p := pagination.GetParams(c)
	result, err := h.service.GetAll(c.Request.Context(), p.Page, p.PerPage)
	if err != nil {
		c.Error(err)
		return
	}

	meta := response.NewPaginationMeta(result.CurrentPage, result.TotalPerPage, int(result.Total))
	resp := ToResponseList(result.Entities)

	response.SuccessWithMeta(c, http.StatusOK, "Tickets retrieved successfully", resp, meta)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperror.NewBadRequest("Invalid ticket ID format", err))
		return
	}

	item, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, "Ticket retrieved successfully", ToResponse(item))
}
