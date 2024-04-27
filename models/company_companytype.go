package models

import "github.com/beego/beego/v2/client/orm"

type CompanyCompanyType struct {
	Id        int           `json:"-" orm:"null"`
	CompanyId *Company      `json:"company_id" orm:"column(company_id);rel(fk);null"`
	TypeId    *CompanyTypes `json:"type_id" orm:"column(type_id);rel(fk);null"`
}

func (t *CompanyCompanyType) TableName() string {
	return "company_type"
}

func CompanyCompanyTypes() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(CompanyCompanyType))
}

func init() {
	orm.RegisterModel(new(CompanyCompanyType))
}

func (t *CompanyCompanyType) Insert(m CompanyCompanyType) (*CompanyCompanyType, error) {
	o := orm.NewOrm()

	if _, err := o.Insert(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (t *CompanyCompanyType) Update(fields ...string) error {
	o := orm.NewOrm()
	if _, err := o.Update(t, fields...); err != nil {
		return err
	}
	return nil
}
