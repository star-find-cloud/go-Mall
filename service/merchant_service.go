package service

import (
	"context"
	"errors"
	"fmt"
	_const "github.com/star-find-cloud/star-mall/const"
	"github.com/star-find-cloud/star-mall/domain"
	"github.com/star-find-cloud/star-mall/pkg/jwt"
	applog "github.com/star-find-cloud/star-mall/pkg/logger"
	"github.com/star-find-cloud/star-mall/repo"
	"github.com/star-find-cloud/star-mall/utils"
)

type MerchantService interface {
	// Create 创建商家
	Create(ctx context.Context, merchant *domain.Merchant) (int64, error)

	// Register 注册商家
	Register(ctx context.Context, merchant *domain.Merchant) (int64, string, error)

	// GetByID 根据ID获取商家
	GetByID(ctx context.Context, id int64) (*domain.Merchant, error)

	// GetByName 根据名称获取商家
	GetByName(ctx context.Context, email string) (*[]domain.Merchant, error)

	// Update 更新商家信息
	Update(ctx context.Context, merchant *domain.Merchant) error

	// UpdateLicenseImage 更新商家营业执照
	UpdateLicenseImage(ctx context.Context, merchantID int64, image *domain.Image) error

	// Delete 删除商家
	Delete(ctx context.Context, id int64) error
}

type MerchantServiceImpl struct {
	repo repo.MerchantRepository
	//ossClient *cos.Client
	imageRepo repo.ImageRepo
}

func NewMerchantService(repo repo.MerchantRepository, imageRepo repo.ImageRepo) *MerchantServiceImpl {
	return &MerchantServiceImpl{
		repo: repo,
		//ossClient: ossClient,
		imageRepo: imageRepo,
	}
}

func (s *MerchantServiceImpl) Create(ctx context.Context, merchant *domain.Merchant) (int64, error) {
	return s.repo.Create(ctx, merchant)
}

func (s *MerchantServiceImpl) GetByID(ctx context.Context, id int64) (*domain.Merchant, error) {
	return s.repo.GetMerchantByID(ctx, id)
}

func (s *MerchantServiceImpl) GetByName(ctx context.Context, email string) (*[]domain.Merchant, error) {
	return s.repo.GetMerchantByName(ctx, email)
}

func (s *MerchantServiceImpl) Register(ctx context.Context, merchant *domain.Merchant) (int64, string, error) {
	if !utils.VerifyEmail(merchant.Email) {
		return 0, "", errors.New("email is not valid")
	}

	// 判断邮箱是否已存在
	exist, err := s.repo.IsExistsByEmail(ctx, merchant.Email)
	if err != nil {
		return 0, "", err
	}
	if exist {
		return 0, "", errors.New("email already exist")
	}

	// 判断手机号是否已存在
	exist, err = s.repo.IsExistsByPhone(ctx, merchant.Name)
	if err != nil {
		return 0, "", err
	}
	if exist {
		return 0, "", errors.New("phone already  exist")
	}

	// 哈希密码
	hashedPassword, err := utils.HashPassword(merchant.Password)
	if err != nil {
		return 0, "", errors.New("failed to hash password")
	}
	merchant.Password = hashedPassword

	id, err := s.repo.Create(ctx, merchant)
	if err != nil {
		return 0, "", errors.New("failed to create merchant")
	}
	token, err := jwt.GenerateTempToken(id, merchant.Name, _const.MerchantRole)

	fmt.Println(id, token)
	return id, token, nil
}

func (s *MerchantServiceImpl) Update(ctx context.Context, merchant *domain.Merchant) error {
	if merchant.ID == 0 {
		return errors.New("invalid merchant")
	}
	exists, err := s.repo.IsExistsByID(ctx, merchant.ID)
	if err != nil {
		return errors.New("merchant not found")
	}
	if exists == false {
		return errors.New("merchant not exist")
	}
	return s.repo.Update(ctx, merchant)
}

func (s *MerchantServiceImpl) UpdateLicenseImage(ctx context.Context, merchantID int64, image *domain.Image) error {

	exists, err := s.repo.IsExistsByID(ctx, merchantID)
	if err != nil {
		return errors.New("merchant inventoryRepo err")
	}
	if exists == false {
		return errors.New("merchant not exist")
	}

	id, err := s.imageRepo.UploadImage(ctx, image)
	if err != nil {
		applog.AppLogger.Errorf("failed to upload image: %v", err)
		return err
	}

	err = s.repo.UpdateLicenseImage(ctx, merchantID, id)
	if err != nil {
		applog.AppLogger.Errorf("failed to update merchant: %v", err)
		return err
	}
	return nil
}

func (s *MerchantServiceImpl) Delete(ctx context.Context, id int64) error {
	exits, err := s.repo.IsExistsByID(ctx, id)
	if err != nil {
		return errors.New("merchant inventoryRepo err")
	}
	if !exits {
		return errors.New("merchant not exist")
	}
	return s.repo.Delete(ctx, id)
}
