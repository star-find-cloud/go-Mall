package service

import (
	"errors"
	"github.com/star-find-cloud/star-mall/domain"
	"github.com/star-find-cloud/star-mall/pkg/jwt"
	log "github.com/star-find-cloud/star-mall/pkg/logger"
	"github.com/star-find-cloud/star-mall/repo"
	"github.com/star-find-cloud/star-mall/utils"
	"golang.org/x/net/context"
)

type UserService interface {
	// GetByID 根据id获取用户元数据
	GetByID(ctx context.Context, id int64) (*domain.User, error)

	// Create 创建用户
	Create(ctx context.Context, user *domain.User) (int64, error)

	// LoginByID 根据id和密码登录
	LoginByID(ctx context.Context, id int64, password string) (string, int64, error)

	// LoginByEmail 根据邮箱和密码登录
	LoginByEmail(ctx context.Context, email, password string) (string, int64, error)

	// Register 注册用户
	Register(ctx context.Context, user *domain.User) (string, int64, error)

	// Update 修改用户信息
	Update(ctx context.Context, name, phone, email string, id int64, sex int) error

	// UpdatePassword 修改密码
	UpdatePassword(ctx context.Context, id int64, newPassword, oldPassword string) error

	// UpdateImage 修改用户头像
	UpdateImage(ctx context.Context, userID int64, image *domain.Image) error

	// UpdateEmail 修改邮箱
	UpdateEmail(ctx context.Context, email string, verificationCode string, userID int64) error

	// Delete 删除用户
	Delete(ctx context.Context, id int64) error

	// CheckEmailVerificationCode 检查邮箱验证码
	CheckEmailVerificationCode(ctx context.Context, email string, verificationCode string) (bool, error)

	// ForgetPassword 忘记密码
	ForgetPassword(ctx context.Context, email string, verificationCode string, newPassword string) error
}

type UserServiceImpl struct {
	repo      repo.UserRepo
	imageRepo repo.ImageRepo
}

func NewUserService(repo repo.UserRepo, imageRepo repo.ImageRepo) *UserServiceImpl {
	return &UserServiceImpl{
		repo: repo,
		//ossClient: oosClient,
		imageRepo: imageRepo,
	}
}

func (s *UserServiceImpl) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	exist, err := s.repo.IsExistsByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("user not found")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *UserServiceImpl) Create(ctx context.Context, user *domain.User) (int64, error) {
	return s.repo.Create(ctx, user)
}

func (s *UserServiceImpl) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *UserServiceImpl) LoginByID(ctx context.Context, id int64, password string) (string, int64, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return "", 0, errors.New("user not found")
	}

	userPassword, err := s.repo.GetPasswordByID(ctx, id)
	err = utils.CheckPasswordHash(password, userPassword)
	if err != nil {
		return "", 0, errors.New("invalid password")
	}

	token, err := jwt.GenerateToken(user.ID, user.Name, user.RoleID)
	if err != nil {
		return "", 0, errors.New("failed to generate token")
	}
	return token, user.RoleID, nil
}

func (s *UserServiceImpl) LoginByEmail(ctx context.Context, email, password string) (string, int64, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", 0, errors.New("user not found")
	}

	userPassword, err := s.repo.GetPasswordByID(ctx, user.ID)
	err = utils.CheckPasswordHash(password, userPassword)
	if err != nil {
		return "", 0, errors.New("invalid password")
	}

	token, err := jwt.GenerateToken(user.ID, user.Name, user.RoleID)
	if err != nil {
		return "", 0, errors.New("failed to generate token")
	}
	return token, user.RoleID, nil
}

// Register 注册用户函数 返回 token和用户id
func (s *UserServiceImpl) Register(ctx context.Context, user *domain.User) (string, int64, error) {
	// 验证用户邮箱是否合法
	if !utils.VerifyEmail(user.Email) {
		return "", 0, errors.New("email is not valid")
	}

	// 检查用户邮箱是否存在
	existingUser, _ := s.repo.GetByEmail(ctx, user.Email)
	if existingUser != nil {
		return "", 0, errors.New("email already exists")
	}

	// 将密码哈希加密
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return "", 0, errors.New("failed to hash password")
	}

	user.Password = hashedPassword

	id, err := s.Create(ctx, user)
	if err != nil {
		return "", 0, errors.New("failed to create user")
	}

	token, err := jwt.GenerateToken(user.ID, user.Name, user.RoleID)
	if err != nil {
		return "", 0, errors.New("failed to generate token")
	}

	return token, id, nil
}

func (s *UserServiceImpl) Update(ctx context.Context, name, phone, email string, id int64, sex int) error {
	if id == 0 {
		return errors.New("invalid user")
	}

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}

	user.Name = name
	user.Phone = phone
	user.Email = email
	user.Sex = sex
	user.UpdateTime = utils.GetTimeNow()

	if err = s.repo.Update(ctx, user); err != nil {
		return errors.New("failed to update user")
	}
	return nil
}

func (s *UserServiceImpl) UpdatePassword(ctx context.Context, id int64, newPassword, oldPassword string) error {
	if id == 0 {
		return errors.New("invalid user")
	}

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}
	if user.Password != oldPassword {
		return errors.New("the passwords don't match")
	}

	user.Password, err = utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}
	user.UpdateTime = utils.GetTimeNow()

	if err = s.repo.UpdatePasswd(ctx, user); err != nil {
		return errors.New("failed to update user")
	}
	return nil
}

func (s *UserServiceImpl) UpdateImage(ctx context.Context, userID int64, image *domain.Image) error {
	exists, err := s.repo.IsExistsByID(ctx, userID)
	if err != nil {
		return errors.New("user inventoryRepo err")
	}
	if exists == false {
		return errors.New("user not found")
	}

	id, err := s.imageRepo.UploadImage(ctx, image)
	if err != nil {
		log.AppLogger.Errorf("failed to upload image: %v", err)
		return err
	}

	err = s.repo.UpdateImage(ctx, userID, id)
	if err != nil {
		log.AppLogger.Errorf("failed to update user: %v", err)
		return err
	}
	return nil
}

func (s *UserServiceImpl) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.New("invalid user")
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		log.AppLogger.Errorf("failed to delete user: %v", err)
		return errors.New("failed to delete user")
	}
	return nil
}

// CheckEmailVerificationCode 验证邮箱验证码
func (s *UserServiceImpl) CheckEmailVerificationCode(ctx context.Context, email string, verificationCode string) (bool, error) {
	if email == "" || verificationCode == "" {
		return false, errors.New("invalid email or verification code")
	}
	return s.repo.CheckEmailVerificationCode(ctx, email, verificationCode)
}

// ForgetPassword 忘记密码修改密码函数
func (s *UserServiceImpl) ForgetPassword(ctx context.Context, email string, verificationCode string, newPassword string) error {
	if email == "" || verificationCode == "" || newPassword == "" {
		return errors.New("invalid email or verification code or new password")
	}

	isInvalid, err := s.repo.CheckEmailVerificationCode(ctx, email, verificationCode)
	if err != nil {
		return errors.New("failed to verify email verification code")
	}
	if isInvalid == false {
		return errors.New("invalid email verification code")
	}

	user, err := s.repo.GetByEmail(ctx, email)
	user.Password = newPassword

	return s.repo.UpdatePasswd(ctx, user)
}

func (s *UserServiceImpl) UpdateEmail(ctx context.Context, email string, verificationCode string, userID int64) error {
	if email == "" || verificationCode == "" || userID == 0 {
		return errors.New("invalid email or verification code or user id")
	}

	isInvalid, err := s.repo.CheckEmailVerificationCode(ctx, email, verificationCode)
	if err != nil {
		return errors.New("failed to verify email verification code")
	}
	if isInvalid == false {
		return errors.New("invalid email verification code")
	}

	return s.repo.UpdateEmail(ctx, email, userID)
}
