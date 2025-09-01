package service

import (
	"context"
	"errors"
	"github.com/star-find-cloud/star-mall/domain"
	"github.com/star-find-cloud/star-mall/pkg/oss"
	"github.com/star-find-cloud/star-mall/repo"
)

type ProductService interface {
	// Create 创建商品
	Create(ctx context.Context, product *domain.Product, merchantID int64) (int64, error)

	// GetByID 根据ID获取商品
	GetByID(ctx context.Context, id int64) (*domain.Product, error)

	// GetByMerchantID 根据商家ID获取商品
	GetByMerchantID(ctx context.Context, merchantID int64) ([]*domain.Product, error)

	// GetByTitleAndKeywords 根据标题和关键词获取商品
	GetByTitleAndKeywords(ctx context.Context, title, keywords string, offset int) ([]*domain.Product, error)

	// Update 更新商品
	Update(ctx context.Context, product *domain.Product, merchantID int64) error

	// Delete 删除商品
	Delete(ctx context.Context, id int64) error

	// Search 根据搜索关键词获取商品
	Search(ctx context.Context, msg string) ([]domain.Product, error)
}

type ProductServiceImpl struct {
	productRepo repo.ProductRepo
	oss         oss.OSS
	imageRepo   repo.ImageRepo
}

func NewProductService(repo repo.ProductRepo, oss oss.OSS, imageRepo repo.ImageRepo) *ProductServiceImpl {
	return &ProductServiceImpl{
		productRepo: repo,
		oss:         oss,
		imageRepo:   imageRepo,
	}
}

func (s *ProductServiceImpl) Create(ctx context.Context, product *domain.Product, merchantID int64) (int64, error) {
	if merchantID == 0 {
		return 0, errors.New("merchantID is empty")
	}
	if product.MerchantID != merchantID {
		return 0, errors.New("商家ID与商品所属商家ID不匹配")
	}

	return s.productRepo.Create(ctx, product)
}

func (s *ProductServiceImpl) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	if id == 0 {
		return nil, errors.New("id is empty")
	}
	return s.productRepo.GetByID(ctx, id)
}

func (s *ProductServiceImpl) GetByMerchantID(ctx context.Context, merchantID int64) ([]*domain.Product, error) {
	if merchantID == 0 {
		return nil, errors.New("merchantID is empty")
	}
	return s.productRepo.GetByMerchantID(ctx, merchantID)
}

func (s *ProductServiceImpl) GetByTitleAndKeywords(ctx context.Context, title, keywords string, offset int) ([]*domain.Product, error) {
	if title == "" && keywords == "" {
		return nil, errors.New("title and keywords is empty")
	}
	return s.productRepo.GetByTitleAndKeywords(ctx, title, keywords, offset)
}

func (s *ProductServiceImpl) Update(ctx context.Context, product *domain.Product, merchantID int64) error {
	storeID, err := s.productRepo.GetMerchantID(ctx, product.ID)
	if err != nil {
		return errors.New("merchantID is empty")
	}
	if !product.ValidateMerchantID(merchantID, storeID) {
		return errors.New("商家ID与商品所属商家ID不匹配")
	}
	if product.ID == 0 {
		return errors.New("id is empty")
	}
	return s.productRepo.Update(ctx, product)
}

func (s *ProductServiceImpl) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.New("id is empty")
	}
	return s.productRepo.Delete(ctx, id)
}

func (s *ProductServiceImpl) Search(ctx context.Context, msg string) ([]domain.Product, error) {
	return s.productRepo.SearchByMsg(ctx, msg)
}
