package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/star-find-cloud/star-mall/domain"
	"github.com/star-find-cloud/star-mall/pkg/database"
	"github.com/star-find-cloud/star-mall/pkg/idgen"
	applog "github.com/star-find-cloud/star-mall/pkg/logger"
	"time"
)

type CartRepositoryImpl struct {
	db    database.Database
	cache *database.Redis
}

func NewCartRepo(db database.Database, cache *database.Redis) *CartRepositoryImpl {
	return &CartRepositoryImpl{db: db, cache: cache}
}

func (r *CartRepositoryImpl) Create(ctx context.Context, cart *domain.Cart) (int64, error) {
	// 生成 uid
	var cartUid, err = idgen.GetUid(ctx, r.cache)
	if err != nil {
		return 0, errors.New("generate uid failed")
	}
	var cartSqlStr = "insert into shop.cart(id, user_id, created_at) values(?,?,?)"

	_, err = r.db.GetDB().ExecContext(ctx, cartSqlStr, cartUid, cart.UserID, time.Now().Unix())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Warnf("cart creat err: %v", err)
			return 0, fmt.Errorf("cart creat err: %v", err)
		}
		applog.AppLogger.Errorf("cart repo error: %v", err)
		return 0, fmt.Errorf("failed to get cart: %w", err)
	}

	return cartUid, nil
}

func (r *CartRepositoryImpl) GetByID(ctx context.Context, id int64) (*domain.Cart, error) {
	//TODO implement me
	panic("implement me")
}

func (r *CartRepositoryImpl) GetByUserID(ctx context.Context, userID int64) (*domain.Cart, error) {
	var cart = &domain.Cart{}
	sqlStr := "select id, user_id from shop.cart where user_id = ?"
	itemSqlStr := "select product_id, product_title, create_price, product_image_oss, specs, added_at from shop.cart_item where cart_id = ?"

	err := r.db.GetDB().GetContext(ctx, cart, sqlStr, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Warnf("cart not found (id: %d)", userID)
			return nil, fmt.Errorf("%w: cart id %d", err, userID)
		}
		applog.AppLogger.Errorf("cart repo error: %v", err)
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	err = r.db.GetDB().SelectContext(ctx, &cart.CartItems, itemSqlStr, cart.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Warnf("cart item not found (id: %d)", userID)
			return nil, fmt.Errorf("%w: cart id %d", err, userID)
		}
		applog.AppLogger.Errorf("cart repo error: %v", err)
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}
	return cart, nil
}

func (r *CartRepositoryImpl) AddProduct(ctx context.Context, cartItem *domain.CartItemVO) error {
	//// 定义并发所需的原子类型
	//var (
	//	successCount int32          // 成功操作计数器（原子类型）
	//	atomicErr    atomic.Value   // 原子存储首个错误
	//	wg           sync.WaitGroup // 等待组控制协程同步
	//)
	specs, err := json.Marshal(cartItem.Specs)
	if err != nil {
		applog.AppLogger.Errorf("cart repo json marshal error: %v", err)
		return fmt.Errorf("cart repo json marshal error: %w", err)
	}

	// 定义SQL语句
	sqlStr := "insert into shop.cart_item(cart_id, product_id, product_title, create_price, now_price, product_image_oss, specs, added_at) value (?,?,?,?,?,?,?,?)"

	_, err = r.db.GetDB().ExecContext(ctx, sqlStr, cartItem.CartID, cartItem.ProductID, cartItem.ProductTitle, cartItem.CreatePrice, cartItem.NowPrice, cartItem.ProductImageOss, specs, time.Now().Unix())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Warnf("cart add product err: %v", err)
			return errors.New("cart add product err")
		}
		applog.AppLogger.Errorf("cart add product err: %v", err)
		return fmt.Errorf("cart add product err: %w", err)
	}
	//// 在启动goroutine前初始化错误存储
	//atomicErr.Store(nil)
	//
	//for i := 0; i < len(CartItems); i++ {
	//	wg.Add(1)
	//
	//	go func(idx int) {
	//		// 获取当前购物车项
	//		item := CartItems[idx]
	//
	//		// 执行数据库操作
	//		_, err := r.db.ExecContext(
	//			ctx,
	//			cartItemSqlStr,
	//			item.CartID,
	//			item.ProductID,
	//			item.ProductTitle,
	//			item.CreatePrice,
	//			item.ProductImageOss,
	//			item.Specs,
	//			time.Now().Unix(),
	//		)
	//
	//		// 错误处理（原子操作保证并发安全）
	//		if err != nil {
	//			// 使用CAS操作确保只记录第一个错误
	//			if atomic.CompareAndSwapPointer(
	//				(*unsafe.Pointer)(unsafe.Pointer(&atomicErr)),
	//				nil,
	//				unsafe.Pointer(&err),
	//			) {
	//				// 错误分类记录
	//				if errors.Is(err, sql.ErrNoRows) {
	//					applog.MySQLLogger.Warnf("Cart item missing: Index=%d | PID=%d | Err=%v", idx, item.ProductID, err)
	//				} else {
	//					applog.AppLogger.Errorf("DB operation failed: Index=%d | PID=%d | Err=%v", idx, item.ProductID, err)
	//				}
	//			}
	//		} else {
	//			// 成功时原子增加计数器
	//			atomic.AddInt32(&successCount, 1)
	//		}
	//	}(i)
	//}
	//
	//if firstErr := atomicErr.Load(); firstErr != nil {
	//	return fmt.Errorf("batch operation failed: %w", *(*error)(firstErr.(unsafe.Pointer)))
	//}
	return nil
}
