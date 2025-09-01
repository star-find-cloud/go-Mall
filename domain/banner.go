package domain

// 轮播图
type Banner struct {
	ID          int
	Title       string
	BannerType  int
	BannerImage string
	Link        string
	Sort        int
	Status      int
	CreateTime  string
	UpdateTime  string
	UpdateUser  string
	DeleteTime  string
}
