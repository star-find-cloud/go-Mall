package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/star-find-cloud/star-mall/model"
	applog "github.com/star-find-cloud/star-mall/pkg/logger"
)

type ProductImpl struct {
	db *sqlx.DB
}

func NewProductRepositoryImpl(db *sqlx.DB) *ProductImpl {
	return &ProductImpl{db: db}
}

func (r ProductImpl) GetByID(ctx context.Context, id int) (*model.Product, error) {
	var product *model.Product
	sqlStr := "select id,title, sub_title, product_sn, cate_id, click_count, product_num, price, market_price, relation_product, product_attr, product_version, product_images, product_gift, product_fitting, product_color, product_keywords, product_desc, product_content, is_deleted, created_at, updated_at, deleted_at, is_hot, is_best, is_new, product_type_id, sort, status from shop.product where id = ? limit 1;"

	err := r.db.GetContext(ctx, product, sqlStr, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.AppLogger.Warnf("product get err: %v", err)
			return nil, fmt.Errorf("product get err: %v", err)
		}
		applog.AppLogger.Errorf("product repo error: %v", err)
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return product, nil
}

func (r ProductImpl) Create(ctx context.Context, product *model.Product) error {
	sqlStr := "insert into shop.product (title, sub_title, product_sn, cate_id, product_num, price, market_price, relation_product, product_attr, product_version, product_images, product_gift, product_fitting, product_color, product_keywords, product_desc, product_content, product_type_id, is_deleted, is_hot, is_best, is_new, created_at) values ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"

	_, err := r.db.ExecContext(ctx, sqlStr,
		product.Title,
		product.SubTitle,
		product.ProductSn,
		product.CateID,
		product.ProductNum,
		product.Price,
		product.MarketPrice,
		product.RelationProduct,
		product.ProductAttr,
		product.ProductVersion,
		product.ProductImages,
		product.ProductGift,
		product.ProductFitting,
		product.ProductColor,
		product.ProductKeywords,
		product.ProductDesc,
		product.ProductContent,
		product.ProductTypeID,
		product.IsDeleted,
		product.IsHot,
		product.IsBest,
		product.IsNew,
		product.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.AppLogger.Warnf("product create err: %v", err)
			return fmt.Errorf("product create err: %v", err)
		}
		applog.AppLogger.Errorf("product repo error: %v", err)
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

func (r ProductImpl) Update(ctx context.Context, product *model.Product) error {
	sqlStr := "update shop.product set title = ?, sub_title = ?, product_sn = ?, cate_id = ?, product_num = ?, price = ?, market_price = ?, relation_product = ?, product_attr = ?, product_version = ?, product_images = ?, product_gift = ?, product_fitting = ?, product_color = ?, product_keywords = ?, product_desc = ?, product_content = ?, product_type_id = ?, is_deleted = ?, is_hot = ?, is_best = ?, product_type_id = ?, status =?;"

	_, err := r.db.ExecContext(ctx, sqlStr,
		product.Title,
		product.SubTitle,
		product.ProductSn,
		product.CateID,
		product.ProductNum,
		product.Price,
		product.MarketPrice,
		product.RelationProduct,
		product.ProductAttr,
		product.ProductVersion,
		product.ProductImages,
		product.ProductGift,
		product.ProductFitting,
		product.ProductColor,
		product.ProductKeywords,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.AppLogger.Warnf("product update err: %v", err)
			return fmt.Errorf("product update err: %v", err)
		}
		applog.AppLogger.Errorf("product repo error: %v", err)
		return fmt.Errorf("failed to update product: %w", err)
	}
	return nil
}

func (r ProductImpl) Delete(ctx context.Context, id int) error {
	sqlStr := "delete from shop.product where id = ?;"

	_, err := r.db.ExecContext(ctx, sqlStr, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.AppLogger.Warnf("product delete err: %v", err)
			return fmt.Errorf("product delete err: %v", err)
		}
		applog.AppLogger.Errorf("product repo error: %v", err)
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}
