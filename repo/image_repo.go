package repo

import (
	"context"
	"github.com/star-find-cloud/star-mall/domain"
)

type ImageRepo interface {
	// UploadImage 上传图片元信息到数据库
	UploadImage(ctx context.Context, image *domain.Image) (int64, error)

	// GetByID 根据 imageID 获取图片元信息
	GetByID(ctx context.Context, imageID int64) (*domain.Image, error)

	// GetByOwner 根据 ownerType 和 ownerID 获取图片元信息
	GetByOwner(ctx context.Context, ownerType, ownerID int64) ([]*domain.Image, error)

	// GetPathByImageID 根据 imageID 获取图片的 ossPath
	GetPathByImageID(ctx context.Context, id int64) (string, error)

	// GetPathAndContentTypeByImageID 根据 imageID 获取图片的 ossPath 和 contentType
	GetPathAndContentTypeByImageID(ctx context.Context, id int64) (string, int64, error)

	// UpdatePath 更新图片的 ossPath
	UpdatePath(ctx context.Context, image domain.Image) error

	// Update 更新图片元信息
	Update(ctx context.Context, oldImageID int64, newImage *domain.Image) (int64, error)

	// Delete 删除图片元信息
	Delete(ctx context.Context, imageID int64) error

	// CheckImageExistsByID 检查图片元信息是否存在
	CheckImageExistsByID(ctx context.Context, id int64) (bool, error)
}
