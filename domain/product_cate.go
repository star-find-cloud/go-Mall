package domain

// 商品分类
type ProductCate struct {
	ID              int           `db:"id"`
	Title           string        `db:"title"`
	CateImg         string        `db:"cate_img"`
	Link            string        `db:"link"`
	Template        string        `db:"template"`    // 模板
	Pid             int           `db:"pid"`         // 父级ID
	SubTitle        string        `db:"sub_title"`   // 子标题
	Keywords        string        `db:"keywords"`    // 关键词
	Description     string        `db:"description"` // 描述
	Sort            int           `db:"sort"`
	Status          int           `db:"status"`
	AddTime         int           `db:"add_time"`
	ProductCateItem []ProductCate `db:"product_cate_item"`
}
