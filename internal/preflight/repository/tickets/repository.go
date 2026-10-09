package tickets

import (
	"context"

	"github.com/yoyojr404/internal/common/datastore"
	"github.com/yoyojr404/internal/common/entity"
)

type TicketWriteRepository interface {
	WithTx(tx datastore.Executor) TicketWriteRepository
	Write(ctx context.Context, data []entity.Ticket) (int64, error)
}

type TicketReadRepository interface {
	Read(ctx context.Context, process func(row entity.Ticket) error) error
}
