package repo

import (
	"context"
	"github.com/star-find-cloud/star-mall/domain"
)

type UserRepo interface {
	Create(ctx context.Context, user *domain.User) (int64, error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetPasswordByID(ctx context.Context, id int64) (string, error)
	Update(ctx context.Context, user *domain.User) error
	UpdatePasswd(ctx context.Context, user *domain.User) error
	UpdateImage(ctx context.Context, userID, imageID int64) error
	IsExistsByID(ctx context.Context, id int64) (bool, error)
	IsExistsByEmail(ctx context.Context, email string) (bool, error)
	Delete(ctx context.Context, id int64) error
	CheckEmailVerificationCode(ctx context.Context, email string, verificationCode string) (bool, error)
	UpdateEmail(ctx context.Context, email string, userID int64) error
	UpdateUserTags(ctx context.Context, userID int64, tag int64) error
	GetUserTags(ctx context.Context, userID int64) ([]byte, error)
}
