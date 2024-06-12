package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type (
	Price struct {
		Id            int        `json:"id" orm:"column(id);auto;pk"`
		EffectiveDate time.Time  `json:"effective_date" orm:"column(effective_date);type(date)"`
		ExpiredDate   *time.Time `json:"expired_date" orm:"column(expired_date);type(date);null"`
		ProductId     int        `json:"product_id" orm:"column(product_id)"`
		ProductCode   string     `json:"product_code" orm:"column(product_code)"`
		NormalPrice   float64    `json:"normal_price" orm:"column(normal_price)"`
		UomIdOne      int        `json:"uom_id_one" orm:"column(uom_id_one)"`
		UomCodeOne    string     `json:"uom_code_one" orm:"column(uom_code_one)"`
		Ratio         float64    `json:"ratio" orm:"column(ratio);digits(12);decimals(4);default(0)"`
		UomIdTwo      int        `json:"uom_id_two" orm:"column(uom_id_two)"`
		UomCodeTwo    string     `json:"uom_code_two" orm:"column(uom_code_two)"`
		DiscOne       float64    `json:"disc_one" orm:"column(disc_one)"`
		DiscTwo       float64    `json:"disc_two" orm:"column(disc_two)"`
		DiscTpr       float64    `json:"disc_type" orm:"column(disc_type)"`
		Price         float64    `json:"price" orm:"column(price)"`
		SureName      string     `json:"sure_name" orm:"column(sure_name)"`
		PriceType     string     `json:"price_type" orm:"column(price_type)"`
		Qty           float64    `json:"qty" orm:"column(qty);digits(12);decimals(4);default(0)"`
		UomId         int        `json:"uom_id" orm:"column(uom_id)"`
		UomCode       string     `json:"uom_code" orm:"column(uom_code)"`
		StatusId      int8       `json:"status_id" orm:"column(status_id)"`
		CreatedAt     time.Time  `json:"created_at" orm:"column(created_at);type(timestamp);auto_now_add"`
		UpdatedAt     time.Time  `json:"updated_at" orm:"column(updated_at);type(timestamp);auto_now"`
		DeletedAt     time.Time  `json:"deleted_at" orm:"column(deleted_at);type(timestamp);null"`
		CreatedBy     string     `json:"created_by" orm:"column(created_by)"`
		UpdatedBy     string     `json:"updated_by" orm:"column(updated_by)"`
		DeletedBy     string     `json:"deleted_by" orm:"column(deleted_by)"`
		Companies     []*Company `json:"-" orm:"reverse(many);rel_through(mikiwa/models.PriceCompany)"`
	}

	PriceProductUom struct {
		Id        int     `json:"id" orm:"column(id);auto;pk"`
		PriceId   int     `json:"price_id" orm:"column(price_id)"`
		CompanyId int     `json:"company_id" orm:"column(company_id)"`
		ProductId int     `json:"product_id" orm:"column(product_id)"`
		ItemNo    int     `json:"item_no" orm:"column(item_no)"`
		UomId     int     `json:"uom_id" orm:"column(uom_id)"`
		Ratio     float64 `json:"ratio" orm:"column(ratio);digits(12);decimals(4);default(0)"`
	}
)

func (t *Price) TableName() string {
	return "prices"
}

func Prices() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(Price))
}

func (t *PriceProductUom) TableName() string {
	return "price_product_uom"
}

func PriceProductUoms() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(PriceProductUom))
}

func init() {
	orm.RegisterModel(new(PriceProductUom), new(Price))
}

func (t *Price) Insert(m Price) (*Price, error) {
	o := orm.NewOrm()

	if _, err := o.Insert(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (t *Price) Update(fields ...string) error {
	o := orm.NewOrm()
	if _, err := o.Update(t, fields...); err != nil {
		return err
	}
	return nil
}

func (t *PriceProductUom) Insert(m PriceProductUom) (*PriceProductUom, error) {
	o := orm.NewOrm()

	if _, err := o.Insert(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (t *PriceProductUom) Update(fields ...string) error {
	o := orm.NewOrm()
	if _, err := o.Update(t, fields...); err != nil {
		return err
	}
	return nil
}
