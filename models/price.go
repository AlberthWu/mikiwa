package models

import (
	"fmt"
	"mikiwa/utils"
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
		CompanyId     int        `json:"company_id" orm:"column(company_id)"`
		CompanyCode   string     `json:"company_code" orm:"column(company_code)"`
		NormalPrice   float64    `json:"normal_price" orm:"column(normal_price)"`
		DiscOne       float64    `json:"disc_one" orm:"column(disc_one)"`
		DiscTwo       float64    `json:"disc_two" orm:"column(disc_two)"`
		DiscTpr       float64    `json:"disc_tpr" orm:"column(disc_tpr)"`
		Price         float64    `json:"price" orm:"column(price)"`
		UomId         int        `json:"uom_id" orm:"column(uom_id)"`
		UomCode       string     `json:"uom_code" orm:"column(uom_code)"`
		Ratio         float64    `json:"ratio" orm:"column(ratio);digits(12);decimals(4);default(0)"`
		SureName      string     `json:"sure_name" orm:"column(sure_name)"`
		PriceType     string     `json:"price_type" orm:"column(price_type)"`
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
		Id             int       `json:"id" orm:"column(id);auto;pk"`
		PriceId        int       `json:"price_id" orm:"column(price_id)"`
		ProductId      int       `json:"product_id" orm:"column(product_id)"`
		ItemNo         int       `json:"item_no" orm:"column(item_no)"`
		UomId          int       `json:"uom_id" orm:"column(uom_id)"`
		IsDefault      int8      `json:"is_default" orm:"column(is_default)"`
		Ratio          float64   `json:"ratio" orm:"column(ratio);digits(12);decimals(4);default(0)"`
		DiscOne        float64   `json:"disc_one" orm:"column(disc_one)"`
		DiscTwo        float64   `json:"disc_two" orm:"column(disc_two)"`
		DiscTpr        float64   `json:"disc_tpr" orm:"column(disc_tpr)"`
		IsDefaultRatio float64   `json:"is_default_ratio" orm:"column(is_default_ratio);digits(8);decimals(4);default(0)"`
		FinalRatio     float64   `json:"final_ratio" orm:"column(final_ratio);digits(8);decimals(4);default(0)"`
		CreatedAt      time.Time `json:"created_at" orm:"column(created_at);type(timestamp);auto_now_add"`
		UpdatedAt      time.Time `json:"updated_at" orm:"column(updated_at);type(timestamp);auto_now"`
		DeletedAt      time.Time `json:"deleted_at" orm:"column(deleted_at);type(timestamp);null"`
		CreatedBy      string    `json:"created_by" orm:"column(created_by)"`
		UpdatedBy      string    `json:"updated_by" orm:"column(updated_by)"`
		DeletedBy      string    `json:"deleted_by" orm:"column(deleted_by)"`
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

type (
	PriceRtn struct {
		Id            int                  `json:"id"`
		EffectiveDate string               `json:"effective_date"`
		ExpiredDate   *string              `json:"expired_date"`
		CompanyId     SimpleCompanyRtnJson `json:"company_id"`
		ProductId     SimpleProductRtn     `json:"product_id"`
		SureName      string               `json:"sure_name"`
		StatusId      int                  `json:"status_id"`
		UomList       []orm.Params         `json:"uom"`
	}

	PriceRtnJson struct {
		Id                  int     `json:"id"`
		EffectiveDate       string  `json:"effective_date"`
		ExpiredDate         *string `json:"expired_date"`
		CompanyId           int     `json:"company_id"`
		CompanyCode         string  `json:"company_code"`
		CompanyName         string  `json:"company_name"`
		ProductId           int     `json:"product_id"`
		ProductCode         string  `json:"product_code"`
		ProductName         string  `json:"product_name"`
		ProductDivisionId   int     `json:"product_division_id"`
		ProductDivisionCode string  `json:"product_division_code"`
		ProductTypeId       int     `json:"product_type_id"`
		ProductTypeName     string  `json:"product_type_name"`
		NormalPrice         float64 `json:"normal_price"`
		DiscOne             float64 `json:"disc_one"`
		DiscOneDesc         string  `json:"disc_one_desc"`
		DiscTwo             float64 `json:"disc_two"`
		DiscTwoDesc         string  `json:"disc_two_desc"`
		DiscTpr             float64 `json:"disc_tpr"`
		DiscTprDesc         string  `json:"disc_tpr_desc"`
		Price               float64 `json:"price"`
		UomId               int     `json:"uom_id"`
		UomCode             string  `json:"uom_code"`
		Ratio               float64 `json:"ratio"`
		SureName            string  `json:"sure_name"`
		StatusId            int     `json:"status_id"`
		StatusData          string  `json:"status_data"`
	}
)

func (t *Price) GetById(id, user_id int) (m *PriceRtn, err error) {
	o := orm.NewOrm()
	cond := orm.NewCondition()
	cond1 := cond.And("deleted_at__isnull", true).And("id", id)
	qs := Prices().SetCond(cond1)
	err = qs.One(t)

	var company SimpleCompanyRtnJson
	o.Raw("select id,code,name  from companies where id  =" + utils.Int2String(t.CompanyId) + " ").QueryRow(&company)

	var product SimpleProductRtn
	o.Raw("select id,product_code,product_name,serial_number,lead_time,uom_id,uom_code from products where id  = " + utils.Int2String(t.ProductId) + " ").QueryRow(&product)

	ulist := t.GetDetail(id, user_id)

	var expireddate *string

	if t.ExpiredDate == nil {
		expireddate = nil
	} else {
		thedate := t.ExpiredDate.Format("2006-01-02")
		expireddate = &thedate
	}

	m = &PriceRtn{
		Id:            t.Id,
		EffectiveDate: t.EffectiveDate.Format("2006-01-02"),
		ExpiredDate:   expireddate,
		CompanyId:     company,
		ProductId:     product,
		StatusId:      int(t.StatusId),
		UomList:       ulist,
	}

	return m, err
}

func (t *Price) GetAll(keyword, field_name, match_mode, value_name string, p, size, allsize, user_id int, id int, division_ids, type_ids, status_ids, price_type string, issue_date, updated_at *string) (u utils.PageDynamic, err error) {
	o := orm.NewOrm()
	var m []orm.Params
	var c int

	// theDate date,updatedAt date,uId int,priceTypeId varchar(50),divisionIds varchar(7),typeIds varchar(7),statusIds varchar(5),reportTypeId int,userId int,keyword varchar(255),in TheField varchar(8000),in MatchMode varchar(8000),in ValueName varchar(8000), in limitVal int, in offsetVal int
	o.Raw("call sp_priceCount(?,?,"+utils.Int2String(id)+",'"+price_type+"','"+division_ids+"','"+type_ids+"','"+status_ids+"',1,"+utils.Int2String(user_id)+",'"+keyword+"','"+field_name+"','"+match_mode+"','"+value_name+"',null,null)", &issue_date, &updated_at).QueryRow(&c)

	if allsize == 1 && c > 0 {
		size = c
	}
	_, err = o.Raw("call sp_price(?,?,"+utils.Int2String(id)+",'"+price_type+"','"+division_ids+"','"+type_ids+"','"+status_ids+"',1,"+utils.Int2String(user_id)+",'"+keyword+"','"+field_name+"','"+match_mode+"','"+value_name+"',"+utils.Int2String(size)+", "+utils.Int2String((p-1)*size)+")", &issue_date, &updated_at).Values(&m)

	if c == 0 && err == nil {
		err = orm.ErrNoRows
		return utils.PaginationDynamic(int(c), p, size, "", "", "", "", "", "", "", m), err
	} else if err != nil {
		return utils.PaginationDynamic(int(c), p, size, "", "", "", "", "", "", "", m), err
	}

	return utils.PaginationDynamic(int(c), p, size, fmt.Sprintf("%v", m[0]["field_key"]), fmt.Sprintf("%v", m[0]["field_label"]), fmt.Sprintf("%v", m[0]["field_int"]), fmt.Sprintf("%v", m[0]["field_level"]), fmt.Sprintf("%v", m[0]["field_export"]), fmt.Sprintf("%v", m[0]["field_export_label"]), fmt.Sprintf("%v", m[0]["field_footer"]), m), err
}

func (c *Price) GetDetail(id, user_id int) (m []orm.Params) {
	o := orm.NewOrm()
	o.Raw("call sp_price(null,null," + utils.Int2String(id) + ",null,null,null,null,0," + utils.Int2String(user_id) + ",'',null,null,null,null,null)").Values(&m)
	return m
}
