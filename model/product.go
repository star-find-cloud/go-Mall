package model

type Product struct {
	ID              int
	Title           string
	SubTitle        string
	ProductSn       string   // 商品编号
	CateID          int      // 分类ID
	ClickCount      int      // 点击量
	ProductNum      int      // 商品数量
	Price           float64  // 价格
	MarketPrice     float64  //  市场价
	RelationProduct string   // 关联商品
	ProductAttr     string   // 商品属性
	ProductVersion  string   // 商品版本
	ProductImages   []string // 商品图片
	ProductGift     string   //  商品赠品
	ProductFitting  string   //  商品配件
	ProductColor    string   // 商品颜色
	ProductKeywords string   // 商品关键词
	ProductDesc     string   // 商品描述
	ProductContent  string   // 商品内容
	IsDeleted       int      // 是否删除
	CreatedAt       int      // 创建时间
	UpdatedAt       int      // 更新时间
	DeletedAt       int      // 删除时间
	IsHot           int      // 是否热门
	IsBest          int      // 是否精品
	IsNew           int      // 是否新品
	ProductTypeID   int      // 商品类型ID
	Sort            int      // 排序
	status          int
}
