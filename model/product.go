package model

type Product struct {
	ID            int64   `db:"id"`
	ShopID        int64   `db:"shop_id"` // 商家ID
	Title         string  `db:"title"`
	SubTitle      string  `db:"sub_title"`
	ProductSn     string  `db:"product_sn"`      // 商品编号
	CateID        int     `db:"cate_id"`         // 分类ID
	ClickCount    int     `db:"click_count"`     // 点击量
	PurchaseCount int64   `db:"purchase_count"`  // 购买量
	ProductNum    int     `db:"product_num"`     // 商品数量
	Price         float64 `db:"price"`           // 价格
	MarketPrice   float64 `db:"market_price"`    // 市场价
	InventoryID   int64   `db:"inventory_id"`    // 库存ID
	Attr          string  `db:"attr"`            // 商品属性
	Version       string  `db:"version"`         // 商品版本
	ImagesID      []int64 `db:"images"`          // 商品图片id
	Keywords      string  `db:"keywords"`        // 商品关键词
	Desc          string  `db:"desc"`            // 商品描述
	Content       string  `db:"content"`         // 商品内容
	IsDeleted     int     `db:"is_deleted"`      // 是否删除
	CreatedAt     int     `db:"created_at"`      // 创建时间
	UpdatedAt     int     `db:"updated_at"`      // 更新时间
	DeletedAt     int     `db:"deleted_at"`      // 删除时间
	IsHot         int     `db:"is_hot"`          // 是否热门
	IsBest        int     `db:"is_best"`         // 是否精品
	IsNew         int     `db:"is_new"`          // 是否新品
	IsBooking     int     `db:"is_booking"`      // 是否预售
	ProductTypeID int     `db:"product_type_id"` // 商品类型ID
	Sort          int     `db:"sort"`            // 排序权重
	status        int     `db:"status"`
}
