package domain

// User 用户模型
// @Description 用户模型
type User struct {
	ID         int64  `db:"id" json:"ID,omitempty"`
	Name       string `db:"name" json:"name,omitempty"`
	Password   string `db:"password" json:"password,omitempty"`
	Email      string `db:"email" json:"email,omitempty"`
	Phone      string `db:"phone" json:"phone,omitempty"`
	Sex        int    `db:"sex" json:"sex,omitempty"`
	Tags       []byte `db:"tags" json:"tags,omitempty"`
	CreateTime int64  `db:"create_time" json:"createTime,omitempty"`
	UpdateTime int64  `db:"update_time" json:"updateTime,omitempty"`
	DeleteTime int64  `db:"delete_time" json:"deleteTime,omitempty"`
	Status     int64  `db:"status" json:"status,omitempty"`
	LastIP     string `db:"last_ip" json:"lastIP,omitempty"`
	ImageID    int64  `db:"image" json:"imageID,omitempty"` // uid
	IsVip      bool   `db:"is_vip" json:"isVip,omitempty"`
	RoleID     int64  `db:"role" json:"role,omitempty"` // 规则ID, 用于分别用户类型
}
