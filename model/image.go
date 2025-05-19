package model

type Image struct {
	ImageID      int64  `db:"imageID"`
	OwnerType    int64  `db:"ownerType"`
	OwnerID      int64  `db:"ownerID"`
	OssPath      string `db:"ossPath"`
	SHA256Hash   string `db:"SHA256Hash"`
	IsCompressed bool   `db:"isCompressed"` // 是否压缩
}
