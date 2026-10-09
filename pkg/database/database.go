package database

import (
	"fmt"
	"log"
	"time"

	"github.com/yoyojr404/pkg/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// InitDB initializes and configures a SQL database connection pool based on configuration.
// If the driver is set to "memory", it safely returns nil to signal in-memory fallback.
func InitDB(cfg *config.Config) (*sqlx.DB, error) {
	if cfg.DBDriver == "memory" || cfg.DBDriver == "" {
		log.Println("[INFO] Database driver set to 'memory' - running with in-memory storage")
		return nil, nil
	}

	db, err := sqlx.Open(cfg.DBDriver, fmt.Sprintf(
		"%s?_journal_mode=WAL&_busy_timeout=5000&_synchronous=NORMAL&_cache_size=-64000&_foreign_keys=ON",
		cfg.DBPath,
	))
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection (%s): %w", cfg.DBDriver, err)
	}

	// Connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping database at %s: %w", cfg.DBPath, err)
	}

	log.Printf("[INFO] Database connection established: %s (%s)", cfg.DBPath, cfg.DBDriver)
	return db, nil
}

// Close gracefully closes the database connection.
func Close(db *sqlx.DB) error {
	if db == nil {
		return nil
	}
	log.Println("[INFO] Closing database connection pool...")
	return db.Close()
}
