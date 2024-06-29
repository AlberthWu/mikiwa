package controllers

import (
	"encoding/json"
	"fmt"
	base "mikiwa/controllers"
	"mikiwa/models"
	"mikiwa/utils"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
)

type PlantController struct {
	base.BaseController
}

func (c *PlantController) Prepare() {
	c.Ctx.Request.Header.Set("token", "No Aut")
	c.BaseController.Prepare()
}

func (c *PlantController) Post() {
	var user_id, form_id int
	var err error
	var deletedAt string

	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}

	form_id = base.FormName(form_plant)

	write_aut := models.CheckPrivileges(user_id, form_id, base.Write)
	write_aut = true
	if !write_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Post not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}

	company_id, _ := c.GetInt("company_id")
	name := strings.TrimSpace(c.GetString("name"))
	pic := strings.TrimSpace(c.GetString("pic"))
	phone := strings.TrimSpace(c.GetString("phone"))
	fax := strings.TrimSpace(c.GetString("fax"))
	address := strings.TrimSpace(c.GetString("address"))
	is_do, _ := c.GetInt8("is_do")
	is_po, _ := c.GetInt8("is_po")
	is_schedule, _ := c.GetInt8("is_schedule")
	is_receipt, _ := c.GetInt8("is_receipt")
	status, _ := c.GetInt8("status")

	var querydata models.Company
	err = models.Companies().Filter("id", company_id).One(&querydata)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Company unregistered/Illegal data")
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
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been DELETED", querydata.Code))
		c.ServeJSON()
		return
	}

	valid := validation.Validation{}
	valid.Required(company_id, "company_id").Message("Company is required")
	valid.Required(name, "name").Message("Name is required")
	valid.MinSize(name, 3, "name").Message("Name min char is 3")
	valid.MaxSize(name, 175, "name").Message("Name max char is 175")

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

	t_plant := models.Plant{
		CompanyId:   company_id,
		Name:        name,
		Pic:         pic,
		Phone:       phone,
		Fax:         fax,
		Address:     address,
		IsDo:        is_do,
		IsPo:        is_po,
		IsSchedule:  is_schedule,
		IsReceipt:   is_receipt,
		PriceMethod: 1,
		Status:      status,
	}
	d, err_ := t_plant.Insert(t_plant)
	errcode, errmessage = base.DecodeErr(err_)
	if err_ != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
	} else {
		v, err_ := t_plant.GetById(d.Id)
		errcode, errmessage = base.DecodeErr(err_)
		if err_ != nil {
			c.Ctx.ResponseWriter.WriteHeader(errcode)
			utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
		} else {
			utils.ReturnHTTPSuccessWithMessage(&c.Controller, errcode, errmessage, v)
		}
	}
	c.ServeJSON()
}

func (c *PlantController) Put() {
	var user_id, form_id int
	var err error
	var deletedAt string

	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}

	form_id = base.FormName(form_plant)

	put_aut := models.CheckPrivileges(user_id, form_id, base.Write)
	put_aut = true
	if !put_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Put not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	company_id, _ := c.GetInt("company_id")
	name := strings.TrimSpace(c.GetString("name"))
	pic := strings.TrimSpace(c.GetString("pic"))
	phone := strings.TrimSpace(c.GetString("phone"))
	fax := strings.TrimSpace(c.GetString("fax"))
	address := strings.TrimSpace(c.GetString("address"))
	is_do, _ := c.GetInt8("is_do")
	is_po, _ := c.GetInt8("is_po")
	is_schedule, _ := c.GetInt8("is_schedule")
	is_receipt, _ := c.GetInt8("is_receipt")
	status, _ := c.GetInt8("status")

	var plants models.Plant
	err = models.Plants().Filter("id", id).Filter("company_id", company_id).One(&plants)
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

	deletedAt = plants.DeletedAt.Format("2006-01-02")
	if deletedAt != "0001-01-01" {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been DELETED", plants.Name))
		c.ServeJSON()
		return
	}

	valid := validation.Validation{}
	valid.Required(company_id, "company_id").Message("Company is required")
	valid.Required(name, "name").Message("Name is required")
	valid.MinSize(name, 3, "name").Message("Name min char is 3")
	valid.MaxSize(name, 175, "name").Message("Name max char is 175")

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

	t_plant.Id = id
	t_plant.CompanyId = company_id
	t_plant.Name = name
	t_plant.Pic = pic
	t_plant.Phone = phone
	t_plant.Fax = fax
	t_plant.Address = address
	t_plant.IsDo = is_do
	t_plant.IsPo = is_po
	t_plant.IsSchedule = is_schedule
	t_plant.IsReceipt = is_receipt
	t_plant.PriceMethod = plants.PriceMethod
	t_plant.Status = status

	err_ := t_plant.Update()
	errcode, errmessage = base.DecodeErr(err_)
	if err_ != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
	} else {
		v, err_ := t_plant.GetById(id)
		errcode, errmessage = base.DecodeErr(err_)
		if err_ != nil {
			c.Ctx.ResponseWriter.WriteHeader(errcode)
			utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
		} else {
			utils.ReturnHTTPSuccessWithMessage(&c.Controller, errcode, errmessage, v)
		}
	}
	c.ServeJSON()
}

func (c *PlantController) Delete() {
	o := orm.NewOrm()
	var user_id, form_id int

	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	delete_aut := models.CheckPrivileges(user_id, form_id, base.Write)
	delete_aut = true
	if !delete_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Delete not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}

	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	var ub DeleteBody
	body := c.Ctx.Input.RequestBody
	json.Unmarshal(body, &ub)
	plants := new(models.Plant)
	var plant models.Plant

	err := o.QueryTable(plants).Filter("company_id", id).Filter("id__in", ub.Id).One(&plant)

	if err == orm.ErrNoRows {
		c.Abort("No data")
	}
	plant.DeletedAt = utils.GetSvrDate()
	o.Update(&plant, "deleted_at")
	utils.ReturnHTTPError(&c.Controller, 200, "soft delete success")
	c.ServeJSON()
}

func (c *PlantController) GetOne() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	v, err := t_plant.GetById(id)
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

func (c *PlantController) GetAllList() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	keyword := strings.TrimSpace(c.GetString("keyword"))

	d, err := t_plant.GetAllList(id, keyword)
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

func (c *PlantController) GetAllListOutlet() {
	keyword := strings.TrimSpace(c.GetString("keyword"))

	company_id := 1
	d, err := t_plant.GetAllList(company_id, keyword)
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
