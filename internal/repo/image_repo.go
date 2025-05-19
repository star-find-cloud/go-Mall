package repo

import (
	"context"
	"github.com/star-find-cloud/star-mall/model"
)

type ImageRepository interface {
	Create(ctx context.Context, image *model.Image) (string, error)
	GetByID(ctx context.Context, imageID string) (*model.Image, error)
	GetByOwner(ctx context.Context, ownerType, ownerID string) ([]*model.Image, error)
	GetByHash(ctx context.Context, hash string) (*model.Image, error)
}
