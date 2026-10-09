package ticketverifications

import (
	"context"

	"github.com/yoyojr404/internal/common/datastore"
	"github.com/yoyojr404/internal/common/entity"
)

type TicketVerificationWriteRepository interface {
	WithTx(tx datastore.Executor) TicketVerificationWriteRepository
	Write(ctx context.Context, data []entity.TicketVerification) (int64, error)
}

type TicketVerificationReadRepository interface {
	Read(ctx context.Context, process func(row entity.TicketVerification) error) error
}
