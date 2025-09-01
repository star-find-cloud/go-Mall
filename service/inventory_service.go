package service

import (
	"context"
	"errors"
	"github.com/star-find-cloud/star-mall/domain"
	"github.com/star-find-cloud/star-mall/repo"
)

type InventoryService interface {
	// Create 创建库存
	Create(ctx context.Context, inventory *domain.Inventory, inputID int64) error

	// GetByID 根据ID查询库存
	GetByID(ctx context.Context, id int64) (*domain.Inventory, error)

	// GetByMerchantID 根据商家ID查询库存
	GetByMerchantID(ctx context.Context, MerchantID int64) ([]*domain.Inventory, error)

	// Update 更新库存
	Update(ctx context.Context, inventory *domain.Inventory, inputID int64) error

	// Delete 删除库存
	Delete(ctx context.Context, id int64) error
}

type InventoryServiceImpl struct {
	inventoryRepo repo.InventoryRepo
	merchantRepo  repo.MerchantRepository
	productRepo   repo.ProductRepo
}

func NewInventoryService(inventoryRepo repo.InventoryRepo, merchantRepo repo.MerchantRepository, productRepo repo.ProductRepo) *InventoryServiceImpl {
	return &InventoryServiceImpl{inventoryRepo: inventoryRepo, merchantRepo: merchantRepo, productRepo: productRepo}
}

func (s *InventoryServiceImpl) Create(ctx context.Context, inventory *domain.Inventory, inputID int64) error {
	var storeID, err = s.productRepo.GetMerchantID(ctx, inventory.ProductID)
	if err != nil {
		return errors.New("商品ID不存在")
	}
	if !inventory.ValidateMerchantID(inputID, storeID) {
		return errors.New("商家ID不匹配")
	}
	return s.inventoryRepo.Create(ctx, inventory)
}

func (s *InventoryServiceImpl) Update(ctx context.Context, inventory *domain.Inventory, inputID int64) error {
	var storeID, err = s.productRepo.GetMerchantID(ctx, inventory.ProductID)
	if err != nil {
		return errors.New("商品ID不存在")
	}
	if !inventory.ValidateMerchantID(inputID, storeID) {
		return errors.New("商家ID不匹配")
	}
	return s.inventoryRepo.Update(ctx, inventory)
}

func (s *InventoryServiceImpl) GetByID(ctx context.Context, id int64) (*domain.Inventory, error) {
	if id == 0 {
		return nil, errors.New("库存ID不存在")
	}
	return s.inventoryRepo.GetByID(ctx, id)
}

func (s *InventoryServiceImpl) GetByMerchantID(ctx context.Context, MerchantID int64) ([]*domain.Inventory, error) {
	if MerchantID == 0 {
		return nil, errors.New("商家ID不存在")
	}
	return s.inventoryRepo.GetByMerchantID(ctx, MerchantID)
}

func (s *InventoryServiceImpl) Delete(ctx context.Context, id int64) error {
	// TODO implement me
	panic("implement me")
}
