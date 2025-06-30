package service

import (
	"context"
	"errors"
	"github.com/star-find-cloud/star-mall/domain"
	"github.com/star-find-cloud/star-mall/repo"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type ProductService struct {
	productRepo repo.ProductRepo
	ossClient   *cos.Client
	imageRepo   repo.ImageRepo
}

func NewProductService(repo repo.ProductRepo, ossClient *cos.Client, imageRepo repo.ImageRepo) *ProductService {
	return &ProductService{
		productRepo: repo,
		ossClient:   ossClient,
		imageRepo:   imageRepo,
	}
}

func (s *ProductService) Create(ctx context.Context, product *domain.Product, merchantID int64) (int64, error) {
	if merchantID == 0 {
		return 0, errors.New("merchantID is empty")
	}
	if product.MerchantID != merchantID {
		return 0, errors.New("商家ID与商品所属商家ID不匹配")
	}

	return s.productRepo.Create(ctx, product)
}

func (s *ProductService) GetProductByID(ctx context.Context, id int64) (*domain.Product, error) {
	if id == 0 {
		return nil, errors.New("id is empty")
	}
	return s.productRepo.GetByID(ctx, id)
}

func (s *ProductService) GetProductByMerchantID(ctx context.Context, merchantID int64) ([]*domain.Product, error) {
	if merchantID == 0 {
		return nil, errors.New("merchantID is empty")
	}
	return s.productRepo.GetByMerchantID(ctx, merchantID)
}

func (s *ProductService) GetByTitleAndKeywords(ctx context.Context, title, keywords string, offset int) ([]*domain.Product, error) {
	if title == "" && keywords == "" {
		return nil, errors.New("title and keywords is empty")
	}
	return s.productRepo.GetByTitleAndKeywords(ctx, title, keywords, offset)
}

func (s *ProductService) Update(ctx context.Context, product *domain.Product, merchantID int64) error {
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

func (s *ProductService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.New("id is empty")
	}
	return s.productRepo.Delete(ctx, id)
}

func (s *ProductService) Search(ctx context.Context, msg string) ([]domain.Product, error) {
	return s.productRepo.SearchByMsg(ctx, msg)
}
