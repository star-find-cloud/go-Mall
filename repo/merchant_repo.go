package repo

import (
	"context"
	"github.com/star-find-cloud/star-mall/domain"
)

type MerchantRepository interface {
	Create(ctx context.Context, merchant *domain.Merchant) (int64, error)
	Update(ctx context.Context, merchant *domain.Merchant) error
	Delete(ctx context.Context, id int64) error
	GetMerchantByID(cxt context.Context, id int64) (*domain.Merchant, error)
	GetMerchantByName(cxt context.Context, name string) (*[]domain.Merchant, error)
	GetMerchantByEmail(ctx context.Context, email string) (*domain.Merchant, error)
	GetMerchantByPhone(ctx context.Context, phone string) (*domain.Merchant, error)
	IsExistsByID(ctx context.Context, id int64) (bool, error)
	IsExistsByEmail(ctx context.Context, email string) (bool, error)
	IsExistsByPhone(ctx context.Context, phone string) (bool, error)
	UpdateLicenseImage(ctx context.Context, merchantID, imageID int64) error
}
