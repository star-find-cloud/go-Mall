package domain

import (
	"errors"
	"time"
)

// Role 角色领域模型
type Role struct {
	ID          int64  `db:"id" json:"id"`                   // 角色ID
	Name        string `db:"name" json:"name"`               // 角色名称
	Code        string `db:"code" json:"code"`               // 角色编码
	Description string `db:"description" json:"description"` // 角色描述
	Status      int    `db:"status" json:"status"`           // 状态：1-启用，0-禁用
	IsSystem    bool   `db:"is_system" json:"isSystem"`      // 是否系统内置角色
	Sort        int    `db:"sort" json:"sort"`               // 排序
	CreatedAt   int64  `db:"created_at" json:"createdAt"`    // 创建时间
	UpdatedAt   int64  `db:"updated_at" json:"updatedAt"`    // 更新时间
	CreatedBy   int64  `db:"created_by" json:"createdBy"`    // 创建人ID
	UpdatedBy   int64  `db:"updated_by" json:"updatedBy"`    // 更新人ID
}

// NewRole 创建新角色
func NewRole(name, code, description string, isSystem bool, operatorID int64) (*Role, error) {
	if err := validateRoleData(name, code); err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	return &Role{
		Name:        name,
		Code:        code,
		Description: description,
		Status:      1, // 默认启用
		IsSystem:    isSystem,
		Sort:        0,
		CreatedAt:   now,
		UpdatedAt:   now,
		CreatedBy:   operatorID,
		UpdatedBy:   operatorID,
	}, nil
}

// validateRoleData 验证角色数据
func validateRoleData(name, code string) error {
	if name == "" {
		return errors.New("角色名称不能为空")
	}
	if code == "" {
		return errors.New("角色编码不能为空")
	}
	return nil
}

// TableName 返回表名
func (r *Role) TableName() string {
	return "sys_role"
}

// Validate 验证角色数据
func (r *Role) Validate() error {
	return validateRoleData(r.Name, r.Code)
}

// IsEnabled 判断角色是否启用
func (r *Role) IsEnabled() bool {
	return r.Status == 1
}

// Enable 启用角色
func (r *Role) Enable(operatorID int64) {
	r.Status = 1
	r.UpdatedAt = time.Now().Unix()
	r.UpdatedBy = operatorID
}

// Disable 禁用角色
func (r *Role) Disable(operatorID int64) {
	r.Status = 0
	r.UpdatedAt = time.Now().Unix()
	r.UpdatedBy = operatorID
}

// Update 更新角色信息
func (r *Role) Update(name, description string, operatorID int64) error {
	if name == "" {
		return errors.New("角色名称不能为空")
	}
	r.Name = name
	r.Description = description
	r.UpdatedAt = time.Now().Unix()
	r.UpdatedBy = operatorID
	return nil
}

// UpdateSort 更新排序
func (r *Role) UpdateSort(sort int, operatorID int64) {
	r.Sort = sort
	r.UpdatedAt = time.Now().Unix()
	r.UpdatedBy = operatorID
}

// RolePermission 角色-权限关联
type RolePermission struct {
	ID           int64 `db:"id" json:"id"`                      // 主键ID
	RoleID       int64 `db:"role_id" json:"roleId"`             // 角色ID
	PermissionID int64 `db:"permission_id" json:"permissionId"` // 权限ID
	CreatedAt    int64 `db:"created_at" json:"createdAt"`       // 创建时间
	CreatedBy    int64 `db:"created_by" json:"createdBy"`       // 创建人ID
}

// TableName 返回表名
func (rp *RolePermission) TableName() string {
	return "sys_role_permission"
}

// NewRolePermission 创建角色-权限关联
func NewRolePermission(roleID, permissionID, operatorID int64) *RolePermission {
	return &RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
		CreatedAt:    time.Now().Unix(),
		CreatedBy:    operatorID,
	}
}

// UserRole 用户-角色关联
type UserRole struct {
	ID        int64 `db:"id" json:"id"`                // 主键ID
	UserID    int64 `db:"user_id" json:"userId"`       // 用户ID
	RoleID    int64 `db:"role_id" json:"roleId"`       // 角色ID
	CreatedAt int64 `db:"created_at" json:"createdAt"` // 创建时间
	CreatedBy int64 `db:"created_by" json:"createdBy"` // 创建人ID
}

// TableName 返回表名
func (ur *UserRole) TableName() string {
	return "sys_user_role"
}

// NewUserRole 创建用户-角色关联
func NewUserRole(userID, roleID, operatorID int64) *UserRole {
	return &UserRole{
		UserID:    userID,
		RoleID:    roleID,
		CreatedAt: time.Now().Unix(),
		CreatedBy: operatorID,
	}
}
