package models

import (
	"strings"

	"github.com/beego/beego/v2/client/orm"
)

type (
	SysRolePermission struct {
		Id           int            `json:"-" orm:"null"`
		RoleId       *SysRole       `json:"role_id" orm:"column(role_id);rel(fk);null"`
		PermissionId *SysPermission `json:"permission_id" orm:"column(permission_id);rel(fk);null"`
	}

	CompanyCompanyType struct {
		Id        int           `json:"-" orm:"null"`
		CompanyId *Company      `json:"company_id" orm:"column(company_id);rel(fk);null"`
		TypeId    *CompanyTypes `json:"type_id" orm:"column(type_id);rel(fk);null"`
	}

	CompanyBusinessUnit struct {
		Id             int           `json:"-" orm:"null"`
		CompanyId      *Company      `json:"company_id" orm:"column(company_id);rel(fk);null"`
		BusinessUnitId *BusinessUnit `json:"business_unit_id" orm:"column(business_unit_id);rel(fk);null"`
	}
)

func (t *SysRolePermission) TableName() string {
	return "sys_role_permission"
}

func SysRolePermissions() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(SysRolePermission))
}

func (t *CompanyCompanyType) TableName() string {
	return "company_type"
}

func CompanyCompanyTypes() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(CompanyCompanyType))
}

func (t *CompanyBusinessUnit) TableName() string {
	return "company_business_unit"
}

func CompanyBusinessUnits() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(CompanyBusinessUnit))
}

func init() {
	orm.RegisterModel(new(SysRolePermission), new(CompanyCompanyType), new(CompanyBusinessUnit))
}

type (
	CompanyTy struct {
		TypeId int    `json:"id"`
		Name   string `json:"name"`
	}

	CompanyBu struct {
		Id               int    `json:"id"`
		BusinessUnitCode string `json:"business_unit_code"`
		BusinessUnitName string `json:"business_unit_name"`
	}
)

func InsertCType(id int, company_type string) (int64, error) {

	o := orm.NewOrm()
	companies := Company{Id: id}
	o.Read(&companies)

	sql := "delete from company_type where company_id =?"
	if _, err := o.Raw(sql, id).Exec(); err != nil {
		return 0, err
	}

	typeidsarray := strings.Split(company_type, ",")

	companytype := new(CompanyTypes)
	var companytypes []CompanyTypes
	o.QueryTable(companytype).Filter("id__in", typeidsarray).All(&companytypes)

	m2m := o.QueryM2M(&companies, "CompanyTypes")
	num, err := m2m.Add(companytypes)
	if err != nil {
		return 0, err
	} else {
		return num, err
	}
}

func (t *CompanyBusinessUnit) InsertM2M(id int, business_type string) (int64, error) {

	o := orm.NewOrm()
	companies := Company{Id: id}
	o.Read(&companies)

	sql := "delete from company_business_unit where company_id =?"
	if _, err := o.Raw(sql, id).Exec(); err != nil {
		return 0, err
	}

	idArrays := strings.Split(business_type, ",")

	var businessUnit []BusinessUnit
	BusinessUnits().Filter("id__in", idArrays).All(&businessUnit)

	m2m := o.QueryM2M(&companies, "BusinessUnit")
	num, err := m2m.Add(businessUnit)
	if err != nil {
		return 0, err
	} else {
		return num, err
	}
}
