package domain

// Merchant 商户
type Merchant struct {
	ID             int64   `db:"id" json:"id" example:"12345"`
	UserID         int64   `db:"user_id" json:"userId" example:"67890"`
	Name           string  `db:"name" json:"name" example:"Star Electronics"`
	Phone          string  `db:"phone" json:"phone" example:"+8613800138000"`
	Email          string  `db:"email" json:"email" example:"contact@starelectronics.com"`
	Password       string  `db:"password" json:"-"`
	RealName       string  `db:"real_name" json:"realName" example:"张三"`                // 法人真实姓名
	RealID         string  `db:"real_id" json:"realId" example:"310000199001011234"`    // 身份证号
	LicenseImageID int64   `db:"license_image_id" json:"licenseImageId" example:"1001"` // 营业执照图片ID
	OldName        string  `db:"old_name" json:"oldName,omitempty" example:"Star Tech"` // 曾用名
	Tag            int     `db:"tag" json:"tag" example:"1"`                            // 商户标签
	CateID         int64   `db:"cate_id" json:"categoryId" example:"100"`               // 店铺分类ID
	BusinessType   []int64 `db:"business_type" json:"businessTypes"`                    // 店铺经营类型
	Score          float64 `db:"score" json:"score" example:"4.8"`                      // 商户评分
	ImageID        int64   `db:"image_id" json:"imageId" example:"2001"`                // 店铺图片ID
	CreateAt       int64   `db:"create_at" json:"createAt" example:"1609459200"`
	UpdateAt       int64   `db:"update_at" json:"updateAt" example:"1609459200"`
	DeleteAt       int64   `db:"delete_at" json:"deleteAt,omitempty" example:"0"`
	Status         int     `db:"status" json:"status" example:"1"`
}
