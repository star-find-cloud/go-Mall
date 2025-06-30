package repo

import (
	"context"
	"github.com/star-find-cloud/star-mall/domain"
)

type CartRepo interface {
	Create(ctx context.Context, cart *domain.Cart) (int64, error)
	GetByID(ctx context.Context, id int64) (*domain.Cart, error)
	GetByUserID(ctx context.Context, userID int64) (*domain.Cart, error)
	AddProduct(ctx context.Context, CartItems *domain.CartItemVO) error
}
