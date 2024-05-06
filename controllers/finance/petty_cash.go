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

type PettyCashHController struct {
	base.BaseController
}

func (c *PettyCashHController) Prepare() {
	c.Ctx.Request.Header.Set("token", "No Aut")
	c.BaseController.Prepare()
}

func (c *PettyCashHController) Post() {
	var err error
	var deletedat, reference_no, voucher_code, period, batch_no string
	var num int
	var user_name string
	var user_id, form_id int

	fmt.Print("Check :", user_id, form_id, "..")
	sess := c.GetSession("profile")
	if sess != nil {
		user_name = sess.(map[string]interface{})["username"].(string)
		user_id = sess.(map[string]interface{})["id"].(int)
	}

	issue_date := strings.TrimSpace(c.GetString("issue_date"))
	company_id, _ := c.GetInt("company_id")
	account_id, _ := c.GetInt("account_id")
	debet, _ := c.GetFloat("debet")
	credit, _ := c.GetFloat("credit")
	transaction_type := strings.TrimSpace(c.GetString("transaction_type"))
	pic := strings.TrimSpace(c.GetString("pic"))
	memo := strings.TrimSpace(c.GetString("memo"))

	valid := validation.Validation{}
	valid.Required(issue_date, "issue_date").Message("Issue date is required")
	valid.Required(company_id, "company_id").Message("Company is required")
	valid.Required(account_id, "account_id").Message("Account is required")
	valid.Required(transaction_type, "transaction_type").Message("Transaction type is required")

	if debet+credit == 0 {
		if transaction_type == "In" && debet == 0 {
			valid.AddError("debet", "Debet is required")
		} else if transaction_type == "Out" && credit == 0 {
			valid.AddError("credit", "Credit is required")
		}
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

	var companies models.Company
	if err = models.Companies().Filter("id", company_id).Filter("CompanyTypes__TypeId__Id", 1).One(&companies); err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Company unregistered/Illegal data")
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
	err = models.ChartOfAccounts().Filter("id", account_id).One(&coa)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Account id unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	deletedat = coa.DeletedAt.Format("2006-01-02")
	if deletedat != "0001-01-01" {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Invalid data", map[string]interface{}{"account_id": "'" + coa.NameCoa + "' has been deleted"})
		c.ServeJSON()
		return
	}

	issuedate, errDate := time.Parse("2006-01-02", issue_date)
	if errDate != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, errDate.Error())
		c.ServeJSON()
		return
	}

	if utils.ToUpper(transaction_type) == "IN" {
		voucher_code = coa.CodeIn
	} else if utils.ToUpper(transaction_type) == "OUT" {
		voucher_code = coa.CodeOut
	}

	num, reference_no = models.GeneratePettyCashNumber(issuedate, company_id, account_id, companies.Code, voucher_code, transaction_type)
	period = string(issuedate.Format("20060102"))

	batch_no = string(issuedate.Format("200601")) + reference_no
	t_pettycashh = models.PettyCashHeader{
		IssueDate:       issuedate,
		CompanyId:       company_id,
		CompanyCode:     companies.Code,
		CompanyName:     companies.Name,
		AccountId:       account_id,
		AccountCode:     coa.CodeCoa,
		AccountName:     coa.NameCoa,
		VoucherSeqNo:    num,
		VoucherCode:     voucher_code,
		VoucherNo:       reference_no,
		Debet:           debet,
		Credit:          credit,
		BatchNo:         batch_no,
		TransactionType: transaction_type,
		Period:          utils.String2Int(period),
		Pic:             pic,
		Memo:            memo,
		CreatedBy:       user_name,
	}

	d, err_ := t_pettycashh.Insert(t_pettycashh)

	errcode, errmessage := base.DecodeErr(err_)

	if err_ != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
	} else {
		v, _ := t_pettycashh.GetById(d.Id, user_id)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, errcode, errmessage, v)
	}
	c.ServeJSON()
	// fmt.Print("check :", utils.ToUpper(transaction_type) == "IN", transaction_type, "..")
	// c.Data["json"] = t_pettycashh
	// c.ServeJSON()
}

func (c *PettyCashHController) Put() {
	var err error
	var deletedat, transaction_type, voucher_no string
	var company_id, account_id int
	var user_name string
	var user_id, form_id int
	fmt.Print("Check :", user_id, form_id, "..")
	sess := c.GetSession("profile")
	if sess != nil {
		user_name = sess.(map[string]interface{})["username"].(string)
		user_id = sess.(map[string]interface{})["id"].(int)
	}

	// form_id = base.FormName("petty_cash")
	// update_aut := models.CheckPrivileges(user_id, form_id, 3)
	// if !update_aut {
	// 	c.Ctx.ResponseWriter.WriteHeader(402)
	// 	utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Update not authorize", map[string]interface{}{"message": "Please contact administrator"})
	// 	c.ServeJSON()
	// 	return
	// }

	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	debet, _ := c.GetFloat("debet")
	credit, _ := c.GetFloat("credit")
	issue_date := strings.TrimSpace(c.GetString("issue_date"))
	pic := strings.TrimSpace(c.GetString("pic"))
	memo := strings.TrimSpace(c.GetString("memo"))

	var querydata models.PettyCashHeader
	err = models.PettyCashHeaders().Filter("id", id).One(&querydata)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, "Invoice id unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	var deletedatData = querydata.DeletedAt.Format("2006-01-02")
	if deletedatData != "0001-01-01" {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("'%v' has been deleted", querydata.VoucherNo))
		c.ServeJSON()
		return
	}

	if querydata.StatusId > 0 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("'%v' has been CLOSED", querydata.VoucherNo))
		c.ServeJSON()
		return
	}

	if querydata.StatusGlId > 0 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("'%v' has been POSTED", querydata.VoucherNo))
		c.ServeJSON()
		return
	}

	if querydata.LoanId > 0 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Unable to edit", map[string]interface{}{"id": "'" + querydata.VoucherNo + "' has been CLAIM by loan '" + querydata.LoanReferenceNo + "'"})
		c.ServeJSON()
		return
	}

	if querydata.ArId > 0 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Unable to edit", map[string]interface{}{"id": "'" + querydata.VoucherNo + "' has been CLAIM by ar '" + querydata.ArReferenceNo + "'"})
		c.ServeJSON()
		return
	}

	if querydata.ApId > 0 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Unable to edit", map[string]interface{}{"id": "'" + querydata.VoucherNo + "' has been CLAIM by ap '" + querydata.ApReferenceNo + "'"})
		c.ServeJSON()
		return
	}

	thedate, errdate := time.Parse("2006-01-02", issue_date)
	if errdate != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, errdate.Error())
		c.ServeJSON()
		return
	}

	if thedate.Month() != querydata.IssueDate.Month() || thedate.Year() != querydata.IssueDate.Year() {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, "Allowed changes date part only")
		c.ServeJSON()
		return
	}

	company_id = querydata.CompanyId
	account_id = querydata.AccountId
	transaction_type = querydata.TransactionType
	voucher_no = querydata.VoucherNo

	valid := validation.Validation{}
	valid.Required(issue_date, "issue_date").Message("Issue date is required")
	valid.Required(company_id, "company_id").Message("Company is required")
	valid.Required(account_id, "account_id").Message("Account is required")
	valid.Required(transaction_type, "transaction_type").Message("Transaction type is required")
	if debet+credit == 0 {
		if transaction_type == "In" && debet == 0 {
			valid.AddError("debet", "Debet is required")
		} else if transaction_type == "Out" && credit == 0 {
			valid.AddError("credit", "Credit is required")
		}
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

	var companies models.Company
	if err = models.Companies().Filter("id", company_id).Filter("CompanyTypes__TypeId__Id", 1).One(&companies); err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Company unregistered/Illegal data")
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
	err = models.ChartOfAccounts().Filter("id", account_id).One(&coa)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Account id unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	deletedat = coa.DeletedAt.Format("2006-01-02")
	if deletedat != "0001-01-01" {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Invalid data", map[string]interface{}{"account_id": "'" + coa.NameCoa + "' has been deleted"})
		c.ServeJSON()
		return
	}

	t_pettycashh.Id = id
	t_pettycashh.IssueDate = thedate
	t_pettycashh.CompanyId = company_id
	t_pettycashh.CompanyCode = companies.Code
	t_pettycashh.CompanyName = companies.Name
	t_pettycashh.AccountId = account_id
	t_pettycashh.AccountCode = coa.CodeCoa
	t_pettycashh.AccountName = coa.NameCoa
	t_pettycashh.VoucherSeqNo = querydata.VoucherSeqNo
	t_pettycashh.VoucherCode = querydata.VoucherCode
	t_pettycashh.VoucherNo = voucher_no
	t_pettycashh.Debet = debet
	t_pettycashh.Credit = credit
	t_pettycashh.BatchNo = querydata.BatchNo
	t_pettycashh.TransactionType = transaction_type
	t_pettycashh.Pic = pic
	t_pettycashh.Memo = memo
	t_pettycashh.ArId = querydata.ArId
	t_pettycashh.ArReferenceNo = querydata.ArReferenceNo
	t_pettycashh.ApId = querydata.ApId
	t_pettycashh.ApReferenceNo = querydata.ApReferenceNo
	t_pettycashh.LoanId = querydata.LoanId
	t_pettycashh.LoanReferenceNo = querydata.LoanReferenceNo
	t_pettycashh.Period = querydata.Period
	t_pettycashh.CreatedBy = querydata.CreatedBy
	t_pettycashh.UpdatedBy = user_name

	err_ := t_pettycashh.Update()
	errcode, errmessage := base.DecodeErr(err_)

	if err_ != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
	} else {
		v, _ := t_pettycashh.GetById(id, user_id)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, errcode, errmessage, v)
	}
	c.ServeJSON()
}

func (c *PettyCashHController) GetAll() {
	var user_id int
	var issuedate *string
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	currentPage, _ := c.GetInt("page")
	if currentPage == 0 {
		currentPage = 1
	}

	pageSize, _ := c.GetInt("pagesize")
	if pageSize == 0 {
		pageSize = 10
	}
	keyword := strings.TrimSpace(c.GetString("keyword"))
	issue_date := strings.TrimSpace(c.GetString("issue_date"))
	status_id, _ := c.GetInt("status_id")
	status_gl_id, _ := c.GetInt("status_gl_id")
	company_id, _ := c.GetInt("status_gl_id")
	account_id, _ := c.GetInt("status_gl_id")
	match_mode := strings.TrimSpace(c.GetString("match_mode"))
	value_name := strings.TrimSpace(c.GetString("value_name"))
	field_name := strings.TrimSpace(c.GetString("field_name"))
	if issue_date == "" {
		issuedate = nil
	} else {
		issuedate = &issue_date
	}
	d, err := t_pettycashh.GetAll(keyword, currentPage, pageSize, issuedate, field_name, match_mode, value_name, user_id, status_id, status_gl_id, company_id, account_id)
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

func (c *PettyCashHController) GetOne() {
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}

	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	v, err := t_pettycashh.GetById(id, user_id)
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

func (c *PettyCashHController) Delete() {
	var err error
	var deletedat string
	var user_name string
	var user_id, form_id int
	fmt.Print("Check :", user_id, form_id, user_name, "..")

	sess := c.GetSession("profile")
	if sess != nil {
		user_name = sess.(map[string]interface{})["username"].(string)
		user_id = sess.(map[string]interface{})["id"].(int)
	}

	// form_id = base.FormName("petty_cash")
	// update_aut := models.CheckPrivileges(user_id, form_id, 3)
	// if !update_aut {
	// 	c.Ctx.ResponseWriter.WriteHeader(402)
	// 	utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Update not authorize", map[string]interface{}{"message": "Please contact administrator"})
	// 	c.ServeJSON()
	// 	return
	// }

	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	var querydata models.PettyCashHeader
	err = models.PettyCashHeaders().Filter("id", id).One(&querydata)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, "petty cash id unregistered/Illegal data")
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
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("voucher_no :'%v' has been deleted", querydata.VoucherNo))
		c.ServeJSON()
		return
	}

	if querydata.StatusGlId > 0 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("voucher_no :'%v' has been POSTED", querydata.VoucherNo))
		c.ServeJSON()
		return
	}

	if querydata.LoanId > 0 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Invalid data", map[string]interface{}{"id": "'" + querydata.VoucherNo + "' has been CLAIM by loan '" + querydata.LoanReferenceNo + "'"})
		c.ServeJSON()
		return
	}

	if querydata.ArId > 0 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Invalid data", map[string]interface{}{"id": "'" + querydata.VoucherNo + "' has been CLAIM by ar '" + querydata.ArReferenceNo + "'"})
		c.ServeJSON()
		return
	}

	if querydata.ApId > 0 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Invalid data", map[string]interface{}{"id": "'" + querydata.VoucherNo + "' has been CLAIM by ap '" + querydata.ApReferenceNo + "'"})
		c.ServeJSON()
		return
	}

	models.PettyCashHeaders().Filter("id", id).Update(orm.Params{"deleted_at": utils.GetSvrDate(), "deleted_by": user_name})
	models.PettyCashs().Filter("voucher_id", id).Update(orm.Params{"deleted_at": utils.GetSvrDate(), "deleted_by": user_name})

	utils.ReturnHTTPError(&c.Controller, 200, "soft delete success")
	c.ServeJSON()
}

func (c *PettyCashHController) GetAllList() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	transaction_type := strings.TrimSpace(c.GetString("transaction_type"))
	issue_date := strings.TrimSpace(c.GetString("issue_date"))
	keyword := strings.TrimSpace(c.GetString("keyword"))
	d, err := t_pettycashh.GetAllList(id, issue_date, utils.ToUpper(transaction_type), keyword)
	code, message := base.DecodeErr(err)
	if err == orm.ErrNoRows {
		code = 200
		c.Ctx.ResponseWriter.WriteHeader(code)
		utils.ReturnHTTPError(&c.Controller, code, "No data")
	} else if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(code)
		utils.ReturnHTTPError(&c.Controller, code, message)
	} else {

		utils.ReturnHTTPSuccessWithMessage(&c.Controller, code, message, d)
	}
	c.ServeJSON()
}

func (c *PettyCashHController) ReOrderNum() {
	var user_name string
	var user_id, form_id int
	fmt.Print("Check :", user_id, form_id, user_name, "..")
	sess := c.GetSession("profile")
	if sess != nil {
		user_name = sess.(map[string]interface{})["username"].(string)
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	action_status, _ := c.GetInt("action_status")
	account_id, _ := c.GetInt("account_id")
	month_id, _ := c.GetInt("month_id")
	year_id, _ := c.GetInt("year_id")
	ids := strings.TrimSpace(c.GetString("ids"))

	valid := validation.Validation{}
	valid.Required(account_id, "account_id").Message("is required")
	valid.Required(ids, "ids").Message("is required")
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

	if year_id == 0 {
		year_id = utils.GetSvrDate().Year()
	}

	if month_id == 0 {
		month_id = int(utils.GetSvrDate().Month())
	}

	d, err := t_pettycashh.ReOrderNum(year_id, month_id, account_id, action_status, ids, user_name)
	code, message := base.DecodeErr(err)
	if err == orm.ErrNoRows {
		code = 200
		c.Ctx.ResponseWriter.WriteHeader(code)
		utils.ReturnHTTPError(&c.Controller, code, "No data")
	} else if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(code)
		utils.ReturnHTTPError(&c.Controller, code, message)
	} else {

		utils.ReturnHTTPSuccessWithMessage(&c.Controller, code, message, d)
	}
	c.ServeJSON()
}

func (c *PettyCashHController) ReOrderNumList() {
	var user_name string
	var user_id, form_id int
	fmt.Print("Check :", user_id, form_id, user_name, "..")
	sess := c.GetSession("profile")
	if sess != nil {
		user_name = sess.(map[string]interface{})["username"].(string)
		user_id = sess.(map[string]interface{})["id"].(int)
	}

	account_id, _ := c.GetInt("account_id")
	month_id, _ := c.GetInt("month_id")
	year_id, _ := c.GetInt("year_id")
	keyword := strings.TrimSpace(c.GetString("keyword"))
	match_mode := strings.TrimSpace(c.GetString("match_mode"))
	value_name := strings.TrimSpace(c.GetString("value_name"))
	field_name := strings.TrimSpace(c.GetString("field_name"))

	valid := validation.Validation{}
	valid.Required(account_id, "account_id").Message("is required")
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

	if year_id == 0 {
		year_id = utils.GetSvrDate().Year()
	}

	if month_id == 0 {
		month_id = int(utils.GetSvrDate().Month())
	}

	d, err := t_pettycashh.ReOrderNumList(keyword, field_name, match_mode, value_name, 0, 0, year_id, month_id, account_id, user_id)
	code, message := base.DecodeErr(err)
	if err == orm.ErrNoRows {
		code = 200
		c.Ctx.ResponseWriter.WriteHeader(code)
		utils.ReturnHTTPError(&c.Controller, code, "No data")
	} else if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(code)
		utils.ReturnHTTPError(&c.Controller, code, message)
	} else {

		utils.ReturnHTTPSuccessWithMessage(&c.Controller, code, message, d)
	}
	c.ServeJSON()
}
