package repo

import (
	"context"
	"github.com/star-find-cloud/star-mall/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	UpdatePasswd(ctx context.Context, user *model.User) error
	UpdateImage(ctx context.Context, userID, imageID int64) error
	Delete(ctx context.Context, id int) error
}
