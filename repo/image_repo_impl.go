package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_const "github.com/star-find-cloud/star-mall/const"
	"github.com/star-find-cloud/star-mall/domain"
	"github.com/star-find-cloud/star-mall/pkg/database"
	log "github.com/star-find-cloud/star-mall/pkg/logger"
)

type ImageRepositoryImpl struct {
	db database.Database
}

func NewImageRepo(db database.Database) *ImageRepositoryImpl {
	return &ImageRepositoryImpl{db: db}
}

// UploadImage 上传图片元信息到数据库
func (r *ImageRepositoryImpl) UploadImage(ctx context.Context, image *domain.Image) (int64, error) {
	sqlStr := "insert into shop.images (imageID, ownerType, ownerID, path, sha256hash, isCompressed, content_type, create_at) values (?,?,?,?,?,?,?,?);"
	_, err := r.db.GetDB().ExecContext(ctx, sqlStr,
		image.ImageID,
		image.OwnerType,
		image.OwnerID,
		image.Path,
		image.IsCompressed,
		image.ContentType,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.MySQLLogger.Warnf("image creat err: %v", err)
			return 0, fmt.Errorf("image creat err: %v", err)
		}
		log.AppLogger.Errorf("image repo error: %v", err)
		return 0, fmt.Errorf("failed to get image: %w", err)
	}
	return image.ImageID, nil
}

// GetByID 根据 imageID 获取图片元信息
func (r *ImageRepositoryImpl) GetByID(ctx context.Context, imageID int64) (*domain.Image, error) {
	var image = &domain.Image{}
	sqlStr := "select imageID, ownerType, ownerID, path, isCompressed, content_type, create_at from shop.images where imageID = ?;"

	err := r.db.GetDB().GetContext(ctx, image, sqlStr, imageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.MySQLLogger.Warnf("image not found (id: %d)", imageID)
			return nil, fmt.Errorf("%w: image id %d", err, imageID)
		}
		log.AppLogger.Errorf("image repo error: %v", err)
		return nil, fmt.Errorf("failed to get image: %w", err)
	}
	return image, nil
}

// GetByOwner 根据 ownerType 和 ownerID 获取图片列表
func (r *ImageRepositoryImpl) GetByOwner(ctx context.Context, ownerType, ownerID int64) ([]*domain.Image, error) {
	var images = []*domain.Image{}
	sqlStr := "select imageID, ownerType, ownerID, path, isCompressed, content_type, create_at from shop.images where ownerType = ? and ownerID = ?;"
	err := r.db.GetDB().SelectContext(ctx, &images, sqlStr, ownerType, ownerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.MySQLLogger.Warnf("image not found (ownerType: %d, ownerID: %d)", ownerType, ownerID)
			return nil, fmt.Errorf("image not found (ownerType: %d, ownerID: %d)", ownerType, ownerID)
		}
		log.AppLogger.Errorf("image repo error: %v", err)
		return nil, fmt.Errorf("failed to get image: %w", err)
	}

	return images, nil
}

func (r *ImageRepositoryImpl) GetPathByImageID(ctx context.Context, id int64) (string, error) {
	sqlStr := "select path from shop.images where imageID = ?;"
	var path string

	err := r.db.GetDB().QueryRowContext(ctx, path, sqlStr, id).Scan()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.MySQLLogger.Warnf("image not found (id: %d)", id)
			return "", fmt.Errorf("image not found (id: %d)", id)
		}
		log.AppLogger.Errorf("image repo error: %v", err)
		return "", fmt.Errorf("failed to get image: %w", err)
	}

	return path, nil
}

func (r *ImageRepositoryImpl) GetPathAndContentTypeByImageID(ctx context.Context, id int64) (string, int64, error) {
	var image = &domain.Image{}
	sqlStr := "select path, content_type from shop.images where imageID = ?;"

	err := r.db.GetDB().GetContext(ctx, image, sqlStr, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.MySQLLogger.Warnf("image not found (id: %d)", id)
			return "", 0, fmt.Errorf("image not found (id: %d)", id)
		}
		log.AppLogger.Errorf("image repo error: %v", err)
		return "", 0, fmt.Errorf("failed to get image: %w", err)
	}

	return image.Path, image.ContentType, nil
}

// UpdatePath 更新图片路径
func (r *ImageRepositoryImpl) UpdatePath(ctx context.Context, image domain.Image) error {
	tx, err := r.db.GetDB().Begin()
	if err != nil {
		return fmt.Errorf("tx begin err: %v", err)
	}

	selectStr := "select path from shop.images where imageID = ?;"
	_, err = tx.QueryContext(ctx, selectStr, image.ImageID)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			log.AppLogger.Warnf("tx rollback err: %v", err)
		}
		return fmt.Errorf("image not found (id: %d)", image.ImageID)
	}

	updateStr := "update shop.images set path=?, isCompressed=?, content_type=? where imageID = ?"
	_, err = tx.ExecContext(ctx, updateStr, image.Path, image.IsCompressed, image.ContentType, image.ImageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.MySQLLogger.Warnf("image update err: %v", err)
			if err = tx.Rollback(); err != nil {
				log.AppLogger.Warnf("tx rollback err: %v", err)
			}
			return fmt.Errorf("image update err: %v", err)
		}
		log.AppLogger.Errorf("image repo error: %v", err)
		if err = tx.Rollback(); err != nil {
			log.AppLogger.Warnf("tx rollback err: %v", err)
		}
		return fmt.Errorf("image update err: %v", err)
	}
	return tx.Commit()
}

// Update 更新图片元信息
func (r *ImageRepositoryImpl) Update(ctx context.Context, oldImageID int64, newImage *domain.Image) (int64, error) {
	// 开启事务
	tx, err := r.db.GetDB().Begin()
	if err != nil {
		return 0, fmt.Errorf("tx begin err: %v", err)
	}

	// 将旧数据标记为删除
	selectStr := "update shop.images set status = ? where imageID = ?"
	_, err = tx.ExecContext(ctx, selectStr, _const.StatusDeleted, oldImageID)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			log.AppLogger.Warnf("tx rollback err: %v", err)
		}
		return 0, fmt.Errorf("image not found (id: %d)", oldImageID)
	}

	// 插入新数据
	insertStr := "insert into shop.images (imageID, ownerType, ownerID, path, sha256hash, isCompressed, content_type, create_at) values (?,?,?,?,?,?,?,?);"
	_, err = tx.ExecContext(ctx, insertStr,
		newImage.ImageID,
		newImage.Path,
		newImage.OwnerType,
		newImage.OwnerID,
		newImage.IsCompressed,
		newImage.ContentType,
	)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.AppLogger.Warnf("tx rollback err: %v", err)
		}
		if errors.Is(err, sql.ErrNoRows) {
			log.MySQLLogger.Warnf("image update err: %v", err)
			return 0, fmt.Errorf("image update err: %v", err)
		}
		log.AppLogger.Errorf("image repo error: %v", err)
		return 0, fmt.Errorf("image update err: %v", err)
	}
	err = tx.Commit()
	if err != nil {
		log.AppLogger.Errorf("image repo error: %v", err)
		return 0, fmt.Errorf("image update err: %v", err)
	}
	return newImage.ImageID, nil
}

// Delete 删除图片元信息
func (r *ImageRepositoryImpl) Delete(ctx context.Context, imageID int64) error {
	exists, err := r.CheckImageExistsByID(ctx, imageID)
	if !exists {
		return fmt.Errorf("image not found (id: %d)", imageID)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.MySQLLogger.Errorf("image check exists error :%v", err)
			return err
		}
		log.AppLogger.Errorf("image repo error: %v", err)
		return fmt.Errorf("image delete err: %v", err)
	}

	deleteStr := "update shop.images set status = ? where imageID = ?;"
	_, err = r.db.GetDB().ExecContext(ctx, deleteStr, imageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.MySQLLogger.Warnf("image delete error :%v", err)
			return fmt.Errorf("image delete err: %v", err)
		}
		log.AppLogger.Errorf("image repo error: %v", err)
		return fmt.Errorf("image delete err: %v", err)
	}

	return nil
}

// CheckImageExistsByID 根据ID判断是否存在根据 id 检查图片是否存在
func (r *ImageRepositoryImpl) CheckImageExistsByID(ctx context.Context, id int64) (bool, error) {
	var sqlStr = "select exists (select 1 from shop.images where imageID = ?)"
	var exist bool

	var err = r.db.GetDB().GetContext(ctx, &exist, sqlStr, id)
	return exist, err
}
