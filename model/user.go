package model

// @Description 用户模型
type User struct {
	ID         int64  `db:"id"`
	Name       string `db:"name"`
	Password   string `db:"password"`
	Email      string `db:"email"`
	Phone      string `db:"phone"`
	Sex        int    `db:"sex"`
	CreateTime int64  `db:"create_time"`
	UpdateTime int64  `db:"update_time"`
	Status     int64  `db:"status"`
	LastIP     string `db:"last_ip"`
	ImageID    string `db:"image"` // uid
	IsVip      bool   `db:"is_vip"`
}
