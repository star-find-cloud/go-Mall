package model

type Shop struct {
	ID       int     `db:"id"`
	Name     string  `db:"name"`
	UserID   int     `db:"user_id"`
	Status   int     `db:"status"`
	Score    float64 `db:"score"` // 商家评分
	ImageId  int64   `db:"image_id"`
	Products []Product
	CreateAt int `db:"create_at"`
	Tag      int `db:"tag"` // 商家荣誉标签
}
