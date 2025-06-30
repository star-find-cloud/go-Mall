package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/star-find-cloud/star-mall/domain"
	"github.com/star-find-cloud/star-mall/pkg/database"
	applog "github.com/star-find-cloud/star-mall/pkg/logger"
)

type ImageRepositoryImpl struct {
	db    *sqlx.DB
	cache *database.Redis
}

func NewImageRepo(db *sqlx.DB) *ImageRepositoryImpl {
	return &ImageRepositoryImpl{db: db}
}

// Save 保存图片
func (r *ImageRepositoryImpl) UploadImage(ctx context.Context, image *domain.Image) (int64, error) {
	sqlStr := "insert into shop.images (imageID,  SHA256Hash, isCompressed, file_path, content_type) values (?,?,?,?,?);"
	_, err := r.db.ExecContext(ctx, sqlStr,
		image.ImageID,
		image.SHA256Hash,
		image.IsCompressed,
		image.FilePath,
		image.ContentType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Warnf("image creat err: %v", err)
			return 0, fmt.Errorf("image creat err: %v", err)
		}
		applog.AppLogger.Errorf("image repo error: %v", err)
		return 0, fmt.Errorf("failed to get image: %w", err)
	}
	return image.ImageID, nil
}

//func (r *ImageRepositoryImpl) Save(ctx context.Context, image *domain.Image) (string, int64, error) {
//	sqlStr := "insert into shop.images (imageID, ownerType, ownerID, ossPath, SHA256Hash, isCompressed, raw) values (?,?,?,?,?,?,?);"
//	_, err := r.db.ExecContext(ctx, sqlStr,
//		image.ImageID,
//		image.OwnerType,
//		image.OwnerID,
//		image.OssPath,
//		image.SHA256Hash,
//		image.IsCompressed)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			applog.MySQLLogger.Warnf("image creat err: %v", err)
//			return "", 0, fmt.Errorf("image creat err: %v", err)
//		}
//		applog.AppLogger.Errorf("image repo error: %v", err)
//		return "", 0, fmt.Errorf("failed to get image: %w", err)
//	}
//	return image.OssPath, image.ImageID, nil
//}
//

func (r *ImageRepositoryImpl) GetByID(ctx context.Context, imageID int64) (*domain.Image, error) {
	var image = &domain.Image{}
	sqlStr := "select ownerID, ownerType, SHA256Hash, isCompressed, file_path, content_type from shop.images where imageID = ?;"

	err := r.db.GetContext(ctx, image, sqlStr, imageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Warnf("image not found (id: %d)", imageID)
			return nil, fmt.Errorf("%w: image id %d", err, imageID)
		}
		applog.AppLogger.Errorf("image repo error: %v", err)
		return nil, fmt.Errorf("failed to get image: %w", err)
	}
	return image, nil
}

//
//func (r *ImageRepositoryImpl) GetByOwner(ctx context.Context, ownerType, ownerID string) ([]*domain.Image, error) {
//	var images = []*domain.Image{}
//	sqlStr := "select ossPath, SHA256Hash, isCompressed, raw from shop.images where ownerType = ? and ownerID = ?;"
//	err := r.db.SelectContext(ctx, &images, sqlStr, ownerType, ownerID)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			applog.MySQLLogger.Warnf("image not found (ownerType: %s, ownerID: %s)", ownerType, ownerID)
//			return nil, errors.New("image not found (ownerType: %s, ownerID: %s)")
//		}
//		applog.AppLogger.Errorf("image repo error: %v", err)
//		return nil, fmt.Errorf("failed to get image: %w", err)
//	}
//
//	return images, nil
//}
//
//func (r *ImageRepositoryImpl) GetByHash(ctx context.Context, hash string) (*domain.Image, error) {
//	var image *domain.Image
//	sqlStr := "select ownerType, ownerID, ossPath, isCompressed, raw from shop.images where SHA256Hash = ?;"
//	err := r.db.GetContext(ctx, image, sqlStr, hash)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			applog.MySQLLogger.Warnf("image not found (hash: %s)", hash)
//			return nil, fmt.Errorf("%w: hash %s", err, hash)
//		}
//		applog.AppLogger.Errorf("image repo error: %v", err)
//		return nil, fmt.Errorf("failed to get image: %w", err)
//	}
//
//	return image, nil
//}
//
//func (r *ImageRepositoryImpl) UpDate(ctx context.Context, image domain.Image) error {
//	tx, _ := r.db.Begin()
//	_, err := tx.QueryContext(ctx, "select ossPath from shop.images where imageID = ?;", image.ImageID)
//	if err != nil {
//		tx.Rollback()
//	}
//
//	sqlStr := "update shop.images set ossPath = ?, SHA256Hash = ?, isCompressed = ? where imageID = ?"
//	_, err = tx.ExecContext(ctx, sqlStr, image.OssPath, image.SHA256Hash, image.IsCompressed)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			applog.MySQLLogger.Warnf("image update err: %v", err)
//			return fmt.Errorf("image update err: %v", err)
//		}
//		applog.AppLogger.Errorf("image repo error: %v", err)
//		return fmt.Errorf("failed to get image: %w", err)
//	}
//	return tx.Commit()
//}
//
//func (r *ImageRepositoryImpl) Delete(ctx context.Context, imageID int) error {
//	sqlStr := "delete from shop.images where imageID = ?;"
//	_, err := r.db.ExecContext(ctx, sqlStr, imageID)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			applog.MySQLLogger.Warnf("image delete err: %v", err)
//			return fmt.Errorf("image delete err: %v", err)
//		}
//		applog.AppLogger.Errorf("image repo error: %v", err)
//		return fmt.Errorf("failed to get image: %w", err)
//	}
//	return nil
//}
