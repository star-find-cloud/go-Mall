package model

type ProductCate struct {
	ID              int
	Title           string
	CateImg         string
	Link            string
	Template        string // 模板
	Pid             int    // 父级ID
	SubTitle        string // 子标题
	Keywords        string // 关键词
	Description     string // 描述
	Sort            int
	Status          int
	AddTime         int
	ProductCateItem []ProductCate
}
