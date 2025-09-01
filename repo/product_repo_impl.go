package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_const "github.com/star-find-cloud/star-mall/const"
	"github.com/star-find-cloud/star-mall/domain"
	"github.com/star-find-cloud/star-mall/pkg/database"
	log "github.com/star-find-cloud/star-mall/pkg/logger"
	"time"
)

type ProductRepoImpl struct {
	db    database.Database
	cache *database.Redis
}

func NewProductRepo(db database.Database, cache *database.Redis) *ProductRepoImpl {
	return &ProductRepoImpl{db: db, cache: cache}
}

func (r *ProductRepoImpl) Create(ctx context.Context, product *domain.Product) (int64, error) {
	//jsonImage, err := json.Marshal(product.ImageID)
	//if err != nil {
	//	log.AppLogger.Errorf("序列化图片id失败: %d", product.ImageID)
	//}

	sqlStr := "insert into shop.product (merchant_id, title, sub_title, brand, product_sn, cate_id, product_num, price, market_price, attr, version, keywords, `desc`, content, created_at, is_best, is_new, is_booking, product_type_id, sort, status) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

	result, err := r.db.GetDB().ExecContext(ctx, sqlStr,
		product.MerchantID,
		product.Title,
		//jsonImage,
		product.SubTitle,
		product.Brand,
		product.ProductSn,
		product.CateID,
		product.ProductNum,
		product.Price,
		product.MarketPrice,
		product.Attr,
		product.Version,
		product.Keywords,
		product.Desc,
		product.Content,
		product.CreatedAt,
		product.IsBest,
		product.IsNew,
		product.IsBooking,
		product.ProductTypeID,
		product.Sort,
		_const.ProductStatusOnSale,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.MySQLLogger.Warnf("product not found (shopID: %d)", product.MerchantID)
			return 0, fmt.Errorf("%w: product shopID %d", err, product.MerchantID)
		}
		log.AppLogger.Errorf("product repo error: %v", err)
		return 0, fmt.Errorf("failed to get product: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.AppLogger.Errorf("product repo error: %v", err)
		return 0, fmt.Errorf("failed to get product: %w", err)
	}

	return id, nil
}

//func (r ProductRepoImpl) CreateMore(ctx context.Context, products *[]model.Product) error {
//	return nil
//}

func (r *ProductRepoImpl) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	var product = &domain.Product{}
	sqlStr := "select id, merchant_id, title, sub_title, brand, product_sn, cate_id, click_count, purchase_count, product_num, price, market_price, attr, version, keywords, `desc`, content, is_deleted, created_at, updated_at, deleted_at, is_hot, is_best, is_new, is_booking, product_type_id, sort, status from shop.product where id = ?"

	err := r.db.GetDB().GetContext(ctx, product, sqlStr, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.MySQLLogger.Errorf("product not found (id: %d)", id)
			return nil, fmt.Errorf("%w: product id %d", err, id)
		}
		log.AppLogger.Errorf("product repo error: %v", err)
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	if product.IsDeleted == 1 {
		return nil, fmt.Errorf("product is deleted")
	}
	return product, nil
}

func (r *ProductRepoImpl) GetByMerchantID(ctx context.Context, MerchantID int64) ([]*domain.Product, error) {
	var products = []*domain.Product{}
	sqlStr := "select id, title, sub_title,  brand, product_sn, cate_id, click_count, purchase_count, product_num, price, market_price, attr, version, keywords, `desc`, content, is_deleted, created_at, updated_at, deleted_at, is_hot, is_best, is_new, is_booking, product_type_id, sort, status from shop.product where merchant_id = ?"

	err := r.db.GetDB().SelectContext(ctx, &products, sqlStr, MerchantID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.MySQLLogger.Warnf("product not found (MerchantID: %d)", MerchantID)
			return nil, fmt.Errorf("%w: product MerchantID %d", err, MerchantID)
		}
		log.AppLogger.Errorf("product repo error: %v", err)
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return products, nil
}

func (r *ProductRepoImpl) GetByTitleAndKeywords(ctx context.Context, title, keywords string, offset int) ([]*domain.Product, error) {
	var products = []*domain.Product{}
	sqlStr := "select id, merchant_id, sub_title,  brand, product_sn, cate_id, click_count, purchase_count, product_num, price, market_price, attr, version, `desc`, content, is_deleted, created_at, updated_at, deleted_at, is_hot, is_best, is_new, is_booking, product_type_id, sort, status from shop.product where title like ? and keywords like ? limit ? offset ?"

	err := r.db.GetDB().SelectContext(ctx, &products, sqlStr, "%"+title+"%", "%"+keywords+"%", 15, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.MySQLLogger.Warnf("product not found (title: %s)", title)
			return nil, fmt.Errorf("%w: product title %s", err, title)
		}
		log.AppLogger.Errorf("product repo error: %v", err)
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return products, nil
}

func (r *ProductRepoImpl) GetMerchantID(ctx context.Context, id int64) (int64, error) {
	var merchantID int64
	sqlStr := "select merchant_id from shop.product where id = ?"

	err := r.db.GetDB().GetContext(ctx, &merchantID, sqlStr, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.MySQLLogger.Warnf("product not found (id: %d)", id)
			return 0, fmt.Errorf("%w: product id %d", err, id)
		}
		log.AppLogger.Errorf("product repo error: %v", err)
		return 0, fmt.Errorf("failed to get product: %w", err)
	}
	return merchantID, nil
}

func (r *ProductRepoImpl) GetByCateIDsAndHot(ctx context.Context, ids []int64) ([]domain.Product, error) {
	var products = []domain.Product{}
	sqlStr := "select id, title, brand, `desc`,content from shop.product where (product_type_id = ? or cate_id = ?) and is_hot = 1 ORDER BY purchase_count DESC limit 30 "

	for _, id := range ids {
		err := r.db.GetDB().SelectContext(ctx, &products, sqlStr, id, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.MySQLLogger.Warnf("product not found (cateIDs: %v)", id)
				return nil, fmt.Errorf("%w: product cateIDs %v", err, id)
			}
			log.AppLogger.Errorf("product repo error: %v", err)
			return nil, fmt.Errorf("failed to get product: %w", err)
		}
	}
	return products, nil
}

// SearchByMsg 搜索商品
func (r *ProductRepoImpl) SearchByMsg(ctx context.Context, msg string) ([]domain.Product, error) {
	var products = []domain.Product{}
	sqlStr := "select id, merchant_id, title, sub_title,  brand, product_sn, cate_id, click_count, purchase_count, product_num, price, market_price, attr, version, `desc`, content, is_deleted, created_at, updated_at, deleted_at, is_hot, is_best, is_new, is_booking, product_type_id, sort, status from shop.product where title like ? or keywords like ?"

	err := r.db.GetDB().SelectContext(ctx, &products, sqlStr, "%"+msg+"%", "%"+msg+"%")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.MySQLLogger.Warnf("product not found (title: %s)", msg)
			return nil, fmt.Errorf("%w: product title %s", err, msg)
		}
		log.AppLogger.Errorf("product repo error: %v", err)
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return products, nil
}

func (r *ProductRepoImpl) Update(ctx context.Context, product *domain.Product) error {
	sqlStr := "update shop.product set title = ?, sub_title = ?, cate_id = ?, price = ?, market_price = ?, attr = ?, version = ?, `desc` = ?, content = ?, updated_at = ?, is_best = ?, is_booking = ?, product_type_id = ?, status = ? where id = ?"

	_, err := r.db.GetDB().ExecContext(ctx, sqlStr, product.Title, product.SubTitle, product.CateID, product.Price, product.MarketPrice, product.Attr, product.Version, product.Desc, product.Content, product.UpdatedAt, product.IsBest, product.IsBooking, product.ProductTypeID, product.Status, product.ID)
	if err != nil {
		log.AppLogger.Errorf("product repo error: %v", err)
		return fmt.Errorf("failed to update product: %w", err)
	}
	return nil
}

func (r *ProductRepoImpl) Delete(ctx context.Context, id int64) error {
	sqlstr := "update shop.product set is_deleted = 1, deleted_at = ?, status = ? where id = ?"

	now := int(time.Now().Unix())
	_, err := r.db.GetDB().ExecContext(ctx, sqlstr, now, _const.StatusDeleted, id)
	if err != nil {
		log.AppLogger.Errorf("product repo error: %v", err)
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}

func (r *ProductRepoImpl) GetCateID(ctx context.Context, id int64) (int64, error) {
	var cateID int64
	sqlStr := "select cate_id from shop.product where id = ?"

	err := r.db.GetDB().GetContext(ctx, &cateID, sqlStr, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.MySQLLogger.Warnf("product not found (id: %d)", id)
			return 0, fmt.Errorf("%w: product id %d", err, id)
		}
		log.AppLogger.Errorf("product repo error: %v", err)
		return 0, fmt.Errorf("failed to get product: %w", err)
	}
	return cateID, nil
}
