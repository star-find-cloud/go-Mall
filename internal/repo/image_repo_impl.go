package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/star-find-cloud/star-mall/model"
	applog "github.com/star-find-cloud/star-mall/pkg/logger"
)

type ImageRepositoryImpl struct {
	db    *sqlx.DB
	cache *redis.Client
}

func NewImageRepositoryImpl(db *sqlx.DB) *ImageRepositoryImpl {
	return &ImageRepositoryImpl{db: db}
}

func (r *ImageRepositoryImpl) Create(ctx context.Context, image *model.Image) (string, error) {
	sqlStr := "insert into shop.images (imageID, ownerType, ownerID, ossPath, SHA256Hash, isCompressed) values (?,?,?,?,?,?);"
	_, err := r.db.ExecContext(ctx, sqlStr,
		image.ImageID,
		image.OwnerType,
		image.OwnerID,
		image.OssPath,
		image.SHA256Hash,
		image.IsCompressed)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Warnf("image creat err: %v", err)
			return "", fmt.Errorf("image creat err: %v", err)
		}
		applog.AppLogger.Errorf("image repo error: %v", err)
		return "", fmt.Errorf("failed to get image: %w", err)
	}
	return image.OssPath, nil
}

func (r *ImageRepositoryImpl) GetByID(ctx context.Context, imageID string) (*model.Image, error) {
	var image *model.Image
	sqlStr := "select ossPath, SHA256Hash, isCompressed from shop.images where imageID = ?;"

	err := r.db.GetContext(ctx, image, sqlStr, imageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Warnf("image not found (id: %s)", imageID)
			return nil, fmt.Errorf("%w: image id %s", err, imageID)
		}
		applog.AppLogger.Errorf("image repo error: %v", err)
		return nil, fmt.Errorf("failed to get image: %w", err)
	}
	return image, nil
}

func (r *ImageRepositoryImpl) GetByOwner(ctx context.Context, ownerType, ownerID string) ([]*model.Image, error) {
	var images []*model.Image
	sqlStr := "select ossPath, SHA256Hash, isCompressed from shop.images where ownerType = ? and ownerID = ?;"
	err := r.db.SelectContext(ctx, &images, sqlStr, ownerType, ownerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Warnf("image not found (ownerType: %s, ownerID: %s)", ownerType, ownerID)
			return nil, fmt.Errorf("%w: ownerType %s, ownerID %s", err, ownerType, ownerID)
		}
		applog.AppLogger.Errorf("image repo error: %v", err)
		return nil, fmt.Errorf("failed to get image: %w", err)
	}

	return images, nil
}

func (r *ImageRepositoryImpl) GetByHash(ctx context.Context, hash string) (*model.Image, error) {
	var image *model.Image
	sqlStr := "select ownerType, ownerID, ossPath, isCompressed from shop.images where SHA256Hash = ?;"
	err := r.db.GetContext(ctx, image, sqlStr, hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Warnf("image not found (hash: %s)", hash)
			return nil, fmt.Errorf("%w: hash %s", err, hash)
		}
		applog.AppLogger.Errorf("image repo error: %v", err)
		return nil, fmt.Errorf("failed to get image: %w", err)
	}

	return image, nil
}
