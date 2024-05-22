package models

import "github.com/beego/beego/v2/client/orm"

type (
	SysRolePermission struct {
		Id           int            `json:"-" orm:"null"`
		RoleId       *SysRole       `json:"role_id" orm:"column(role_id);rel(fk);null"`
		PermissionId *SysPermission `json:"permission_id" orm:"column(permission_id);rel(fk);null"`
	}
)

func (t *SysRolePermission) TableName() string {
	return "sys_role_permission"
}

func SysRolePermissions() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(SysRolePermission))
}

func init() {
	orm.RegisterModel(new(SysRolePermission))
}
