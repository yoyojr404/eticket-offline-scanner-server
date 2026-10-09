package tickets

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type Repository interface {
	FindAll(ctx context.Context, offset, limit int) ([]Ticket, int, error)
	FindByID(ctx context.Context, id int64) (*Ticket, error)
}

type sqlRepository struct {
	db *sqlx.DB
}

func NewSQLRepository(db *sqlx.DB) Repository {
	return &sqlRepository{db: db}
}

func (r *sqlRepository) FindAll(ctx context.Context, offset, limit int) ([]Ticket, int, error) {
	eg, gCtx := errgroup.WithContext(ctx)

	var total int
	eg.Go(func() error {
		err := r.db.GetContext(gCtx, &total, "SELECT COUNT(*) FROM tickets")
		if err != nil {
			return err
		}
		return nil
	})

	var list []Ticket
	eg.Go(func() error {
		query := `
		SELECT 
			t.id as id,
			t.code as code,
			t.category_id as category_id,
			c.name as category_name,
			tv.id as verification_id,
			tv.verified as verification_verified,
			tv.verified_at as verification_verified_at
		FROM tickets t
		JOIN ticket_categories c ON t.category_id = c.id
		JOIN ticket_verifications tv ON t.id = tv.ticket_id
		LIMIT ? OFFSET ?
	`
		if err := r.db.SelectContext(gCtx, &list, query, limit, offset); err != nil {
			return err
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *sqlRepository) FindByID(ctx context.Context, id int64) (*Ticket, error) {
	query := `
		SELECT 
			t.id as id, 
			t.code as code, 
			t.category_id as category_id,
			c.name as category_name,
			tv.id as verification_id,
			tv.verified as verification_verified,
			tv.verified_at as verification_verified_at
		FROM tickets t
		JOIN ticket_categories c ON t.category_id = c.id
		JOIN ticket_verifications tv ON t.id = tv.ticket_id
		WHERE t.id = ?
	`
	var t Ticket
	err := r.db.GetContext(ctx, &t, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}
