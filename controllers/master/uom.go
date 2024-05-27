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

type UomController struct {
	base.BaseController
}

func (c *UomController) Prepare() {
	c.Ctx.Request.Header.Set("token", "No Aut")
	c.BaseController.Prepare()
}

func (c *UomController) Post() {
	var user_id, form_id int
	fmt.Print("Check :", user_id, form_id, "..")
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}

	form_id = base.FormName(form_uom)

	write_aut := models.CheckPrivileges(user_id, form_id, base.Write)
	write_aut = true
	if !write_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Post not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}

	uom_code := strings.TrimSpace(c.GetString("uom_code"))
	uom_name := strings.TrimSpace(c.GetString("uom_name"))
	status_id, _ := c.GetInt8("status_id")

	valid := validation.Validation{}
	valid.Required(uom_code, "uom_code").Message("is required")
	valid.Required(uom_name, "uom_name").Message("is required")

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

	if t_uom.CheckCode(0, uom_code) {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("uom_code : '%v' has been REGISTERED", uom_code))
		c.ServeJSON()
		return
	}

	t_uom = models.Uom{
		UomCode:  uom_code,
		UomName:  uom_name,
		StatusId: status_id,
	}

	d, err_ := t_uom.Insert(t_uom)
	errcode, errmessage := base.DecodeErr(err_)

	if err_ != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
	} else {
		v, err := t_uom.GetById(d.Id)
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

func (c *UomController) Put() {
	var err error
	var user_id, form_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}

	form_id = base.FormName(form_uom)

	put_aut := models.CheckPrivileges(user_id, form_id, base.Update)
	put_aut = true
	if !put_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Put not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	uom_code := strings.TrimSpace(c.GetString("uom_code"))
	uom_name := strings.TrimSpace(c.GetString("uom_name"))
	status_id, _ := c.GetInt8("status_id")

	var querydata models.Uom
	err = models.Uoms().Filter("id", id).One(&querydata)
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
	valid.Required(uom_code, "uom_code").Message("is required")
	valid.Required(uom_name, "uom_name").Message("is required")

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

	if t_uom.CheckCode(id, uom_code) {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("uom_code : '%v' has been REGISTERED", uom_code))
		c.ServeJSON()
		return
	}

	t_uom.Id = id
	t_uom.UomCode = uom_code
	t_uom.UomName = uom_name
	t_uom.StatusId = status_id
	err_ := t_uom.Update()
	errcode, errmessage := base.DecodeErr(err_)
	if err_ != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
	} else {
		v, err := t_uom.GetById(id)
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

func (c *UomController) Delete() {}

func (c *UomController) GetOne() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	v, err := t_uom.GetById(id)
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

func (c *UomController) GetAll() {
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

	d, err := t_uom.GetAll(keyword, field_name, match_mode, value_name, currentPage, pageSize, status_id)
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

func (c *UomController) GetAllList() {
	keyword := strings.TrimSpace(c.GetString("keyword"))

	d, err := t_uom.GetAllList(keyword)
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
