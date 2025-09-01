package domain

// 商品属性
type ProductAttr struct {
	ID              int
	ProductID       int
	AttributeCateID int    // 属性分类ID
	AttributeID     int    // 属性ID
	AttributeTitle  string // 属性标题
	AttributeValue  string // 属性值
	Sort            int
	AddTime         int
	DeleteTime      int
	UpdateTime      int
	Status          int
}
