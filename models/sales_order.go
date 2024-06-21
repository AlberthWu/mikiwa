package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type (
	SalesOrder struct {
		Id              int       `json:"id" orm:"column(id);auto;pk"`
		ReferenceNo     string    `json:"reference_no" orm:"column(reference_no)"`
		SeqNo           int       `json:"seq_no" orm:"column(seq_no)"`
		IssueDate       time.Time `json:"issue_date" orm:"column(issue_date);type(date)"`
		DueDate         time.Time `json:"due_date" orm:"column(due_date);type(date)"`
		PoolId          int       `json:"pool_id" orm:"column(pool_id)"`
		PoolName        string    `json:"pool_name" orm:"column(pool_name)"`
		CustomerId      int       `json:"customer_id" orm:"column(customer_id)"`
		CustomerName    string    `json:"customer_name" orm:"column(customer_name)"`
		Terms           int       `json:"terms" orm:"column(terms)"`
		DeliveryAddress string    `json:"delivery_address" orm:"column(delivery_address)"`
		EmployeeId      int       `json:"employee_id" orm:"column(employee_id)"`
		EmployeeName    string    `json:"employee_name" orm:"column(employee_name)"`
		LeadTime        int       `json:"lead_time" orm:"column(lead_time)"`
		Subtotal        float64   `json:"subtotal" orm:"column(subtotal);digits(18);decimals(2);default(0)"`
		TotalDisc       float64   `json:"total_disc" orm:"column(total_disc);digits(18);decimals(2);default(0)"`
		Dpp             float64   `json:"dpp" orm:"column(dpp);digits(18);decimals(2);default(0)"`
		Ppn             int       `json:"ppn" orm:"column(ppn)"`
		PpnAmount       float64   `json:"ppn_amount" orm:"column(ppn_amount);digits(18);decimals(2);default(0)"`
		Total           float64   `json:"total" orm:"column(total);digits(18);decimals(2);default(0)"`
		StatusId        int8      `json:"status_id" orm:"column(status_id)"`
		CreatedAt       time.Time `json:"created_at" orm:"column(created_at);type(timestamp);auto_now_add"`
		UpdatedAt       time.Time `json:"updated_at" orm:"column(updated_at);type(timestamp);auto_now"`
		DeletedAt       time.Time `json:"deleted_at" orm:"column(deleted_at);type(timestamp);null"`
		CreatedBy       string    `json:"created_by" orm:"column(created_by)"`
		UpdatedBy       string    `json:"updated_by" orm:"column(updated_by)"`
		DeletedBy       string    `json:"deleted_by" orm:"column(deleted_by)"`
	}

	SalesOrderDetail struct {
		Id                int       `json:"id" orm:"column(id);auto;pk"`
		SalesOrderId      int       `json:"sales_order_id" orm:"column(sales_order_id)"`
		ReferenceNo       string    `json:"reference_no" orm:"column(reference_no)"`
		IssueDate         time.Time `json:"issue_date" orm:"column(issue_date);type(date)"`
		DueDate           time.Time `json:"due_date" orm:"column(due_date);type(date)"`
		ItemNo            int       `json:"item_no" orm:"column(item_no)"`
		ProductId         int       `json:"product_id" orm:"column(product_id)"`
		ProductCode       string    `json:"product_code" orm:"column(product_code)"`
		Qty               float64   `json:"qty" orm:"column(qty);digits(12);decimals(2);default(0)"`
		UomId             int       `json:"uom_id" orm:"column(uom_id)"`
		UomCode           string    `json:"uom_code" orm:"column(uom_code)"`
		Ratio             float64   `json:"ratio" orm:"column(ratio);digits(12);decimals(2);default(0)"`
		PackagingId       int       `json:"packaging_id" orm:"column(packaging_id)"`
		PackagingCode     string    `json:"packaging_code" orm:"column(packaging_code)"`
		FinalQty          float64   `json:"final_qty" orm:"column(final_qty);digits(12);decimals(2);default(0)"`
		FinalUomId        int       `json:"final_uom_id" orm:"column(final_uom_id)"`
		FinalUomCode      string    `json:"final_uom_code" orm:"column(final_uom_code)"`
		NormalPrice       float64   `json:"normal_price" orm:"column(normal_price);digits(18);decimals(2);default(0)"`
		PriceId           int       `json:"price_id" orm:"column(price_id)"`
		Price             float64   `json:"price" orm:"column(price);digits(18);decimals(2);default(0)"`
		Disc1             float64   `json:"disc1" orm:"column(disc1);digits(5);decimals(2);default(0)"`
		Disc2             float64   `json:"disc2" orm:"column(disc2);digits(5);decimals(2);default(0)"`
		DiscTpr           float64   `json:"disc_tpr" orm:"column(disc_tpr);digits(18);decimals(2);default(0)"`
		TotalDisc         float64   `json:"total_disc" orm:"column(total_disc);digits(18);decimals(2);default(0)"`
		NettPrice         float64   `json:"nett_price" orm:"column(nett_price);digits(18);decimals(2);default(0)"`
		Total             float64   `json:"total" orm:"column(total);digits(18);decimals(2);default(0)"`
		LeadTime          int       `json:"lead_time" orm:"column(lead_time)"`
		ConversionQty     float64   `json:"conversion_qty" orm:"column(conversion_qty);digits(12);decimals(2);default(0)"`
		ConversionUomId   int       `json:"conversion_uom_id" orm:"column(conversion_uom_id)"`
		ConversionUomCode string    `json:"conversion_uom_code" orm:"column(conversion_uom_code)"`
		CreatedAt         time.Time `json:"created_at" orm:"column(created_at);type(timestamp);auto_now_add"`
		UpdatedAt         time.Time `json:"updated_at" orm:"column(updated_at);type(timestamp);auto_now"`
		DeletedAt         time.Time `json:"deleted_at" orm:"column(deleted_at);type(timestamp);null"`
		CreatedBy         string    `json:"created_by" orm:"column(created_by)"`
		UpdatedBy         string    `json:"updated_by" orm:"column(updated_by)"`
		DeletedBy         string    `json:"deleted_by" orm:"column(deleted_by)"`
	}
)

func (t *SalesOrder) TableName() string {
	return "sales_order"
}

func SalesOrders() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(SalesOrder))
}

func (t *SalesOrderDetail) TableName() string {
	return "sales_order_detail"
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
