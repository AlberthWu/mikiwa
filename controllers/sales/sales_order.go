package controllers

import (
	base "mikiwa/controllers"
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
		DeliveryAddress string                  `json:"delivery_address"`
		Detail          []InputDetailSalesOrder `json:"detail"`
	}

	InputDetailSalesOrder struct {
		Id           int `json:"id"`
		ProductId    int `json:"product_id"`
		QtyFormulaId int `json:"qty_formula_id"`
		Qty          int `json:"qty"`
		LeadTime     int `json:"lead_time"`
	}
)

func (c *SalesOrderController) Post() {
	// o := orm.NewOrm()
	// var user_id, form_id int
	// var user_name string
	// var err error
	// var folderName string = "sales_order"
	// sess := c.GetSession("profile")
	// if sess != nil {
	// 	user_id = sess.(map[string]interface{})["id"].(int)
	// 	user_name = sess.(map[string]interface{})["username"].(string)
	// }

	// form_id = base.FormName(form_sales_order)
	// write_aut := models.CheckPrivileges(user_id, form_id, base.Write)
	// write_aut = true
	// if !write_aut {
	// 	c.Ctx.ResponseWriter.WriteHeader(402)
	// 	utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Post not authorize", map[string]interface{}{"message": "Please contact administrator"})
	// 	c.ServeJSON()
	// 	return
	// }

	// var i int = 0
	// var ob InputHeaderSalesOrder
	// var inputDetail []models.SalesOrderDetail

	// body := c.Ctx.Input.RequestBody
	// json.Unmarshal(body, &ob)
	// valid := validation.Validation{}
	// valid.Required(strings.TrimSpace(ob.IssueDate), "issue_date").Message("Is required")
	// valid.Required(ob.CustomerId, "customer_id").Message("Is required")
	// valid.Required(strings.TrimSpace(ob.DeliveryAddress), "delivery_address").Message("Is required")

	// if len(ob.Detail) == 0 {
	// 	valid.AddError("detail", "Detail list is required")
	// }

	// if valid.HasErrors() {
	// 	out := make([]utils.ApiError, len(valid.Errors))
	// 	for i, err := range valid.Errors {
	// 		out[i] = utils.ApiError{Param: err.Key, Message: err.Message}
	// 	}
	// 	c.Ctx.ResponseWriter.WriteHeader(400)
	// 	utils.ReturnHTTPSuccessWithMessage(&c.Controller, 400, "Invalid input field", out)
	// 	c.ServeJSON()
	// 	return
	// }

	// var customers models.Company
	// err = models.Companies().Filter("id", ob.CustomerId).Filter("deleted_at__isnull", true).Filter("CompanyTypes__TypeId__Id", base.Customer).One(&customers)
	// if err == orm.ErrNoRows {
	// 	c.Ctx.ResponseWriter.WriteHeader(401)
	// 	utils.ReturnHTTPError(&c.Controller, 401, "Customer unregistered/Illegal data")
	// 	c.ServeJSON()
	// 	return
	// }

	// if err != nil {
	// 	c.Ctx.ResponseWriter.WriteHeader(401)
	// 	utils.ReturnHTTPError(&c.Controller, 401, err.Error())
	// 	c.ServeJSON()
	// 	return
	// }

	// if customers.Status == 0 {
	// 	c.Ctx.ResponseWriter.WriteHeader(402)
	// 	utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been set as INACTIVE", customers.Code))
	// 	c.ServeJSON()
	// 	return
	// }

	// var products models.Product
	// for _, v := range ob.Detail {
	// 	err = models.Products().Filter("id", v.ProductId).Filter("deleted_at__isnull", true).Filter("product_type_id", 3).One(&products)
	// 	if err == orm.ErrNoRows {
	// 		c.Ctx.ResponseWriter.WriteHeader(401)
	// 		utils.ReturnHTTPError(&c.Controller, 401, "Product unregistered/Illegal data")
	// 		c.ServeJSON()
	// 		return
	// 	}

	// 	if err != nil {
	// 		c.Ctx.ResponseWriter.WriteHeader(401)
	// 		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
	// 		c.ServeJSON()
	// 		return
	// 	}

	// 	if products.StatusId == 0 {
	// 		c.Ctx.ResponseWriter.WriteHeader(402)
	// 		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been set as INACTIVE", products.ProductName))
	// 		c.ServeJSON()
	// 		return
	// 	}
	// }

}

func (c *SalesOrderController) Put()    {}
func (c *SalesOrderController) Delete() {}
func (c *SalesOrderController) GetOne() {}
func (c *SalesOrderController) GetAll() {}
