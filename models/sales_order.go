package models

import (
	"fmt"
	"mikiwa/utils"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type (
	SalesOrder struct {
		Id                int       `json:"id" orm:"column(id);auto;pk"`
		ReferenceNo       string    `json:"reference_no" orm:"column(reference_no)"`
		SeqNo             int       `json:"seq_no" orm:"column(seq_no)"`
		IssueDate         time.Time `json:"issue_date" orm:"column(issue_date);type(date)"`
		DueDate           time.Time `json:"due_date" orm:"column(due_date);type(date)"`
		PoolId            int       `json:"pool_id" orm:"column(pool_id)"`
		PoolName          string    `json:"pool_name" orm:"column(pool_name)"`
		OutletId          int       `json:"outlet_id" orm:"column(outlet_id)"`
		OutletName        string    `json:"outlet_name" orm:"column(outlet_name)"`
		CustomerId        int       `json:"customer_id" orm:"column(customer_id)"`
		CustomerCode      string    `json:"customer_code" orm:"column(customer_code)"`
		PlantId           int       `json:"plant_id" orm:"column(plant_id)"`
		PlantName         string    `json:"plant_name" orm:"column(plant_name)"`
		Terms             int       `json:"terms" orm:"column(terms)"`
		DeliveryAddress   string    `json:"delivery_address" orm:"column(delivery_address)"`
		TransporterId     int       `json:"transporter_id" orm:"column(transporter_id)"`
		TransporterCode   string    `json:"transporter_code" orm:"column(transporter_code)"`
		EmployeeId        int       `json:"employee_id" orm:"column(employee_id)"`
		EmployeeName      string    `json:"employee_name" orm:"column(employee_name)"`
		LeadTime          int       `json:"lead_time" orm:"column(lead_time)"`
		Subtotal          float64   `json:"subtotal" orm:"column(subtotal);digits(18);decimals(2);default(0)"`
		TotalDisc         float64   `json:"total_disc" orm:"column(total_disc);digits(18);decimals(2);default(0)"`
		Dpp               float64   `json:"dpp" orm:"column(dpp);digits(18);decimals(2);default(0)"`
		Ppn               int       `json:"ppn" orm:"column(ppn)"`
		PpnAmount         float64   `json:"ppn_amount" orm:"column(ppn_amount);digits(18);decimals(2);default(0)"`
		Total             float64   `json:"total" orm:"column(total);digits(18);decimals(2);default(0)"`
		StatusId          int8      `json:"status_id" orm:"column(status_id)"`
		StatusDescription string    `json:"status_description" orm:"column(status_description)"`
		CreatedAt         time.Time `json:"created_at" orm:"column(created_at);type(timestamp);auto_now_add"`
		UpdatedAt         time.Time `json:"updated_at" orm:"column(updated_at);type(timestamp);auto_now"`
		DeletedAt         time.Time `json:"deleted_at" orm:"column(deleted_at);type(timestamp);null"`
		CreatedBy         string    `json:"created_by" orm:"column(created_by)"`
		UpdatedBy         string    `json:"updated_by" orm:"column(updated_by)"`
		DeletedBy         string    `json:"deleted_by" orm:"column(deleted_by)"`
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
		Disc1Amount       float64   `json:"disc1_amount" orm:"column(disc1_amount);digits(18);decimals(2);default(0)"`
		Disc2             float64   `json:"disc2" orm:"column(disc2);digits(5);decimals(2);default(0)"`
		Disc2Amount       float64   `json:"disc2_amount" orm:"column(disc2_amount);digits(18);decimals(2);default(0)"`
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

func (t *SalesOrder) InsertWithDetail(m SalesOrder, d []SalesOrderDetail) (data *SalesOrder, err error) {
	o := orm.NewOrm()

	if _, err = o.Insert(&m); err != nil {
		return nil, err
	}

	for i := range d {
		d[i].SalesOrderId = m.Id
	}

	_, err = o.InsertMulti(len(d), d)
	if err != nil {
		o.Raw("update sales_order set deleted_at = now(),deleted_by = 'Failed' where id = " + utils.Int2String(m.Id)).Exec()
		return nil, err
	}

	return &m, nil
}

func (t *SalesOrder) UpdateWithDetail(m SalesOrder, data_post, data_put []SalesOrderDetail, user_name string) error {
	o := orm.NewOrm()
	tx, err := o.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	if _, err := tx.Update(&m); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update order: %v", err)
	}

	// Update existing SalesOrderDetail (Details) and delete
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

	if len(deleteIds) == 0 {
		joinId = "0"
	} else {
		joinId = strings.Join(deleteIds, ",")
	}
	_, err = o.Raw("update sales_order_detail set deleted_at = now(), deleted_by = '" + user_name + "' where deleted_at is null and sales_order_id = " + utils.Int2String(m.Id) + " and id not in (" + joinId + ") ").Exec()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete existing details: %v", err)
	}

	// Insert new SalesOrderDetail (Details)
	for i := range data_post {
		data_post[i].SalesOrderId = m.Id
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

func (t *SalesOrder) InsertWithDetailBeginCommit(m SalesOrder, d []SalesOrderDetail) error {
	o := orm.NewOrm()
	tx, err := o.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Insert(&m); err != nil {
		_ = tx.Rollback()
		return err
	}

	for i := range d {
		d[i].SalesOrderId = m.Id
	}

	_, err = tx.InsertMulti(len(d), d)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("failed to rollback transaction: %v after insert details failed: %v", rbErr, err)
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}

type (
	SalesOrderRtn struct {
		Id                int                  `json:"id"`
		ReferenceNo       string               `json:"reference_no"`
		IssueDate         string               `json:"issue_date"`
		DueDate           string               `json:"due_date"`
		LeadTime          int                  `json:"lead_time"`
		PoolId            PoolRtnJson          `json:"pool_id"`
		OutletId          SimplePlantRtnJson   `json:"outlet_id"`
		CustomerId        SimpleCompanyRtnJson `json:"customer_id"`
		PlantId           SimplePlantRtnJson   `json:"plant_id"`
		TransporterId     SimpleCompanyRtnJson `json:"transporter_id"`
		Terms             int                  `json:"terms"`
		DeliveryAddress   string               `json:"delivery_address"`
		EmployeeId        int                  `json:"employee_id"`
		Subtotal          float64              `json:"subtotal"`
		TotalDisc         float64              `json:"total_disc"`
		Dpp               float64              `json:"dpp"`
		Ppn               int                  `json:"ppn" `
		PpnAmount         float64              `json:"ppn_amount"`
		Total             float64              `json:"total"`
		StatusId          int8                 `json:"status_id"`
		StatusDescription string               `json:"status_description"`
		Detail            []orm.Params         `json:"detail"`
	}

	SalesOrderRtnJson struct {
		Id                int     `json:"id"`
		ReferenceNo       string  `json:"reference_no"`
		IssueDate         string  `json:"issue_date"`
		DueDate           string  `json:"due_date"`
		LeadTime          int     `json:"lead_time"`
		PoolId            int     `json:"pool_id"`
		PoolName          string  `json:"pool_name"`
		OutletId          int     `json:"outlet_id"`
		OutletName        string  `json:"outlet_name"`
		CustomerId        int     `json:"customer_id"`
		CustomerCode      int     `json:"customer_code"`
		CustomerName      int     `json:"customer_name"`
		PlantId           int     `json:"plant_id"`
		PlantName         string  `json:"plant_name"`
		FullName          string  `json:"full_name"`
		Terms             int     `json:"terms"`
		DeliveryAddress   string  `json:"delivery_address"`
		EmployeeId        int     `json:"employee_id"`
		EmployeeName      string  `json:"employee_name"`
		TransporterId     int     `json:"transporter_id"`
		TransporterCode   int     `json:"transporter_code"`
		TransporterName   int     `json:"transporter_name"`
		Subtotal          float64 `json:"subtotal"`
		TotalDisc         float64 `json:"total_disc"`
		Dpp               float64 `json:"dpp"`
		Ppn               int     `json:"ppn" `
		PpnAmount         float64 `json:"ppn_amount"`
		Total             float64 `json:"total"`
		StatusId          int8    `json:"status_id"`
		StatusDescription string  `json:"status_description"`
		StatusData        string  `json:"status_data"`
	}

	SalesOrderDetailRtnJson struct {
		Id                int     `json:"id"`
		SalesOrderId      int     `json:"sales_order_id"`
		ItemNo            int     `json:"item_no"`
		ProductId         int     `json:"product_id"`
		ProductCode       string  `json:"product_code"`
		ProductName       string  `json:"product_name"`
		Qty               float64 `json:"qty"`
		UomId             int     `json:"uom_id"`
		UomCode           string  `json:"uom_code"`
		Ratio             float64 `json:"ratio"`
		PackagingId       int     `json:"packaging_id"`
		PackagingCode     string  `json:"packaging_code"`
		FinalQty          float64 `json:"final_qty"`
		FinalUomId        int     `json:"final_uom_id"`
		FinalUomCode      string  `json:"final_uom_code"`
		ConversionQty     float64 `json:"conversion_qty"`
		ConversionUomId   int     `json:"convertsion_uom_id"`
		ConversionUomCode string  `json:"convertsion_uom_code"`
		NormalPrice       float64 `json:"normal_price"`
		PriceId           int     `json:"price_id"`
		Price             float64 `json:"price"`
		StatusData        string  `json:"status_data"`
	}
)

func (t *SalesOrder) GetById(id, user_id int) (m *SalesOrderRtn, err error) {
	o := orm.NewOrm()
	cond := orm.NewCondition()
	cond1 := cond.And("deleted_at__isnull", true).And("id", id)
	qs := SalesOrders().SetCond(cond1)
	err = qs.One(t)

	var customer SimpleCompanyRtnJson
	o.Raw("select id,code,name  from companies where id  =" + utils.Int2String(t.CustomerId) + " ").QueryRow(&customer)

	var pool PoolRtnJson
	o.Raw("select id,name,status  from pools where id  = " + utils.Int2String(t.PoolId) + " ").QueryRow(&pool)

	var outlet SimplePlantRtnJson
	o.Raw("select t0.id,t0.code,name,concat(t1.code,' - ',t0.name) full_name,company_id,t1.code company_code,status from plants t0 left join (select id,`code` from companies) t1 on t1.id = t0.company_id where t0.id = " + utils.Int2String(t.OutletId) + "").QueryRow(&outlet)

	var plant SimplePlantRtnJson
	o.Raw("select t0.id,t0.code,name,concat(t1.code,' - ',t0.name) full_name,company_id,t1.code company_code,status from plants t0 left join (select id,`code` from companies) t1 on t1.id = t0.company_id where t0.id = " + utils.Int2String(t.PlantId) + "").QueryRow(&plant)

	var transporter SimpleCompanyRtnJson
	o.Raw("select id,code,name  from companies where id  =" + utils.Int2String(t.TransporterId) + " ").QueryRow(&transporter)

	dlist := t.GetDetail(id, user_id)

	m = &SalesOrderRtn{
		Id:                t.Id,
		ReferenceNo:       t.ReferenceNo,
		IssueDate:         t.IssueDate.Format("2006-01-02"),
		DueDate:           t.DueDate.Format("2006-01-02"),
		LeadTime:          t.LeadTime,
		PoolId:            pool,
		OutletId:          outlet,
		CustomerId:        customer,
		PlantId:           plant,
		Terms:             t.Terms,
		DeliveryAddress:   t.DeliveryAddress,
		EmployeeId:        t.EmployeeId,
		Subtotal:          t.Subtotal,
		TotalDisc:         t.TotalDisc,
		Dpp:               t.Dpp,
		Ppn:               t.Ppn,
		PpnAmount:         t.PpnAmount,
		Total:             t.Total,
		StatusId:          t.StatusId,
		StatusDescription: t.StatusDescription,
		TransporterId:     transporter,
		Detail:            dlist,
	}

	return m, err
}

func (t *SalesOrder) GetAll(keyword, field_name, match_mode, value_name string, p, size, allsize, user_id, id, search_detail int, plant_id int, employee_ids, outlet_ids, customer_ids, status_ids, product_ids, lead_time_ids string, issue_date, due_date, updated_at *string) (u utils.PageDynamic, err error) {
	o := orm.NewOrm()
	var m []orm.Params
	var c int

	//  sp_SalesOrder(theDate date,dueDate date,updatedAt date,uId int,employeeIds varchar(15),outletIds varchar(15),customerIds varchar(15),plantId int,productIds varchar(15),statusIds varchar(15),in leadTime varchar(5), reportTypeId int,userId int,keyword varchar(255),in searchDetail int,in TheField varchar(8000),in MatchMode varchar(8000),in ValueName varchar(8000), in limitVal int, in offsetVal int )
	o.Raw("call sp_SalesOrderCount(?,?,?,"+utils.Int2String(id)+",'"+employee_ids+"','"+outlet_ids+"','"+customer_ids+"',"+utils.Int2String(plant_id)+",'"+product_ids+"','"+status_ids+"','"+lead_time_ids+"',1,"+utils.Int2String(user_id)+",'"+keyword+"',"+utils.Int2String(search_detail)+",'"+field_name+"','"+match_mode+"','"+value_name+"',null,null)", &issue_date, &due_date, &updated_at).QueryRow(&c)

	if allsize == 1 && c > 0 {
		size = c
	}
	_, err = o.Raw("call sp_SalesOrder(?,?,?,"+utils.Int2String(id)+",'"+employee_ids+"','"+outlet_ids+"','"+customer_ids+"',"+utils.Int2String(plant_id)+",'"+product_ids+"','"+status_ids+"','"+lead_time_ids+"',1,"+utils.Int2String(user_id)+",'"+keyword+"',"+utils.Int2String(search_detail)+",'"+field_name+"','"+match_mode+"','"+value_name+"',"+utils.Int2String(size)+", "+utils.Int2String((p-1)*size)+")", &issue_date, &due_date, &updated_at).Values(&m)

	if c == 0 && err == nil {
		err = orm.ErrNoRows
		return utils.PaginationDynamic(int(c), p, size, "", "", "", "", "", "", "", m), err
	} else if err != nil {
		return utils.PaginationDynamic(int(c), p, size, "", "", "", "", "", "", "", m), err
	}
	return utils.PaginationDynamic(int(c), p, size, fmt.Sprintf("%v", m[0]["field_key"]), fmt.Sprintf("%v", m[0]["field_label"]), fmt.Sprintf("%v", m[0]["field_int"]), fmt.Sprintf("%v", m[0]["field_level"]), fmt.Sprintf("%v", m[0]["field_export"]), fmt.Sprintf("%v", m[0]["field_export_label"]), fmt.Sprintf("%v", m[0]["field_footer"]), m), err
}

func (t *SalesOrder) GetAllDetail(keyword, field_name, match_mode, value_name string, p, size, allsize, user_id, id, search_detail int, plant_id int, employee_ids, outlet_ids, customer_ids, status_ids, product_ids, lead_time_ids string, issue_date, due_date, updated_at *string) (m []orm.Params, err error) {
	o := orm.NewOrm()
	var c int64
	c, err = o.Raw("call sp_SalesOrder(?,?,?,"+utils.Int2String(id)+",'"+employee_ids+"','"+outlet_ids+"','"+customer_ids+"',"+utils.Int2String(plant_id)+",'"+product_ids+"','"+status_ids+"','"+lead_time_ids+"',0,"+utils.Int2String(user_id)+",'"+keyword+"',"+utils.Int2String(search_detail)+",'"+field_name+"','"+match_mode+"','"+value_name+"',null,null)", &issue_date, &due_date, &updated_at).Values(&m)
	if c == 0 && err == nil {
		err = orm.ErrNoRows
	}
	return m, err
}

func (c *SalesOrder) GetDetail(id, user_id int) (m []orm.Params) {
	o := orm.NewOrm()
	// theDate date,dueDate date,updatedAt date,uId int,employeeIds varchar(15),outletIds varchar(15),customerIds varchar(15),plantId int,productIds varchar(15),statusIds varchar(15),in leadTime varchar(5), reportTypeId int,userId int,keyword varchar(255),in searchDetail int,in TheField varchar(8000),in MatchMode varchar(8000),in ValueName varchar(8000), in limitVal int, in offsetVal int
	o.Raw("call sp_SalesOrder(null,null,null," + utils.Int2String(id) + ",null,null,null,0,null,null,null,0," + utils.Int2String(user_id) + ",'',null,null,null,null,null,null)").Values(&m)
	return m
}
