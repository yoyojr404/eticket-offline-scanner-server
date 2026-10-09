package ticketcategories

import (
	"context"

	"github.com/yoyojr404/internal/common/datastore"
	"github.com/yoyojr404/internal/common/entity"
)

type TicketCategoriesWriteRepository interface {
	WithTx(tx datastore.Executor) TicketCategoriesWriteRepository
	Write(ctx context.Context, data []entity.TicketCategory) (int64, error)
}

type TicketCategoriesReadRepository interface {
	Read(ctx context.Context, process func(row entity.TicketCategory) error) error
}
