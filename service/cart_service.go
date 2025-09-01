package service

import (
	"context"
	"github.com/star-find-cloud/star-mall/domain"
	"github.com/star-find-cloud/star-mall/repo"
)

type CartService interface {
	// Create 创建购物车
	Create(ctx context.Context, cart *domain.Cart) (int64, error)

	// GetByID 通过id获取购物车获取购物车
	GetByID(ctx context.Context, id int64) (*domain.Cart, error)

	// GetByUserID 通过用户id获取购物车
	GetByUserID(ctx context.Context, userID int64) (*domain.Cart, error)

	// AddProduct 添加商品
	AddProduct(ctx context.Context, cartItems *domain.CartItemVO) error
}

type CartServiceImpl struct {
	CartRepo    repo.CartRepo
	ProductRepo repo.ProductRepo
}

func NewCartService(cartRepo repo.CartRepo, productRepo repo.ProductRepo) *CartServiceImpl {
	return &CartServiceImpl{
		CartRepo:    cartRepo,
		ProductRepo: productRepo,
	}
}

func (s *CartServiceImpl) Create(ctx context.Context, cart *domain.Cart) (int64, error) {
	return s.CartRepo.Create(ctx, cart)
}

func (s *CartServiceImpl) GetByID(ctx context.Context, id int64) (*domain.Cart, error) {
	return s.CartRepo.GetByID(ctx, id)
}

func (s *CartServiceImpl) GetByUserID(ctx context.Context, userID int64) (*domain.Cart, error) {
	return s.CartRepo.GetByUserID(ctx, userID)
}

func (s *CartServiceImpl) AddProduct(ctx context.Context, cartItems *domain.CartItemVO) error {
	return s.CartRepo.AddProduct(ctx, cartItems)
}
