package models

import (
	"mikiwa/utils"
	"strconv"

	"github.com/beego/beego/v2/client/orm"
)

type (
	Bank struct {
		Id         int    `json:"id" orm:"column(id);auto;pk"`
		Name       string `json:"name" orm:"column(name);size(200)"`
		Swift      string `json:"swift" orm:"column(swift);size(200)"`
		Bi         string `json:"bi" orm:"column(bi);size(200)"`
		DistrictId int    `json:"district_id"  orm:"null;column(district_id)"` //null --> left join,not null --> inner join
		StateId    int    `json:"state_id"  orm:"null;column(state_id)"`
		CityId     int    `json:"city_id"  orm:"null;column(city_id)"` //omitempty hasil json akan tidak muncul jika tidak ada data (blank)
	}
)

func (t *Bank) TableName() string {
	return "banks"
}

func Banks() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(Bank))
}

func init() {
	orm.RegisterModel(new(Bank))
}

func (t *Bank) Insert(m Bank) (*Bank, error) {
	o := orm.NewOrm()

	if _, err := o.Insert(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (t *Bank) Update(fields ...string) error {
	o := orm.NewOrm()
	if _, err := o.Update(t, fields...); err != nil {
		return err
	}
	return nil
}

type (
	BankReturn struct {
		Id    int    `json:"id"`
		Name  string `json:"name"`
		Swift string `json:"swift"`
	}
)

func (t *Bank) GetById(id int) (m *Bank, err error) {
	m = &Bank{}
	cond := orm.NewCondition()
	cond1 := cond.And("id", id)
	qs := Banks().SetCond(cond1)

	if err = qs.One(m); err != nil {
		return nil, err
	}
	return m, err
}

func (t *Bank) GetAll(keyword, field_name, match_mode, value_name string, p, size int) (u utils.Page, err error) {

	var details []Bank
	var d int64
	cond := orm.NewCondition()
	cond1 := cond.AndCond(cond.Or("name__icontains", keyword).Or("swift__icontains", keyword))
	qs := Banks().SetCond(cond1)

	d, err = qs.Limit(size).Offset((p - 1) * size).OrderBy("-id").All(&details)
	count, _ := qs.Limit(-1).Count()
	c, _ := strconv.Atoi(strconv.FormatInt(count, 10))

	if err == nil && d == 0 {
		err = orm.ErrNoRows
	}
	return utils.Pagination(c, p, size, details), err
}

func (t *Bank) GetAllList(keyword string) (m []Bank, err error) {
	var num int64
	cond := orm.NewCondition()
	cond = cond.AndCond(cond.Or("name__icontains", keyword).Or("swift__icontains", keyword))

	qs := Banks().SetCond(cond).OrderBy("name")
	num, err = qs.Limit(100).Offset(0).All(&m)

	if num == 0 && err == nil {
		err = orm.ErrNoRows
	}
	return m, err
}
