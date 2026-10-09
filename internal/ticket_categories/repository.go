package ticket_categories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/yoyojr404/internal/common/entity"
	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

// Repository interface defines data access operations for ticket categories.
type Repository interface {
	FindAll(ctx context.Context, offset, limit int) ([]entity.TicketCategory, int, error)
	FindByID(ctx context.Context, id int64) (*entity.TicketCategory, error)
}

// sqlRepository implements Repository using standard database/sql (*sql.DB).
type sqlRepository struct {
	db *sqlx.DB
}

// NewSQLRepository creates a new SQL-based repository.
func NewSQLRepository(db *sqlx.DB) Repository {
	return &sqlRepository{db: db}
}

func (r *sqlRepository) FindAll(ctx context.Context, offset, limit int) ([]entity.TicketCategory, int, error) {
	eg, egCtx := errgroup.WithContext(ctx)

	queryTotal := "SELECT COUNT(*) FROM ticket_categories"
	querySelect := "SELECT id, name, created_at FROM ticket_categories ORDER BY id DESC LIMIT ? OFFSET ?"

	var total int
	eg.Go(func() error {
		err := r.db.GetContext(egCtx, &total, queryTotal)
		if err != nil {
			return err
		}
		return nil
	})

	var categories []entity.TicketCategory
	eg.Go(func() error {
		err := r.db.SelectContext(egCtx, &categories, querySelect, limit, offset)
		if err != nil {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *sqlRepository) FindByID(ctx context.Context, id int64) (*entity.TicketCategory, error) {
	var cat entity.TicketCategory
	err := r.db.GetContext(ctx, &cat, "SELECT id, name, created_at FROM ticket_categories WHERE id = ?", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &cat, nil
}
