package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/star-find-cloud/star-mall/domain"
	applog "github.com/star-find-cloud/star-mall/pkg/logger"
	"github.com/star-find-cloud/star-mall/repo"
	"github.com/star-find-cloud/star-mall/utils"
	"time"
)

type OrderService interface {
	// Create 创建订单
	Create(ctx context.Context, order *domain.Order, orderItem *domain.OrderItem, userID int64) (int64, int64, error)

	// GetByID 获取订单
	GetByID(ctx context.Context, id int64) (*domain.Order, *domain.OrderItem, error)

	// Delete 删除订单
	Delete(ctx context.Context, id int64) error
}

type OrderServiceImpl struct {
	OrderRepo     repo.OrderRepo
	productRepo   repo.ProductRepo
	userRepo      repo.UserRepo
	inventoryRepo repo.InventoryRepo
}

func NewOrderService(orderRepo repo.OrderRepo, productRepo repo.ProductRepo, userRepo repo.UserRepo, inventoryRepo repo.InventoryRepo) *OrderServiceImpl {
	return &OrderServiceImpl{
		OrderRepo:     orderRepo,
		productRepo:   productRepo,
		userRepo:      userRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (s *OrderServiceImpl) Create(ctx context.Context, order *domain.Order, orderItem *domain.OrderItem, userID int64) (int64, int64, error) {
	if userID == 0 {
		return 0, 0, errors.New("userID is empty")
	}
	if order.UserID != userID {
		return 0, 0, errors.New("用户ID与订单所属用户ID不匹配")
	}

	tag, err := s.productRepo.GetCateID(ctx, orderItem.ProductID)
	err = s.userRepo.UpdateUserTags(ctx, userID, tag)
	id, err := utils.GenerateOrderID()
	if err != nil {
		applog.AppLogger.Errorf("生成订单号失败: %v", err)
		return 0, 0, fmt.Errorf("生成订单号失败: %v", err)
	}
	order.ID = id
	order.CreatedAt = time.Now().Unix()

	err = s.OrderRepo.Create(ctx, order, orderItem)
	if err != nil {
		applog.AppLogger.Errorf("创建订单失败: %v", err)
		return 0, 0, fmt.Errorf("创建订单失败: %v", err)
	}

	return order.ID, order.CreatedAt, s.inventoryRepo.Deduction(ctx, orderItem.ProductID, orderItem.Quantity)
}

func (s *OrderServiceImpl) GetByID(ctx context.Context, id int64) (*domain.Order, *domain.OrderItem, error) {
	if id == 0 {
		return nil, nil, errors.New("id is empty")
	}
	return s.OrderRepo.GetByID(ctx, id)
}

func (s *OrderServiceImpl) Delete(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}
