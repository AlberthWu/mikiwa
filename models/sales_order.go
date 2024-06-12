package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type (
	SalesOrder struct {
		Id              int       `json:"id" orm:"column(id);auto;pk"`
		ReferenceNo     string    `json:"reference_no" orm:"column(reference_no)"`
		IssueDate       time.Time `json:"issue_date" orm:"column(issue_date);type(date)"`
		CustomerId      int       `json:"customer_id" orm:"column(customer_id)"`
		CustomerName    string    `json:"customer_name" orm:"column(customer_name)"`
		Terms           int       `json:"terms" orm:"column(terms)"`
		DeliveryAddress string    `json:"delivery_address" orm:"column(delivery_address)"`
		LeadTime        int       `json:"lead_time" orm:"column(lead_time)"`
		Subtotal        float64   `json:"subtotal" orm:"column(subtotal);digits(18);decimals(2);default(0)"`
		TotalDisc       float64   `json:"total_disc" orm:"column(total_disc);digits(18);decimals(2);default(0)"`
		Dpp             float64   `json:"dpp" orm:"column(dpp);digits(18);decimals(2);default(0)"`
		Ppn             int       `json:"ppn" orm:"column(ppn)"`
		PpnAmount       float64   `json:"ppn_amount" orm:"column(ppn_amount);digits(18);decimals(2);default(0)"`
		Total           float64   `json:"total" orm:"column(total);digits(18);decimals(2);default(0)"`
		CreatedAt       time.Time `json:"created_at" orm:"column(created_at);type(timestamp);auto_now_add"`
		UpdatedAt       time.Time `json:"updated_at" orm:"column(updated_at);type(timestamp);auto_now"`
		DeletedAt       time.Time `json:"deleted_at" orm:"column(deleted_at);type(timestamp);null"`
		CreatedBy       string    `json:"created_by" orm:"column(created_by)"`
		UpdatedBy       string    `json:"updated_by" orm:"column(updated_by)"`
		DeletedBy       string    `json:"deleted_by" orm:"column(deleted_by)"`
	}

	SalesOrderDetail struct {
		Id           int       `json:"id" orm:"column(id);auto;pk"`
		SalesOrderId int       `json:"sales_order_id" orm:"column(sales_order_id)"`
		ReferenceNo  string    `json:"reference_no" orm:"column(reference_no)"`
		IssueDate    time.Time `json:"issue_date" orm:"column(issue_date);type(date)"`
		ItemNo       int       `json:"item_no" orm:"column(item_no)"`
		ProductId    int       `json:"product_id" orm:"column(product_id)"`
		ProductCode  string    `json:"product_code" orm:"column(product_code)"`
		QtyFormulaId int       `json:"qty_formula_id" orm:"column(qty_formula_id)"`
		Qty1         float64   `json:"qty1" orm:"column(qty1);digits(12);decimals(2);default(0)"`
		UomId        int       `json:"uom_id" orm:"column(uom_id)"`
		UomCode      string    `json:"uom_code" orm:"column(uom_code)"`
		Ratio        float64   `json:"ratio" orm:"column(ratio);digits(12);decimals(2);default(0)"`
		Qty2         float64   `json:"qty2" orm:"column(qty2);digits(12);decimals(2);default(0)"`
		UomId2       int       `json:"uom_id2" orm:"column(uom_id2)"`
		UomCode2     string    `json:"uom_code2" orm:"column(uom_code2)"`
		CreatedAt    time.Time `json:"created_at" orm:"column(created_at);type(timestamp);auto_now_add"`
		UpdatedAt    time.Time `json:"updated_at" orm:"column(updated_at);type(timestamp);auto_now"`
		DeletedAt    time.Time `json:"deleted_at" orm:"column(deleted_at);type(timestamp);null"`
		CreatedBy    string    `json:"created_by" orm:"column(created_by)"`
		UpdatedBy    string    `json:"updated_by" orm:"column(updated_by)"`
		DeletedBy    string    `json:"deleted_by" orm:"column(deleted_by)"`
	}
)

func (t *SalesOrder) TableName() string {
	return "uoms"
}

func SalesOrders() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(SalesOrder))
}

func (t *SalesOrderDetail) TableName() string {
	return "product_divisions"
}

func SalesOrderDetails() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(SalesOrderDetail))
}

func init() {
	orm.RegisterModel(new(SalesOrder), new(SalesOrderDetail))
}

func (t *SalesOrder) Insert(m SalesOrder) (*SalesOrder, error) {
	o := orm.NewOrm()

	if _, err := o.Insert(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (t *SalesOrder) Update(fields ...string) error {
	o := orm.NewOrm()
	if _, err := o.Update(t, fields...); err != nil {
		return err
	}
	return nil
}

func (t *SalesOrderDetail) Insert(m SalesOrderDetail) (*SalesOrderDetail, error) {
	o := orm.NewOrm()

	if _, err := o.Insert(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (t *SalesOrderDetail) Update(fields ...string) error {
	o := orm.NewOrm()
	if _, err := o.Update(t, fields...); err != nil {
		return err
	}
	return nil
}
