package ticket_categories

import (
	"time"

	"github.com/yoyojr404/internal/common/entity"
)

// CategoryResponse defines the response structure for category data.
type CategoryResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse maps a TicketCategory entity to CategoryResponse DTO.
func ToResponse(e *entity.TicketCategory) CategoryResponse {
	return CategoryResponse{
		ID:        e.ID,
		Name:      e.Name,
		CreatedAt: e.CreatedAt,
	}
}

// ToResponseList maps a slice of TicketCategory entities to a slice of CategoryResponse DTOs.
func ToResponseList(entities []entity.TicketCategory) []CategoryResponse {
	res := make([]CategoryResponse, len(entities))
	for i, e := range entities {
		res[i] = ToResponse(&e)
	}
	return res
}
