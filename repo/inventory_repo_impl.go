package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/star-find-cloud/star-mall/domain"
	"github.com/star-find-cloud/star-mall/pkg/database"
	applog "github.com/star-find-cloud/star-mall/pkg/logger"
	"time"
)

type InventoryRepoImpl struct {
	db    database.Database
	cache *database.Redis
}

func NewInventoryRepo(db database.Database, cache *database.Redis) *InventoryRepoImpl {
	return &InventoryRepoImpl{
		db:    db,
		cache: cache,
	}
}

func (r *InventoryRepoImpl) Create(ctx context.Context, inventory *domain.Inventory) error {
	sqlStr := "insert into shop.inventory (product_id, available_stock, low_stock_threshold, create_at) values (?, ?, ?, ?) "

	_, err := r.db.GetDB().ExecContext(ctx, sqlStr, inventory.ProductID, inventory.AvailableStock, time.Now().Unix(), inventory.UpdateAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Warnf("product not found id: %d", inventory.ProductID)
			return fmt.Errorf("product not found id: %d", inventory.ProductID)
		}
		applog.AppLogger.Errorf("inventory repo error: %v", err)
		return fmt.Errorf("failed to get inventory: %w", err)
	}

	return nil
}

func (r *InventoryRepoImpl) Update(ctx context.Context, inventory *domain.Inventory) error {
	sqlStr := "update shop.inventory set available_stock = ?, low_stock_threshold = ?, update_at = ? where product_id = ?"

	_, err := r.db.GetDB().ExecContext(ctx, sqlStr, inventory.AvailableStock, inventory.LowStockThreshold, time.Now().Unix(), inventory.ProductID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Warnf("product not found id: %d", inventory.ProductID)
			return fmt.Errorf("product not found id: %d", inventory.ProductID)
		}
		applog.AppLogger.Errorf("inventory repo error: %v", err)
		return fmt.Errorf("failed to get inventory: %w", err)
	}
	return nil
}

func (r *InventoryRepoImpl) GetByID(ctx context.Context, id int64) (*domain.Inventory, error) {
	sqlStr := "select available_stock, low_stock_threshold, create_at, update_at from shop.inventory where product_id = ?"

	var inventory = &domain.Inventory{}
	err := r.db.GetDB().GetContext(ctx, inventory, sqlStr, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Warnf("product not found id: %d", id)
			return nil, fmt.Errorf("product not found id: %d", id)
		}
		applog.AppLogger.Errorf("inventory repo error: %v", err)
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}
	return inventory, nil
}

func (r *InventoryRepoImpl) GetByMerchantID(ctx context.Context, MerchantID int64) ([]*domain.Inventory, error) {
	sqlStr := "select product_id, available_stock, low_stock_threshold, create_at, update_at from shop.inventory where product_id in (select product_id from shop.product where merchant_id = ?)"

	var inventories = make([]*domain.Inventory, 0)
	err := r.db.GetDB().SelectContext(ctx, &inventories, sqlStr, MerchantID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Warnf("product not found id: %d", MerchantID)
			return nil, fmt.Errorf("product not found id: %d", MerchantID)
		}
		applog.AppLogger.Errorf("inventory repo error: %v", err)
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}
	return inventories, nil
}

// Deduction 扣减库存
func (r *InventoryRepoImpl) Deduction(ctx context.Context, ProductID int64, count int64) error {
	sqlStr := "update shop.inventory set available_stock = available_stock - ? where product_id = ?"

	_, err := r.db.GetDB().ExecContext(ctx, sqlStr, count, ProductID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Warnf("product not found id: %d", ProductID)
			return fmt.Errorf("product not found id: %d", ProductID)
		}
		applog.AppLogger.Errorf("inventory repo error: %v", err)
		return fmt.Errorf("failed to get inventory: %w", err)
	}

	return nil
}
