package model

type Product struct {
	ID              int      `db:"id"`
	Title           string   `db:"title"`
	SubTitle        string   `db:"sub_title"`
	ProductSn       string   `db:"product_sn"`       // 商品编号
	CateID          int      `db:"cate_id"`          // 分类ID
	ClickCount      int      `db:"click_count"`      // 点击量
	ProductNum      int      `db:"product_num"`      // 商品数量
	Price           float64  `db:"price"`            // 价格
	MarketPrice     float64  `db:"market_price"`     //  市场价
	RelationProduct string   `db:"relation_product"` // 关联商品
	ProductAttr     string   `db:"product_attr"`     // 商品属性
	ProductVersion  string   `db:"product_version"`  // 商品版本
	ProductImages   []string `db:"product_images"`   // 商品图片
	ProductGift     string   `db:"product_gift"`     //  商品赠品
	ProductFitting  string   `db:"product_fitting"`  //  商品配件
	ProductColor    string   `db:"product_color"`    // 商品颜色
	ProductKeywords string   `db:"product_keywords"` // 商品关键词
	ProductDesc     string   `db:"product_desc"`     // 商品描述
	ProductContent  string   `db:"product_content"`  // 商品内容
	IsDeleted       int      `db:"is_deleted"`       // 是否删除
	CreatedAt       int      `db:"created_at"`       // 创建时间
	UpdatedAt       int      `db:"updated_at"`       // 更新时间
	DeletedAt       int      `db:"deleted_at"`       // 删除时间
	IsHot           int      `db:"is_hot"`           // 是否热门
	IsBest          int      `db:"is_best"`          // 是否精品
	IsNew           int      `db:"is_new"`           // 是否新品
	ProductTypeID   int      `db:"product_type_id"`  // 商品类型ID
	Sort            int      `db:"sort"`             // 排序
	status          int      `db:"status"`
}
