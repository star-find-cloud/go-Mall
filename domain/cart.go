package domain

type Cart struct {
	ID        int64        `db:"id"`
	UserID    int64        `db:"user_id"`
	CartItems []CartItemVO `db:"cart_items"`
}

// CartItemVO 购物车商品详情
type CartItemVO struct {
	CartID          int64                  `db:"cartId"`
	ProductID       int64                  `db:"product_id"`
	ProductTitle    string                 `db:"product_title"`
	CreatePrice     float64                `db:"create_price"`
	NowPrice        float64                `db:"now_price"`
	ProductImageOss string                 `db:"product_image_oss"`
	Quantity        int64                  `db:"quantity"` // 商品数量
	Specs           map[string]interface{} `db:"specs"`    // 商品规格
	AddedAt         int64                  `db:"added_at"`
}

// CartHasData 判断购物车是否有相同商品
// 如果有相同的商品，则返回true
func (m *Cart) CartHasData(cartList []Cart, currentData Cart) bool {
	for i := 0; i < len(cartList); i++ {
		if cartList[i].ID == currentData.ID && m.equalCartItemVO(cartList[i].CartItems[0], currentData.CartItems[0]) {
			return true
		}
	}
	return false
}

func (m *Cart) equalCartItemVO(a, b CartItemVO) bool {
	if a.ProductID != b.ProductID {
		return false
	}
	if a.ProductTitle != b.ProductTitle {
		return false
	}
	if a.CreatePrice != b.CreatePrice {
		return false
	}
	if a.NowPrice != b.NowPrice {
		return false
	}
	if a.ProductImageOss != b.ProductImageOss {
		return false
	}
	if a.Quantity != b.Quantity {
		return false
	}

	// 比较 specs map[string]interface{}
	if len(a.Specs) != len(b.Specs) {
		return false
	}
	for k, v := range a.Specs {
		if bv, ok := b.Specs[k]; !ok || bv != v {
			return false
		}
	}
	return true
}
