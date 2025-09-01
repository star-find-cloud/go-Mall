package repo

import (
	"context"
	"github.com/star-find-cloud/star-mall/domain"
)

type UserRepo interface {
	// Create 创建用户
	Create(ctx context.Context, user *domain.User) (int64, error)

	// GetByID 根据ID获取用户
	GetByID(ctx context.Context, id int64) (*domain.User, error)

	// GetByEmail 根据邮箱获取用户
	GetByEmail(ctx context.Context, email string) (*domain.User, error)

	// GetPasswordByID 根据ID获取用户密码
	GetPasswordByID(ctx context.Context, id int64) (string, error)

	// GetUserTags 获取用户标签
	GetUserTags(ctx context.Context, userID int64) ([]byte, error)

	// Update 更新用户
	Update(ctx context.Context, user *domain.User) error

	// UpdatePasswd 更新用户密码
	UpdatePasswd(ctx context.Context, user *domain.User) error

	// UpdateImage 更新用户头像
	UpdateImage(ctx context.Context, userID, imageID int64) error

	// UpdateEmail 更新用户邮箱
	UpdateEmail(ctx context.Context, email string, userID int64) error

	// UpdateUserTags 更新用户标签
	UpdateUserTags(ctx context.Context, userID int64, tag int64) error

	// IsExistsByID 根据ID判断用户是否存在
	IsExistsByID(ctx context.Context, id int64) (bool, error)

	// IsExistsByEmail 根据邮箱判断用户是否存在
	IsExistsByEmail(ctx context.Context, email string) (bool, error)

	// Delete 删除用户
	Delete(ctx context.Context, id int64) error

	// CheckEmailVerificationCode 检查邮箱验证码
	CheckEmailVerificationCode(ctx context.Context, email string, verificationCode string) (bool, error)
}
