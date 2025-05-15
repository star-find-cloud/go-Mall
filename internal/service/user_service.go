package service

import (
	"errors"
	"github.com/star-find-cloud/star-mall/internal/repo"
	"github.com/star-find-cloud/star-mall/model"
	appjwt "github.com/star-find-cloud/star-mall/pkg/jwt"
	"github.com/star-find-cloud/star-mall/utils"
	"golang.org/x/net/context"
)

type UserService struct {
	repo repo.UserRepository
}

func NewUserService(repo repo.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s UserService) GetUsrByID(ctx context.Context, id int) (*model.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s UserService) Create(ctx context.Context, user *model.User) error {
	return s.repo.Create(ctx, user)
}

func (s UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *UserService) LoginByID(ctx context.Context, id int, password string) (string, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return "", errors.New("user not found")
	}

	err = utils.CheckPasswordHash(password, user.Password)
	if err != nil {
		return "", errors.New("invalid password")
	}

	token, err := appjwt.GenerateToken(user.ID, user.Name)
	if err != nil {
		return "", errors.New("failed to generate token")
	}
	return token, nil
}
