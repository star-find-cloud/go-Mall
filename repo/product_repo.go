package repo

import (
	"context"
	"github.com/star-find-cloud/star-mall/model"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	GetByID(ctx context.Context, id int) (*model.Product, error)
	GetByShopID(ctx context.Context, shopID int) ([]*model.Product, error)

	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id int) error
}
