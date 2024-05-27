package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type (
	BusinessUnit struct {
		Id               int        `json:"id" orm:"column(id);auto;pk"`
		BusinessUnitCode string     `json:"business_unit_code" orm:"column(business_unit_code)"`
		BusinessUnitName string     `json:"business_unit_name" orm:"column(business_unit_name)"`
		Position         int8       `json:"position" orm:"column(position)"`
		StatusId         int8       `json:"status_id" orm:"column(status_id)"`
		CreatedAt        time.Time  `json:"created_at" orm:"column(created_at);type(timestamp);auto_now_add"`
		UpdatedAt        time.Time  `json:"updated_at" orm:"column(updated_at);type(timestamp);auto_now"`
		DeletedAt        time.Time  `json:"deleted_at" orm:"column(deleted_at);type(timestamp);null"`
		CreatedBy        string     `json:"created_by" orm:"column(created_by)"`
		UpdatedBy        string     `json:"updated_by" orm:"column(updated_by)"`
		DeletedBy        string     `json:"deleted_by" orm:"column(deleted_by)"`
		Companies        []*Company `json:"-" orm:"rel(m2m);rel_through(mikiwa/models.CompanyBusinessUnit)"`
	}
)

func (t *BusinessUnit) TableName() string {
	return "business_units"
}

func BusinessUnits() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(BusinessUnit))
}

func init() {
	orm.RegisterModel(new(BusinessUnit))
}

func (t *BusinessUnit) Insert(m BusinessUnit) (*BusinessUnit, error) {
	o := orm.NewOrm()

	if _, err := o.Insert(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (t *BusinessUnit) Update(fields ...string) error {
	o := orm.NewOrm()
	if _, err := o.Update(t, fields...); err != nil {
		return err
	}
	return nil
}

func (t *BusinessUnit) GetAllList(keyword string) (m []BusinessUnit, err error) {
	var num int64
	cond := orm.NewCondition()
	cond = cond.AndCond(cond.Or("business_unit_code__icontains", keyword).Or("business_unit_name__icontains", keyword)).And("status_id", 1).And("deleted_at__isnull", true)

	qs := BusinessUnits().SetCond(cond).OrderBy("business_unit_code")
	num, err = qs.Limit(100).Offset(0).All(&m)

	if num == 0 && err == nil {
		err = orm.ErrNoRows
	}
	return m, err
}
