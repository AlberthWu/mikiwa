package controllers

import (
	"fmt"
	"strconv"
	"strings"

	base "mikiwa/controllers"
	"mikiwa/models"
	"mikiwa/utils"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
)

type AccountTypeController struct {
	base.BaseController
}

func (c *AccountTypeController) Prepare() {
	c.Ctx.Request.Header.Set("token", "No Aut")
	c.BaseController.Prepare()
}

func (c *AccountTypeController) Post() {
	var user_name string
	var user_id, form_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_name = sess.(map[string]interface{})["username"].(string)
		user_id = sess.(map[string]interface{})["id"].(int)
	}

	form_id = base.FormName(form_account_type)

	write_aut := models.CheckPrivileges(user_id, form_id, base.Write)
	write_aut = true
	if !write_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Post not authorization", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}

	name := strings.TrimSpace(c.GetString("name"))
	account_code := strings.TrimSpace(c.GetString("account_code"))
	journal_position := strings.TrimSpace(c.GetString("journal_position"))
	component_account := strings.TrimSpace(c.GetString("component_account"))

	valid := validation.Validation{}
	valid.Required(account_code, "account_code").Message("is required")
	valid.Required(name, "name").Message("is required")
	valid.Required(journal_position, "journal_position").Message("is required")
	valid.Required(component_account, "component_account").Message("is required")

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

	if exist := models.GlAccountTypes().Filter("deleted_at__isnull", true).Filter("account_code", account_code).Exist(); exist == true {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("account_code : '%v' has been REGISTERED", account_code))
		c.ServeJSON()
		return
	}

	t_account_type = models.GlAccountType{
		Name:             name,
		AccountCode:      account_code,
		JournalPosition:  journal_position,
		ComponentAccount: component_account,
		CreatedBy:        user_name,
	}

	d, err_ := t_account_type.Insert(t_account_type)
	errcode, errmessage := base.DecodeErr(err_)

	if err_ != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
	} else {
		v, _ := t_account_type.GetById(d.Id)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, errcode, errmessage, v)
	}
	c.ServeJSON()
}

func (c *AccountTypeController) Put() {

	var err error
	var deletedat string
	var user_name string
	var user_id, form_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_name = sess.(map[string]interface{})["username"].(string)
		user_id = sess.(map[string]interface{})["id"].(int)
	}

	form_id = base.FormName(form_account_type)

	put_aut := models.CheckPrivileges(user_id, form_id, base.Update)
	put_aut = true
	if !put_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Put not authorization", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}

	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	name := strings.TrimSpace(c.GetString("name"))
	account_code := strings.TrimSpace(c.GetString("account_code"))
	journal_position := strings.TrimSpace(c.GetString("journal_position"))
	component_account := strings.TrimSpace(c.GetString("component_account"))

	var querydata models.GlAccountType
	err = models.GlAccountTypes().Filter("id", id).One(&querydata)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, "Account type id unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	deletedat = querydata.DeletedAt.Format("2006-01-02")
	if deletedat != "0001-01-01" {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("account_type_id :'%v' has been deleted", querydata.Name))
		c.ServeJSON()
		return
	}

	valid := validation.Validation{}
	valid.Required(account_code, "account_code").Message("is required")
	valid.Required(name, "name").Message("is required")
	valid.Required(journal_position, "journal_position").Message("is required")
	valid.Required(component_account, "component_account").Message("is required")

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

	if exist := models.GlAccountTypes().Exclude("id", id).Filter("deleted_at__isnull", true).Filter("account_code", account_code).Exist(); exist == true {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("account_code : '%v' has been REGISTERED", account_code))
		c.ServeJSON()
		return
	}

	t_account_type.Id = id
	t_account_type.AccountCode = account_code
	t_account_type.Name = name
	t_account_type.JournalPosition = journal_position
	t_account_type.ComponentAccount = component_account
	t_account_type.CreatedBy = querydata.CreatedBy
	t_account_type.UpdatedBy = user_name

	err_ := t_account_type.Update()
	errcode, errmessage := base.DecodeErr(err_)

	if err_ != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
	} else {
		v, _ := t_account_type.GetById(id)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, errcode, errmessage, v)
	}
	c.ServeJSON()
}

func (c *AccountTypeController) Delete() {
	var err error
	var user_id, form_id int = 0, 0
	var user_name string

	sess := c.GetSession("profile")
	if sess != nil {
		user_name = sess.(map[string]interface{})["username"].(string)
		user_id = sess.(map[string]interface{})["id"].(int)
	}

	form_id = base.FormName(form_account_type)

	delete_aut := models.CheckPrivileges(user_id, form_id, base.Delete)
	if !delete_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Delete not authorization", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}

	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	var querydata models.GlAccountType
	err = models.GlAccountTypes().Filter("id", id).One(&querydata)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, "Account type id unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	var deletedat = querydata.DeletedAt.Format("2006-01-02")
	if deletedat != "0001-01-01" {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("account_type_id :'%v' has been deleted", querydata.Name))
		c.ServeJSON()
		return
	}

	models.GlAccountTypes().Filter("id", id).Filter("deleted_at__isnull", true).Update(orm.Params{"deleted_at": utils.GetSvrDate(), "deleted_by": user_name})

	utils.ReturnHTTPError(&c.Controller, 200, "soft delete success")
	c.ServeJSON()
}

func (c *AccountTypeController) GetOne() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	v, err := t_account_type.GetById(id)
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

func (c *AccountTypeController) GetAll() {

	currentPage, _ := c.GetInt("page")
	if currentPage == 0 {
		currentPage = 1
	}

	pageSize, _ := c.GetInt("pagesize")
	if pageSize == 0 {
		pageSize = 10
	}
	keyword := strings.TrimSpace(c.GetString("keyword"))

	d, err := t_account_type.GetAll(keyword, currentPage, pageSize)
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

func (c *AccountTypeController) GetAllList() {
	keyword := strings.TrimSpace(c.GetString("keyword"))
	compaonent_account := strings.TrimSpace(c.GetString("compaonent_account"))
	d, err := t_account_type.GetAllList(keyword, compaonent_account)
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

func (c *AccountTypeController) GetAllListAssets() {
	keyword := strings.TrimSpace(c.GetString("keyword"))
	d, err := t_account_type.GetAllListAsset(keyword)
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

func (c *AccountTypeController) GetAllListExpenses() {
	keyword := strings.TrimSpace(c.GetString("keyword"))
	d, err := t_account_type.GetAllListExpenses(keyword)
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

func (c *AccountTypeController) GetAllListLiability() {
	keyword := strings.TrimSpace(c.GetString("keyword"))
	d, err := t_account_type.GetAllListLiability(keyword)
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

func (c *AccountTypeController) GetAllListEquity() {
	keyword := strings.TrimSpace(c.GetString("keyword"))
	d, err := t_account_type.GetAllListEquity(keyword)
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

func (c *AccountTypeController) GetAllListRevenue() {
	keyword := strings.TrimSpace(c.GetString("keyword"))
	d, err := t_account_type.GetAllListRevenue(keyword)
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

func (c *AccountTypeController) GetAllListCogs() {
	keyword := strings.TrimSpace(c.GetString("keyword"))
	d, err := t_account_type.GetAllListCogs(keyword)
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
