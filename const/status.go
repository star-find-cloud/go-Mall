package _const

const (
	// 存在状态
	StatusNotDeleted = 60 + iota // 存在
	StatusDeleted                // 标记为删除

	// 商品状态
	ProductStatusPending = 70 + iota // 待上架
	ProductStatusOnSale              // 上架
	ProductStatusOffSale             // 下架
	ProductStatusDeleted             // 已删除

)
