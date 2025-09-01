package repo

import (
	"context"
	"github.com/star-find-cloud/star-mall/domain"
)

type ProductRepo interface {
	Create(ctx context.Context, product *domain.Product) (int64, error)
	//CreateMore(ctx context.Context, products *[]domain.Product) error
	GetByID(ctx context.Context, id int64) (*domain.Product, error)
	GetByMerchantID(ctx context.Context, shopID int64) ([]*domain.Product, error)
	GetCateID(ctx context.Context, id int64) (int64, error)
	GetByCateIDsAndHot(ctx context.Context, ids []int64) ([]domain.Product, error)
	//GetByShopName(ctx context.Context, shopName string) ([]*domain.Product, error)
	GetByTitleAndKeywords(ctx context.Context, title, keywords string, offset int) ([]*domain.Product, error)
	GetMerchantID(ctx context.Context, id int64) (int64, error)
	SearchByMsg(ctx context.Context, msg string) ([]domain.Product, error)
	//GetByKeywords(ctx context.Context, keywords string) ([]*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id int64) error
}
