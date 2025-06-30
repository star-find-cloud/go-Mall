package domain

type Admin struct {
	ID       int
	UserName string
	Password string
	PhoneNum string
	Email    string
	Status   int
	RoleID   int
	AddTime  int
	IsSuper  int
	Role     Role
}
