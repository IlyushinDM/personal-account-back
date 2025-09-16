package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Dbtx реализует интерфейс Transactor.
type dbtx struct {
	db *sqlx.DB
}

// NewTransactor создает новый экземпляр Transactor.
func NewTransactor(db *sqlx.DB) Transactor {
	return &dbtx{db: db}
}

// WithinTransaction выполняет функцию fn внутри транзакции БД.
// Если fn возвращает ошибку, транзакция откатывается.
// Если fn завершается без ошибок, транзакция коммитится.
func (d *dbtx) WithinTransaction(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
			}
		} else {
			err = tx.Commit()
			if err != nil {
				err = fmt.Errorf("failed to commit transaction: %w", err)
			}
		}
	}()

	err = fn(tx)
	return err
}
