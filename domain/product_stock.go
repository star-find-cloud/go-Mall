package domain

type ProductSock struct {
	ID                int64  `db:"id"`
	ProductID         int64  `db:"product_id"`
	SKU               string `json:"sku"`                // 商品SKU
	Name              string `json:"name"`               // 商品名称
	CategoryID        string `json:"category_id"`        // 分类ID
	CurrentQuantity   int    `json:"current_quantity"`   // 当前库存数量
	LockedQuantity    int    `json:"locked_quantity"`    // 锁定数量
	AvailableQuantity int    `json:"available_quantity"` // 可用数量(当前-锁定)
	SafetyStock       int    `json:"safety_stock"`       // 安全库存
	MaxStock          int    `json:"max_stock"`          // 最大库存
	MinStock          int    `json:"min_stock"`          // 最小库存
	WarehouseID       string `json:"warehouse_id"`       // 仓库ID
	LocationCode      string `json:"location_code"`      // 库位编码
	BatchNumber       string `json:"batch_number"`       // 批次号
	ProductionDate    int64  `json:"production_date"`    // 生产日期
	ExpiryDate        int64  `json:"expiry_date"`        // 过期日期
	Status            string `json:"status"`             // 状态(正常/临期/过期)
	CreatedAt         int64  `json:"created_at"`         // 创建时间
	UpdatedAt         int64  `json:"updated_at"`         // 更新时间
}
