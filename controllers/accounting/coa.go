package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	base "mikiwa/controllers"
	"mikiwa/models"
	"mikiwa/utils"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
)

type CoaController struct {
	base.BaseController
}

func (c *CoaController) Prepare() {
	c.Ctx.Request.Header.Set("token", "No Aut")
	c.BaseController.Prepare()
}

func (c *CoaController) Post() {
	var err error
	var deletedat string
	var user_name string
	var user_id, form_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_name = sess.(map[string]interface{})["username"].(string)
		user_id = sess.(map[string]interface{})["id"].(int)
	}

	form_id = base.FormName(form_coa)

	write_aut := models.CheckPrivileges(user_id, form_id, base.Write)
	if !write_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Post not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}
	effective_date := strings.TrimSpace(c.GetString("effective_date"))
	expired_date := strings.TrimSpace(c.GetString("expired_date"))
	account_type_id, _ := c.GetInt("account_type_id")
	company_id, _ := c.GetInt("company_id")
	level_no, _ := c.GetInt("level_no")
	parent_id, _ := c.GetInt("parent_id")
	code_coa := strings.TrimSpace(c.GetString("code_coa"))
	name_coa := strings.TrimSpace(c.GetString("name_coa"))
	code_in := strings.TrimSpace(c.GetString("code_in"))
	code_out := strings.TrimSpace(c.GetString("code_out"))
	status_id, _ := c.GetInt("status_id")
	is_header, _ := c.GetInt8("is_header")
	// sales_type_id := strings.TrimSpace(c.GetString("sales_type_id"))

	valid := validation.Validation{}
	valid.Required(effective_date, "effective_date").Message("is required")
	// valid.Required(company_id, "company_id").Message("is required")
	valid.Required(code_coa, "code_coa").Message("is required")
	valid.Required(name_coa, "name_coa").Message("is required")
	// valid.Required(sales_type_id, "sales_type_id").Message("is required")

	if level_no == 1 {
		valid.AddError("level_no", "Level no : 1 not allowed")
	}

	if level_no > 2 {
		valid.Required(parent_id, "parent_id").Message("is required")
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

	if t_coa.CheckCode(0, company_id, code_coa) {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("code_coa : '%v' has been REGISTERED", code_coa))
		c.ServeJSON()
		return
	}

	var accounttype models.GlAccountType
	if err = models.GlAccountTypes().Filter("id", account_type_id).Filter("deleted_at__isnull", true).One(&accounttype); err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Account type unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	deletedat = accounttype.DeletedAt.Format("2006-01-02")
	if deletedat != "0001-01-01" {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("account_type_id :'%v' has been deleted", accounttype.Name))
		c.ServeJSON()
		return
	}

	var companies models.Company
	if err = models.Companies().Filter("id", company_id).Filter("CompanyTypes__TypeId__Id", 1).One(&companies); err == orm.ErrNoRows {
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

	deletedat = companies.DeletedAt.Format("2006-01-02")
	if deletedat != "0001-01-01" {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("company_id :'%v' has been deleted", companies.Code))
		c.ServeJSON()
		return
	}

	var coa models.CharOfAccount
	if parent_id > 0 {
		if err = models.ChartOfAccounts().Filter("account_type_id", account_type_id).Filter("id", parent_id).One(&coa); err == orm.ErrNoRows {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, "Parent account unregistered/Illegal data")
			c.ServeJSON()
			return
		}

		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, err.Error())
			c.ServeJSON()
			return
		}

		deletedat = coa.DeletedAt.Format("2006-01-02")
		if deletedat != "0001-01-01" {
			c.Ctx.ResponseWriter.WriteHeader(402)
			utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("chart_of_account_id :'%v' has been deleted", coa.CodeCoa))
			c.ServeJSON()
			return
		}

		if level_no-1 != coa.LevelNo {
			c.Ctx.ResponseWriter.WriteHeader(402)
			utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("parent_id : '%v' Invalid sub level", coa.CodeCoa))
			c.ServeJSON()
			return
		}
	}

	issuedate, errDate := time.Parse("2006-01-02", effective_date)
	if errDate != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, errDate.Error())
		c.ServeJSON()
		return

	}

	var expireddate *time.Time
	if expired_date == "" {
		expireddate = nil
	} else {
		expiredthedate, err_date := time.Parse("2006-01-02", expired_date)
		if err_date != nil {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprint("expired_date: ", err_date.Error()))
			c.ServeJSON()
			return
		}
		expireddate = &expiredthedate
	}

	t_coa = models.CharOfAccount{
		EffectiveDate:   issuedate,
		ExpiredDate:     expireddate,
		AccountTypeId:   account_type_id,
		AccountTypeName: accounttype.Name,
		CompanyId:       company_id,
		CompanyCode:     companies.Code,
		CompanyName:     companies.Name,
		LevelNo:         level_no,
		ParentId:        parent_id,
		ParentCode:      coa.CodeCoa,
		CodeCoa:         code_coa,
		NameCoa:         name_coa,
		CodeIn:          code_in,
		CodeOut:         code_out,
		StatusId:        int8(status_id),
		JournalPosition: accounttype.JournalPosition,
		IsHeader:        is_header,
		CreatedBy:       user_name,
	}

	d, err_ := t_coa.Insert(t_coa)

	errcode, errmessage := base.DecodeErr(err_)

	if err_ != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
	} else {
		v, _ := t_coa.GetById(d.Id)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, errcode, errmessage, v)
	}
	c.ServeJSON()
}

func (c *CoaController) Put() {
	var err error
	var deletedat string
	var user_name string
	var user_id, form_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_name = sess.(map[string]interface{})["username"].(string)
		user_id = sess.(map[string]interface{})["id"].(int)
	}

	form_id = base.FormName(form_coa)

	put_aut := models.CheckPrivileges(user_id, form_id, base.Update)
	if !put_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Put not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}

	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	effective_date := strings.TrimSpace(c.GetString("effective_date"))
	expired_date := strings.TrimSpace(c.GetString("expired_date"))
	account_type_id, _ := c.GetInt("account_type_id")
	company_id, _ := c.GetInt("company_id")
	level_no, _ := c.GetInt("level_no")
	parent_id, _ := c.GetInt("parent_id")
	code_coa := strings.TrimSpace(c.GetString("code_coa"))
	name_coa := strings.TrimSpace(c.GetString("name_coa"))
	code_in := strings.TrimSpace(c.GetString("code_in"))
	code_out := strings.TrimSpace(c.GetString("code_out"))
	status_id, _ := c.GetInt("status_id")
	is_header, _ := c.GetInt8("is_header")
	// sales_type_id := strings.TrimSpace(c.GetString("sales_type_id"))

	var querydata models.CharOfAccount
	err = models.ChartOfAccounts().Filter("id", id).One(&querydata)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, "Chart of account id unregistered/Illegal data")
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
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("chart_of_account_id :'%v' has been deleted", querydata.CodeCoa))
		c.ServeJSON()
		return
	}

	valid := validation.Validation{}
	valid.Required(effective_date, "effective_date").Message("is required")
	// valid.Required(company_id, "company_id").Message("is required")
	valid.Required(code_coa, "code_coa").Message("is required")
	valid.Required(name_coa, "name_coa").Message("is required")
	// valid.Required(sales_type_id, "sales_type_id").Message("is required")

	if level_no == 1 {
		valid.AddError("level_no", "Level no : 1 not allowed")
	}

	if level_no > 2 {
		valid.Required(parent_id, "parent_id").Message("is required")
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

	if t_coa.CheckCode(id, company_id, code_coa) {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("code_coa : '%v' has been REGISTERED", code_coa))
		c.ServeJSON()
		return
	}

	var accounttype models.GlAccountType
	if err = models.GlAccountTypes().Filter("id", account_type_id).Filter("deleted_at__isnull", true).One(&accounttype); err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Account type unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	deletedat = accounttype.DeletedAt.Format("2006-01-02")
	if deletedat != "0001-01-01" {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("account_type_id :'%v' has been deleted", accounttype.Name))
		c.ServeJSON()
		return
	}

	var companies models.Company
	if err = models.Companies().Filter("id", company_id).Filter("CompanyTypes__TypeId__Id", 1).One(&companies); err == orm.ErrNoRows {
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

	deletedat = companies.DeletedAt.Format("2006-01-02")
	if deletedat != "0001-01-01" {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("company_id :'%v' has been deleted", companies.Code))
		c.ServeJSON()
		return
	}

	var coa models.CharOfAccount
	if parent_id > 0 {
		if err = models.ChartOfAccounts().Filter("account_type_id", account_type_id).Filter("id", parent_id).One(&coa); err == orm.ErrNoRows {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, "Parent account unregistered/Illegal data")
			c.ServeJSON()
			return
		}

		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, err.Error())
			c.ServeJSON()
			return
		}

		deletedat = coa.DeletedAt.Format("2006-01-02")
		if deletedat != "0001-01-01" {
			c.Ctx.ResponseWriter.WriteHeader(402)
			utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("chart_of_account_id :'%v' has been deleted", coa.CodeCoa))
			c.ServeJSON()
			return
		}

		if level_no-1 != coa.LevelNo {
			c.Ctx.ResponseWriter.WriteHeader(402)
			utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("parent_id : '%v' Invalid sub level", coa.CodeCoa))
			c.ServeJSON()
			return
		}
	}

	issuedate, errDate := time.Parse("2006-01-02", effective_date)
	if errDate != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, errDate.Error())
		c.ServeJSON()
		return

	}

	var expireddate *time.Time
	if expired_date == "" {
		expireddate = nil
	} else {
		expiredthedate, err_date := time.Parse("2006-01-02", expired_date)
		if err_date != nil {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprint("expired_date: ", err_date.Error()))
			c.ServeJSON()
			return
		}
		expireddate = &expiredthedate
	}

	t_coa.Id = id
	t_coa.EffectiveDate = issuedate
	t_coa.ExpiredDate = expireddate
	t_coa.AccountTypeId = account_type_id
	t_coa.AccountTypeName = accounttype.Name
	t_coa.CompanyId = company_id
	t_coa.CompanyCode = companies.Code
	t_coa.CompanyName = companies.Name
	t_coa.LevelNo = level_no
	t_coa.ParentId = parent_id
	t_coa.CodeCoa = code_coa
	t_coa.NameCoa = name_coa
	t_coa.CodeIn = code_in
	t_coa.CodeOut = code_out
	t_coa.StatusId = int8(status_id)
	t_coa.JournalPosition = accounttype.JournalPosition
	t_coa.IsHeader = is_header
	t_coa.CreatedBy = querydata.CreatedBy
	t_coa.UpdatedBy = user_name
	err_ := t_coa.Update()
	errcode, errmessage := base.DecodeErr(err_)

	if err_ != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
	} else {
		v, _ := t_coa.GetById(id)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, errcode, errmessage, v)
	}
	c.ServeJSON()
}

func (c *CoaController) Delete() {
	var err error
	var user_id, form_id int = 0, 0
	var user_name string

	sess := c.GetSession("profile")
	if sess != nil {
		user_name = sess.(map[string]interface{})["username"].(string)
		user_id = sess.(map[string]interface{})["id"].(int)
	}

	form_id = base.FormName(form_coa)

	delete_aut := models.CheckPrivileges(user_id, form_id, base.Delete)
	if !delete_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Delete not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}

	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	var querydata models.CharOfAccount
	err = models.ChartOfAccounts().Filter("id", id).One(&querydata)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, "Char of account unregistered/Illegal data")
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
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("char_of_account_id :'%v' has been deleted", querydata.CodeCoa))
		c.ServeJSON()
		return
	}

	models.ChartOfAccounts().Filter("id", id).Filter("deleted_at__isnull", true).Update(orm.Params{"deleted_at": utils.GetSvrDate(), "deleted_by": user_name})

	utils.ReturnHTTPError(&c.Controller, 200, "soft delete success")
	c.ServeJSON()
}

func (c *CoaController) GetOne() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	v, err := t_coa.GetById(id)
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

func (c *CoaController) GetAll() {
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	user_id = 1
	var issuedate, updatedat *string

	currentPage, _ := c.GetInt("page")
	if currentPage == 0 {
		currentPage = 1
	}

	pageSize, _ := c.GetInt("pagesize")
	if pageSize == 0 {
		pageSize = 10
	}

	issue_date := strings.TrimSpace(c.GetString("issue_date"))
	updated_at := strings.TrimSpace(c.GetString("updated_at"))
	keyword := strings.TrimSpace(c.GetString("keyword"))
	match_mode := strings.TrimSpace(c.GetString("match_mode"))
	value_name := strings.TrimSpace(c.GetString("value_name"))
	field_name := strings.TrimSpace(c.GetString("field_name"))
	account_type_id, _ := c.GetInt("account_type_id")
	level_no, _ := c.GetInt("level_no")
	parent_level_no, _ := c.GetInt("parent_level_no")
	company_id, _ := c.GetInt("company_id")
	sales_type_id, _ := c.GetInt("sales_type_id")
	if issue_date == "" {
		issuedate = nil
	} else {
		issuedate = &issue_date
	}

	if updated_at == "" {
		updatedat = nil
	} else {
		updatedat = &updated_at
	}
	d, err := t_coa.GetAll(keyword, field_name, match_mode, value_name, currentPage, pageSize, parent_level_no, level_no, account_type_id, company_id, sales_type_id, user_id, issuedate, updatedat)
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

func (c *CoaController) GetAllLimit() {
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	user_id = 1
	keyword := strings.TrimSpace(c.GetString("keyword"))
	account_type_id, _ := c.GetInt("account_type_id")
	level_no, _ := c.GetInt("level_no")
	parent_level_no, _ := c.GetInt("parent_level_no")
	sales_type_id, _ := c.GetInt("sales_type_id")

	company_id := 1

	d, err := t_coa.GetAllLimit(keyword, parent_level_no, level_no, account_type_id, company_id, sales_type_id, user_id)
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

func (c *CoaController) GetAllLimiChildByCompany() {
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	user_id = 1
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	fmt.Println(id)
	issue_date := strings.TrimSpace(c.GetString("issue_date"))
	keyword := strings.TrimSpace(c.GetString("keyword"))
	account_type_id, _ := c.GetInt("account_type_id")
	sales_type_id, _ := c.GetInt("sales_type_id")

	if issue_date == "" {
		issue_date = utils.GetSvrDate().Format("2006-01-02")
	}

	d, err := t_coa.GetAllLimitChild(issue_date, keyword, "", 1, account_type_id, sales_type_id, user_id)
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

func (c *CoaController) GetAllLimiChildByCompanyAssets() {
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	user_id = 1

	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	fmt.Println(id)
	issue_date := strings.TrimSpace(c.GetString("issue_date"))
	keyword := strings.TrimSpace(c.GetString("keyword"))
	sales_type_id, _ := c.GetInt("sales_type_id")
	account_type_id, _ := c.GetInt("account_type_id")

	if issue_date == "" {
		issue_date = utils.GetSvrDate().Format("2006-01-02")
	}

	d, err := t_coa.GetAllLimitChild(issue_date, keyword, "Assets", 1, account_type_id, sales_type_id, user_id)
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
