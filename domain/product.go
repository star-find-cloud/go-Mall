package domain

type Product struct {
	ID         int64 `db:"id"`
	MerchantID int64 `db:"merchant_id"` // 商家ID
	//ImageID       string  `db:"image_id"`
	Title         string  `db:"title"`
	SubTitle      string  `db:"sub_title"`
	Brand         string  `db:"brand"`           // 品牌
	ProductSn     string  `db:"product_sn"`      // 商品编号
	ProductTypeID int     `db:"product_type_id"` // 商品大类ID
	CateID        int64   `db:"cate_id"`         // 小类ID
	ClickCount    int     `db:"click_count"`     // 点击量
	PurchaseCount int64   `db:"purchase_count"`  // 购买量
	ProductNum    int     `db:"product_num"`     // 商品数量
	Price         float64 `db:"price"`           // 价格
	MarketPrice   float64 `db:"market_price"`    // 市场价
	Attr          string  `db:"attr"`            // 商品属性
	Version       string  `db:"version"`         // 商品版本
	Keywords      string  `db:"keywords"`        // 商品关键词
	Desc          string  `db:"desc"`            // 商品描述
	Content       string  `db:"content"`         // 商品内容
	Specs         string  `db:"specs"`           // 商品规格
	IsDeleted     int     `db:"is_deleted"`      // 是否删除
	CreatedAt     int     `db:"created_at"`      // 创建时间
	UpdatedAt     int     `db:"updated_at"`      // 更新时间
	DeletedAt     int     `db:"deleted_at"`      // 删除时间
	IsHot         int     `db:"is_hot"`          // 是否热门
	IsBest        bool    `db:"is_best"`         // 是否精品
	IsNew         bool    `db:"is_new"`          // 是否新品
	IsBooking     bool    `db:"is_booking"`      // 是否预售
	BookingTime   int64   `db:"booking_time"`    // 预售时间
	Sort          int     `db:"sort"`            // 排序权重
	Status        int     `db:"status"`
}

func (d *Product) ValidateMerchantID(inputID, storeID int64) bool {
	return inputID == storeID
}
