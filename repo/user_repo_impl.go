package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	_const "github.com/star-find-cloud/star-mall/const"
	"github.com/star-find-cloud/star-mall/domain"
	"github.com/star-find-cloud/star-mall/pkg/database"
	log "github.com/star-find-cloud/star-mall/pkg/logger"
	"github.com/star-find-cloud/star-mall/utils"
)

type UserRepoImpl struct {
	db    database.Database
	cache *database.Redis
}

func NewUserRepo(db database.Database, cache *database.Redis) *UserRepoImpl {
	return &UserRepoImpl{db: db, cache: cache}
}

func (r *UserRepoImpl) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	var user = &domain.User{}
	sqlStr := "select id,name,image,sex,tags,last_ip,image,is_vip, role from shop.user where id = ? LIMIT 1;"

	err := r.db.GetDB().GetContext(ctx, user, sqlStr, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.MySQLLogger.Warnf("user not found (id: %d)", id)
			return nil, fmt.Errorf("%w: user id %d", err, id)
		}
		log.AppLogger.Errorf("user repo error: %v", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user.Status == _const.StatusDeleted {
		sqlStr = "update shop.user set status = ? where id = ?;"
		_, err = r.db.GetDB().ExecContext(ctx, sqlStr, _const.StatusNotDeleted, id)
		if err != nil {
			log.AppLogger.Errorf("user repo error: %v", err)
			return nil, fmt.Errorf("failed to get user: %w", err)
		}
	}

	return user, nil
}

func (r *UserRepoImpl) GetPasswordByID(ctx context.Context, id int64) (string, error) {
	var password string
	sqlStr := "select password from shop.user where id = ? LIMIT 1;"

	err := r.db.GetDB().GetContext(ctx, &password, sqlStr, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.MySQLLogger.Warnf("user not found (id: %d)", id)
			return "", fmt.Errorf("%w: user id %d", err, id)
		}
		log.AppLogger.Errorf("user repo error: %v", err)
		return "", fmt.Errorf("failed to get user: %w", err)
	}
	return password, nil
}

func (r *UserRepoImpl) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user = &domain.User{}
	sqlStr := "select id,name, sex,image,last_ip,image,is_vip, role from shop.user where email = ? LIMIT 1;"

	err := r.db.GetDB().GetContext(ctx, user, sqlStr, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.MySQLLogger.Warnf("user not found (email: %s)", email)
			return nil, fmt.Errorf("%w: user id %s", err, email)
		}
		log.AppLogger.Errorf("user repo error: %v", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r *UserRepoImpl) Create(ctx context.Context, user *domain.User) (int64, error) {
	sqlStr := "insert into shop.user (name,password,email,phone,sex,create_time, status, last_ip, image, is_vip, role) values (?,?,?,?,?,?,?,?,?,?,?)"

	result, err := r.db.GetDB().ExecContext(ctx, sqlStr,
		user.Name,
		user.Password,
		user.Email,
		user.Phone,
		user.Sex,
		utils.GetTimeNow(),
		_const.StatusNotDeleted,
		user.LastIP,
		user.ImageID,
		user.IsVip,
		user.RoleID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.MySQLLogger.Warnf("user creat err: %v", err)
			return 0, fmt.Errorf("user creat err: %v", err)
		}
		log.AppLogger.Errorf("user repo error: %v", err)
		return 0, fmt.Errorf("failed to get user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.AppLogger.Errorf("user repo error: %v", err)
		return 0, fmt.Errorf("failed to get user: %w", err)
	}

	return id, nil
}

func (r *UserRepoImpl) Update(ctx context.Context, user *domain.User) error {
	sqlStr := "update shop.user set name = ?,email = ?,phone = ?, sex = ?, update_time = ?, image = ? where id = ?;"

	_, err := r.db.GetDB().ExecContext(ctx, sqlStr, user.Name, user.Email, user.Phone, user.UpdateTime, user.ImageID, user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.MySQLLogger.Warnf("user update err: %v", err)
			return fmt.Errorf("user update err: %v", err)
		}
		log.AppLogger.Errorf("user repo error: %v", err)
		return fmt.Errorf("failed to get user: %w", err)
	}
	return nil
}

func (r *UserRepoImpl) UpdateImage(ctx context.Context, userID, imageID int64) error {
	sqlStr := "update shop.user set image = ?, update_time = ? where id = ?"

	now := utils.GetTimeNow()
	_, err := r.db.GetDB().ExecContext(ctx, sqlStr, imageID, now, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.MySQLLogger.Warnf("user update image err: %v", err)
			return fmt.Errorf("user update image err: %v", err)
		}
		log.AppLogger.Errorf("user repo error: %v", err)
		return err
	}
	return nil
}

func (r *UserRepoImpl) UpdatePasswd(ctx context.Context, user *domain.User) error {
	sqlStr := "update shop.user set password = ?, update_time = ? where id = ?"

	_, err := r.db.GetDB().ExecContext(ctx, sqlStr, user.Password, utils.GetTimeNow(), user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.MySQLLogger.Warnf("user update err: %v", err)
			return fmt.Errorf("user update err: %v", err)
		}
		log.AppLogger.Errorf("user repo error: %v", err)
		return fmt.Errorf("failed to get user: %w", err)
	}
	return nil
}

func (r *UserRepoImpl) Delete(ctx context.Context, id int64) error {
	sqlStr := "UPDATE shop.user SET status = ?, update_time = ? WHERE id = ?"

	now := utils.GetTimeNow()
	_, err := r.db.GetDB().ExecContext(ctx, sqlStr, _const.StatusDeleted, now, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.MySQLLogger.Warnf("user delete err: %v", err)
			return fmt.Errorf("user delete err: %v", err)
		}
		log.AppLogger.Errorf("user repo error: %v", err)
		return fmt.Errorf("failed to get user: %w", err)
	}
	return nil
}

func (r *UserRepoImpl) IsExistsByID(ctx context.Context, id int64) (bool, error) {
	var exists bool
	sqlStr := "select exists(select 1 from shop.user where id = ?)"

	err := r.db.GetDB().QueryRowxContext(ctx, sqlStr, id).Scan(&exists)
	if err != nil {
		return false, errors.New("failed to check if user exists")
	}
	return exists, nil
}

func (r *UserRepoImpl) IsExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	sqlStr := "select exists(select 1 from shop.user where email = ?)"

	err := r.db.GetDB().QueryRowxContext(ctx, sqlStr, email).Scan(&exists)
	if err != nil {
		return false, errors.New("failed to check if user exists")
	}
	return exists, nil
}

// CheckEmailVerificationCode 验证邮箱验证码
func (r *UserRepoImpl) CheckEmailVerificationCode(ctx context.Context, email string, verificationCode string) (bool, error) {
	// 从Redis缓存中获取存储的验证码
	storedCode, err := r.cache.GetCache().Get(ctx, email).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.AppLogger.Warnf("email:%s err: %v", email, err)
			return false, err
		}
		log.AppLogger.Errorf("获取验证码失败: %v", err)
		return false, fmt.Errorf("获取验证码失败: %w", err)
	}

	// 比较验证码是否匹配
	isValid, err := domain.ValidateCode(verificationCode, storedCode)
	if err != nil {
		log.AppLogger.Errorln(err)
		return false, err
	}

	// 如果验证成功，删除验证码防止重复使用
	if isValid {
		_, err = r.cache.GetCache().Del(ctx, email).Result()
		if err != nil {
			log.AppLogger.Warnf("删除已使用的验证码失败: %v", err)
		}
	}

	return isValid, nil
}

func (r *UserRepoImpl) UpdateEmail(ctx context.Context, email string, userID int64) error {
	sqlStr := "update shop.user set email = ?, update_time = ? where id = ?"

	now := utils.GetTimeNow()
	_, err := r.db.GetDB().ExecContext(ctx, sqlStr, email, now, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.MySQLLogger.Warnf("user update email err: %v", err)
			return fmt.Errorf("user update email err: %v", err)
		}
		log.AppLogger.Errorf("user repo error: %v", err)
		return fmt.Errorf("failed to get user: %w", err)
	}
	return nil
}

func (r *UserRepoImpl) UpdateUserTags(ctx context.Context, userID int64, tag int64) error {
	sqlStr := "update shop.user set tags = case when json_length(tags) < 3 then JSON_ARRAY_APPEND(tags,'$',?) else json_array_append(json_remove(tags, '$[0]'), '$', ?) end where id = ?"

	_, err := r.db.GetDB().ExecContext(ctx, sqlStr, tag, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.MySQLLogger.Warnf("user update tags err: %v", err)
			return fmt.Errorf("user update tags err: %v", err)
		}
		log.AppLogger.Errorf("user repo error: %v", err)
		return fmt.Errorf("failed to get user: %w", err)
	}
	return nil
}

func (r *UserRepoImpl) GetUserTags(ctx context.Context, userID int64) ([]byte, error) {
	var tags []byte
	sqlStr := "select tags from shop.user where id = ?"

	err := r.db.GetDB().QueryRowxContext(ctx, sqlStr, userID).Scan(&tags)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.MySQLLogger.Warnf("user get tags err: %v", err)
			return nil, fmt.Errorf("user get tags err: %v", err)
		}
		log.AppLogger.Errorf("user repo error: %v", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return tags, nil
}
