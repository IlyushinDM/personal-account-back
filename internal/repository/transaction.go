package repository

import (
	"context"

	"gorm.io/gorm"
)

// Dbtx реализует интерфейс Transactor.
type dbtx struct {
	db *gorm.DB
}

// NewTransactor создает новый экземпляр Transactor.
func NewTransactor(db *gorm.DB) Transactor {
	return &dbtx{db: db}
}

// WithinTransaction выполняет функцию fn внутри транзакции БД.
func (d *dbtx) WithinTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return d.db.WithContext(ctx).Transaction(fn)
}
