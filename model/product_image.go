package model

type ProductImage struct {
	ID        int
	ProductID int
	ImageUrl  string
	ColorID   int
	Sort      int // 排序
	CreatedAt int
	UpdatedAt int
	DeletedAt int
	Status    int
}
