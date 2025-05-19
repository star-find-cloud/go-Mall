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

func (s *UserService) GetUsrByID(ctx context.Context, id int) (*model.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) Create(ctx context.Context, user *model.User) error {
	return s.repo.Create(ctx, user)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
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

func (s *UserService) LoginByEmail(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
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

func (s *UserService) Register(ctx context.Context, name, password, email, phone string, sex int) (string, error) {
	// 验证用户邮箱是否合法
	if !utils.VerifyEmail(email) {
		return "", errors.New("email is not valid.")
	}

	// 检查用户邮箱是否存在
	existingUser, _ := s.repo.GetByEmail(ctx, email)
	if existingUser != nil {
		return "", errors.New("email already exists")
	}

	// 将密码哈希加密
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return "", errors.New("failed to hash password")
	}
	password = hashedPassword

	user := &model.User{
		Name:     name,
		Password: hashedPassword,
		Email:    email,
		Phone:    phone,
		Sex:      sex,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return "", errors.New("failed to create user")
	}

	token, err := appjwt.GenerateToken(user.ID, user.Name)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}
