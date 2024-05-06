package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type (
	SysMenu struct {
		Id          int       `json:"id" orm:"column(id);auto;pk"`
		LevelNo     int       `json:"level_no" orm:"column(level_no)"`
		SeqNo       float64   `json:"seq_no" orm:"column(seq_no);digits(8);decimals(5);default(0)"`
		ParentId    int       `json:"parent_id" orm:"column(parent_id)"`
		Label       string    `json:"label" orm:"column(label)"`
		Icon        string    `json:"icon" orm:"column(icon)"`
		To          string    `json:"to" orm:"column(to)"`
		Status      int8      `json:"status" orm:"column(status); type(tinyint)"`
		FormDefault int8      `json:"form_default" orm:"column(form_default)"`
		FormName    string    `json:"form_name" orm:"column(form_name)"`
		CreatedAt   time.Time `json:"created_at" orm:"column(created_at);type(timestamp);auto_now_add"`
		UpdatedAt   time.Time `json:"updated_at" orm:"column(updated_at);type(timestamp);auto_now"`
		DeletedAt   time.Time `json:"deleted_at" orm:"column(deleted_at);type(timestamp);null"`
		// Users     []*Users  `json:"-" orm:"rel(m2m);rel_through(mikiwa/models.UsersSysMenu)"`
	}

	SysPermission struct {
		Id             int       `json:"id" orm:"column(id);auto;pk"`
		PermissionName string    `json:"permission_name" orm:"column(permission_name)"`
		CreatedAt      time.Time `json:"created_at" orm:"column(created_at);type(timestamp);auto_now_add"`
		UpdatedAt      time.Time `json:"updated_at" orm:"column(updated_at);type(timestamp);auto_now"`
		DeletedAt      time.Time `json:"deleted_at" orm:"column(deleted_at);type(timestamp);null"`
	}

	SysRole struct {
		Id        int       `json:"id" orm:"column(id);auto;pk"`
		ParentId  int       `json:"parent_id" orm:"column(parent_id)"`
		RoleName  string    `json:"role_name" orm:"column(role_name)"`
		Status    int8      `json:"status" orm:"column(status); type(tinyint)"`
		CreatedAt time.Time `json:"created_at" orm:"column(created_at);type(timestamp);auto_now_add"`
		UpdatedAt time.Time `json:"updated_at" orm:"column(updated_at);type(timestamp);auto_now"`
		DeletedAt time.Time `json:"deleted_at" orm:"column(deleted_at);type(timestamp);null"`
		// Users     []*Users  `json:"-" orm:"rel(m2m);rel_through(mikiwa/models.SysUserRole)"`
	}

	SysRoleMenuPermission struct {
		Id           int `json:"-" orm:"null"`
		RoleId       int `json:"role_id" orm:"column(role_id)"`
		MenuId       int `json:"menu_id" orm:"column(menu_id)"`
		PermissionId int `json:"permission_id" orm:"column(permission_id)"`
	}
)

func (t *SysRole) TableName() string {
	return "sys_roles"
}

func SysRoles() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(SysRole))
}

func (t *SysPermission) TableName() string {
	return "sys_permissions"
}

func SysPermissions() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(SysPermission))
}

func (t *SysMenu) TableName() string {
	return "sys_menus"
}

func SysMenus() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(SysMenu))
}

func (t *SysRoleMenuPermission) TableName() string {
	return "sys_role_menu_permission"
}

func SysRoleMenuPermissions() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(SysRoleMenuPermission))
}

func init() {
	orm.RegisterModel(new(SysRole), new(SysPermission), new(SysMenu), new(SysRoleMenuPermission))
}

func (t *SysRole) Insert(fields ...string) (*SysRole, error) {
	o := orm.NewOrm()

	if _, err := o.Insert(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (t *SysRole) Update(fields ...string) error {
	o := orm.NewOrm()
	if _, err := o.Update(t, fields...); err != nil {
		return err
	}
	return nil
}

func (t *SysPermission) Insert(fields ...string) (*SysPermission, error) {
	o := orm.NewOrm()

	if _, err := o.Insert(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (t *SysPermission) Update(fields ...string) error {
	o := orm.NewOrm()
	if _, err := o.Update(t, fields...); err != nil {
		return err
	}
	return nil
}

func (t *SysMenu) Insert(fields ...string) (*SysMenu, error) {
	o := orm.NewOrm()

	if _, err := o.Insert(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (t *SysMenu) Update(fields ...string) error {
	o := orm.NewOrm()
	if _, err := o.Update(t, fields...); err != nil {
		return err
	}
	return nil
}

func (t *SysRolePermission) Insert(fields ...string) (*SysRolePermission, error) {
	o := orm.NewOrm()

	if _, err := o.Insert(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (t *SysRolePermission) Update(fields ...string) error {
	o := orm.NewOrm()
	if _, err := o.Update(t, fields...); err != nil {
		return err
	}
	return nil
}

func CheckPrivileges(user_id, form_id, permmission int) bool {
	o := orm.NewOrm()
	status := false
	type RolePrivileges struct {
		RoleId         int
		RoleName       string
		MenuId         int
		FormName       string
		PermissionId   int
		PermissionName string
	}
	var role RolePrivileges
	o.Raw("call sp_SysUserPrivileges(?,?,?)", user_id, form_id, permmission).QueryRow(&role)
	if role.PermissionId > 0 {
		status = true
	}
	return status
}

type (
	SysMenuHeaderRtn struct {
		Id    int                `json:"id"`
		Label string             `json:"label"`
		Items []SysMenuParentRtn `json:"items"`
	}

	SysMenuParentRtn struct {
		Id    int               `json:"id"`
		Label string            `json:"label"`
		Icon  string            `json:"icon"`
		To    string            `json:"to"`
		Items []SysMenuChildRtn `json:"items,omitempty"`
	}

	SysMenuChildRtn struct {
		Id    int    `json:"id"`
		Label string `json:"label"`
		Icon  string `json:"icon"`
		To    string `json:"to"`
	}
)

func (t *SysMenu) GetAllMenu(id int) ([]SysMenuHeaderRtn, error) {
	var detail []SysMenu
	var m SysMenu
	// var []detail int
	o := orm.NewOrm()
	num, err := o.Raw("call sp_SysUserMenu(?,1,0)", id).QueryRows(&detail)
	var detailrtn []SysMenuHeaderRtn
	for _, list := range detail {

		plist := m.GetParentList(list.Id, id)
		detailrtn = append(detailrtn, SysMenuHeaderRtn{
			Id:    list.Id,
			Label: list.Label,
			Items: plist,
		})
	}

	if num == 0 {
		return nil, orm.ErrNoRows
	}
	return detailrtn, err
}

func (t *SysMenu) GetParentList(id, user_id int) []SysMenuParentRtn {
	var detail []SysMenu
	var m SysMenu
	// num, _ := SysMenus().Filter("parent_id", id).Filter("deleted_at__isnull", true).All(&detail)
	o := orm.NewOrm()
	num, _ := o.Raw("call sp_SysUserMenu(?,2,?)", user_id, id).QueryRows(&detail)

	var detailrtn []SysMenuParentRtn
	for _, list := range detail {
		clist := m.GetChildList(list.Id, user_id)
		detailrtn = append(detailrtn, SysMenuParentRtn{
			Id:    list.Id,
			Label: list.Label,
			Icon:  list.Icon,
			To:    list.To,
			Items: clist,
		})
	}

	if num == 0 {
		return nil
	}
	return detailrtn
}

func (t *SysMenu) GetChildList(id, user_id int) []SysMenuChildRtn {
	var detail []SysMenu

	o := orm.NewOrm()
	num, _ := o.Raw("call sp_SysUserMenu(?,3,?)", user_id, id).QueryRows(&detail)

	var detailrtn []SysMenuChildRtn
	for _, list := range detail {
		detailrtn = append(detailrtn, SysMenuChildRtn{
			Id:    list.Id,
			Label: list.Label,
			Icon:  list.Icon,
			To:    list.To,
		})
	}

	if num == 0 {
		return nil
	}
	return detailrtn
}
