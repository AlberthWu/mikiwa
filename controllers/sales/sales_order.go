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

type SalesOrderController struct {
	base.BaseController
}

func (c *SalesOrderController) Prepare() {
	c.Ctx.Request.Header.Set("token", "No Aut")
	c.BaseController.Prepare()
}

type (
	InputHeaderSalesOrder struct {
		IssueDate       string                  `json:"issue_date"`
		OutletId        int                     `json:"outlet_id"`
		CustomerId      int                     `json:"customer_id"`
		PlantId         int                     `json:"plant_id"`
		EmployeeId      int                     `json:"employee_id"`
		LeadTime        int                     `json:"lead_time"`
		StatusId        int8                    `json:"status_id"`
		TransporterId   int                     `json:"transporter_id"`
		DeliveryAddress string                  `json:"delivery_address"`
		UploadFile      models.DocumentList     `json:"upload_file"`
		Detail          []InputDetailSalesOrder `json:"detail"`
	}

	InputDetailSalesOrder struct {
		Id        int     `json:"id"`
		ProductId int     `json:"product_id"`
		Qty       float64 `json:"qty"`
		UomId     int     `json:"uom_id"`
		LeadTime  int     `json:"lead_time"`
		Disc1     float64 `json:"disc1"`
		Disc2     float64 `json:"disc2"`
		DiscTpr   float64 `json:"disc_tpr"`
		CreatedBy string  `json:"created_by"`
	}
)

func (c *SalesOrderController) Post() {
	o := orm.NewOrm()
	var user_id, form_id int
	var user_name string
	var err error
	var folderName string = "sales_order"
	var status_id int8 = base.OpenSo
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
		user_name = sess.(map[string]interface{})["username"].(string)
	}

	form_id = base.FormName(form_sales_order)
	write_aut := models.CheckPrivileges(user_id, form_id, base.Write)
	write_aut = true
	if !write_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Post not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}

	var i int = 0
	var ob InputHeaderSalesOrder
	var inputDetail []models.SalesOrderDetail

	body := c.Ctx.Input.RequestBody
	err = json.Unmarshal(body, &ob)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	ob.EmployeeId = 1
	ob.StatusId = status_id
	valid := validation.Validation{}
	valid.Required(strings.TrimSpace(ob.IssueDate), "issue_date").Message("Is required")
	valid.Required(ob.LeadTime, "lead_time").Message("Is required")
	valid.Required(ob.CustomerId, "customer_id").Message("Is required")
	valid.Required(ob.EmployeeId, "employee_id").Message("Is required")
	valid.Required(ob.TransporterId, "transporter_id").Message("Is required")
	valid.Required(ob.OutletId, "outlet_id").Message("Is required")
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

	var outlet models.Plant
	o.Raw("select * from plants where deleted_at is null and id = " + utils.Int2String(ob.OutletId) + " and company_id in (select company_id from company_type where type_id = " + utils.Int2String(base.Internal) + " )").QueryRow(&outlet)
	if err == orm.ErrNoRows {
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

	var transporter models.Company
	err = models.Companies().Filter("id", ob.TransporterId).Filter("deleted_at__isnull", true).Filter("CompanyTypes__TypeId__Id", base.Transporter).One(&transporter)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Transporter unregistered/Illegal data")
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

	var products models.Product
	var productUom models.ProductUom
	var priceRtn *models.ProductConversionRtnJson
	var subtotal, disc1, disc2, disctpr, totalDisc, dpp, price, normal_price, nettprice, subtotal_, totalDisc_ float64
	var ppn int
	for _, v := range ob.Detail {
		if err = o.Raw("select * from products where deleted_at is null and product_type_id = " + utils.Int2String(base.ProductFinishing) + " and product_division_id in (select business_unit_id from company_business_unit where company_id = " + utils.Int2String(ob.CustomerId) + ") and id = " + utils.Int2String(v.ProductId)).QueryRow(&products); err == orm.ErrNoRows {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, "Product unregistered/Illegal data")
			c.ServeJSON()
			return
		}

		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, err.Error())
			c.ServeJSON()
			return
		}

		if products.StatusId == 0 {
			c.Ctx.ResponseWriter.WriteHeader(402)
			utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been set as INACTIVE", products.ProductName))
			c.ServeJSON()
			return
		}

		exist := models.CompanyBusinessUnits().Filter("business_unit_id", products.ProductDivisionId).Exist()
		if !exist {
			c.Ctx.ResponseWriter.WriteHeader(402)
			utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' product division", products.ProductName))
			c.ServeJSON()
			return
		}

		err = o.Raw("select * from product_uom where product_id = " + utils.Int2String(v.ProductId) + " and uom_id = " + utils.Int2String(v.UomId)).QueryRow(&productUom)
		if err == orm.ErrNoRows {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, "Product uom unregistered/Illegal data")
			c.ServeJSON()
			return
		}

		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, err.Error())
			c.ServeJSON()
			return
		}

		priceRtn = products.GetConversion(ob.IssueDate, v.Qty, ob.CustomerId, v.ProductId, v.UomId, user_id)
		if priceRtn == nil {
			subtotal = 0
			disc1 = 0
			disc2 = 0
			disctpr = 0
			nettprice = 0
			totalDisc = 0
		} else {
			if priceRtn.Price == 0 {
				subtotal = 0
				disc1 = 0
				disc2 = 0
				disctpr = 0
				nettprice = 0
				totalDisc = 0
			} else {
				subtotal = priceRtn.FinalQty * priceRtn.Price
				disc1 = (priceRtn.Price * v.Disc1 / 100) * -1
				disc2 = ((priceRtn.Price + disc1) * v.Disc2 / 100) * -1
				disctpr = v.DiscTpr * -1
				nettprice = priceRtn.Price + disc1 + disc2 + disctpr
				totalDisc = (disc1 + disc2 + disctpr) * priceRtn.FinalQty

			}

		}

		subtotal_ += subtotal
		totalDisc_ += totalDisc
		dpp = subtotal_ + totalDisc_
		fmt.Println(v.ProductId, subtotal, totalDisc, subtotal_, totalDisc_)
	}

	if customers.IsTax == 1 {
		ppn = 11
	}

	dpp_amount, _, _, _, _, ppn_amount, total := utils.GetDppPpnTotal(ob.IssueDate, ppn, 0, 0, 0, 0, dpp)
	thedate, errDate := time.Parse("2006-01-02", ob.IssueDate)
	if errDate != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprint("issue_date: ", errDate.Error()))
		c.ServeJSON()
		return
	}

	dueDate := thedate.AddDate(0, 0, ob.LeadTime)

	seqno, referenceno := models.GenerateNumber(thedate, 1, ob.CustomerId)

	t_sales_order = models.SalesOrder{
		IssueDate:       thedate,
		ReferenceNo:     referenceno,
		SeqNo:           seqno,
		DueDate:         dueDate,
		PoolId:          1,
		OutletId:        ob.OutletId,
		OutletName:      outlet.Name,
		CustomerId:      ob.CustomerId,
		CustomerCode:    customers.Code,
		PlantId:         ob.PlantId,
		PlantName:       plants.Name,
		Terms:           customers.Terms,
		DeliveryAddress: ob.DeliveryAddress,
		EmployeeId:      ob.EmployeeId,
		EmployeeName:    "",
		TransporterId:   ob.TransporterId,
		TransporterCode: transporter.Code,
		LeadTime:        ob.LeadTime,
		Subtotal:        subtotal_,
		TotalDisc:       totalDisc_,
		Dpp:             dpp_amount,
		Ppn:             ppn,
		PpnAmount:       ppn_amount,
		Total:           total,
		StatusId:        ob.StatusId,
		CreatedBy:       user_name,
		UpdatedBy:       user_name,
	}

	wg = new(sync.WaitGroup)
	var mutex sync.Mutex
	for k, v := range ob.Detail {
		i = 0
		wg.Add(1)
		go func(k int, v InputDetailSalesOrder) {
			priceRtn = products.GetConversion(ob.IssueDate, v.Qty, ob.CustomerId, v.ProductId, v.UomId, user_id)
			if priceRtn == nil {
				disc1 = 0
				disc2 = 0
				disctpr = 0
				nettprice = 0
				subtotal = 0
				price = 0
				normal_price = 0
			} else {
				if priceRtn.Price == 0 {
					disc1 = 0
					disc2 = 0
					disctpr = 0
					nettprice = 0
					subtotal = 0
					price = 0
					normal_price = 0
				} else {
					disc1 = (priceRtn.Price * v.Disc1 / 100) * -1
					disc2 = ((priceRtn.Price + disc1) * v.Disc2 / 100) * -1
					disctpr = v.DiscTpr * -1
					price = priceRtn.Price
					normal_price = priceRtn.NormalPrice
					nettprice = price + disc1 + disc2 + disctpr
					subtotal = priceRtn.FinalQty * nettprice
				}

			}
			defer wg.Done()
			mutex.Lock()
			if v.Id == 0 {
				inputDetail = append(inputDetail, models.SalesOrderDetail{
					SalesOrderId:      t_sales_order.Id,
					ReferenceNo:       referenceno,
					IssueDate:         thedate,
					DueDate:           dueDate,
					ItemNo:            k + 1,
					ProductId:         v.ProductId,
					ProductCode:       priceRtn.ProductCode,
					Qty:               v.Qty,
					UomId:             v.UomId,
					UomCode:           priceRtn.UomCode,
					Ratio:             priceRtn.Ratio,
					PackagingId:       priceRtn.PackagingId,
					PackagingCode:     priceRtn.PackagingCode,
					FinalQty:          priceRtn.FinalQty,
					FinalUomId:        priceRtn.FinalUomId,
					FinalUomCode:      priceRtn.FinalUomCode,
					NormalPrice:       normal_price,
					PriceId:           priceRtn.PriceId,
					Price:             price,
					Disc1:             v.Disc1,
					Disc1Amount:       disc1,
					Disc2:             v.Disc2,
					Disc2Amount:       disc2,
					DiscTpr:           disctpr,
					TotalDisc:         disc1 + disc2 + disctpr,
					NettPrice:         nettprice,
					Total:             subtotal,
					LeadTime:          v.LeadTime,
					ConversionQty:     priceRtn.ConversionQty,
					ConversionUomId:   priceRtn.ConversionUomId,
					ConversionUomCode: priceRtn.ConversionUomCode,
					CreatedBy:         user_name,
					UpdatedBy:         user_name,
				})
				i += 1
			}
			mutex.Unlock()
		}(k, v)
	}
	wg.Wait()

	d, err_ := t_sales_order.InsertWithDetail(t_sales_order, inputDetail)
	errcode, errmessage = base.DecodeErr(err_)
	if err_ != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
		c.ServeJSON()
		return
	} else {
		if err = base.PostFirebaseRaw(ob.UploadFile, user_name, d.Id, folderName+"/"+utils.Int2String(d.Id), folderName+"/"+utils.Int2String(d.Id), folderName); err != nil {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprint("Error processing data and uploading to Firebase: ", err.Error()))
		} else {
			v, err := t_sales_order.GetById(d.Id, user_id)
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

func (c *SalesOrderController) Put() {
	o := orm.NewOrm()
	var user_id, form_id int
	var user_name string
	var err error
	var deletedAt string
	var folderName string = "sales_order"
	var status_id int8 = base.OpenSo
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
		user_name = sess.(map[string]interface{})["username"].(string)
	}

	form_id = base.FormName(form_sales_order)
	put_aut := models.CheckPrivileges(user_id, form_id, base.Update)
	put_aut = true
	if !put_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Post not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}

	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	var i int = 0
	var ob InputHeaderSalesOrder
	var inputDetail []models.SalesOrderDetail
	var putDetail []models.SalesOrderDetail

	body := c.Ctx.Input.RequestBody
	err = json.Unmarshal(body, &ob)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	var querydata models.SalesOrder
	err = models.SalesOrders().Filter("id", id).One(&querydata)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, "Sales order id unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	deletedAt = querydata.DeletedAt.Format("2006-01-02")
	if deletedAt != "0001-01-01" {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("'%v' has been deleted", querydata.ReferenceNo))
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

	if querydata.StatusId == base.ConfirmSo {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("'%v' has been CONFIRMED", querydata.ReferenceNo))
		c.ServeJSON()
		return
	}

	if querydata.StatusId == base.ProgressSo {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("'%v' is ON PROGRESS", querydata.ReferenceNo))
		c.ServeJSON()
		return
	}

	if querydata.StatusId == base.DoneSo {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("'%v' has been DONE", querydata.ReferenceNo))
		c.ServeJSON()
		return
	}

	ob.EmployeeId = 1
	ob.StatusId = status_id
	valid := validation.Validation{}
	valid.Required(strings.TrimSpace(ob.IssueDate), "issue_date").Message("Is required")
	valid.Required(ob.LeadTime, "lead_time").Message("Is required")
	valid.Required(ob.CustomerId, "customer_id").Message("Is required")
	valid.Required(ob.EmployeeId, "employee_id").Message("Is required")
	valid.Required(ob.TransporterId, "transporter_id").Message("Is required")
	valid.Required(ob.OutletId, "outlet_id").Message("Is required")
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

	var outlet models.Plant
	o.Raw("select * from plants where deleted_at is null and id = " + utils.Int2String(ob.OutletId) + " and company_id in (select company_id from company_type where type_id = " + utils.Int2String(base.Internal) + " )").QueryRow(&outlet)
	if err == orm.ErrNoRows {
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

	var transporter models.Company
	err = models.Companies().Filter("id", ob.TransporterId).Filter("deleted_at__isnull", true).Filter("CompanyTypes__TypeId__Id", base.Transporter).One(&transporter)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Transporter unregistered/Illegal data")
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

	var querydetail models.SalesOrderDetail
	var products models.Product
	var productUom models.ProductUom
	var priceRtn *models.ProductConversionRtnJson
	var subtotal, disc1, disc2, disctpr, totalDisc, dpp, price, normal_price, nettprice, subtotal_, totalDisc_ float64
	var ppn int
	wg = new(sync.WaitGroup)
	var mutex sync.Mutex
	resultChan := make(chan utils.ResultChan, len(ob.Detail))
	var queryResults []utils.ResultChan
	wg.Add(len(ob.Detail))
	for _, v := range ob.Detail {
		go func(v InputDetailSalesOrder) {
			defer wg.Done()
			mutex.Lock()
			defer mutex.Unlock()
			if v.Id != 0 {
				if err = models.SalesOrderDetails().Filter("deleted_at__isnull", true).Filter("sales_order_id", id).Filter("id", v.Id).One(&querydetail); err == orm.ErrNoRows {
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

			exist := models.CompanyBusinessUnits().Filter("business_unit_id", products.ProductDivisionId).Exist()
			if !exist {
				resultChan <- utils.ResultChan{Id: v.ProductId, Data: products.ProductCode, Message: fmt.Sprintf("error '%v' product division", products.ProductName)}
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

		priceRtn = products.GetConversion(ob.IssueDate, v.Qty, ob.CustomerId, v.ProductId, v.UomId, user_id)
		if priceRtn == nil {
			subtotal = 0
			disc1 = 0
			disc2 = 0
			disctpr = 0
			nettprice = 0
			totalDisc = 0
		} else {
			if priceRtn.Price == 0 {
				subtotal = 0
				disc1 = 0
				disc2 = 0
				disctpr = 0
				nettprice = 0
				totalDisc = 0
			} else {
				subtotal = priceRtn.FinalQty * priceRtn.Price
				disc1 = (priceRtn.Price * v.Disc1 / 100) * -1
				disc2 = ((priceRtn.Price + disc1) * v.Disc2 / 100) * -1
				disctpr = v.DiscTpr * -1
				nettprice = priceRtn.Price + disc1 + disc2 + disctpr
				totalDisc = (disc1 + disc2 + disctpr) * priceRtn.FinalQty

			}

		}

		subtotal_ += subtotal
		totalDisc_ += totalDisc
		dpp = subtotal_ + totalDisc_
	}

	// Use goroutine to wait until all goroutines are finished
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
	if customers.IsTax == 1 {
		ppn = 11
	}

	dpp_amount, _, _, _, _, ppn_amount, total := utils.GetDppPpnTotal(ob.IssueDate, ppn, 0, 0, 0, 0, dpp)
	thedate, errDate := time.Parse("2006-01-02", ob.IssueDate)
	if errDate != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprint("issue_date: ", errDate.Error()))
		c.ServeJSON()
		return
	}

	dueDate := thedate.AddDate(0, 0, ob.LeadTime)

	seqno := querydata.SeqNo
	referenceno := querydata.ReferenceNo

	t_sales_order.Id = id
	t_sales_order.IssueDate = thedate
	t_sales_order.ReferenceNo = referenceno
	t_sales_order.SeqNo = seqno
	t_sales_order.DueDate = dueDate
	t_sales_order.PoolId = 1
	t_sales_order.OutletId = ob.OutletId
	t_sales_order.OutletName = outlet.Name
	t_sales_order.CustomerId = ob.CustomerId
	t_sales_order.CustomerCode = customers.Code
	t_sales_order.PlantId = ob.PlantId
	t_sales_order.PlantName = plants.Name
	t_sales_order.Terms = customers.Terms
	t_sales_order.DeliveryAddress = ob.DeliveryAddress
	t_sales_order.EmployeeId = ob.EmployeeId
	t_sales_order.EmployeeName = ""
	t_sales_order.TransporterId = ob.TransporterId
	t_sales_order.TransporterCode = transporter.Code
	t_sales_order.LeadTime = ob.LeadTime
	t_sales_order.Subtotal = subtotal_
	t_sales_order.TotalDisc = totalDisc_
	t_sales_order.Dpp = dpp_amount
	t_sales_order.Ppn = ppn
	t_sales_order.PpnAmount = ppn_amount
	t_sales_order.Total = total
	t_sales_order.StatusId = ob.StatusId
	t_sales_order.CreatedBy = querydata.CreatedBy
	t_sales_order.UpdatedBy = user_name

	for k, v := range ob.Detail {
		i = 0
		wg.Add(1)
		go func(k int, v InputDetailSalesOrder) {
			priceRtn = products.GetConversion(ob.IssueDate, v.Qty, ob.CustomerId, v.ProductId, v.UomId, user_id)
			if priceRtn == nil {
				disc1 = 0
				disc2 = 0
				disctpr = 0
				nettprice = 0
				subtotal = 0
				price = 0
				normal_price = 0
			} else {
				if priceRtn.Price == 0 {
					disc1 = 0
					disc2 = 0
					disctpr = 0
					nettprice = 0
					subtotal = 0
					price = 0
					normal_price = 0
				} else {
					disc1 = (priceRtn.Price * v.Disc1 / 100) * -1
					disc2 = ((priceRtn.Price + disc1) * v.Disc2 / 100) * -1
					disctpr = v.DiscTpr * -1
					price = priceRtn.Price
					normal_price = priceRtn.NormalPrice
					nettprice = price + disc1 + disc2 + disctpr
					subtotal = priceRtn.FinalQty * nettprice
				}

			}
			defer wg.Done()
			mutex.Lock()
			if v.Id == 0 {
				inputDetail = append(inputDetail, models.SalesOrderDetail{
					SalesOrderId:      t_sales_order.Id,
					ReferenceNo:       referenceno,
					IssueDate:         thedate,
					DueDate:           dueDate,
					ItemNo:            k + 1,
					ProductId:         v.ProductId,
					ProductCode:       priceRtn.ProductCode,
					Qty:               v.Qty,
					UomId:             v.UomId,
					UomCode:           priceRtn.UomCode,
					Ratio:             priceRtn.Ratio,
					PackagingId:       priceRtn.PackagingId,
					PackagingCode:     priceRtn.PackagingCode,
					FinalQty:          priceRtn.FinalQty,
					FinalUomId:        priceRtn.FinalUomId,
					FinalUomCode:      priceRtn.FinalUomCode,
					NormalPrice:       normal_price,
					PriceId:           priceRtn.PriceId,
					Price:             price,
					Disc1:             v.Disc1,
					Disc1Amount:       disc1,
					Disc2:             v.Disc2,
					Disc2Amount:       disc2,
					DiscTpr:           disctpr,
					TotalDisc:         disc1 + disc2 + disctpr,
					NettPrice:         nettprice,
					Total:             subtotal,
					LeadTime:          v.LeadTime,
					ConversionQty:     priceRtn.ConversionQty,
					ConversionUomId:   priceRtn.ConversionUomId,
					ConversionUomCode: priceRtn.ConversionUomCode,
					CreatedBy:         user_name,
					UpdatedBy:         user_name,
				})
				i += 1
			} else {
				putDetail = append(putDetail, models.SalesOrderDetail{
					Id:                v.Id,
					SalesOrderId:      id,
					ReferenceNo:       referenceno,
					IssueDate:         thedate,
					DueDate:           dueDate,
					ItemNo:            k + 1,
					ProductId:         v.ProductId,
					ProductCode:       priceRtn.ProductCode,
					Qty:               v.Qty,
					UomId:             v.UomId,
					UomCode:           priceRtn.UomCode,
					Ratio:             priceRtn.Ratio,
					PackagingId:       priceRtn.PackagingId,
					PackagingCode:     priceRtn.PackagingCode,
					FinalQty:          priceRtn.FinalQty,
					FinalUomId:        priceRtn.FinalUomId,
					FinalUomCode:      priceRtn.FinalUomCode,
					NormalPrice:       normal_price,
					PriceId:           priceRtn.PriceId,
					Price:             price,
					Disc1:             v.Disc1,
					Disc1Amount:       disc1,
					Disc2:             v.Disc2,
					Disc2Amount:       disc2,
					DiscTpr:           disctpr,
					TotalDisc:         disc1 + disc2 + disctpr,
					NettPrice:         nettprice,
					Total:             subtotal,
					LeadTime:          v.LeadTime,
					ConversionQty:     priceRtn.ConversionQty,
					ConversionUomId:   priceRtn.ConversionUomId,
					ConversionUomCode: priceRtn.ConversionUomCode,
					CreatedBy:         v.CreatedBy,
					UpdatedBy:         user_name,
				})
			}
			mutex.Unlock()
		}(k, v)
	}
	wg.Wait()

	err_ := t_sales_order.UpdateWithDetail(t_sales_order, inputDetail, putDetail, user_name)
	errcode, errmessage = base.DecodeErr(err_)
	if err_ != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
		c.ServeJSON()
		return
	} else {
		if err = base.PutFirebaseRaw(ob.UploadFile, user_name, id, folderName+"/"+utils.Int2String(id), folderName+"/"+utils.Int2String(id), folderName); err != nil {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprint("Error processing data and uploading to Firebase: ", err.Error()))
		} else {
			v, err := t_sales_order.GetById(id, user_id)
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

func (c *SalesOrderController) Delete() {
	var user_id, form_id int
	var err error
	var user_name string
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
		user_name = sess.(map[string]interface{})["username"].(string)
	}
	form_id = base.FormName(form_sales_order)
	delete_aut := models.CheckPrivileges(user_id, form_id, base.Delete)
	if !delete_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Delete not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	var querydata models.SalesOrder
	err = models.SalesOrders().Filter("id", id).Filter("deleted_at__isnull", true).One(&querydata)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, "Sales order id unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	if querydata.StatusId == base.ConfirmSo {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Unable to delete", fmt.Sprintf("'%v' has been CONFIRMED", querydata.ReferenceNo))
		c.ServeJSON()
		return
	}

	if querydata.StatusId == base.ProgressSo {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Unable to delete", fmt.Sprintf("'%v' is ON PROGRESS", querydata.ReferenceNo))
		c.ServeJSON()
		return
	}

	if querydata.StatusId == base.DoneSo {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Unable to delete", fmt.Sprintf("'%v' has been DONE", querydata.ReferenceNo))
		c.ServeJSON()
		return
	}

	// to check invoiced

	models.SalesOrders().Filter("id", id).Filter("deleted_at__isnull", true).Update(orm.Params{"deleted_at": utils.GetSvrDate(), "deleted_by": user_name})
	models.SalesOrderDetails().Filter("sales_order_id", id).Filter("deleted_at__isnull", true).Update(orm.Params{"deleted_at": utils.GetSvrDate(), "deleted_by": user_name})

	utils.ReturnHTTPError(&c.Controller, 200, "soft delete success")
	c.ServeJSON()
}

func (c *SalesOrderController) GetOne() {
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	v, err := t_sales_order.GetById(id, user_id)
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

func (c *SalesOrderController) GetAll() {
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	var issueDate, updatedat, dueDate *string

	currentPage, _ := c.GetInt("page")
	if currentPage == 0 {
		currentPage = 1
	}

	pageSize, _ := c.GetInt("pagesize")
	if pageSize == 0 {
		pageSize = 10
	}

	keyword := strings.TrimSpace(c.GetString("keyword"))
	match_mode := strings.TrimSpace(c.GetString("match_mode"))
	value_name := strings.TrimSpace(c.GetString("value_name"))
	field_name := strings.TrimSpace(c.GetString("field_name"))
	allsize, _ := c.GetInt("allsize")
	search_detail, _ := c.GetInt("search_detail")

	status_ids := strings.TrimSpace(c.GetString("status_ids"))
	issue_date := strings.TrimSpace(c.GetString("issue_date"))
	due_date := strings.TrimSpace(c.GetString("due_date"))
	updated_at := strings.TrimSpace(c.GetString("updated_at"))
	outlet_ids := strings.TrimSpace(c.GetString("outlet_ids"))
	employee_ids := strings.TrimSpace(c.GetString("employee_ids"))
	customer_ids := strings.TrimSpace(c.GetString("customer_ids"))
	plant_id, _ := c.GetInt("plant_id")
	product_ids := strings.TrimSpace(c.GetString("product_ids"))
	lead_time_ids := strings.TrimSpace(c.GetString("lead_time_ids"))

	if issue_date == "" {
		issueDate = nil

	} else {
		issueDate = &issue_date
	}

	if due_date == "" {
		dueDate = nil

	} else {
		dueDate = &due_date
	}

	if updated_at == "" {
		updatedat = nil

	} else {
		updatedat = &updated_at
	}

	d, err := t_sales_order.GetAll(keyword, field_name, match_mode, value_name, currentPage, pageSize, allsize, user_id, 0, search_detail, plant_id, employee_ids, outlet_ids, customer_ids, status_ids, product_ids, lead_time_ids, issueDate, dueDate, updatedat)
	code, message := base.DecodeErr(err)
	if err == orm.ErrNoRows {
		code = 200
		c.Ctx.ResponseWriter.WriteHeader(code)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, code, "No data", nil)
	} else if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(code)
		utils.ReturnHTTPError(&c.Controller, code, message)
	} else {
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, code, message, d)
	}
	c.ServeJSON()
}

func (c *SalesOrderController) GetAllDetail() {
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	var issueDate, updatedat, dueDate *string

	currentPage, _ := c.GetInt("page")
	if currentPage == 0 {
		currentPage = 1
	}

	pageSize, _ := c.GetInt("pagesize")
	if pageSize == 0 {
		pageSize = 10
	}

	keyword := strings.TrimSpace(c.GetString("keyword"))
	match_mode := strings.TrimSpace(c.GetString("match_mode"))
	value_name := strings.TrimSpace(c.GetString("value_name"))
	field_name := strings.TrimSpace(c.GetString("field_name"))
	allsize, _ := c.GetInt("allsize")
	search_detail, _ := c.GetInt("search_detail")

	status_ids := strings.TrimSpace(c.GetString("status_ids"))
	issue_date := strings.TrimSpace(c.GetString("issue_date"))
	due_date := strings.TrimSpace(c.GetString("due_date"))
	updated_at := strings.TrimSpace(c.GetString("updated_at"))
	outlet_ids := strings.TrimSpace(c.GetString("outlet_ids"))
	employee_ids := strings.TrimSpace(c.GetString("employee_ids"))
	customer_ids := strings.TrimSpace(c.GetString("customer_ids"))
	plant_id, _ := c.GetInt("plant_id")
	product_ids := strings.TrimSpace(c.GetString("product_ids"))
	lead_time_ids := strings.TrimSpace(c.GetString("lead_time_ids"))

	if issue_date == "" {
		issueDate = nil

	} else {
		issueDate = &issue_date
	}

	if due_date == "" {
		dueDate = nil

	} else {
		dueDate = &due_date
	}

	if updated_at == "" {
		updatedat = nil

	} else {
		updatedat = &updated_at
	}

	d, err := t_sales_order.GetAllDetail(keyword, field_name, match_mode, value_name, currentPage, pageSize, allsize, user_id, 0, search_detail, plant_id, employee_ids, outlet_ids, customer_ids, status_ids, product_ids, lead_time_ids, issueDate, dueDate, updatedat)
	code, message := base.DecodeErr(err)
	if err == orm.ErrNoRows {
		code = 200
		c.Ctx.ResponseWriter.WriteHeader(code)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, code, "No data", nil)
	} else if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(code)
		utils.ReturnHTTPError(&c.Controller, code, message)
	} else {
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, code, message, d)
	}
	c.ServeJSON()
}
