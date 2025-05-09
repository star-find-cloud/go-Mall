package model

type Menu struct {
	ID          int
	Title       string
	Link        string
	Position    int
	IsOpenNew   int
	Relation    string
	Sort        int
	Status      int
	AddTime     int
	ProductItem []Product
}
