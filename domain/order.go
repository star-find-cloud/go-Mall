package domain

type Order struct {
	ID                int64  `db:"id" json:"id,omitempty"`
	UserID            int64  `db:"user_id" json:"userID,omitempty"`
	OrderStatus       string `db:"order_status" json:"orderStatus,omitempty"`
	TotalPrice        int64  `db:"total_price" json:"totalPrice,omitempty"` // 总价
	PayPrice          int64  `db:"pay_price" json:"payPrice,omitempty"`     // 实际支付价格
	CreatedAt         int64  `db:"created_at" json:"createdAt,omitempty"`
	UpdatedAt         int64  `db:"updated_at" json:"updatedAt,omitempty"`
	PaymentMethodID   int64  `db:"payment_id" json:"paymentID,omitempty"`   // 支付方式ID
	ShippingAddressID int64  `db:"shipping_id" json:"shippingID,omitempty"` // 配送地址ID
}

type OrderItem struct {
	ItemID       int64  `db:"item_id" json:"itemID,omitempty"`
	OrderID      int64  `db:"order_id" json:"orderID,omitempty"`
	ProductID    int64  `db:"product_id" json:"productID,omitempty"`
	ProductTitle string `db:"product_title" json:"productTitle,omitempty"`
	UnitPrice    int64  `db:"unit_price" json:"unitPrice,omitempty"` // 单价
	Quantity     int64  `db:"quantity" json:"quantity,omitempty"`    // 数量
	Subtotal     int64  `db:"subtotal" json:"subtotal,omitempty"`    // 小计
}
