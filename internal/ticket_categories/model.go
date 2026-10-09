package ticket_categories

import "github.com/yoyojr404/internal/common/entity"

type TicketCategoryResult struct {
	Entities     []entity.TicketCategory
	CurrentPage  int
	TotalPerPage int
	Total        int64
}
