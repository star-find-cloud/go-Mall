package model

type Cart struct {
	ID             int
	Title          string
	Price          float64
	ProductVersion string
	Number         int
	ProductGifts   []string // 赠品
	ProductFitting string   // 搭配
	ProductColor   string
	ProductImage   string
	ProductAttr    string // 属性
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
