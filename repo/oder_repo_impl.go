package repo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/star-find-cloud/star-mall/domain"
	"github.com/star-find-cloud/star-mall/pkg/database"
	applog "github.com/star-find-cloud/star-mall/pkg/logger"
)

type OrderRepoImpl struct {
	db    database.Database
	cache *database.Redis
}

func NewOrderRepo(db database.Database, cache *database.Redis) *OrderRepoImpl {
	return &OrderRepoImpl{
		db:    db,
		cache: cache,
	}
}

func (r *OrderRepoImpl) Create(ctx context.Context, order *domain.Order, orderItem *domain.OrderItem) error {
	sqlStr := "insert into shop.orders (id, user_id, order_status, total_price, pay_price, created_at, shipping_id) values (?, ?, ?, ?, ?, ?, ?)"
	itemSqlStr := "insert into shop.order_items (order_id, product_id, product_title, unit_price, quantity, subtotal) values (?, ?, ?, ?, ?, ?)"

	_, err := r.db.GetDB().ExecContext(ctx, sqlStr, order.ID, order.UserID, order.OrderStatus, order.TotalPrice, order.PayPrice, order.CreatedAt, order.ShippingAddressID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Warnf("create order failed, err: %v", err)
			return errors.New("create order failed")
		}
		applog.AppLogger.Errorf("create order failed, err: %v", err)
		return errors.New("create order failed")
	}

	_, err = r.db.GetDB().ExecContext(ctx, itemSqlStr, orderItem.OrderID, orderItem.ProductID, orderItem.ProductTitle, orderItem.UnitPrice, orderItem.Quantity, orderItem.Subtotal)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Warnf("create order item failed, err: %v", err)
			return errors.New("create order item failed")
		}
		applog.AppLogger.Errorf("create order item failed, err: %v", err)
		return errors.New("create order item failed")
	}
	return nil
}

func (r *OrderRepoImpl) GetByID(ctx context.Context, orderID int64) (*domain.Order, *domain.OrderItem, error) {
	var order = &domain.Order{}
	var orderItem = &domain.OrderItem{}

	sqlStr := "select id, user_id, order_status, total_price, pay_price, created_at, shipping_id from shop.orders where id = ?"
	ItemSqlStr := "select order_id, product_id, product_title, unit_price, quantity, subtotal from shop.order_items where order_id = ?"

	err := r.db.GetDB().GetContext(ctx, order, sqlStr, orderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Warnf("get order by id failed, err: %v", err)
			return nil, nil, errors.New("get order by id failed")
		}
		applog.AppLogger.Errorf("get order by id failed, err: %v", err)
		return nil, nil, errors.New("get order by id failed")
	}

	err = r.db.GetDB().GetContext(ctx, orderItem, ItemSqlStr, orderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Warnf("get order item by id failed, err: %v", err)
			return nil, nil, errors.New("get order item by id failed")
		}
		applog.AppLogger.Errorf("get order item by id failed, err: %v", err)
		return nil, nil, errors.New("get order item by id failed")
	}

	return order, orderItem, nil
}

func (r *OrderRepoImpl) GetByCreatedAt(ctx context.Context, createdAt int64) ([]*domain.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (r *OrderRepoImpl) UpdateStatus(ctx context.Context, order *domain.Order) error {
	//TODO implement me
	panic("implement me")
}

func (r *OrderRepoImpl) Delete(ctx context.Context, order *domain.Order) error {
	sqlStr := "delete from shop.orders where id = ?"
	itemSqlStr := "delete from shop.order_items where order_id = ?"

	_, err := r.db.GetDB().ExecContext(ctx, sqlStr, order.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Warnf("delete order failed, err: %v", err)
			return errors.New("delete order failed")
		}
		applog.AppLogger.Errorf("delete order failed, err: %v", err)
		return errors.New("delete order failed")
	}

	_, err = r.db.GetDB().ExecContext(ctx, itemSqlStr, order.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			applog.MySQLLogger.Warnf("delete order item failed, err: %v", err)
			return errors.New("delete order item failed")
		}
		applog.AppLogger.Errorf("delete order item failed, err: %v", err)
		return errors.New("delete order item failed")
	}

	return nil
}
