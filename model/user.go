package model

type User struct {
	ID         int    `db:"id"`
	Name       string `db:"name"`
	Password   string `db:"password"`
	Email      string `db:"email"`
	Phone      string `db:"phone"`
	CreateTime int    `db:"create_time"`
	UpdateTime int    `db:"update_time"`
	Status     int    `db:"status"`
	LastIP     string `db:"last_ip"`
	Image      string `db:"image"`
	IsVip      bool   `db:"is_vip"`
}
