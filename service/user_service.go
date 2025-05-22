package service

import (
	"errors"
	"github.com/star-find-cloud/star-mall/model"
	appjwt "github.com/star-find-cloud/star-mall/pkg/jwt"
	applog "github.com/star-find-cloud/star-mall/pkg/logger"
	proto "github.com/star-find-cloud/star-mall/protobuf/pb"
	repo2 "github.com/star-find-cloud/star-mall/repo"
	"github.com/star-find-cloud/star-mall/utils"
	"github.com/tencentyun/cos-go-sdk-v5"
	"golang.org/x/net/context"
)

type UserService struct {
	repo      repo2.UserRepository
	ossClient *cos.Client
	imageRepo repo2.ImageRepository
}

func NewUserService(repo repo2.UserRepository, oosClient *cos.Client, imageRepo repo2.ImageRepository) *UserService {
	return &UserService{
		repo:      repo,
		ossClient: oosClient,
		imageRepo: imageRepo,
	}
}

func (s *UserService) GetUsrByID(ctx context.Context, id int64) (*model.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) Create(ctx context.Context, user *model.User) (int64, error) {
	return s.repo.Create(ctx, user)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *UserService) LoginByID(ctx context.Context, id int64, password string) (string, error) {
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

func (s *UserService) Register(ctx context.Context, name, password, email, phone string, sex int) (string, int64, error) {
	// 验证用户邮箱是否合法
	if !utils.VerifyEmail(email) {
		return "", 0, errors.New("email is not valid.")
	}

	// 检查用户邮箱是否存在
	existingUser, _ := s.repo.GetByEmail(ctx, email)
	if existingUser != nil {
		return "", 0, errors.New("email already exists")
	}

	// 将密码哈希加密
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return "", 0, errors.New("failed to hash password")
	}
	password = hashedPassword

	user := &model.User{
		Name:     name,
		Password: hashedPassword,
		Email:    email,
		Phone:    phone,
		Sex:      sex,
	}

	id, err := s.Create(ctx, user)
	if err != nil {
		return "", 0, errors.New("failed to create user")
	}

	token, err := appjwt.GenerateToken(user.ID, user.Name)
	if err != nil {
		return "", 0, errors.New("failed to generate token")
	}

	return token, id, nil
}

//func (s *UserService) Update(ctx context.Context, name, email, phone string, id int64, sex int) error {
//	if id == 0 {
//		return errors.New("invalid user")
//	}
//
//	user, err := s.repo.GetByID(ctx, id)
//	if err != nil {
//		return errors.New("user not found")
//	}
//
//	if err := s.repo.Update(ctx, user); err != nil {
//		return errors.New("failed to update user")
//	}
//	return nil
//}

func (s UserService) UpdateImage(ctx context.Context, userID int64, newImage *proto.ImageProto) error {
	ossService := NewOSSService(s.ossClient, s.imageRepo)

	_, id, err := ossService.UploadImage(ctx, newImage)
	if err != nil {
		applog.AppLogger.Errorf("failed to upload image: %v", err)
		return err
	}

	user, err := s.repo.GetByID(ctx, userID)
	user.ImageID = id

	err = s.repo.UpdateImage(ctx, user.ID, user.ImageID)
	if err != nil {
		applog.AppLogger.Errorf("failed to update user: %v", err)
		return err
	}
	return nil
}
