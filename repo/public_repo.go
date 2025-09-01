package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/star-find-cloud/star-mall/pkg/database"
	applog "github.com/star-find-cloud/star-mall/pkg/logger"
	"time"
)

type PublicRepository interface {
	SetVerificationCodeCache(ctx context.Context, email, code string) error
}

type PublicRepo struct {
	db    database.Database
	cache *database.Redis
}

func NewPublicRepo(db database.Database, cache *database.Redis) *PublicRepo {
	return &PublicRepo{
		db:    db,
		cache: cache,
	}
}

// 发送验证码
func (r *PublicRepo) SetVerificationCodeCache(ctx context.Context, email, code string) error {
	// 先检查是否存在未过期的验证码
	ttl, err := r.cache.GetCache().TTL(ctx, email).Result()

	if err != nil && err.Error() != "redis: nil" {
		applog.AppLogger.Errorln("cache err")
		applog.RedisLogger.Errorf("检查验证码缓存失败: %v", err)
		return errors.New("检查验证码缓存失败")
	}

	// 如果验证码存在且未过期
	if ttl > 0 {
		remainingMinutes := int(ttl.Minutes())
		return errors.New(fmt.Sprintf("请等待 %d 分钟后再重新发送验证码", remainingMinutes))
	}

	// 设置新的验证码
	err = r.cache.GetCache().Set(ctx, email, code, 5*time.Minute).Err()
	if err != nil {
		applog.AppLogger.Errorln("cache err")
		applog.RedisLogger.Errorf("创建验证码缓存失败: %v", err)
		return errors.New("创建验证码缓存失败")
	}

	return nil
}
