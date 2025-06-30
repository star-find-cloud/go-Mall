package domain

// 送货地址
type Address struct {
	ID             int    `db:"id"`
	Uid            int    `db:"uid"`
	Phone          string `db:"phone"`
	Name           string `db:"name"`
	Address        string `db:"address"`
	DefaultAddress int    `db:"default_address"`
	AddTime        int    `db:"add_time"`
}
