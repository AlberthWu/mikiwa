package controllers

import (
	"fmt"
	base "mikiwa/controllers"
	"mikiwa/models"
	"mikiwa/utils"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
)

type ProductTypeController struct {
	base.BaseController
}

func (c *ProductTypeController) Prepare() {
	c.Ctx.Request.Header.Set("token", "No Aut")
	c.BaseController.Prepare()
}

func (c *ProductTypeController) Post() {
	var user_id, form_id int
	fmt.Print("Check :", user_id, form_id, "..")
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}

	form_id = base.FormName(form_product_type)

	write_aut := models.CheckPrivileges(user_id, form_id, base.Write)
	write_aut = true
	if !write_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Post not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}

	type_name := strings.TrimSpace(c.GetString("type_name"))
	is_purchase, _ := c.GetInt8("is_purchase")
	is_sales, _ := c.GetInt8("is_sales")
	is_production, _ := c.GetInt8("is_production")

	valid := validation.Validation{}
	valid.Required(type_name, "type_name").Message("is required")

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

	if t_product_type.CheckCode(0, type_name) {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("type_name : '%v' has been REGISTERED", type_name))
		c.ServeJSON()
		return
	}

	t_product_type = models.ProductType{
		TypeName:     type_name,
		IsPurchase:   is_purchase,
		IsSales:      is_sales,
		IsProduction: is_production,
	}

	d, err_ := t_product_type.Insert(t_product_type)
	errcode, errmessage := base.DecodeErr(err_)

	if err_ != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
	} else {
		v, err := t_product_type.GetById(d.Id)
		errcode, errmessage := base.DecodeErr(err)
		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(errcode)
			utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
		} else {
			utils.ReturnHTTPSuccessWithMessage(&c.Controller, errcode, errmessage, v)
		}
	}
	c.ServeJSON()
}
func (c *ProductTypeController) Put() {
	var user_id, form_id int
	fmt.Print("Check :", user_id, form_id, "..")
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}

	form_id = base.FormName(form_product_type)

	put_aut := models.CheckPrivileges(user_id, form_id, base.Update)
	put_aut = true
	if !put_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Put not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}

	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	type_name := strings.TrimSpace(c.GetString("type_name"))
	is_purchase, _ := c.GetInt8("is_purchase")
	is_sales, _ := c.GetInt8("is_sales")
	is_production, _ := c.GetInt8("is_production")

	valid := validation.Validation{}
	valid.Required(type_name, "type_name").Message("is required")

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

	if t_product_type.CheckCode(id, type_name) {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("type_name : '%v' has been REGISTERED", type_name))
		c.ServeJSON()
		return
	}

	t_product_type.Id = id
	t_product_type.TypeName = type_name
	t_product_type.IsPurchase = is_purchase
	t_product_type.IsSales = is_sales
	t_product_type.IsProduction = is_production
	err_ := t_product_type.Update()
	errcode, errmessage := base.DecodeErr(err_)

	if err_ != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
	} else {
		v, err := t_product_type.GetById(id)
		errcode, errmessage := base.DecodeErr(err)
		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(errcode)
			utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
		} else {
			utils.ReturnHTTPSuccessWithMessage(&c.Controller, errcode, errmessage, v)
		}
	}
	c.ServeJSON()
}
func (c *ProductTypeController) Delete() {}

func (c *ProductTypeController) GetOne() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	v, err := t_product_type.GetById(id)
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
func (c *ProductTypeController) GetAll() {
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
	is_purchase := strings.TrimSpace(c.GetString("is_purchase"))
	is_sales := strings.TrimSpace(c.GetString("is_sales"))
	is_production := strings.TrimSpace(c.GetString("is_production"))

	d, err := t_product_type.GetAll(keyword, field_name, match_mode, value_name, currentPage, pageSize, is_purchase, is_sales, is_production)
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

func (c *ProductTypeController) GetAllList() {
	keyword := strings.TrimSpace(c.GetString("keyword"))
	is_purchase := strings.TrimSpace(c.GetString("is_purchase"))
	is_sales := strings.TrimSpace(c.GetString("is_sales"))
	is_production := strings.TrimSpace(c.GetString("is_production"))

	d, err := t_product_type.GetAllList(keyword, is_purchase, is_sales, is_production)
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
