package models

type (
	SysRolePermission struct {
		Id           int            `json:"-" orm:"null"`
		RoleId       *SysRole       `json:"role_id" orm:"column(role_id);rel(fk);null"`
		PermissionId *SysPermission `json:"permission_id" orm:"column(permission_id);rel(fk);null"`
	}
)