package repo

import (
	"context"
	"github.com/star-find-cloud/star-mall/domain"
)

type ImageRepo interface {
	UploadImage(ctx context.Context, image *domain.Image) (int64, error)
	GetByID(ctx context.Context, imageID int64) (*domain.Image, error)
	//GetByOwner(ctx context.Context, ownerType, ownerID string) ([]*domain.Image, error)
	//GetByHash(ctx context.Context, hash string) (*domain.Image, error)
	//UpDate(ctx context.Context, image domain.Image) error
	//Delete(ctx context.Context, imageID int) error
}
