package models

import (
	"mikiwa/utils"

	"github.com/beego/beego/v2/client/orm"
)

type City struct {
	Id       int    `json:"id" orm:"column(id);auto;pk" `
	ParentId int    `json:"parent_id" orm:"column(parent_id)"`
	Name     string `json:"name" orm:"column(name);size(200);null" valid:"Required;MinSize(5);MaxSize(10)"`
	Zip      string `json:"zip" orm:"column(zip);size(5);null" valid:"Required;range(5,10)"`
}

func (t *City) TableName() string {
	return "cities"
}

func Cities() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(City))
}

func init() {
	orm.RegisterModel(new(City))
}

func (t *City) Insert(m City) (*City, error) {
	o := orm.NewOrm()

	if _, err := o.Insert(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (t *City) Update(fields ...string) error {
	o := orm.NewOrm()
	if _, err := o.Update(t, fields...); err != nil {
		return err
	}
	return nil
}

type CityRtnJson struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (t *City) GetById(id int) (m *City, err error) {
	m = &City{}
	cond := orm.NewCondition()
	cond1 := cond.And("id", id)
	qs := Cities().SetCond(cond1)

	if err = qs.One(m); err != nil {
		return nil, err
	}
	return m, err
}

func (t *City) GetAll(keyword, field_name, match_mode, value_name string, p, size int, status_id string) (u utils.Page, err error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	var cities []CityRtnJson
	var pagedata []CityRtnJson
	qb.Select("c1.id", "concat(c1.name,', ',c2.name,', ', c3.name) name").
		From("cities c1").
		LeftJoin("cities c2").On("c2.id=c1.parent_id").
		LeftJoin("cities c3").On("c3.id=c2.parent_id").
		Where("c2.parent_id > 0").And("c1.zip is null").And("(c1.name like '%"+keyword+"%' or c2.name like '%"+keyword+"%' or c3.name like '%"+keyword+"%')").
		OrderBy("c2.parent_id", "c1.parent_id", "c1.id")
	sql := qb.String()
	c, err := o.Raw(sql).QueryRows(&cities)
	for idx := (p - 1) * size; idx < p*size && idx < int(c); idx++ {
		pagedata = append(pagedata, cities[idx])
	}

	if len(pagedata) == 0 && err == nil {
		err = orm.ErrNoRows
	}
	return utils.Pagination(int(c), p, size, pagedata), err
}

func (t *City) GetAllList(keyword string) (m []City, err error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("c1.id", "concat(c1.name,', ',c2.name,', ', c3.name) name").
		From("cities c1").
		LeftJoin("cities c2").On("c2.id=c1.parent_id").
		LeftJoin("cities c3").On("c3.id=c2.parent_id").
		Where("c2.parent_id > 0").And("c1.zip is null").And("(c1.name like  '%"+keyword+"%' or c2.name like  '%"+keyword+"%' or c3.name like '%"+keyword+"%')").
		OrderBy("c2.parent_id", "c1.parent_id", "c1.id").Limit(100).Offset(0)
	sql := qb.String()

	_, err = o.Raw(sql).QueryRows(&m)

	if len(m) == 0 && err == nil {
		err = orm.ErrNoRows
	}
	return m, err

}
