package models

// RoleDept 角色和部门关联表。
type RoleDept struct {
	RoleID uint `gorm:"column:role_id;primaryKey" json:"roleId"`
	DeptID uint `gorm:"column:dept_id;primaryKey" json:"deptId"`
}

func (RoleDept) TableName() string { return "sys_role_dept" }
