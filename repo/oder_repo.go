package repo

import (
	"context"
	"github.com/star-find-cloud/star-mall/domain"
)

type OrderRepo interface {
	Create(ctx context.Context, order *domain.Order, orderItem *domain.OrderItem) error
	GetByID(ctx context.Context, orderID int64) (*domain.Order, *domain.OrderItem, error)
	GetByCreatedAt(ctx context.Context, createdAt int64) ([]*domain.Order, error)
	UpdateStatus(ctx context.Context, order *domain.Order) error
	Delete(ctx context.Context, order *domain.Order) error
}
