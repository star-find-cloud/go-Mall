package service

import (
	"context"
	"errors"
	"github.com/star-find-cloud/star-mall/domain"
	"github.com/star-find-cloud/star-mall/repo"
)

type InventoryService struct {
	inventoryRepo repo.InventoryRepo
	merchantRepo  repo.MerchantRepository
	productRepo   repo.ProductRepo
}

func NewInventoryService(inventoryRepo repo.InventoryRepo, merchantRepo repo.MerchantRepository, productRepo repo.ProductRepo) *InventoryService {
	return &InventoryService{inventoryRepo: inventoryRepo, merchantRepo: merchantRepo, productRepo: productRepo}
}

// Create 创建库存
func (s *InventoryService) Create(ctx context.Context, inventory *domain.Inventory, inputID int64) error {
	var storeID, err = s.productRepo.GetMerchantID(ctx, inventory.ProductID)
	if err != nil {
		return errors.New("商品ID不存在")
	}
	if !inventory.ValidateMerchantID(inputID, storeID) {
		return errors.New("商家ID不匹配")
	}
	return s.inventoryRepo.Create(ctx, inventory)
}

// Update 更新库存
func (s *InventoryService) Update(ctx context.Context, inventory *domain.Inventory, inputID int64) error {
	var storeID, err = s.productRepo.GetMerchantID(ctx, inventory.ProductID)
	if err != nil {
		return errors.New("商品ID不存在")
	}
	if !inventory.ValidateMerchantID(inputID, storeID) {
		return errors.New("商家ID不匹配")
	}
	return s.inventoryRepo.Update(ctx, inventory)
}

func (s *InventoryService) GetInventory(ctx context.Context, id int64) (*domain.Inventory, error) {
	if id == 0 {
		return nil, errors.New("库存ID不存在")
	}
	return s.inventoryRepo.GetByID(ctx, id)
}

func (s *InventoryService) GetByMerchantID(ctx context.Context, MerchantID int64) ([]*domain.Inventory, error) {
	if MerchantID == 0 {
		return nil, errors.New("商家ID不存在")
	}
	return s.inventoryRepo.GetByMerchantID(ctx, MerchantID)
}
