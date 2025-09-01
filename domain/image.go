package domain

type Image struct {
	ImageID      int64  `db:"imageID"`
	OwnerID      int64  `db:"ownerID"`
	OwnerType    int64  `db:"ownerType"`
	Path         string `db:"path"`
	SHA256Hash   string `db:"sha256hash"`   // 哈希值, 用于重复图片识别, 前端处理
	IsCompressed bool   `db:"isCompressed"` // 是否压缩
	ContentType  int64  `db:"content_type"` // 图片格式
	CreateAt     int64  `db:"create_at"`    // 创建时间
	Status       int64  `db:"status"`       // 状态
}
