package ticketverification

import (
	"context"
	"database/sql"
	"errors"

	"github.com/yoyojr404/internal/common/datastore"
	"github.com/yoyojr404/internal/common/entity"
)

type Repository interface {
	WithTx(tx datastore.Executor) Repository
	FindTicketVerificationByID(ctx context.Context, id int64) (*entity.TicketVerification, error)
	UpdateVerifiedTicketVerification(ctx context.Context, id int64) (bool, error)
}

type sqlRepository struct {
	db datastore.Executor
}

func NewSQLRepository(db datastore.Executor) Repository {
	return &sqlRepository{db: db}
}

func (d *sqlRepository) WithTx(tx datastore.Executor) Repository {
	if d.db == nil {
		return d
	}

	return &sqlRepository{db: tx}
}

func (s *sqlRepository) UpdateVerifiedTicketVerification(ctx context.Context, id int64) (bool, error) {
	result, err := s.db.ExecContext(
		ctx,
		"UPDATE ticket_verifications SET verified = ?, verified_at = CURRENT_TIMESTAMP WHERE id = ?",
		true,
		id,
	)

	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func (d *sqlRepository) FindTicketVerificationByID(ctx context.Context, id int64) (*entity.TicketVerification, error) {
	query := `
		SELECT 
			tv.id,
			tv.ticket_id,
			tv.verified,
			tv.created_at,
			tv.verified_at
		FROM ticket_verifications tv 
		WHERE tv.id = ?
	`
	var t entity.TicketVerification
	err := d.db.GetContext(ctx, &t, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}
