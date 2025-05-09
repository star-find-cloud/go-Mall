package model

type User struct {
	ID         int
	Name       string
	Password   string
	Email      string
	Phone      string
	CreateTime int
	UpdateTime int
	Status     int
	LastIP     string
	Image      string
	IsVip      bool
}
