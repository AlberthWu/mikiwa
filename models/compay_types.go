package models

import (
	"mikiwa/utils"
	"strconv"

	"github.com/beego/beego/v2/client/orm"
)

type CompanyTypes struct {
	Id        int        `json:"id" orm:"column(id);auto;pk"`
	Name      string     `json:"name" orm:"column(name);size(200)"`
	Position  int        `json:"position" orm:"column(position)"`
	AliasName string     `json:"alias_name" orm:"column(alias_name);size(100)"`
	IsAp      int        `json:"is_ap" orm:"column(is_ap)"`
	IsAr      int        `json:"is_ar" orm:"column(is_ar)"`
	Companies []*Company `json:"-" orm:"rel(m2m);rel_through(mikiwa/models.CompanyCompanyType)"`
}

func (t *CompanyTypes) TableName() string {
	return "company_types"
}

func CompanyTypess() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(CompanyTypes))
}

func init() {
	orm.RegisterModel(new(CompanyTypes))
}

func (t *CompanyTypes) Insert(m CompanyTypes) (*CompanyTypes, error) {
	o := orm.NewOrm()

	if _, err := o.Insert(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (t *CompanyTypes) Update(fields ...string) error {
	o := orm.NewOrm()
	if _, err := o.Update(t, fields...); err != nil {
		return err
	}
	return nil
}

type (
	CompanyTypeRtnJson struct {
		Id        int    `json:"id"`
		Name      string `json:"name"`
		AliasName string `json:"alias_name"`
	}
)

func (t *CompanyTypes) GetAll(keyword string, p, size int) (u utils.Page, err error) {

	var ctypes []CompanyTypes
	qs := CompanyTypess()
	cond := orm.NewCondition()
	cond = cond.And("name__icontains", keyword)
	qs = qs.SetCond(cond)
	count, _ := qs.Limit(-1).Count()
	_, err = qs.RelatedSel().Limit(size).Offset((p - 1) * size).All(&ctypes)
	c, _ := strconv.Atoi(strconv.FormatInt(count, 10))
	return utils.Pagination(c, p, size, ctypes), err
}

func (t *CompanyTypes) GetAllList(keyword string) (m []CompanyTypeRtnJson, err error) {
	o := orm.NewOrm()
	var num int64
	num, err = o.Raw("select id,`name`,alias_name from company_types where `name` like '%" + keyword + "%'").QueryRows(&m)

	if num == 0 && err == nil {
		err = orm.ErrNoRows
	}
	return m, err
}
