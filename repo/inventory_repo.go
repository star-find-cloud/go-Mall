package repo

import (
	"context"
	"github.com/star-find-cloud/star-mall/domain"
)

// InventoryRepo 库存数据库接口
type InventoryRepo interface {
	Create(ctx context.Context, inventory *domain.Inventory) error
	Update(ctx context.Context, inventory *domain.Inventory) error
	GetByID(ctx context.Context, id int64) (*domain.Inventory, error)
	GetByMerchantID(ctx context.Context, MerchantID int64) ([]*domain.Inventory, error)
	Deduction(ctx context.Context, ProductID int64, count int64) error
}
