package controllers

import (
	"encoding/json"
	"fmt"
	base "mikiwa/controllers"
	"mikiwa/models"
	"mikiwa/utils"
	"strings"
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
		CustomerId      int                     `json:"customer_id"`
		EmployeeId      int                     `json:"employee_id"`
		LeadTime        int                     `json:"lead_time"`
		StatusId        int8                    `json:"status_id"`
		DeliveryAddress string                  `json:"delivery_address"`
		Detail          []InputDetailSalesOrder `json:"detail"`
	}

	InputDetailSalesOrder struct {
		Id        int     `json:"id"`
		ProductId int     `json:"product_id"`
		PriceId   int     `json:"price_id"`
		Qty       float64 `json:"qty"`
		UomId     int     `json:"uom_id"`
		LeadTime  int     `json:"lead_time"`
		Disc1     float64 `json:"disc1"`
		Disc2     float64 `json:"disc2"`
		DiscTpr   float64 `json:"disc_tpr"`
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
	json.Unmarshal(body, &ob)

	ob.EmployeeId = 1
	ob.StatusId = status_id
	valid := validation.Validation{}
	valid.Required(strings.TrimSpace(ob.IssueDate), "issue_date").Message("Is required")
	valid.Required(ob.LeadTime, "lead_time").Message("Is required")
	valid.Required(ob.CustomerId, "customer_id").Message("Is required")
	valid.Required(ob.EmployeeId, "employee_id").Message("Is required")
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

	var products models.Product
	var productUom models.ProductUom
	var priceRtn *models.ProductConversionRtnJson
	var disc1, disc2, discTpr, subtotal, nettprice, totalDisc, total float64
	for _, v := range ob.Detail {
		err = models.Products().Filter("id", v.ProductId).Filter("deleted_at__isnull", true).Filter("product_type_id", 3).One(&products)
		if err == orm.ErrNoRows {
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

		err = models.ProductUoms().Filter("product_id", v.ProductId).Filter("uom_id", v.UomId).One(&productUom)
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

		priceRtn = products.GetConversion(ob.IssueDate, v.Qty, v.ProductId, ob.CustomerId, v.UomId, user_id)
		disc1 = (priceRtn.Price * v.Disc1 / 100) * -1
		disc2 = ((priceRtn.Price + disc1) * v.Disc2 / 100) * -1
		discTpr = priceRtn.Price + disc1 + disc2 - v.DiscTpr
		nettprice = priceRtn.Price + discTpr
		subtotal = v.Qty * nettprice
		totalDisc += totalDisc + disc1 + disc2 + (v.DiscTpr * -1)
		total += total + subtotal
	}

	thedate, errDate := time.Parse("2006-01-02", ob.IssueDate)
	if errDate != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprint("issue_date: ", errDate.Error()))
		c.ServeJSON()
		return
	}

	dueDate := thedate.AddDate(0, 0, ob.LeadTime)

	seqno, referenceno := models.GenerateNumber(thedate, 1, ob.CustomerId)

	tx, errTrans := o.Begin()
	if errTrans != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, errTrans.Error())
		c.ServeJSON()
		return
	}

	t_sales_order = models.SalesOrder{
		IssueDate:       thedate,
		ReferenceNo:     referenceno,
		SeqNo:           seqno,
		DueDate:         dueDate,
		PoolId:          1,
		CustomerId:      ob.CustomerId,
		CustomerName:    customers.Name,
		Terms:           customers.Terms,
		DeliveryAddress: ob.DeliveryAddress,
		EmployeeId:      ob.EmployeeId,
		EmployeeName:    "",
		LeadTime:        ob.LeadTime,
		StatusId:        ob.StatusId,
		Subtotal:        total,
		TotalDisc:       totalDisc,
		CreatedBy:       user_name,
		UpdatedBy:       user_name,
	}
	_, err_ := t_sales_order.Insert(t_sales_order)
	errcode, errmessage = base.DecodeErr(err_)

	if err_ != nil {
		tx.Rollback()
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
		c.ServeJSON()
		return
	} else {

	}

	errTrans = tx.Commit()
	errcode, errmessage = base.DecodeErr(errTrans)
	if errTrans != nil {
		tx.Rollback()
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
	} else {
	}
	c.ServeJSON()
}

func (c *SalesOrderController) Put()    {}
func (c *SalesOrderController) Delete() {}
func (c *SalesOrderController) GetOne() {}
func (c *SalesOrderController) GetAll() {}
