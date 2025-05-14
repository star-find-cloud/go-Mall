package model

type Cart struct {
	ID             int64    `db:"id"`
	Title          string   `db:"title"`
	Price          float64  `db:"price"`
	ProductVersion string   `db:"product_version"`
	Number         int64    `db:"number"`
	ProductGifts   []string `db:"product_gifts"`   // 赠品
	ProductFitting string   `db:"product_fitting"` // 搭配
	ProductColor   string   `db:"product_color"`
	ProductImage   string   `db:"product_image"`
	ProductAttr    string   `db:"product_attr"` // 属性
}

//// 验证购物车数据是否添加
//// 如果有相同的商品，则返回true
//func CartHasData(cartList []Cart, currentData Cart) bool {
//	for i := 0; i < len(cartList); i++ {
//		if cartList[i].ID == currentData.ID && cartList[i].ProductColor == currentData.ProductColor && cartList[i].ProductAttr == currentData.ProductAttr {
//			return true
//		}
//	}
//	return false
//}
