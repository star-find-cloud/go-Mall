package domain

// Inventory 商品库存
type Inventory struct {
	ProductID         int64 `db:"product_id"`
	AvailableStock    int64 `db:"available_stock"`     // 可用库存
	ReservedStock     int64 `db:"reserved_stock"`      // 已锁定库存
	LowStockThreshold int64 `db:"low_stock_threshold"` // 低库存阈值
	Version           int64 `db:"version"`             // 乐观锁版本
	CreateAt          int64 `db:"create_at"`
	UpdateAt          int64 `db:"update_at"`
}

// ValidateMerchantID 验证商户ID
func (d *Inventory) ValidateMerchantID(inputID, storeID int64) bool {
	return inputID == storeID
}
