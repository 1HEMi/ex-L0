package service

import (
	"Ex-L0/internal/domain"
	"Ex-L0/internal/repository"
	"context"
)

type OrdersService struct{ Repository repository.OrdersRepository }

func (s *OrdersService) Upsert(ctx context.Context, o *domain.Order) error {
	return s.Repository.UpsertOrderTx(ctx, o)
}

func (s *OrdersService) GetByID(ctx context.Context, uid string) (*domain.Order, error) {
	return s.Repository.GetByID(ctx, uid)
}
