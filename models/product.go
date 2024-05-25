package models

import (
	"fmt"
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
		Id           int    `json:"id" orm:"column(id);auto;pk"`
		TypeName     string `json:"type_name" orm:"column(type_name)"`
		IsPurchase   int8   `json:"is_purchase" orm:"column(is_purchase)"`
		IsSales      int8   `json:"is_sales" orm:"column(is_sales)"`
		IsProduction int8   `json:"is_production" orm:"column(is_production)"`
	}

	Product struct {
		Id                  int       `json:"id" orm:"column(id);auto;pk"`
		ProductCode         string    `json:"product_code" orm:"column(product_code)"`
		ProductName         string    `json:"product_name" orm:"column(product_name)"`
		ProductTypeId       int       `json:"product_type_id" orm:"column(product_type_id)"`
		ProductTypeName     string    `json:"product_type_name" orm:"column(product_type_name)"`
		ProductDivisionId   int       `json:"product_division_id" orm:"column(product_division_id)"`
		ProductDivisionCode string    `json:"product_division_code" orm:"column(product_division_code)"`
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

	ProductUom struct {
		Id        int     `json:"-" orm:"null"`
		ProductId int     `json:"product_id" orm:"column(product_id)"`
		ItemNo    int     `json:"item_no" orm:"column(item_no)"`
		UomId     int     `json:"uom_id" orm:"column(uom_id)"`
		Ratio     float64 `json:"ratio" orm:"column(ratio);digits(8);decimals(4);default(0)"`
		IsDefault int8    `json:"is_default" orm:"column(is_default)"`
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
func (t *ProductUom) TableName() string {
	return "product_uom"
}

func ProductUoms() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(ProductUom))
}

func init() {
	orm.RegisterModel(new(Uom), new(ProductDivision), new(ProductType), new(Product), new(ProductUom))
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
	exist := Uoms().Exclude("id", id).Filter("uom_code", uom).Exist()
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
	exist := Products().Exclude("id", id).Filter("product_code", code).Filter("deleted_at__isnull", true).Exist()
	return exist
}

type (
	ProductTypeRtn struct {
		Id           int    `json:"id"`
		TypeName     string `json:"type_name"`
		IsPurchase   int8   `json:"is_purchase"`
		IsSales      int8   `json:"is_sales"`
		IsProduction int8   `json:"is_production"`
	}

	ProductDivisionRtn struct {
		Id           int    `json:"id"`
		DivisionCode string `json:"division_code"`
		DivisionName string `json:"division_name"`
	}

	ProductRtn struct {
		Id                int                 `json:"id"`
		ProductCode       string              `json:"product_code"`
		ProductName       string              `json:"product_name"`
		ProductTypeId     ProductTypeRtn      `json:"product_type_id"`
		ProductDivisionId ProductDivisionRtn  `json:"product_division_id"`
		SerialNumber      string              `json:"serial_number"`
		LeadTime          int                 `json:"lead_time"`
		UomId             int                 `json:"uom_id"`
		UomCode           string              `json:"uom_code"`
		StatusId          int                 `json:"status_id"`
		UomList           []ProductUomRtnJson `json:"uom"`
		Document          []DocumentRtn       `json:"document"`
	}

	ProductUomRtnJson struct {
		ItemNo      int     `json:"item_no"`
		ProductId   int     `json:"product_id"`
		ProductCode string  `json:"product_code"`
		UomId       int     `json:"uom_id"`
		UomCode     string  `json:"uom_code"`
		Ratio       float64 `json:"ratio"`
		IsDefault   int8    `json:"is_default"`
	}

	SimpleProductRtn struct {
		Id           int    `json:"id"`
		ProductCode  string `json:"product_code"`
		ProductName  string `json:"product_name"`
		SerialNumber string `json:"serial_number"`
		LeadTime     int    `json:"lead_time"`
		UomId        int    `json:"uom_id"`
		UomCode      string `json:"uom_code"`
	}
)

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
	cond1 := cond.AndCond(cond.Or("division_name__icontains", keyword).Or("division_code__icontains", keyword))

	if status_id != "" {
		var joinId []string
		ids := strings.Split(status_id, ",")
		joinId = append(joinId, ids...)
		cond1 = cond1.And("status_id__in", joinId)
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

func (t *ProductDivision) GetAllList(keyword string) (m []ProductDivision, err error) {
	var num int64
	cond := orm.NewCondition()
	cond = cond.AndCond(cond.Or("division_code__icontains", keyword).Or("division_name__icontains", keyword)).And("status_id", 1)

	qs := ProductDivisions().SetCond(cond).OrderBy("division_code")
	num, err = qs.Limit(100).Offset(0).All(&m)

	if num == 0 && err == nil {
		err = orm.ErrNoRows
	}
	return m, err
}

func (t *ProductType) GetById(id int) (m *ProductType, err error) {
	m = &ProductType{}
	cond := orm.NewCondition()
	cond1 := cond.And("id", id)
	qs := ProductTypes().SetCond(cond1)

	if err = qs.One(m); err != nil {
		return nil, err
	}
	return m, err
}

func (t *ProductType) GetAll(keyword, field_name, match_mode, value_name string, p, size int, is_purchase, is_sales, is_production string) (u utils.Page, err error) {

	var details []ProductType
	var d int64
	cond := orm.NewCondition()
	cond1 := cond.And("type_name__icontains", keyword)

	if is_purchase != "" {
		cond1 = cond1.And("is_purchase", is_purchase)
	}
	if is_sales != "" {
		cond1 = cond1.And("is_sales", is_sales)
	}
	if is_production != "" {
		cond1 = cond1.And("is_production", is_production)
	}

	qs := ProductTypes().SetCond(cond1)

	d, err = qs.Limit(size).Offset((p - 1) * size).OrderBy("-id").All(&details)
	count, _ := qs.Limit(-1).Count()
	c, _ := strconv.Atoi(strconv.FormatInt(count, 10))

	if err == nil && d == 0 {
		err = orm.ErrNoRows
	}
	return utils.Pagination(c, p, size, details), err
}

func (t *ProductType) GetAllList(keyword, is_purchase, is_sales, is_production string) (m []ProductType, err error) {
	var num int64
	cond := orm.NewCondition()
	cond = cond.And("type_name__icontains", keyword)
	if is_purchase != "" {
		cond = cond.And("is_purchase", is_purchase)
	}

	if is_sales != "" {
		cond = cond.And("is_sales", is_sales)
	}

	if is_production != "" {
		cond = cond.And("is_production", is_production)
	}

	qs := ProductTypes().SetCond(cond).OrderBy("type_name")
	num, err = qs.Limit(100).Offset(0).All(&m)

	if num == 0 && err == nil {
		err = orm.ErrNoRows
	}
	return m, err
}

func (t *Uom) GetById(id int) (m *Uom, err error) {
	m = &Uom{}
	cond := orm.NewCondition()
	cond1 := cond.And("id", id)
	qs := Uoms().SetCond(cond1)

	if err = qs.One(m); err != nil {
		return nil, err
	}
	return m, err
}

func (t *Uom) GetAll(keyword, field_name, match_mode, value_name string, p, size int, status_id string) (u utils.Page, err error) {

	var details []Uom
	var d int64
	cond := orm.NewCondition()
	cond1 := cond.AndCond(cond.Or("uom_code__icontains", keyword).Or("uom_name__icontains", keyword))

	if status_id != "" {
		var joinId []string
		ids := strings.Split(status_id, ",")
		joinId = append(joinId, ids...)
		cond1 = cond1.And("status_id__in", joinId)
	}

	qs := Uoms().SetCond(cond1)

	d, err = qs.Limit(size).Offset((p - 1) * size).OrderBy("-id").All(&details)
	count, _ := qs.Limit(-1).Count()
	c, _ := strconv.Atoi(strconv.FormatInt(count, 10))

	if err == nil && d == 0 {
		err = orm.ErrNoRows
	}
	return utils.Pagination(c, p, size, details), err
}

func (t *Uom) GetAllList(keyword string) (m []Uom, err error) {
	var num int64
	cond := orm.NewCondition()
	cond = cond.AndCond(cond.Or("uom_code__icontains", keyword).Or("uom_name__icontains", keyword)).And("status_id", 1)

	qs := Uoms().SetCond(cond).OrderBy("uom_code")
	num, err = qs.Limit(100).Offset(0).All(&m)

	if num == 0 && err == nil {
		err = orm.ErrNoRows
	}
	return m, err
}

func (t *Product) GetById(id, user_id int) (m *ProductRtn, err error) {
	o := orm.NewOrm()
	var d Product
	var doc Document
	cond := orm.NewCondition()
	cond1 := cond.And("deleted_at__isnull", true).And("id", id)
	qs := Products().SetCond(cond1)
	err = qs.One(&d)

	var producttype ProductTypeRtn
	o.Raw("select id,type_name,is_purchase,is_sales,is_production from product_types where id  =" + utils.Int2String(d.ProductTypeId) + " ").QueryRow(&producttype)

	var productdivision ProductDivisionRtn
	o.Raw("select id,division_code,division_name from product_divisions where id  = " + utils.Int2String(d.ProductDivisionId) + " ").QueryRow(&productdivision)

	ulist := d.GetDetail(id, user_id)
	dlist := doc.GetDocument(id, "product")

	m = &ProductRtn{
		Id:                d.Id,
		ProductCode:       d.ProductCode,
		ProductName:       d.ProductName,
		ProductTypeId:     producttype,
		ProductDivisionId: productdivision,
		SerialNumber:      d.SerialNumber,
		LeadTime:          d.LeadTime,
		UomId:             d.UomId,
		UomCode:           d.UomCode,
		StatusId:          int(d.StatusId),
		UomList:           ulist,
		Document:          dlist,
	}

	return m, err
}

func (c *Product) GetDetail(id, user_id int) (m []ProductUomRtnJson) {
	o := orm.NewOrm()
	o.Raw("call sp_Product(null," + utils.Int2String(id) + ",'','','','','',null," + utils.Int2String(user_id) + ",0,'',null,null,null,null, null)").QueryRows(&m)
	return m
}

func (t *Product) GetAll(keyword, field_name, match_mode, value_name string, p, size, allsize, user_id int, division_ids, type_ids, production_ids, purchase_ids, sales_ids, status_ids string, updated_at *string) (u utils.PageDynamic, err error) {
	o := orm.NewOrm()
	var m []orm.Params
	var c int

	o.Raw("call sp_ProductCount(?,0,'"+division_ids+"','"+type_ids+"','"+purchase_ids+"','"+sales_ids+"','"+production_ids+"','"+status_ids+"',"+utils.Int2String(user_id)+",1,'"+keyword+"','"+field_name+"','"+match_mode+"','"+value_name+"',null,null)", &updated_at).QueryRow(&c)

	if allsize == 1 && c > 0 {
		size = c
	}
	_, err = o.Raw("call sp_Product(?,0,'"+division_ids+"','"+type_ids+"','"+purchase_ids+"','"+sales_ids+"','"+production_ids+"','"+status_ids+"',"+utils.Int2String(user_id)+",1,'"+keyword+"','"+field_name+"','"+match_mode+"','"+value_name+"',"+utils.Int2String(size)+", "+utils.Int2String((p-1)*size)+")", &updated_at).Values(&m)

	if c == 0 && err == nil {
		err = orm.ErrNoRows
		return utils.PaginationDynamic(int(c), p, size, "", "", "", "", "", "", "", m), err
	} else if err != nil {
		return utils.PaginationDynamic(int(c), p, size, "", "", "", "", "", "", "", m), err
	}

	return utils.PaginationDynamic(int(c), p, size, fmt.Sprintf("%v", m[0]["field_key"]), fmt.Sprintf("%v", m[0]["field_label"]), fmt.Sprintf("%v", m[0]["field_int"]), fmt.Sprintf("%v", m[0]["field_level"]), fmt.Sprintf("%v", m[0]["field_export"]), fmt.Sprintf("%v", m[0]["field_export_label"]), fmt.Sprintf("%v", m[0]["field_footer"]), m), err
}

func (t *Product) GetAllDetail(keyword, field_name, match_mode, value_name string, user_id int, division_ids, type_ids, production_ids, purchase_ids, sales_ids, status_ids string, updated_at *string) (m []orm.Params, err error) {
	o := orm.NewOrm()
	var num int64
	if num, err = o.Raw("call sp_Product(?,0,'"+division_ids+"','"+type_ids+"','"+purchase_ids+"','"+sales_ids+"','"+production_ids+"','"+status_ids+"',"+utils.Int2String(user_id)+",0,'"+keyword+"','"+field_name+"','"+match_mode+"','"+value_name+"',null, null)", &updated_at).Values(&m); num == 0 && err == nil {
		err = orm.ErrNoRows
	}
	return m, err
}

func (t *Product) GetAllListRaw(keyword string) (m []SimpleProductRtn, err error) {
	o := orm.NewOrm()
	d, err := o.Raw("select id,product_code,product_name,serial_number,lead_time,uom_id,uom_code from products where deleted_at is null and status_id = 1 and product_type_id = 1 and (product_code like '%" + keyword + "%' or product_name like '%" + keyword + "%' or serial_number like '%" + keyword + "%') ").QueryRows(&m)

	if d == 0 && err == nil {
		err = orm.ErrNoRows
	}
	return m, err
}

func (t *Product) GetAllListWip(keyword string) (m []SimpleProductRtn, err error) {
	o := orm.NewOrm()
	d, err := o.Raw("select id,product_code,product_name,serial_number,lead_time,uom_id,uom_code from products where deleted_at is null and status_id = 1 and product_type_id = 2 and (product_code like '%" + keyword + "%' or product_name like '%" + keyword + "%' or serial_number like '%" + keyword + "%') ").QueryRows(&m)

	if d == 0 && err == nil {
		err = orm.ErrNoRows
	}
	return m, err
}

func (t *Product) GetAllListFinishing(keyword string) (m []SimpleProductRtn, err error) {
	o := orm.NewOrm()
	d, err := o.Raw("select id,product_code,product_name,serial_number,lead_time,uom_id,uom_code from products where deleted_at is null and status_id = 1 and product_type_id = 3 and (product_code like '%" + keyword + "%' or product_name like '%" + keyword + "%' or serial_number like '%" + keyword + "%') ").QueryRows(&m)

	if d == 0 && err == nil {
		err = orm.ErrNoRows
	}
	return m, err
}

func (t *Product) GetAllListAcc(keyword string) (m []SimpleProductRtn, err error) {
	o := orm.NewOrm()
	d, err := o.Raw("select id,product_code,product_name,serial_number,lead_time,uom_id,uom_code from products where deleted_at is null and status_id = 1 and product_type_id = 4 and (product_code like '%" + keyword + "%' or product_name like '%" + keyword + "%' or serial_number like '%" + keyword + "%') ").QueryRows(&m)

	if d == 0 && err == nil {
		err = orm.ErrNoRows
	}
	return m, err
}

func (t *Product) GetAllListOthers(keyword string) (m []SimpleProductRtn, err error) {
	o := orm.NewOrm()
	d, err := o.Raw("select id,product_code,product_name,serial_number,lead_time,uom_id,uom_code from products where deleted_at is null and status_id = 1 and product_type_id = 5 and (product_code like '%" + keyword + "%' or product_name like '%" + keyword + "%' or serial_number like '%" + keyword + "%') ").QueryRows(&m)

	if d == 0 && err == nil {
		err = orm.ErrNoRows
	}
	return m, err
}

func (t *Product) GetAllListRecycle(keyword string) (m []SimpleProductRtn, err error) {
	o := orm.NewOrm()
	d, err := o.Raw("select id,product_code,product_name,serial_number,lead_time,uom_id,uom_code from products where deleted_at is null and status_id = 1 and product_type_id = 6 and (product_code like '%" + keyword + "%' or product_name like '%" + keyword + "%' or serial_number like '%" + keyword + "%') ").QueryRows(&m)

	if d == 0 && err == nil {
		err = orm.ErrNoRows
	}
	return m, err
}

func (t *Product) Document(id, user_id int, folder_name string) (m []DocumentRtn) {
	var doc Document
	m = doc.GetDocument(id, folder_name)
	return m
}
