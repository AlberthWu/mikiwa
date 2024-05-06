package models

import (
	"errors"
	"fmt"
	"mikiwa/utils"
	"strconv"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Pool struct {
	Id        int       `json:"id"  orm:"column(id);auto;pk"`
	Name      string    `json:"name" orm:"column(name)"`
	Status    int8      `json:"status" orm:"column(status);default(0)"`
	CreatedAt time.Time `json:"created_at" orm:"column(created_at);type(timestamp);auto_now_add"`
	UpdatedAt time.Time `json:"updated_at" orm:"column(updated_at);type(timestamp);auto_now"`
	DeletedAt time.Time `json:"deleted_at" orm:"column(deleted_at);type(timestamp);null"`
}

func (t *Pool) TableName() string {
	return "pools"
}

func Pools() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(Pool))
}

func init() {
	orm.RegisterModel(new(Pool))
}

func CheckPoolName(name string) bool {
	exist := Pools().Filter("name", name).Exist()
	return exist
}

func CheckPoolNamePut(id int, name string) bool {
	exist := Pools().Filter("name", name).Exclude("id", id).Exist()
	return exist
}

func InsertPool(m Pool) (*Pool, error) {
	o := orm.NewOrm()
	if _, err := o.Insert(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func UpdateByIdPool(m *Pool) (err error) {
	o := orm.NewOrm()
	v := Pool{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

type PoolRtnJson struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Status int    `json:"status"`
}

func GetAllPool(keyword string, p, size int) (u utils.Page, err error) {

	var pools []Pool
	cond := orm.NewCondition()
	cond1 := cond.And("deleted_at__isnull", true)
	qs := Pools().SetCond(cond1)
	cond2 := cond.AndCond(cond1).AndCond(cond.Or("name__icontains", keyword))
	qs = qs.SetCond(cond2)
	_, err = qs.Limit(size).Offset((p - 1) * size).All(&pools)
	count, _ := qs.Limit(-1).Count()
	c, _ := strconv.Atoi(strconv.FormatInt(count, 10))

	var poolrtn []PoolRtnJson
	for _, val := range pools {
		poolrtn = append(poolrtn, PoolRtnJson{
			Id:     val.Id,
			Name:   val.Name,
			Status: int(val.Status),
		})
	}

	if len(poolrtn) == 0 {
		return utils.Pagination(c, p, size, nil), errors.New("No data")
	}
	return utils.Pagination(c, p, size, poolrtn), err

}

func GetAllPoolLimit(keyword string) (m []PoolRtnJson, err error) {
	var pools []Pool
	cond := orm.NewCondition()
	cond1 := cond.And("deleted_at__isnull", true)
	qs := Pools().SetCond(cond1)
	cond2 := cond.AndCond(cond1).AndCond(cond.Or("name__icontains", keyword))
	qs = qs.SetCond(cond2)
	_, err = qs.Limit(100).Offset(0).All(&pools)

	var poolrtn []PoolRtnJson
	for _, val := range pools {
		poolrtn = append(poolrtn, PoolRtnJson{
			Id:     val.Id,
			Name:   val.Name,
			Status: int(val.Status),
		})
	}

	if len(poolrtn) == 0 {
		return poolrtn, errors.New("No data")
	}
	return poolrtn, err

}

func GetByIdPool(id int) (m *PoolRtnJson, err error) {
	var pools *Pool
	pools = &Pool{Id: id}
	if err := Pools().Filter("id", id).Filter("deleted_at__isnull", true).One(pools); err == orm.ErrNoRows {
		return nil, errors.New("No data")
	}
	m = &PoolRtnJson{
		Id:     pools.Id,
		Name:   pools.Name,
		Status: int(pools.Status),
	}
	return m, err
}
