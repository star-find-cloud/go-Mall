package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/star-find-cloud/star-mall/model"
	applog "github.com/star-find-cloud/star-mall/pkg/logger"
)

type UserRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepositoryImpl(db *sqlx.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (r UserRepositoryImpl) GetByID(ctx context.Context, id int) (*model.User, error) {
	var user *model.User
	sqlStr := "select id,name,image,last_ip,image,is_vip from user.user where id = ? LIMIT 1;"

	err := r.db.GetContext(ctx, user, sqlStr, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.AppLogger.Warnf("user not found (id: %d)", id)
			return nil, fmt.Errorf("%w: user id %d", err, id)
		}
		applog.AppLogger.Errorf("user repo error: %v", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r UserRepositoryImpl) Create(ctx context.Context, user *model.User) error {
	sqlStr := "insert into user.user (name,password,email,phone,create_time, update_time, status, last_ip, image, is_vip) values (?,?,?,?,?,?,?,?,?,?)"

	_, err := r.db.ExecContext(ctx, sqlStr,
		user.Name,
		user.Password,
		user.Email,
		user.Phone,
		user.CreateTime,
		user.UpdateTime,
		user.Status,
		user.LastIP,
		user.Image,
		user.IsVip)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.AppLogger.Warnf("user creat err: %v", err)
			return fmt.Errorf("user creat err: %v", err)
		}
		applog.AppLogger.Errorf("user repo error: %v", err)
		return fmt.Errorf("failed to get user: %w", err)
	}

	//id, _ := result.LastInsertId()
	//user.ID = int(id)
	return nil
}

func (r UserRepositoryImpl) Update(ctx context.Context, user *model.User) error {
	sqlStr := "update user.user set name = ?,email = ?,phone = ?, update_time = ?, image = ? where id = ?;"

	_, err := r.db.ExecContext(ctx, sqlStr, user.Name, user.Email, user.Phone, user.UpdateTime, user.Image, user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.AppLogger.Warnf("user update err: %v", err)
			return fmt.Errorf("user update err: %v", err)
		}
		applog.AppLogger.Errorf("user repo error: %v", err)
		return fmt.Errorf("failed to get user: %w", err)
	}
	return nil
}

func (r UserRepositoryImpl) UpdatePasswd(ctx context.Context, user *model.User) error {
	sqlStr := "update user.user set password = ?, update_time = ? where id = ?"

	_, err := r.db.ExecContext(ctx, sqlStr, user.Password, user.UpdateTime, user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.AppLogger.Warnf("user update err: %v", err)
			return fmt.Errorf("user update err: %v", err)
		}
		applog.AppLogger.Errorf("user repo error: %v", err)
		return fmt.Errorf("failed to get user: %w", err)
	}
	return nil
}

func (r UserRepositoryImpl) Delete(ctx context.Context, id int) error {
	sqlStr := "DELETE FROM user.user WHERE id = ? "

	_, err := r.db.ExecContext(ctx, sqlStr, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.AppLogger.Warnf("user delete err: %v", err)
			return fmt.Errorf("user delete err: %v", err)
		}
		applog.AppLogger.Errorf("user repo error: %v", err)
		return fmt.Errorf("failed to get user: %w", err)
	}
	return nil
}
