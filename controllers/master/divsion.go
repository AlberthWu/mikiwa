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

type ProductDivisionController struct {
	base.BaseController
}

func (c *ProductDivisionController) Post() {
	var user_id, form_id int
	var user_name string
	fmt.Print("Check :", user_id, form_id, user_name, "..")
	sess := c.GetSession("profile")
	if sess != nil {
		user_name = sess.(map[string]interface{})["username"].(string)
		user_id = sess.(map[string]interface{})["id"].(int)
	}

	form_id = base.FormName(form_product_division)

	write_aut := models.CheckPrivileges(user_id, form_id, base.Write)
	if !write_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Post not authorization", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}

	division_code := strings.TrimSpace(c.GetString("division_code"))
	division_name := strings.TrimSpace(c.GetString("division_name"))
	status_id, _ := c.GetInt8("status_id")

	valid := validation.Validation{}
	valid.Required(division_code, "division_code").Message("is required")
	valid.Required(division_name, "division_name").Message("is required")

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

	if t_product_division.CheckCode(0, division_code) {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("divison_code : '%v' has been REGISTERED", division_code))
		c.ServeJSON()
		return
	}

	t_product_division = models.ProductDivision{
		DivisionCode: division_code,
		DivisionName: division_name,
		StatusId:     status_id,
	}

	d, err_ := t_product_division.Insert(t_product_division)
	errcode, errmessage := base.DecodeErr(err_)

	if err_ != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
	} else {
		v, err := t_product_division.GetById(d.Id)
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
func (c *ProductDivisionController) Put() {
	var err error
	var user_id, form_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}

	form_id = base.FormName(form_product_division)

	put_aut := models.CheckPrivileges(user_id, form_id, base.Update)
	if !put_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Put not authorization", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	division_code := strings.TrimSpace(c.GetString("division_code"))
	division_name := strings.TrimSpace(c.GetString("division_name"))
	status_id, _ := c.GetInt8("status_id")

	var querydata models.ProductDivision
	err = models.ProductDivisions().Filter("id", id).One(&querydata)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, "Id unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	valid := validation.Validation{}
	valid.Required(division_code, "division_code").Message("is required")
	valid.Required(division_name, "division_name").Message("is required")

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

	if t_product_division.CheckCode(id, division_code) {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("divison_code : '%v' has been REGISTERED", division_code))
		c.ServeJSON()
		return
	}

	t_product_division.Id = id
	t_product_division.DivisionCode = division_code
	t_product_division.DivisionName = division_name
	t_product_division.StatusId = status_id
	err_ := t_product_division.Update()
	errcode, errmessage := base.DecodeErr(err_)
	if err_ != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
	} else {
		v, err := t_product_division.GetById(id)
		errcode, errmessage := base.DecodeErr(err)
		if err_ != nil {
			c.Ctx.ResponseWriter.WriteHeader(errcode)
			utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
		} else {
			utils.ReturnHTTPSuccessWithMessage(&c.Controller, errcode, errmessage, v)
		}
	}
	c.ServeJSON()

}

func (c *ProductDivisionController) Delete() {}

func (c *ProductDivisionController) GetOne() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	v, err := t_product_division.GetById(id)
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
func (c *ProductDivisionController) GetAll() {
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
	status_id := strings.TrimSpace(c.GetString("status_id"))

	d, err := t_product_division.GetAll(keyword, field_name, match_mode, value_name, currentPage, pageSize, status_id)
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
