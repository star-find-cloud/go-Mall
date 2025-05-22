package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/star-find-cloud/star-mall/model"
	applog "github.com/star-find-cloud/star-mall/pkg/logger"
	"time"
)

type UserRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepositoryImpl(db *sqlx.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (r UserRepositoryImpl) GetByID(ctx context.Context, id int64) (*model.User, error) {
	var user *model.User
	sqlStr := "select id,name,image,sex,last_ip,image,is_vip from shop.user where id = ? LIMIT 1;"

	err := r.db.GetContext(ctx, user, sqlStr, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Warnf("user not found (id: %d)", id)
			return nil, fmt.Errorf("%w: user id %d", err, id)
		}
		applog.AppLogger.Errorf("user repo error: %v", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user *model.User
	sqlStr := "select id,name, sex,image,last_ip,image,is_vip from shop.user where email = ? LIMIT 1;"

	err := r.db.GetContext(ctx, user, sqlStr, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Warnf("user not found (email: %s)", email)
			return nil, fmt.Errorf("%w: user id %s", err, email)
		}
		applog.AppLogger.Errorf("user repo error: %v", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r UserRepositoryImpl) Create(ctx context.Context, user *model.User) (int64, error) {
	sqlStr := "insert into shop.user (name,password,email,phone,sex,create_time, update_time, status, last_ip, image, is_vip) values (?,?,?,?,?,?,?,?,?,?,?)"

	result, err := r.db.ExecContext(ctx, sqlStr,
		user.Name,
		user.Password,
		user.Email,
		user.Phone,
		user.Sex,
		user.CreateTime,
		user.UpdateTime,
		user.Status,
		user.LastIP,
		user.ImageID,
		user.IsVip)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Warnf("user creat err: %v", err)
			return 0, fmt.Errorf("user creat err: %v", err)
		}
		applog.AppLogger.Errorf("user repo error: %v", err)
		return 0, fmt.Errorf("failed to get user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		applog.AppLogger.Errorf("user repo error: %v", err)
		return 0, fmt.Errorf("failed to get user: %w", err)
	}

	return id, nil
}

func (r UserRepositoryImpl) Update(ctx context.Context, user *model.User) error {
	sqlStr := "update shop.user set name = ?,email = ?,phone = ?, sex = ?, update_time = ?, image = ? where id = ?;"

	_, err := r.db.ExecContext(ctx, sqlStr, user.Name, user.Email, user.Phone, user.UpdateTime, user.ImageID, user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Warnf("user update err: %v", err)
			return fmt.Errorf("user update err: %v", err)
		}
		applog.AppLogger.Errorf("user repo error: %v", err)
		return fmt.Errorf("failed to get user: %w", err)
	}
	return nil
}

func (r UserRepositoryImpl) UpdateImage(ctx context.Context, userID, imageID int64) error {
	sqlStr := "update shop.user set image = ?, update_time = ? where id = ?"

	_, err := r.db.ExecContext(ctx, sqlStr, imageID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Warnf("user update image err: %v", err)
			return fmt.Errorf("user update image err: %v", err)
		}
		applog.AppLogger.Errorf("user repo error: %v", err)
		return err
	}
	return nil
}

func (r UserRepositoryImpl) UpdatePasswd(ctx context.Context, user *model.User) error {
	sqlStr := "update shop.user set password = ?, update_time = ? where id = ?"

	_, err := r.db.ExecContext(ctx, sqlStr, user.Password, time.Now().String(), user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Warnf("user update err: %v", err)
			return fmt.Errorf("user update err: %v", err)
		}
		applog.AppLogger.Errorf("user repo error: %v", err)
		return fmt.Errorf("failed to get user: %w", err)
	}
	return nil
}

func (r UserRepositoryImpl) Delete(ctx context.Context, id int64) error {
	sqlStr := "DELETE FROM shop.user WHERE id = ? "

	_, err := r.db.ExecContext(ctx, sqlStr, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Warnf("user delete err: %v", err)
			return fmt.Errorf("user delete err: %v", err)
		}
		applog.AppLogger.Errorf("user repo error: %v", err)
		return fmt.Errorf("failed to get user: %w", err)
	}
	return nil
}
