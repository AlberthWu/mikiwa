package controllers

import (
	"encoding/json"
	"fmt"
	base "mikiwa/controllers"
	"mikiwa/models"
	"mikiwa/utils"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
)

type DoController struct {
	base.BaseController
}

func (c *DoController) Prepare() {
	c.Ctx.Request.Header.Set("token", "No Aut")
	c.BaseController.Prepare()
}

type (
	InputHeaderDo struct {
		SalesOrderId     int                   `json:"sales_order_id"`
		IssueDate        string                `json:"issue_date"`
		OutletId         int                   `json:"outlet_id"`
		WarehousePlantId int                   `json:"warehouse_plant_id"`
		CustomerId       int                   `json:"customer_id"`
		PlantId          int                   `json:"plant_id"`
		TransporterId    int                   `json:"transporter_id"`
		CourierId        int                   `json:"courier_id"`
		DeliveryAddress  string                `json:"delivery_address"`
		PlateNo          string                `json:"plate_no"`
		Notes            string                `json:"notes"`
		UploadFile       models.DocumentList   `json:"upload_file"`
		Detail           []InputDetailDoDetail `json:"detail"`
	}

	InputDetailDoDetail struct {
		Id                  int     `json:"id"`
		CategoryId          int     `json:"category_id"`
		CategoryDescription string  `json:"category_description"`
		ProductId           int     `json:"product_id"`
		Qty                 float64 `json:"qty"`
		UomId               int     `json:"uom_id"`
		Memo                string  `json:"memo"`
		CreatedBy           string  `json:"created_by"`
	}
)

func (c *DoController) Post() {
	o := orm.NewOrm()
	var user_id, form_id int
	var user_name string
	var err error
	var folderName string = "delivery_order"
	var so_reference_no string
	var status_id int8 = base.OpenDo
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
		user_name = sess.(map[string]interface{})["username"].(string)
	}

	form_id = base.FormName(form_do)
	write_aut := models.CheckPrivileges(user_id, form_id, base.Write)
	write_aut = true
	if !write_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Post not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}

	var i int = 0
	var ob InputHeaderDo
	var inputDetail []models.DoDetail

	body := c.Ctx.Input.RequestBody
	err = json.Unmarshal(body, &ob)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	var salesorder models.SalesOrder
	if ob.SalesOrderId != 0 {
		err = models.SalesOrders().One(&salesorder)
		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, err.Error())
			c.ServeJSON()
			return
		}

		if salesorder.StatusId == base.OpenSo {
			c.Ctx.ResponseWriter.WriteHeader(402)
			utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("'%v' not CONFIRMED yet", salesorder.ReferenceNo))
			c.ServeJSON()
			return
		}

		if salesorder.StatusId == base.DoneSo {
			c.Ctx.ResponseWriter.WriteHeader(402)
			utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("'%v' has been DONE", salesorder.ReferenceNo))
			c.ServeJSON()
			return
		}
		so_reference_no = salesorder.ReferenceNo
	}

	ob.CourierId = 1
	valid := validation.Validation{}
	valid.Required(strings.TrimSpace(ob.IssueDate), "issue_date").Message("Is required")
	valid.Required(ob.WarehousePlantId, "warehouse_plant_id").Message("Is required")
	valid.Required(ob.CustomerId, "customer_id").Message("Is required")
	valid.Required(ob.TransporterId, "transporter_id").Message("Is required")
	valid.Required(ob.PlateNo, "plate_no").Message("Is required")
	valid.Required(strings.TrimSpace(ob.DeliveryAddress), "delivery_address").Message("Is required")

	if len(ob.Detail) == 0 {
		valid.AddError("detail", "Detail list is required")
	}

	if valid.HasErrors() {
		out := make([]utils.ApiError, len(valid.Errors))
		for i, err := range valid.Errors {
			out[i] = utils.ApiError{Param: err.Key, Message: err.Message}
		}
		c.Ctx.ResponseWriter.WriteHeader(400)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 400, "Invalid input field", out)
		c.ServeJSON()
		return
	}

	// check customer
	var customers models.Company
	var plants models.Plant
	err = models.Companies().Filter("id", ob.CustomerId).Filter("deleted_at__isnull", true).Filter("CompanyTypes__TypeId__Id", base.Customer).One(&customers)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Customer unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	if customers.Status == 0 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been set as INACTIVE", customers.Code))
		c.ServeJSON()
		return
	}

	// check customer plant
	if ob.PlantId != 0 {
		err = models.Plants().Filter("deleted_at__isnull", true).Filter("id", ob.PlantId).Filter("company_id", ob.CustomerId).One(&plants)
		if err == orm.ErrNoRows {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, "Plant unregistered/Illegal data")
			c.ServeJSON()
			return
		}

		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, err.Error())
			c.ServeJSON()
			return
		}

		if plants.Status == 0 {
			c.Ctx.ResponseWriter.WriteHeader(402)
			utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been set as INACTIVE", plants.Name))
			c.ServeJSON()
			return
		}
	}

	// check warehouse
	var warehouse models.SimplePlantRtnJson
	if err = o.Raw("select t0.id,t0.code,name,concat(t1.code,' - ',t0.name) full_name,company_id,t1.code company_code,status from plants t0 left join (select id,`code` from companies) t1 on t1.id = t0.company_id where t0.id = " + utils.Int2String(ob.WarehousePlantId) + "").QueryRow(&warehouse); err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Warehouse unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	if warehouse.Status == 0 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been set as INACTIVE", warehouse.Code))
		c.ServeJSON()
		return
	}

	// check transporter
	var transporter models.Company
	err = models.Companies().Filter("id", ob.TransporterId).Filter("deleted_at__isnull", true).Filter("CompanyTypes__TypeId__Id", base.Transporter).One(&transporter)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Logistic unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	if transporter.Status == 0 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been set as INACTIVE", transporter.Code))
		c.ServeJSON()
		return
	}

	// check outlet
	var outlet models.Company
	if err = o.Raw("select * from plants where deleted_at is null and id = " + utils.Int2String(ob.OutletId) + " and company_id in (select company_id from company_type where type_id = " + utils.Int2String(base.Internal) + " )").QueryRow(&outlet); err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Outlet unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	if outlet.Status == 0 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been set as INACTIVE", outlet.Name))
		c.ServeJSON()
		return
	}

	// check courier

	var products models.Product
	var productUom models.ProductUom
	var priceRtn *models.ProductConversionRtnJson
	wg = new(sync.WaitGroup)
	var mutex sync.Mutex
	resultChan := make(chan utils.ResultChan, len(ob.Detail))
	var queryResults []utils.ResultChan
	wg.Add(len(ob.Detail))
	for _, v := range ob.Detail {
		go func(v InputDetailDoDetail) {
			defer wg.Done()
			mutex.Lock()
			defer mutex.Unlock()
			if err = o.Raw("select * from products where deleted_at is null and product_type_id = " + utils.Int2String(base.ProductFinishing) + " and product_division_id in (select business_unit_id from company_business_unit where company_id = " + utils.Int2String(ob.CustomerId) + ") and id = " + utils.Int2String(v.ProductId)).QueryRow(&products); err == orm.ErrNoRows {
				resultChan <- utils.ResultChan{Id: v.ProductId, Data: "Invalid product", Message: "product unregistered/Illegal data"}
				return
			}

			if err != nil {
				resultChan <- utils.ResultChan{Id: v.ProductId, Data: products.ProductCode, Message: err.Error()}
				return
			}

			if products.StatusId == 0 {
				resultChan <- utils.ResultChan{Id: v.ProductId, Data: products.ProductCode, Message: fmt.Sprintf("'%v' has been set as inactive", products.ProductCode)}
				return
			}

			err = o.Raw("select * from product_uom where product_id = " + utils.Int2String(v.ProductId) + " and uom_id = " + utils.Int2String(v.UomId)).QueryRow(&productUom)
			if err == orm.ErrNoRows {
				resultChan <- utils.ResultChan{Id: v.ProductId, Data: products.ProductCode, Message: "product uom unregistered/Illegal data"}
				return
			}

			if err != nil {
				resultChan <- utils.ResultChan{Id: v.ProductId, Data: products.ProductCode, Message: err.Error()}
				return
			}
		}(v)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		if result.Message != "" {
			queryResults = append(queryResults, utils.ResultChan{
				Id:      result.Id,
				Data:    result.Data,
				Message: result.Message,
			})
		}
	}

	if len(queryResults) != 0 {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 401, "Error", map[string]interface{}{"Invalid field": queryResults})
		c.ServeJSON()
		return
	}

	thedate, errDate := time.Parse("2006-01-02", ob.IssueDate)
	if errDate != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprint("issue_date: ", errDate.Error()))
		c.ServeJSON()
		return
	}

	seqno, referenceno := models.GenerateBatchNumber(thedate, warehouse.CompanyId, ob.OutletId, ob.CustomerId, "DeliveryOrder")

	t_delivery_order = models.Do{
		SalesOrderId:       ob.SalesOrderId,
		SalesOrderNo:       so_reference_no,
		IssueDate:          thedate,
		ReferenceNo:        referenceno,
		SeqNo:              seqno,
		WarehouseId:        warehouse.CompanyId,
		WarehouseCode:      warehouse.CompanyCode,
		WarehousePlantId:   ob.WarehousePlantId,
		WarehousePlantCode: warehouse.Code,
		CustomerId:         ob.CustomerId,
		CustomerCode:       customers.Code,
		PlantId:            ob.PlantId,
		PlantCode:          plants.Code,
		DeliveryAddress:    ob.DeliveryAddress,
		TransporterId:      ob.TransporterId,
		TransporterCode:    transporter.Code,
		CourierId:          ob.CourierId,
		PlateNo:            ob.PlateNo,
		Notes:              ob.Notes,
		StatusId:           status_id,
		StatusDescription:  base.GetStatusDo(int(status_id)),
		CreatedBy:          user_name,
		UpdatedBy:          user_name,
	}

	for k, v := range ob.Detail {
		i = 0
		wg.Add(1)
		go func(k int, v InputDetailDoDetail) {
			priceRtn = products.GetConversion(ob.IssueDate, v.Qty, ob.CustomerId, v.ProductId, v.UomId, user_id)

			defer wg.Done()
			mutex.Lock()
			if v.Id == 0 {
				inputDetail = append(inputDetail, models.DoDetail{
					SalesOrderId:        ob.SalesOrderId,
					SalesOrderNo:        so_reference_no,
					ReferenceNo:         referenceno,
					IssueDate:           thedate,
					WarehouseId:         warehouse.CompanyId,
					WarehouseCode:       warehouse.Code,
					WarehousePlantId:    ob.WarehousePlantId,
					WarehousePlantCode:  warehouse.Code,
					CategoryId:          v.CategoryId,
					CategoryDescription: v.CategoryDescription,
					ItemNo:              k + 1,
					ProductId:           v.ProductId,
					ProductCode:         priceRtn.ProductCode,
					Qty:                 v.Qty,
					UomId:               v.UomId,
					UomCode:             priceRtn.UomCode,
					Ratio:               priceRtn.Ratio,
					PackagingId:         priceRtn.PackagingId,
					PackagingCode:       priceRtn.PackagingCode,
					FinalQty:            priceRtn.FinalQty,
					FinalUomId:          priceRtn.FinalUomId,
					FinalUomCode:        priceRtn.FinalUomCode,
					Memo:                v.Memo,
					ConversionQty:       priceRtn.ConversionQty,
					ConversionUomId:     priceRtn.ConversionUomId,
					ConversionUomCode:   priceRtn.ConversionUomCode,
					CreatedBy:           user_name,
					UpdatedBy:           user_name,
				})
				i += 1
			}
			mutex.Unlock()
		}(k, v)
	}
	wg.Wait()

	d, err_ := t_delivery_order.InsertWithDetail(t_delivery_order, inputDetail)
	errcode, errmessage = base.DecodeErr(err_)
	if err_ != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
		c.ServeJSON()
		return
	} else {
		if err = base.PostFirebaseRaw(ob.UploadFile, user_name, d.Id, folderName+"/"+utils.Int2String(d.Id), folderName+"/"+utils.Int2String(d.Id)); err != nil {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprint("Error processing data and uploading to Firebase: ", err.Error()))
		} else {
			v, err := t_delivery_order.GetById(d.Id, user_id)
			errcode, errmessage = base.DecodeErr(err)
			if err != nil {
				c.Ctx.ResponseWriter.WriteHeader(errcode)
				utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
			} else {
				utils.ReturnHTTPSuccessWithMessage(&c.Controller, errcode, errmessage, v)
			}
		}
	}

	c.ServeJSON()
}

func (c *DoController) Put() {
	o := orm.NewOrm()
	var user_id, form_id int
	var user_name string
	var err error
	var folderName string = "delivery_order"
	var so_reference_no string
	var status_id int8
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
		user_name = sess.(map[string]interface{})["username"].(string)
	}

	form_id = base.FormName(form_do)
	put_aut := models.CheckPrivileges(user_id, form_id, base.Write)
	put_aut = true
	if !put_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Put not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}

	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	var i int = 0
	var ob InputHeaderDo
	var inputDetail []models.DoDetail
	var putDetail []models.DoDetail

	body := c.Ctx.Input.RequestBody
	err = json.Unmarshal(body, &ob)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	var querydata models.Do
	err = models.Dos().Filter("id", id).Filter("deleted_at__isnull", true).One(&querydata)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, "Delivery order id unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	issuedate, errDate := time.Parse("2006-01-02", ob.IssueDate)
	if errDate != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, errDate.Error())
		c.ServeJSON()
		return

	}

	if issuedate.Year() != querydata.IssueDate.Year() {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Allowed changes date part or month part only")
		c.ServeJSON()
		return
	}

	var salesorder models.SalesOrder
	if ob.SalesOrderId != 0 {
		err = models.SalesOrders().One(&salesorder)
		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, err.Error())
			c.ServeJSON()
			return
		}

		if salesorder.StatusId == base.OpenSo {
			c.Ctx.ResponseWriter.WriteHeader(402)
			utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("'%v' not CONFIRMED yet", salesorder.ReferenceNo))
			c.ServeJSON()
			return
		}

		if salesorder.StatusId == base.DoneSo {
			c.Ctx.ResponseWriter.WriteHeader(402)
			utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("'%v' has been DONE", salesorder.ReferenceNo))
			c.ServeJSON()
			return
		}
		so_reference_no = salesorder.ReferenceNo
	}

	ob.CourierId = 1
	status_id = querydata.StatusId
	valid := validation.Validation{}
	valid.Required(strings.TrimSpace(ob.IssueDate), "issue_date").Message("Is required")
	valid.Required(ob.WarehousePlantId, "warehouse_plant_id").Message("Is required")
	valid.Required(ob.CustomerId, "customer_id").Message("Is required")
	valid.Required(ob.TransporterId, "transporter_id").Message("Is required")
	valid.Required(ob.PlateNo, "plate_no").Message("Is required")
	valid.Required(strings.TrimSpace(ob.DeliveryAddress), "delivery_address").Message("Is required")

	if len(ob.Detail) == 0 {
		valid.AddError("detail", "Detail list is required")
	}

	if valid.HasErrors() {
		out := make([]utils.ApiError, len(valid.Errors))
		for i, err := range valid.Errors {
			out[i] = utils.ApiError{Param: err.Key, Message: err.Message}
		}
		c.Ctx.ResponseWriter.WriteHeader(400)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 400, "Invalid input field", out)
		c.ServeJSON()
		return
	}

	// check customer
	var customers models.Company
	var plants models.Plant
	err = models.Companies().Filter("id", ob.CustomerId).Filter("deleted_at__isnull", true).Filter("CompanyTypes__TypeId__Id", base.Customer).One(&customers)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Customer unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	if customers.Status == 0 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been set as INACTIVE", customers.Code))
		c.ServeJSON()
		return
	}

	// check customer plant
	if ob.PlantId != 0 {
		err = models.Plants().Filter("deleted_at__isnull", true).Filter("id", ob.PlantId).Filter("company_id", ob.CustomerId).One(&plants)
		if err == orm.ErrNoRows {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, "Plant unregistered/Illegal data")
			c.ServeJSON()
			return
		}

		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, err.Error())
			c.ServeJSON()
			return
		}

		if plants.Status == 0 {
			c.Ctx.ResponseWriter.WriteHeader(402)
			utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been set as INACTIVE", plants.Name))
			c.ServeJSON()
			return
		}
	}

	// check warehouse
	var warehouse models.SimplePlantRtnJson
	if err = o.Raw("select t0.id,t0.code,name,concat(t1.code,' - ',t0.name) full_name,company_id,t1.code company_code,status from plants t0 left join (select id,`code` from companies) t1 on t1.id = t0.company_id where t0.id = " + utils.Int2String(ob.WarehousePlantId) + "").QueryRow(&warehouse); err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Warehouse unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	if warehouse.Status == 0 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been set as INACTIVE", warehouse.Code))
		c.ServeJSON()
		return
	}

	// check transporter
	var transporter models.Company
	err = models.Companies().Filter("id", ob.TransporterId).Filter("deleted_at__isnull", true).Filter("CompanyTypes__TypeId__Id", base.Transporter).One(&transporter)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Logistic unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	if transporter.Status == 0 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been set as INACTIVE", transporter.Code))
		c.ServeJSON()
		return
	}

	// check outlet
	var outlet models.Company
	if err = o.Raw("select * from plants where deleted_at is null and id = " + utils.Int2String(ob.OutletId) + " and company_id in (select company_id from company_type where type_id = " + utils.Int2String(base.Internal) + " )").QueryRow(&outlet); err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Outlet unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	if outlet.Status == 0 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been set as INACTIVE", outlet.Name))
		c.ServeJSON()
		return
	}

	// check courier

	var products models.Product
	var productUom models.ProductUom
	var priceRtn *models.ProductConversionRtnJson
	var querydetail models.DoDetail
	wg = new(sync.WaitGroup)
	var mutex sync.Mutex
	resultChan := make(chan utils.ResultChan, len(ob.Detail))
	var queryResults []utils.ResultChan
	wg.Add(len(ob.Detail))
	for _, v := range ob.Detail {
		go func(v InputDetailDoDetail) {
			defer wg.Done()
			mutex.Lock()
			defer mutex.Unlock()
			if v.Id != 0 {
				if err = models.DoDetails().Filter("deleted_at__isnull", true).Filter("do_id", id).Filter("id", v.Id).One(&querydetail); err == orm.ErrNoRows {
					resultChan <- utils.ResultChan{Id: v.ProductId, Data: "Invalid detail id", Message: "detail unregistered/Illegal data"}
					return
				}

				if err != nil {
					resultChan <- utils.ResultChan{Id: v.ProductId, Data: products.ProductCode, Message: err.Error()}
					return
				}
			}

			if err = o.Raw("select * from products where deleted_at is null and product_type_id = " + utils.Int2String(base.ProductFinishing) + " and product_division_id in (select business_unit_id from company_business_unit where company_id = " + utils.Int2String(ob.CustomerId) + ") and id = " + utils.Int2String(v.ProductId)).QueryRow(&products); err == orm.ErrNoRows {
				resultChan <- utils.ResultChan{Id: v.ProductId, Data: "Invalid product", Message: "product unregistered/Illegal data"}
				return
			}

			if err != nil {
				resultChan <- utils.ResultChan{Id: v.ProductId, Data: products.ProductCode, Message: err.Error()}
				return
			}

			if products.StatusId == 0 {
				resultChan <- utils.ResultChan{Id: v.ProductId, Data: products.ProductCode, Message: fmt.Sprintf("'%v' has been set as inactive", products.ProductCode)}
				return
			}

			err = o.Raw("select * from product_uom where product_id = " + utils.Int2String(v.ProductId) + " and uom_id = " + utils.Int2String(v.UomId)).QueryRow(&productUom)
			if err == orm.ErrNoRows {
				resultChan <- utils.ResultChan{Id: v.ProductId, Data: products.ProductCode, Message: "product uom unregistered/Illegal data"}
				return
			}

			if err != nil {
				resultChan <- utils.ResultChan{Id: v.ProductId, Data: products.ProductCode, Message: err.Error()}
				return
			}
		}(v)
		v.CreatedBy = querydetail.CreatedBy
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		if result.Message != "" {
			queryResults = append(queryResults, utils.ResultChan{
				Id:      result.Id,
				Data:    result.Data,
				Message: result.Message,
			})
		}
	}

	if len(queryResults) != 0 {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 401, "Error", map[string]interface{}{"Invalid field": queryResults})
		c.ServeJSON()
		return
	}

	thedate, errDate := time.Parse("2006-01-02", ob.IssueDate)
	if errDate != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprint("issue_date: ", errDate.Error()))
		c.ServeJSON()
		return
	}

	seqno := querydata.SeqNo
	referenceno := querydata.ReferenceNo

	t_delivery_order.Id = id
	t_delivery_order.SalesOrderId = ob.SalesOrderId
	t_delivery_order.SalesOrderNo = so_reference_no
	t_delivery_order.IssueDate = thedate
	t_delivery_order.ReferenceNo = referenceno
	t_delivery_order.SeqNo = seqno
	t_delivery_order.WarehouseId = warehouse.CompanyId
	t_delivery_order.WarehouseCode = warehouse.CompanyCode
	t_delivery_order.WarehousePlantId = ob.WarehousePlantId
	t_delivery_order.WarehousePlantCode = warehouse.Code
	t_delivery_order.CustomerId = ob.CustomerId
	t_delivery_order.CustomerCode = customers.Code
	t_delivery_order.PlantId = ob.PlantId
	t_delivery_order.PlantCode = plants.Code
	t_delivery_order.DeliveryAddress = ob.DeliveryAddress
	t_delivery_order.TransporterId = ob.TransporterId
	t_delivery_order.TransporterCode = transporter.Code
	t_delivery_order.CourierId = ob.CourierId
	t_delivery_order.PlateNo = ob.PlateNo
	t_delivery_order.Notes = ob.Notes
	t_delivery_order.StatusId = status_id
	t_delivery_order.StatusDescription = base.GetStatusDo(int(status_id))
	t_delivery_order.CreatedBy = querydata.CreatedBy
	t_delivery_order.UpdatedBy = user_name

	for k, v := range ob.Detail {
		i = 0
		wg.Add(1)
		go func(k int, v InputDetailDoDetail) {
			priceRtn = products.GetConversion(ob.IssueDate, v.Qty, ob.CustomerId, v.ProductId, v.UomId, user_id)

			defer wg.Done()
			mutex.Lock()
			if v.Id == 0 {
				inputDetail = append(inputDetail, models.DoDetail{
					SalesOrderId:        ob.SalesOrderId,
					SalesOrderNo:        so_reference_no,
					ReferenceNo:         referenceno,
					IssueDate:           thedate,
					WarehouseId:         warehouse.CompanyId,
					WarehouseCode:       warehouse.CompanyCode,
					WarehousePlantId:    ob.WarehousePlantId,
					WarehousePlantCode:  warehouse.Code,
					CategoryId:          v.CategoryId,
					CategoryDescription: v.CategoryDescription,
					ItemNo:              k + 1,
					ProductId:           v.ProductId,
					ProductCode:         priceRtn.ProductCode,
					Qty:                 v.Qty,
					UomId:               v.UomId,
					UomCode:             priceRtn.UomCode,
					Ratio:               priceRtn.Ratio,
					PackagingId:         priceRtn.PackagingId,
					PackagingCode:       priceRtn.PackagingCode,
					FinalQty:            priceRtn.FinalQty,
					FinalUomId:          priceRtn.FinalUomId,
					FinalUomCode:        priceRtn.FinalUomCode,
					Memo:                v.Memo,
					ConversionQty:       priceRtn.ConversionQty,
					ConversionUomId:     priceRtn.ConversionUomId,
					ConversionUomCode:   priceRtn.ConversionUomCode,
					CreatedBy:           user_name,
					UpdatedBy:           user_name,
				})
				i += 1
			} else {
				putDetail = append(putDetail, models.DoDetail{
					Id:                  v.Id,
					SalesOrderId:        ob.SalesOrderId,
					SalesOrderNo:        so_reference_no,
					ReferenceNo:         referenceno,
					IssueDate:           thedate,
					WarehouseId:         warehouse.CompanyId,
					WarehouseCode:       warehouse.CompanyCode,
					WarehousePlantId:    ob.WarehousePlantId,
					WarehousePlantCode:  warehouse.Code,
					CategoryId:          v.CategoryId,
					CategoryDescription: v.CategoryDescription,
					ItemNo:              k + 1,
					ProductId:           v.ProductId,
					ProductCode:         priceRtn.ProductCode,
					Qty:                 v.Qty,
					UomId:               v.UomId,
					UomCode:             priceRtn.UomCode,
					Ratio:               priceRtn.Ratio,
					PackagingId:         priceRtn.PackagingId,
					PackagingCode:       priceRtn.PackagingCode,
					FinalQty:            priceRtn.FinalQty,
					FinalUomId:          priceRtn.FinalUomId,
					FinalUomCode:        priceRtn.FinalUomCode,
					Memo:                v.Memo,
					ConversionQty:       priceRtn.ConversionQty,
					ConversionUomId:     priceRtn.ConversionUomId,
					ConversionUomCode:   priceRtn.ConversionUomCode,
					CreatedBy:           v.CreatedBy,
					UpdatedBy:           user_name,
				})
			}
			mutex.Unlock()
		}(k, v)
	}
	wg.Wait()

	err_ := t_delivery_order.UpdateWithDetail(t_delivery_order, inputDetail, putDetail, user_name)
	errcode, errmessage = base.DecodeErr(err_)
	if err_ != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
		c.ServeJSON()
		return
	} else {
		if err = base.PutFirebaseRaw(ob.UploadFile, user_name, id, folderName+"/"+utils.Int2String(id), folderName+"/"+utils.Int2String(id)); err != nil {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprint("Error processing data and uploading to Firebase: ", err.Error()))
		} else {
			v, err := t_delivery_order.GetById(id, user_id)
			errcode, errmessage = base.DecodeErr(err)
			if err != nil {
				c.Ctx.ResponseWriter.WriteHeader(errcode)
				utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
			} else {
				utils.ReturnHTTPSuccessWithMessage(&c.Controller, errcode, errmessage, v)
			}
		}
	}

	c.ServeJSON()
}

func (c *DoController) Delete() {
	var user_id, form_id int
	var err error
	var user_name string
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
		user_name = sess.(map[string]interface{})["username"].(string)
	}
	form_id = base.FormName(form_do)
	delete_aut := models.CheckPrivileges(user_id, form_id, base.Delete)
	if !delete_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Delete not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	var querydata models.Do
	err = models.Dos().Filter("id", id).Filter("deleted_at__isnull", true).One(&querydata)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, "Delivery order id unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	if querydata.StatusId == base.ProgressDo {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Unable to delete", fmt.Sprintf("'%v' is on LOADING", querydata.ReferenceNo))
		c.ServeJSON()
		return
	}

	if querydata.StatusId == base.ShippedDo {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Unable to delete", fmt.Sprintf("'%v' is ON SHIPPING", querydata.ReferenceNo))
		c.ServeJSON()
		return
	}

	if querydata.StatusId == base.CompleteDo {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Unable to delete", fmt.Sprintf("'%v' has been COMPELETED", querydata.ReferenceNo))
		c.ServeJSON()
		return
	}

	// to check confirm_shipment

	models.Dos().Filter("id", id).Filter("deleted_at__isnull", true).Update(orm.Params{"deleted_at": utils.GetSvrDate(), "deleted_by": user_name})
	models.DoDetails().Filter("sales_order_id", id).Filter("deleted_at__isnull", true).Update(orm.Params{"deleted_at": utils.GetSvrDate(), "deleted_by": user_name})

	utils.ReturnHTTPError(&c.Controller, 200, "soft delete success")
	c.ServeJSON()
}
func (c *DoController) GetOne() {
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	v, err := t_delivery_order.GetById(id, user_id)
	code, message := base.DecodeErr(err)
	if err == orm.ErrNoRows {
		code = 200
		c.Ctx.ResponseWriter.WriteHeader(code)
		utils.ReturnHTTPError(&c.Controller, code, "No data")
	} else if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(code)
		utils.ReturnHTTPError(&c.Controller, code, message)
	} else {

		utils.ReturnHTTPSuccessWithMessage(&c.Controller, code, message, v)
	}
	c.ServeJSON()
}
func (c *DoController) GetAll() {}
func (c *DoController) Confirm() {
	var user_id, form_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	form_id = base.FormName(form_do)
	confirm_aut := models.CheckPrivileges(user_id, form_id, base.Approval)
	confirm_aut = true
	if !confirm_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Confirm not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	fmt.Print(id)
}

func (c *DoController) Delivery() {
	var user_id, form_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	form_id = base.FormName(form_do)
	confirm_aut := models.CheckPrivileges(user_id, form_id, base.Approval)
	confirm_aut = true
	if !confirm_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Delivery not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	fmt.Print(id)
}

func (c *DoController) Cancel() {
	var user_id, form_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	form_id = base.FormName(form_do)
	confirm_aut := models.CheckPrivileges(user_id, form_id, base.Approval)
	confirm_aut = true
	if !confirm_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Cancel not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	fmt.Print(id)
}

func (c *DoController) JobList() {
	// var user_id, form_id int
	// sess := c.GetSession("profile")
	// if sess != nil {
	// 	user_id = sess.(map[string]interface{})["id"].(int)
	// }
	// form_id = base.FormName(form_do)
	// warehouse_plant_id, _ := c.GetInt("warehouse_plant_id")

}
