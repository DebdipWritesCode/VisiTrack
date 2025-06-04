package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier // Embeds all generated SQLC methods
	// Add transactional methods here, like:
	// BookAppointmentTx(ctx context.Context, arg BookAppointmentTxParams) (BookAppointmentTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store instance
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx runs a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
