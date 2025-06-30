package domain

// 认证
type Auth struct {
	Id          int
	ModuleName  string // 模块名称
	ActionName  string // 操作名称
	Type        int    // 类型
	Url         string
	ModuleID    int // 模块ID
	Sort        int
	Description string // 描述
	Status      int
	CreateTime  int64
	UpdateTime  int64
	UpdateUser  string
}
