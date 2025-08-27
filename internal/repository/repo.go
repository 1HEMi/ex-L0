package repository

import (
	"Ex-L0/internal/domain"
	"context"
	"errors"
)

var ErrNotFound = errors.New("not found")

type OrdersReader interface {
	GetByID(ctx context.Context, uid string) (*domain.Order, error)
	SelectRecentIDS(ctx context.Context, n int) ([]string, error)
}

type OrdersWriter interface {
	UpsertOrderTx(ctx context.Context, o *domain.Order) error
}

type OrdersRepository interface {
	OrdersReader
	OrdersWriter
}
