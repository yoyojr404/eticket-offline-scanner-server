package seed

import (
	"context"
	"fmt"

	"github.com/yoyojr404/internal/common/datastore"
	"github.com/yoyojr404/internal/preflight/repository"
	"github.com/yoyojr404/internal/preflight/repository/ticketcategories"
	"github.com/yoyojr404/internal/preflight/repository/tickets"
	"github.com/yoyojr404/internal/preflight/repository/ticketverifications"
)

type ServiceRepoBucket struct {
	DDLRepo                   repository.DDLRepo
	TicketCategoriesReader    ticketcategories.TicketCategoriesReadRepository
	TicketCategoriesWriter    ticketcategories.TicketCategoriesWriteRepository
	TicketsReader             tickets.TicketReadRepository
	TicketsWriter             tickets.TicketWriteRepository
	TicketVerificationsReader ticketverifications.TicketVerificationReadRepository
	TicketVerificationsWriter ticketverifications.TicketVerificationWriteRepository
}

type seedServiceImpl struct {
	batch      int
	transactor datastore.Transactor
	ddlRepo    repository.DDLRepo
	repoBucket ServiceRepoBucket
}

func NewService(transactor datastore.Transactor, batch int, cfg ServiceRepoBucket) SeedService {
	return &seedServiceImpl{
		batch:      batch,
		transactor: transactor,
		ddlRepo:    cfg.DDLRepo,
		repoBucket: cfg,
	}
}

func (s *seedServiceImpl) Seed(ctx context.Context) error {
	return s.transactor.WithinTx(ctx, func(tx datastore.Executor) error {
		ddlRepo := s.repoBucket.DDLRepo.WithTx(tx)
		ticketCategoryWriter := s.repoBucket.TicketCategoriesWriter.WithTx(tx)
		ticketWriter := s.repoBucket.TicketsWriter.WithTx(tx)
		ticketVerificationWriter := s.repoBucket.TicketVerificationsWriter.WithTx(tx)

		if err := ddlRepo.DefineSchema(ctx); err != nil {
			return err
		}

		if err := batchSync(ctx, s.batch, s.repoBucket.TicketCategoriesReader, ticketCategoryWriter); err != nil {
			return fmt.Errorf("ticket categories sync error: %w", err)
		}

		if err := batchSync(ctx, s.batch, s.repoBucket.TicketsReader, ticketWriter); err != nil {
			return fmt.Errorf("tickets sync error: %w", err)
		}

		if err := batchSync(ctx, s.batch, s.repoBucket.TicketVerificationsReader, ticketVerificationWriter); err != nil {
			return fmt.Errorf("ticket verifications sync error: %w", err)
		}

		return nil
	})
}

type BatchReader[T any] interface {
	Read(ctx context.Context, process func(row T) error) error
}

type BatchWriter[T any] interface {
	Write(ctx context.Context, data []T) (int64, error)
}

func batchSync[T any](ctx context.Context, batchSize int, reader BatchReader[T], writer BatchWriter[T]) error {
	if batchSize <= 0 {
		batchSize = 500
	}

	bulk := make([]T, 0, batchSize)

	err := reader.Read(ctx, func(row T) error {
		bulk = append(bulk, row)
		if len(bulk) >= batchSize {
			if _, err := writer.Write(ctx, bulk); err != nil {
				return err
			}
			bulk = bulk[:0]
		}
		return nil
	})
	if err != nil {
		return err
	}

	if len(bulk) > 0 {
		if _, err := writer.Write(ctx, bulk); err != nil {
			return err
		}
	}

	return nil
}
