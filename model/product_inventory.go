package model

type ProductInventory struct {
	InventoryID      int `db:"inventory_id"`
	SellableQuantity int `db:"sellable_quantity"`
	OccupyQuantity   int `db:"occupy_quantity"`
}
