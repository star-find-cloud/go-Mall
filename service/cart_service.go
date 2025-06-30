package service

import (
	"context"
	"github.com/star-find-cloud/star-mall/domain"
	"github.com/star-find-cloud/star-mall/repo"
)

type CartService struct {
	CartRepo    repo.CartRepo
	ProductRepo repo.ProductRepo
}

func NewCartService(cartRepo repo.CartRepo, productRepo repo.ProductRepo) *CartService {
	return &CartService{
		CartRepo:    cartRepo,
		ProductRepo: productRepo,
	}
}

func (s *CartService) Create(ctx context.Context, cart *domain.Cart) (int64, error) {
	return s.CartRepo.Create(ctx, cart)
}

func (s *CartService) GetByID(ctx context.Context, id int64) (*domain.Cart, error) {
	return s.CartRepo.GetByID(ctx, id)
}

func (s *CartService) GetByUserID(ctx context.Context, userID int64) (*domain.Cart, error) {
	return s.CartRepo.GetByUserID(ctx, userID)
}

func (s *CartService) AddProduct(ctx context.Context, cartItems *domain.CartItemVO) error {
	return s.CartRepo.AddProduct(ctx, cartItems)
}
