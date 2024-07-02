package models

import (
	"fmt"
	"mikiwa/utils"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type (
	Do struct {
		Id                 int       `json:"id" orm:"column(id);auto;pk"`
		SalesOrderId       int       `json:"sales_order_id" orm:"column(sales_order_id)"`
		SalesOrderNo       string    `json:"sales_order_no" orm:"column(sales_order_no)"`
		ReferenceNo        string    `json:"reference_no" orm:"column(reference_no)"`
		SeqNo              int       `json:"seq_no" orm:"column(seq_no)"`
		WarehouseId        int       `json:"warehouse_id" orm:"column(warehouse_id)"`
		WarehouseCode      string    `json:"warehouse_code" orm:"column(warehouse_code)"`
		WarehousePlantId   int       `json:"warehouse_plant_id" orm:"column(warehouse_plant_id)"`
		WarehousePlantCode string    `json:"warehouse_plant_code" orm:"column(warehouse_plant_code)"`
		IssueDate          time.Time `json:"issue_date" orm:"column(issue_date);type(date)"`
		CustomerId         int       `json:"customer_id" orm:"column(customer_id)"`
		CustomerCode       string    `json:"customer_code" orm:"column(customer_code)"`
		PlantId            int       `json:"plant_id" orm:"column(plant_id)"`
		PlantCode          string    `json:"plant_code" orm:"column(plant_code)"`
		DeliveryAddress    string    `json:"delivery_address" orm:"column(delivery_address)"`
		TransporterId      int       `json:"transporter_id" orm:"column(transporter_id)"`
		TransporterCode    string    `json:"transporter_code" orm:"column(transporter_code)"`
		CourierId          int       `json:"courier_id" orm:"column(courier_id)"`
		CourierName        string    `json:"courier_name" orm:"column(courier_name)"`
		PlateNo            string    `json:"plate_no" orm:"column(plate_no)"`
		Notes              string    `json:"Notes" orm:"column(Notes)"`
		StatusId           int8      `json:"status_id" orm:"column(status_id)"`
		StatusDescription  string    `json:"status_description" orm:"column(status_description)"`
		CreatedAt          time.Time `json:"created_at" orm:"column(created_at);type(timestamp);auto_now_add"`
		UpdatedAt          time.Time `json:"updated_at" orm:"column(updated_at);type(timestamp);auto_now"`
		DeletedAt          time.Time `json:"deleted_at" orm:"column(deleted_at);type(timestamp);null"`
		CreatedBy          string    `json:"created_by" orm:"column(created_by)"`
		UpdatedBy          string    `json:"updated_by" orm:"column(updated_by)"`
		DeletedBy          string    `json:"deleted_by" orm:"column(deleted_by)"`
	}

	DoDetail struct {
		Id                  int       `json:"id" orm:"column(id);auto;pk"`
		SalesOrderId        int       `json:"sales_order_id" orm:"column(sales_order_id)"`
		SalesOrderNo        string    `json:"sales_order_no" orm:"column(sales_order_no)"`
		DoId                int       `json:"do_id" orm:"column(do_id)"`
		ReferenceNo         string    `json:"reference_no" orm:"column(reference_no)"`
		WarehouseId         int       `json:"warehouse_id" orm:"column(warehouse_id)"`
		WarehouseCode       string    `json:"warehouse_code" orm:"column(warehouse_code)"`
		WarehousePlantId    int       `json:"warehouse_plant_id" orm:"column(warehouse_plant_id)"`
		WarehousePlantCode  string    `json:"warehouse_plant_code" orm:"column(warehouse_plant_code)"`
		IssueDate           time.Time `json:"issue_date" orm:"column(issue_date);type(date)"`
		CategoryId          int       `json:"category_id" orm:"column(category_id)"`
		CategoryDescription string    `json:"category_description" orm:"column(category_description)"`
		ItemNo              int       `json:"item_no" orm:"column(item_no)"`
		ProductId           int       `json:"product_id" orm:"column(product_id)"`
		ProductCode         string    `json:"product_code" orm:"column(product_code)"`
		Qty                 float64   `json:"qty" orm:"column(qty);digits(12);decimals(2);default(0)"`
		UomId               int       `json:"uom_id" orm:"column(uom_id)"`
		UomCode             string    `json:"uom_code" orm:"column(uom_code)"`
		Ratio               float64   `json:"ratio" orm:"column(ratio);digits(12);decimals(2);default(0)"`
		PackagingId         int       `json:"packaging_id" orm:"column(packaging_id)"`
		PackagingCode       string    `json:"packaging_code" orm:"column(packaging_code)"`
		FinalQty            float64   `json:"final_qty" orm:"column(final_qty);digits(12);decimals(2);default(0)"`
		FinalUomId          int       `json:"final_uom_id" orm:"column(final_uom_id)"`
		FinalUomCode        string    `json:"final_uom_code" orm:"column(final_uom_code)"`
		Memo                string    `json:"memo" orm:"column(memo)"`
		ConversionQty       float64   `json:"conversion_qty" orm:"column(conversion_qty);digits(12);decimals(2);default(0)"`
		ConversionUomId     int       `json:"conversion_uom_id" orm:"column(conversion_uom_id)"`
		ConversionUomCode   string    `json:"conversion_uom_code" orm:"column(conversion_uom_code)"`
		CreatedAt           time.Time `json:"created_at" orm:"column(created_at);type(timestamp);auto_now_add"`
		UpdatedAt           time.Time `json:"updated_at" orm:"column(updated_at);type(timestamp);auto_now"`
		DeletedAt           time.Time `json:"deleted_at" orm:"column(deleted_at);type(timestamp);null"`
		CreatedBy           string    `json:"created_by" orm:"column(created_by)"`
		UpdatedBy           string    `json:"updated_by" orm:"column(updated_by)"`
		DeletedBy           string    `json:"deleted_by" orm:"column(deleted_by)"`
	}

	DoConfirm struct {
		Id                  int       `json:"id" orm:"column(id);auto;pk"`
		SalesOrderId        int       `json:"sales_order_id" orm:"column(sales_order_id)"`
		SalesOrderNo        string    `json:"sales_order_no" orm:"column(sales_order_no)"`
		DoId                int       `json:"do_id" orm:"column(do_id)"`
		ReferenceNo         string    `json:"reference_no" orm:"column(reference_no)"`
		IssueDate           time.Time `json:"issue_date" orm:"column(issue_date);type(date)"`
		WarehouseId         int       `json:"warehouse_id" orm:"column(warehouse_id)"`
		WarehouseCode       string    `json:"warehouse_code" orm:"column(warehouse_code)"`
		WarehousePlantId    int       `json:"warehouse_plant_id" orm:"column(warehouse_plant_id)"`
		WarehousePlantCode  string    `json:"warehouse_plant_code" orm:"column(warehouse_plant_code)"`
		CategoryId          int       `json:"category_id" orm:"column(category_id)"`
		CategoryDescription string    `json:"category_description" orm:"column(category_description)"`
		ItemNo              int       `json:"item_no" orm:"column(item_no)"`
		ProductId           int       `json:"product_id" orm:"column(product_id)"`
		ProductCode         string    `json:"product_code" orm:"column(product_code)"`
		Qty                 float64   `json:"qty" orm:"column(qty);digits(12);decimals(2);default(0)"`
		UomId               int       `json:"uom_id" orm:"column(uom_id)"`
		UomCode             string    `json:"uom_code" orm:"column(uom_code)"`
		Ratio               float64   `json:"ratio" orm:"column(ratio);digits(12);decimals(2);default(0)"`
		PackagingId         int       `json:"packaging_id" orm:"column(packaging_id)"`
		PackagingCode       string    `json:"packaging_code" orm:"column(packaging_code)"`
		FinalQty            float64   `json:"final_qty" orm:"column(final_qty);digits(12);decimals(2);default(0)"`
		FinalUomId          int       `json:"final_uom_id" orm:"column(final_uom_id)"`
		FinalUomCode        string    `json:"final_uom_code" orm:"column(final_uom_code)"`
		Memo                string    `json:"memo" orm:"column(memo)"`
		ConversionQty       float64   `json:"conversion_qty" orm:"column(conversion_qty);digits(12);decimals(2);default(0)"`
		ConversionUomId     int       `json:"conversion_uom_id" orm:"column(conversion_uom_id)"`
		ConversionUomCode   string    `json:"conversion_uom_code" orm:"column(conversion_uom_code)"`
		CreatedAt           time.Time `json:"created_at" orm:"column(created_at);type(timestamp);auto_now_add"`
		UpdatedAt           time.Time `json:"updated_at" orm:"column(updated_at);type(timestamp);auto_now"`
		DeletedAt           time.Time `json:"deleted_at" orm:"column(deleted_at);type(timestamp);null"`
		CreatedBy           string    `json:"created_by" orm:"column(created_by)"`
		UpdatedBy           string    `json:"updated_by" orm:"column(updated_by)"`
		DeletedBy           string    `json:"deleted_by" orm:"column(deleted_by)"`
	}
)

func (t *Do) TableName() string {
	return "dos"
}

func Dos() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(Do))
}

func (t *DoDetail) TableName() string {
	return "do_detail"
}

func DoDetails() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(DoDetail))
}

func (t *DoConfirm) TableName() string {
	return "do_confirm"
}

func DoConfirms() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(DoConfirm))
}

func init() {
	orm.RegisterModel(new(Do), new(DoDetail), new(DoConfirm))
}

func (t *Do) InsertWithDetail(m Do, d []DoDetail) (data *Do, err error) {
	o := orm.NewOrm()

	if _, err = o.Insert(&m); err != nil {
		return nil, err
	}

	for i := range d {
		d[i].DoId = m.Id
	}

	_, err = o.InsertMulti(len(d), d)
	if err != nil {
		o.Raw("update dos set deleted_at = now(),deleted_by = 'Failed' where id = " + utils.Int2String(m.Id)).Exec()
		return nil, err
	}

	return &m, nil
}

func (t *Do) UpdateWithDetail(m Do, data_post, data_put []DoDetail, user_name string) error {
	o := orm.NewOrm()
	tx, err := o.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	if _, err := tx.Update(&m); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update order: %v", err)
	}

	// Update existing DoDetail (Details) and delete
	var deleteIds []string
	var joinId string
	for _, detail := range data_put {
		if detail.Id != 0 {
			if _, err := tx.Update(&detail); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update detail: %v", err)
			}
		}

		deleteIds = append(deleteIds, utils.Int2String(detail.Id))
	}

	joinId = strings.Join(deleteIds, ",")
	_, err = o.Raw("update do_detail set deleted_at = now(), deleted_by = '" + user_name + "' where deleted_at is null and do_id = " + utils.Int2String(m.Id) + " and id not in (" + joinId + ") ").Exec()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete existing details: %v", err)
	}

	// Insert new DoDetail (Details)
	for i := range data_post {
		data_post[i].DoId = m.Id
	}

	if len(data_post) > 0 {
		_, err = o.InsertMulti(len(data_post), data_post)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert new details: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("failed to rollback transaction: %v after commit failed: %v", rbErr, err)
		}
		return fmt.Errorf("failed to commit transaction: %v", err)
	}
	return nil
}

func (t *DoConfirm) Insert(m DoConfirm) (*DoConfirm, error) {
	o := orm.NewOrm()

	if _, err := o.Insert(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (t *DoConfirm) Update(fields ...string) error {
	o := orm.NewOrm()
	if _, err := o.Update(t, fields...); err != nil {
		return err
	}
	return nil
}

type (
	DoRtn struct {
	}

	DoRtnJson struct {
	}

	DoDetailRtn struct {
	}
)

func (t *Do) GetById(id, user_id int) (m *DoRtn, err error) {
	// o := orm.NewOrm()
	// cond := orm.NewCondition()
	// cond1 := cond.And("deleted_at__isnull", true).And("id", id)
	// qs := Dos().SetCond(cond1)
	// err = qs.One(t)

	// var customer SimpleCompanyRtnJson
	// o.Raw("select id,code,name  from companies where id  =" + utils.Int2String(t.CustomerId) + " ").QueryRow(&customer)

	// var pool PoolRtnJson
	// o.Raw("select id,name,status  from pools where id  = " + utils.Int2String(t.PoolId) + " ").QueryRow(&pool)

	// var outlet SimplePlantRtnJson
	// o.Raw("select t0.id,name,concat(t1.code,' - ',t0.name) full_name,company_id from plants t0 left join (select id,`code` from companies) t1 on t1.id = t0.company_id where t0.id = " + utils.Int2String(t.OutletId) + "").QueryRow(&outlet)

	// var plant SimplePlantRtnJson
	// o.Raw("select t0.id,name,concat(t1.code,' - ',t0.name) full_name,company_id from plants t0 left join (select id,`code` from companies) t1 on t1.id = t0.company_id where t0.id = " + utils.Int2String(t.PlantId) + "").QueryRow(&plant)

	// dlist := t.GetDetail(id, user_id)

	// m = &DoRtn{
	// 	Id:                t.Id,
	// 	ReferenceNo:       t.ReferenceNo,
	// 	IssueDate:         t.IssueDate.Format("2006-01-02"),
	// 	DueDate:           t.DueDate.Format("2006-01-02"),
	// 	LeadTime:          t.LeadTime,
	// 	PoolId:            pool,
	// 	OutletId:          outlet,
	// 	CustomerId:        customer,
	// 	PlantId:           plant,
	// 	Terms:             t.Terms,
	// 	DeliveryAddress:   t.DeliveryAddress,
	// 	EmployeeId:        t.EmployeeId,
	// 	Subtotal:          t.Subtotal,
	// 	TotalDisc:         t.TotalDisc,
	// 	Dpp:               t.Dpp,
	// 	Ppn:               t.Ppn,
	// 	PpnAmount:         t.PpnAmount,
	// 	Total:             t.Total,
	// 	StatusId:          t.StatusId,
	// 	StatusDescription: t.StatusDescription,
	// 	Detail:            dlist,
	// }

	// return m, err
	return nil, err
}
func (t *Do) GetAll(id int)    {}
func (t *Do) GetDetail(id int) {}
