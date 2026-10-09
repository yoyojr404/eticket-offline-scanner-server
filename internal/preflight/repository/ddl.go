package repository

import (
	"context"

	"github.com/yoyojr404/internal/common/datastore"
	"github.com/jmoiron/sqlx"
)

type DDLRepo interface {
	WithTx(tx datastore.Executor) DDLRepo
	DefineSchema(ctx context.Context) error
}

type ddlRepoImpl struct {
	db datastore.Executor
}

func NewDDLRepo(db *sqlx.DB) DDLRepo {
	return &ddlRepoImpl{db: db}
}

func (d *ddlRepoImpl) WithTx(tx datastore.Executor) DDLRepo {
	if d.db == nil {
		return d
	}

	return &ddlRepoImpl{db: tx}
}

func (d *ddlRepoImpl) DefineSchema(ctx context.Context) error {
	ticketCategorySchema := `
CREATE TABLE IF NOT EXISTS ticket_categories (
  id INT AUTO_INCREMENT PRIMARY KEY, 
  name VARCHAR(255) NOT NULL, 
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`

	ticketSchema := `
CREATE TABLE IF NOT EXISTS tickets (
  id INT AUTO_INCREMENT PRIMARY KEY, 
  code VARCHAR(255) UNIQUE NOT NULL, 
  category_id INT FK REFERENCES ticket_categories(id) ON UPDATE CASCADE ON DELETE RESTRICT
);
`

	ticketVerificationSchema := `
CREATE TABLE IF NOT EXISTS ticket_verifications(
  id INT AUTO_INCREMENT PRIMARY KEY, 
  ticket_id INT FK REFERENCES tickets(id) ON UPDATE CASCADE ON DELETE RESTRICT, 
  verified TINYINT NOT NULL DEFAULT 0, 
  verified_at DATETIME NULL, 
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`

	if _, err := d.db.ExecContext(ctx, ticketCategorySchema); err != nil {
		return err
	}
	if _, err := d.db.ExecContext(ctx, ticketSchema); err != nil {
		return err
	}
	if _, err := d.db.ExecContext(ctx, ticketVerificationSchema); err != nil {
		return err
	}

	return nil
}
