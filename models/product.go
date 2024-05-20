package models

import (
	"mikiwa/utils"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type (
	Uom struct {
		Id       int    `json:"id" orm:"column(id);auto;pk"`
		UomCode  string `json:"uom_code" orm:"column(uom_code)"`
		UomName  string `json:"uom_name" orm:"column(uom_name)"`
		StatusId int8   `json:"status_id" orm:"column(status_id)"`
	}

	ProductDivision struct {
		Id           int    `json:"id" orm:"column(id);auto;pk"`
		DivisionCode string `json:"division_code" orm:"column(division_code)"`
		DivisionName string `json:"division_name" orm:"column(division_name)"`
		StatusId     int8   `json:"status_id" orm:"column(status_id)"`
	}

	ProductType struct {
		Id              int    `json:"id" orm:"column(id);auto;pk"`
		ProductTypeName string `json:"product_type_name" orm:"column(product_type_name)"`
		IsPurchase      int8   `json:"is_purchase" orm:"column(is_purchase)"`
		IsSales         int8   `json:"is_sales" orm:"column(is_sales)"`
		IsProduction    int8   `json:"is_production" orm:"column(is_production)"`
		StatusId        int8   `json:"status_id" orm:"column(status_id)"`
	}

	Product struct {
		Id                  int       `json:"id" orm:"column(id);auto;pk"`
		ProductCode         string    `json:"product_code" orm:"column(product_code)"`
		ProductName         string    `json:"product_name" orm:"column(product_name)"`
		ProductTypeId       int       `json:"product_type_id" orm:"column(product_type_id)"`
		ProductTypeName     string    `json:"product_type_name" orm:"column(product_type_name)"`
		ProductDivisionId   int       `json:"product_division_id" orm:"column(product_division_id)"`
		ProductDivisionName string    `json:"product_division_name" orm:"column(product_division_name)"`
		SerialNumber        string    `json:"serial_number" orm:"column(serial_number)"`
		UomId               int       `json:"uom_id" orm:"column(uom_id)"`
		UomCode             string    `json:"uom_code" orm:"column(uom_code)"`
		LeadTime            int       `json:"lead_time" orm:"column(lead_time)"`
		StatusId            int8      `json:"status_id" orm:"column(status_id)"`
		CreatedAt           time.Time `json:"created_at" orm:"column(created_at);type(timestamp);auto_now_add"`
		UpdatedAt           time.Time `json:"updated_at" orm:"column(updated_at);type(timestamp);auto_now"`
		DeletedAt           time.Time `json:"deleted_at" orm:"column(deleted_at);type(timestamp);null"`
		CreatedBy           string    `json:"created_by" orm:"column(created_by)"`
		UpdatedBy           string    `json:"updated_by" orm:"column(updated_by)"`
		DeletedBy           string    `json:"deleted_by" orm:"column(deleted_by)"`
	}
)

func (t *Uom) TableName() string {
	return "uoms"
}

func Uoms() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(Uom))
}

func (t *ProductDivision) TableName() string {
	return "product_divisions"
}

func ProductDivisions() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(ProductDivision))
}

func (t *ProductType) TableName() string {
	return "product_types"
}

func ProductTypes() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(ProductType))
}

func (t *Product) TableName() string {
	return "products"
}

func Products() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(Product))
}

func init() {
	orm.RegisterModel(new(Uom), new(ProductDivision), new(ProductType), new(Product))
}

func (t *Uom) Insert(m Uom) (*Uom, error) {
	o := orm.NewOrm()

	if _, err := o.Insert(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (t *Uom) Update(fields ...string) error {
	o := orm.NewOrm()
	if _, err := o.Update(t, fields...); err != nil {
		return err
	}
	return nil
}

func (t *ProductDivision) Insert(m ProductDivision) (*ProductDivision, error) {
	o := orm.NewOrm()

	if _, err := o.Insert(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (t *ProductDivision) Update(fields ...string) error {
	o := orm.NewOrm()
	if _, err := o.Update(t, fields...); err != nil {
		return err
	}
	return nil
}

func (t *ProductType) Insert(m ProductType) (*ProductType, error) {
	o := orm.NewOrm()

	if _, err := o.Insert(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (t *ProductType) Update(fields ...string) error {
	o := orm.NewOrm()
	if _, err := o.Update(t, fields...); err != nil {
		return err
	}
	return nil
}
func (t *Product) Insert(m Product) (*Product, error) {
	o := orm.NewOrm()

	if _, err := o.Insert(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (t *Product) Update(fields ...string) error {
	o := orm.NewOrm()
	if _, err := o.Update(t, fields...); err != nil {
		return err
	}
	return nil
}

func (t *Uom) CheckCode(id int, uom string) bool {
	exist := ProductTypes().Exclude("id", id).Filter("uom_code", uom).Exist()
	return exist
}

func (t *ProductDivision) CheckCode(id int, code string) bool {
	exist := ProductDivisions().Exclude("id", id).Filter("division_code", code).Exist()
	return exist
}

func (t *ProductType) CheckCode(id int, name string) bool {
	exist := ProductTypes().Exclude("id", id).Filter("type_name", name).Exist()
	return exist
}

func (t *Product) CheckCode(id int, code string) bool {
	exist := ProductTypes().Exclude("id", id).Filter("product_code", code).Filter("deleted_at__isnull", true).Exist()
	return exist
}

func (t *ProductDivision) GetById(id int) (m *ProductDivision, err error) {
	m = &ProductDivision{}
	cond := orm.NewCondition()
	cond1 := cond.And("id", id)
	qs := ProductDivisions().SetCond(cond1)

	if err = qs.One(m); err != nil {
		return nil, err
	}
	return m, err
}

func (t *ProductDivision) GetAll(keyword, field_name, match_mode, value_name string, p, size int, status_id string) (u utils.Page, err error) {

	var details []ProductDivision
	var d int64
	cond := orm.NewCondition()
	cond1 := cond.And("division_name__icontains", keyword).Or("division_code__icontains", keyword)

	if status_id != "" {
		var joinId []string
		ids := strings.Split(status_id, ",")
		for _, st := range ids {
			joinId = append(joinId, st)
		}
		cond1 = cond1.And("status_id", joinId)
	}

	qs := ProductDivisions().SetCond(cond1)

	d, err = qs.Limit(size).Offset((p - 1) * size).OrderBy("-id").All(&details)
	count, _ := qs.Limit(-1).Count()
	c, _ := strconv.Atoi(strconv.FormatInt(count, 10))

	if err == nil && d == 0 {
		err = orm.ErrNoRows
	}
	return utils.Pagination(c, p, size, details), err
}
