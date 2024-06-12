package controllers

import (
	base "mikiwa/controllers"
)

type PriceController struct {
	base.BaseController
}

func (c *PriceController) Prepare() {
	c.Ctx.Request.Header.Set("token", "No Aut")
	c.BaseController.Prepare()
}

func (c *PriceController) Post() {
	// var user_id, form_id int
	// fmt.Print("Check :", user_id, form_id, "..")
	// sess := c.GetSession("profile")
	// if sess != nil {
	// 	user_id = sess.(map[string]interface{})["id"].(int)
	// }

	// form_id = base.FormName(form_price)

	// write_aut := models.CheckPrivileges(user_id, form_id, base.Write)
	// write_aut = true
	// if !write_aut {
	// 	c.Ctx.ResponseWriter.WriteHeader(402)
	// 	utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Post not authorize", map[string]interface{}{"message": "Please contact administrator"})
	// 	c.ServeJSON()
	// 	return
	// }
	// price_type := "sales"
	// effective_date := strings.TrimSpace(c.GetString("effective_date"))
	// expired_date := strings.TrimSpace(c.GetString("expired_date"))
	// product_id, _ := c.GetInt("product_id")
	// price, _ := c.GetFloat("price")
	// uom_id_one, _ := c.GetInt("uom_id_one")
	// ratio, _ := c.GetFloat("ratio")
	// disc_one, _ := c.GetFloat("disc_one")
	// disc_two, _ := c.GetFloat("disc_two")
	// disc_tpr, _ := c.GetFloat("disc_tpr")
	// sure_name := strings.TrimSpace(c.GetString("sure_name"))
	// status_id, _ := c.GetInt8("status_id")

	// valid := validation.Validation{}
	// valid.Required(effective_date, "effective_date").Message("is required")
	// valid.Required(product_id, "product_id").Message("is required")
	// valid.Required(uom_id_one, "uom_id_one").Message("is required")
	// if price == 0 {
	// 	valid.AddError("price", "0 not allowed")
	// }
	// if ratio == 0 {
	// 	valid.AddError("ratio", "0 not allowed")
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

	// t_price = models.Price{
	// 	StatusId: status_id,
	// }

	// d, err_ := t_price.Insert(t_price)
	// errcode, errmessage := base.DecodeErr(err_)

	// if err_ != nil {
	// 	c.Ctx.ResponseWriter.WriteHeader(errcode)
	// 	utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
	// } else {
	// 	v, err := t_price.GetById(d.Id)
	// 	errcode, errmessage := base.DecodeErr(err)
	// 	if err != nil {
	// 		c.Ctx.ResponseWriter.WriteHeader(errcode)
	// 		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
	// 	} else {
	// 		utils.ReturnHTTPSuccessWithMessage(&c.Controller, errcode, errmessage, v)
	// 	}
	// }
	// c.ServeJSON()
}
func (c *PriceController) Put()      {}
func (c *PriceController) Delete()   {}
func (c *PriceController) GetOne()   {}
func (c *PriceController) GetAll()   {}
func (c *PriceController) GetPrice() {}
