package ticketcategories

import (
	"context"
	"fmt"

	"github.com/yoyojr404/internal/common/datastore"
	"github.com/yoyojr404/internal/common/entity"
	"github.com/jmoiron/sqlx"
)

type DBWriterRepository struct {
	db datastore.Executor
}

func NewDBWriterRepository(db *sqlx.DB) *DBWriterRepository {
	return &DBWriterRepository{db: db}
}

func (d *DBWriterRepository) WithTx(tx datastore.Executor) TicketCategoriesWriteRepository {
	if d.db == nil {
		return d
	}

	return &DBWriterRepository{db: tx}
}

func (d *DBWriterRepository) Write(ctx context.Context, data []entity.TicketCategory) (int64, error) {
	if len(data) == 0 {
		return 0, nil
	}

	query := `
        INSERT INTO ticket_categories (id, name)
        VALUES (:id, :name)
        ON CONFLICT (id) DO NOTHING;
    `

	res, err := d.db.NamedExecContext(ctx, query, data)
	if err != nil {
		return 0, fmt.Errorf("d.db.NamedExecContext: %w", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return 0, fmt.Errorf("res.RowsAffected: %w", err)
	}

	return rowsAffected, nil
}
