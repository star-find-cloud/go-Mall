package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_const "github.com/star-find-cloud/star-mall/const"
	"github.com/star-find-cloud/star-mall/domain"
	"github.com/star-find-cloud/star-mall/pkg/database"
	applog "github.com/star-find-cloud/star-mall/pkg/logger"
	"time"
)

type MerchantRepo struct {
	db    database.Database
	cache *database.Redis
}

func NewMerchantRepo(db database.Database, cache *database.Redis) *MerchantRepo {
	return &MerchantRepo{db: db, cache: cache}
}

func (r *MerchantRepo) Create(ctx context.Context, merchant *domain.Merchant) (int64, error) {
	businessTypeJSON, err := json.Marshal(merchant.BusinessType)
	if err != nil {
		applog.AppLogger.Error("CreateMerchant", "err", err)
		return 0, err
	}

	sqlStr := "insert into shop.merchant (userID, name, phone, email, password, real_name, real_id, license_image_id, tag, cate_id, business_type, score, create_at, status) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	result, err := r.db.GetDB().ExecContext(ctx, sqlStr, merchant.UserID, merchant.Name, merchant.Phone, merchant.Email, merchant.Password, merchant.RealName, merchant.RealID, merchant.LicenseImageID, merchant.Tag, merchant.CateID, businessTypeJSON, merchant.Score, merchant.CreateAt, _const.StatusNotDeleted)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Error("GetMerchantByID", "err", err)
			return 0, err
		}
		applog.AppLogger.Error("CreateMerchant", "err", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Error("GetMerchantByID", "err", err)
			return 0, err
		}
		applog.AppLogger.Error("CreateMerchant", "err", err)
		return 0, err
	}
	return id, nil
}

func (r *MerchantRepo) GetMerchantByID(cxt context.Context, id int64) (*domain.Merchant, error) {
	var merchant = &domain.Merchant{}
	sqlStr := "select  name, phone, email, license_image_id, status, score, image_id, create_at, old_name, tag, cate_id  from shop.merchant where id=?"

	err := r.db.GetDB().GetContext(cxt, &merchant, sqlStr, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Error("GetMerchantByID", "err", err)
			return nil, err
		}
		applog.AppLogger.Error("GetMerchantByID", "err", err)
		return nil, err
	}
	return merchant, nil
}

func (r *MerchantRepo) GetMerchantByName(cxt context.Context, name string) (*[]domain.Merchant, error) {
	var merchant = &[]domain.Merchant{}
	sqlStr := "select name, phone, email, license_image_id, status, score, image_id, create_at, old_name, tag, cate_id  from shop.merchant where name like ? limit 15 offset ?"

	sqlName := "%" + name + "%"
	err := r.db.GetDB().SelectContext(cxt, &merchant, sqlStr, sqlName, 0)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Error("GetMerchantByName", "err", err)
			return nil, err
		}
		applog.AppLogger.Error("GetMerchantByName", "err", err)
		return nil, err
	}
	return merchant, nil
}

func (r *MerchantRepo) GetMerchantByEmail(ctx context.Context, email string) (*domain.Merchant, error) {
	var merchant = &domain.Merchant{}
	sqlStr := "select  name, phone, email, license_image_id, status, score, image_id, create_at, old_name, tag, cate_id  from shop.merchant where email=?"

	err := r.db.GetDB().GetContext(ctx, &merchant, sqlStr, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Error("GetMerchantByEmail", "err", err)
			return nil, err
		}
		applog.AppLogger.Error("GetMerchantByEmail", "err", err)
		return nil, err
	}
	return merchant, nil
}

func (r *MerchantRepo) GetMerchantByPhone(ctx context.Context, phone string) (*domain.Merchant, error) {
	var merchant = &domain.Merchant{}
	sqlStr := "select  name, phone, email, license_image_id, status, score, image_id, create_at, old_name, tag, cate_id  from shop.merchant where phone=?"

	err := r.db.GetDB().GetContext(ctx, &merchant, sqlStr, phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Error("GetMerchantByPhone", "err", err)
			return nil, err
		}
		applog.AppLogger.Error("GetMerchantByPhone", "err", err)
		return nil, err
	}
	return merchant, nil
}

func (r *MerchantRepo) Update(ctx context.Context, merchant *domain.Merchant) error {
	originalMerchant, err := r.GetMerchantByID(ctx, merchant.ID)
	if err != nil {
		applog.AppLogger.Error("GetOriginalMerchantFailed", "err", err)
		return err
	}

	if originalMerchant.Name != merchant.Name {
		merchant.OldName = originalMerchant.Name
	} else {
		merchant.OldName = originalMerchant.OldName
	}

	sqlStr := "update shop.merchant set name=?, old_name=?, phone=?, email=?, cate_id=?, business_type=? where id=?"

	_, err = r.db.GetDB().ExecContext(ctx, sqlStr, merchant.Name, merchant.OldName, merchant.Phone, merchant.Email, merchant.CateID, merchant.BusinessType, merchant.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Error("GetMerchantByID", "err", err)
			return err
		}
		applog.AppLogger.Error("UpdateMerchant", "err", err)
		return err
	}
	return nil
}

func (r *MerchantRepo) UpdateLicenseImage(ctx context.Context, merchantID, imageID int64) error {
	sqlStr := "update shop.merchant set license_image_id = ?, update_at = ? where id = ?"

	now := time.Now().Unix()
	_, err := r.db.GetDB().ExecContext(ctx, sqlStr, imageID, now, merchantID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Warnf("merchant update image err: %v", err)
			return fmt.Errorf("merchant update image err: %v", err)
		}
		applog.AppLogger.Errorf("merchant repo error: %v", err)
		return err
	}
	return nil
}

func (r *MerchantRepo) Delete(ctx context.Context, id int64) error {
	sqlStr := "update shop.merchant set status = ?, delete_at = ? where id = ?"

	now := time.Now().Unix()
	_, err := r.db.GetDB().ExecContext(ctx, sqlStr, _const.StatusDeleted, now, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Error("GetMerchantByID", "err", err)
			return err
		}
		applog.AppLogger.Error("DeleteMerchant", "err", err)
		return err
	}
	return nil
}

func (r *MerchantRepo) IsExistsByID(ctx context.Context, id int64) (bool, error) {
	var exists bool
	sqlStr := "select exists(select 1 from shop.merchant where id = ?)"

	err := r.db.GetDB().QueryRowxContext(ctx, sqlStr, id).Scan(&exists)
	if err != nil {
		return false, errors.New("failed to check if merchant exists")
	}

	return exists, nil
}

func (r *MerchantRepo) IsExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	sqlStr := "select exists(select 1 from shop.merchant where id = ?)"

	err := r.db.GetDB().QueryRowxContext(ctx, sqlStr, email).Scan(&exists)
	if err != nil {
		return false, errors.New("failed to check if merchant exists")
	}
	return exists, nil
}

func (r *MerchantRepo) IsExistsByPhone(ctx context.Context, phone string) (bool, error) {
	var exists bool
	sqlStr := "select exists(select 1 from shop.merchant where id = ?)"

	err := r.db.GetDB().QueryRowxContext(ctx, sqlStr, phone).Scan(&exists)
	if err != nil {
		return false, errors.New("failed to check if merchant exists")
	}
	return exists, nil
}
