package controllers

import (
	"strconv"
	"strings"

	base "mikiwa/controllers"
	"mikiwa/utils"

	"github.com/beego/beego/v2/client/orm"
)

type FinanceReportController struct {
	base.BaseController
}

func (c *FinanceReportController) Prepare() {
	c.Ctx.Request.Header.Set("token", "No Aut")
	c.BaseController.Prepare()
}

func (c *FinanceReportController) ReportPettyCashSummaryDaily() {
	var issuedate, issuedate2 *string
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	report_grup := 1
	report_type := 1
	keyword := strings.TrimSpace(c.GetString("keyword"))
	match_mode := strings.TrimSpace(c.GetString("match_mode"))
	value_name := strings.TrimSpace(c.GetString("value_name"))
	field_name := strings.TrimSpace(c.GetString("field_name"))
	field_name_top := strings.TrimSpace(c.GetString("field_name_top"))
	issue_date := strings.TrimSpace(c.GetString("issue_date"))
	issue_date2 := strings.TrimSpace(c.GetString("issue_date2"))
	status_id, _ := c.GetInt("status_id")
	status_gl_id, _ := c.GetInt("status_gl_id")
	company_id, _ := c.GetInt("company_id")
	account_id, _ := c.GetInt("account_id")
	sales_type_id, _ := c.GetInt("sales_type_id")
	search_detail, _ := c.GetInt("search_detail")
	currentPage, _ := c.GetInt("page")

	if currentPage == 0 {
		currentPage = 1
	}
	pageSize, _ := c.GetInt("pagesize")
	if pageSize == 0 {
		pageSize = 10
	}
	allsize, _ := c.GetInt("allsize")
	if issue_date == "" {
		currdate := utils.GetSvrDate().Format("2006-01-02")
		issuedate = &currdate
	} else {
		issuedate = &issue_date
	}
	if issue_date2 == "" {
		issuedate2 = issuedate
	} else {
		issuedate2 = &issue_date2
	}

	d, _, err := t_pettycashh.GetPettyCashHeader(keyword, field_name, match_mode, value_name, currentPage, pageSize, allsize, issuedate, issuedate2, nil, company_id, account_id, sales_type_id, status_id, status_gl_id, user_id, report_grup, report_type, search_detail, field_name_top)
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

func (c *FinanceReportController) ReportPettyCashSummaryMonthly() {
	var issuedate, issuedate2 *string
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	report_grup := 1
	report_type := 2
	keyword := strings.TrimSpace(c.GetString("keyword"))
	match_mode := strings.TrimSpace(c.GetString("match_mode"))
	value_name := strings.TrimSpace(c.GetString("value_name"))
	field_name := strings.TrimSpace(c.GetString("field_name"))
	field_name_top := strings.TrimSpace(c.GetString("field_name_top"))
	issue_date := strings.TrimSpace(c.GetString("issue_date"))
	issue_date2 := strings.TrimSpace(c.GetString("issue_date2"))
	status_id, _ := c.GetInt("status_id")
	status_gl_id, _ := c.GetInt("status_gl_id")
	company_id, _ := c.GetInt("company_id")
	account_id, _ := c.GetInt("account_id")
	sales_type_id, _ := c.GetInt("sales_type_id")
	search_detail, _ := c.GetInt("search_detail")
	currentPage, _ := c.GetInt("page")

	if currentPage == 0 {
		currentPage = 1
	}
	pageSize, _ := c.GetInt("pagesize")
	if pageSize == 0 {
		pageSize = 10
	}
	allsize, _ := c.GetInt("allsize")
	if issue_date == "" {
		currdate := utils.GetSvrDate().Format("2006-01-02")
		issuedate = &currdate
	} else {
		issuedate = &issue_date
	}
	if issue_date2 == "" {
		issuedate2 = issuedate
	} else {
		issuedate2 = &issue_date2
	}

	d, _, err := t_pettycashh.GetPettyCashHeader(keyword, field_name, match_mode, value_name, currentPage, pageSize, allsize, issuedate, issuedate2, nil, company_id, account_id, sales_type_id, status_id, status_gl_id, user_id, report_grup, report_type, search_detail, field_name_top)
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

func (c *FinanceReportController) ReportPettyCashSummaryYearly() {
	var issuedate, issuedate2 *string
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	report_grup := 1
	report_type := 3
	keyword := strings.TrimSpace(c.GetString("keyword"))
	match_mode := strings.TrimSpace(c.GetString("match_mode"))
	value_name := strings.TrimSpace(c.GetString("value_name"))
	field_name := strings.TrimSpace(c.GetString("field_name"))
	field_name_top := strings.TrimSpace(c.GetString("field_name_top"))
	issue_date := strings.TrimSpace(c.GetString("issue_date"))
	issue_date2 := strings.TrimSpace(c.GetString("issue_date2"))
	status_id, _ := c.GetInt("status_id")
	status_gl_id, _ := c.GetInt("status_gl_id")
	company_id, _ := c.GetInt("company_id")
	account_id, _ := c.GetInt("account_id")
	sales_type_id, _ := c.GetInt("sales_type_id")
	search_detail, _ := c.GetInt("search_detail")
	currentPage, _ := c.GetInt("page")

	if currentPage == 0 {
		currentPage = 1
	}
	pageSize, _ := c.GetInt("pagesize")
	if pageSize == 0 {
		pageSize = 10
	}
	allsize, _ := c.GetInt("allsize")
	if issue_date == "" {
		currdate := utils.GetSvrDate().Format("2006-01-02")
		issuedate = &currdate
	} else {
		issuedate = &issue_date
	}
	if issue_date2 == "" {
		issuedate2 = issuedate
	} else {
		issuedate2 = &issue_date2
	}

	d, _, err := t_pettycashh.GetPettyCashHeader(keyword, field_name, match_mode, value_name, currentPage, pageSize, allsize, issuedate, issuedate2, nil, company_id, account_id, sales_type_id, status_id, status_gl_id, user_id, report_grup, report_type, search_detail, field_name_top)
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

func (c *FinanceReportController) ReportPettyCashDaily() {
	var issuedate, issuedate2 *string
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	report_grup := 0
	report_type := 1
	keyword := strings.TrimSpace(c.GetString("keyword"))
	match_mode := strings.TrimSpace(c.GetString("match_mode"))
	value_name := strings.TrimSpace(c.GetString("value_name"))
	field_name := strings.TrimSpace(c.GetString("field_name"))
	field_name_top := strings.TrimSpace(c.GetString("field_name_top"))
	issue_date := strings.TrimSpace(c.GetString("issue_date"))
	issue_date2 := strings.TrimSpace(c.GetString("issue_date2"))
	status_id, _ := c.GetInt("status_id")
	status_gl_id, _ := c.GetInt("status_gl_id")
	company_id, _ := c.GetInt("company_id")
	account_id, _ := c.GetInt("account_id")
	sales_type_id, _ := c.GetInt("sales_type_id")
	search_detail, _ := c.GetInt("search_detail")
	currentPage, _ := c.GetInt("page")

	if currentPage == 0 {
		currentPage = 1
	}
	pageSize, _ := c.GetInt("pagesize")
	if pageSize == 0 {
		pageSize = 10
	}
	allsize, _ := c.GetInt("allsize")

	if issue_date == "" {
		currdate := utils.GetSvrDate().Format("2006-01-02")
		issuedate = &currdate
	} else {
		issuedate = &issue_date
	}
	if issue_date2 == "" {
		issuedate2 = issuedate
	} else {
		issuedate2 = &issue_date2
	}

	d, _, err := t_pettycashh.GetPettyCash(keyword, field_name, match_mode, value_name, currentPage, pageSize, allsize, issuedate, issuedate2, nil, company_id, account_id, sales_type_id, status_id, status_gl_id, user_id, report_grup, report_type, search_detail, field_name_top)
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

func (c *FinanceReportController) ReportPettyCashMonthly() {
	var issuedate, issuedate2 *string
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	report_grup := 0
	report_type := 2
	keyword := strings.TrimSpace(c.GetString("keyword"))
	match_mode := strings.TrimSpace(c.GetString("match_mode"))
	value_name := strings.TrimSpace(c.GetString("value_name"))
	field_name := strings.TrimSpace(c.GetString("field_name"))
	field_name_top := strings.TrimSpace(c.GetString("field_name_top"))
	issue_date := strings.TrimSpace(c.GetString("issue_date"))
	issue_date2 := strings.TrimSpace(c.GetString("issue_date2"))
	status_id, _ := c.GetInt("status_id")
	status_gl_id, _ := c.GetInt("status_gl_id")
	company_id, _ := c.GetInt("company_id")
	account_id, _ := c.GetInt("account_id")
	sales_type_id, _ := c.GetInt("sales_type_id")
	search_detail, _ := c.GetInt("search_detail")
	currentPage, _ := c.GetInt("page")

	if currentPage == 0 {
		currentPage = 1
	}
	pageSize, _ := c.GetInt("pagesize")
	if pageSize == 0 {
		pageSize = 10
	}
	allsize, _ := c.GetInt("allsize")

	if issue_date == "" {
		currdate := utils.GetSvrDate().Format("2006-01-02")
		issuedate = &currdate
	} else {
		issuedate = &issue_date
	}
	if issue_date2 == "" {
		issuedate2 = issuedate
	} else {
		issuedate2 = &issue_date2
	}

	d, _, err := t_pettycashh.GetPettyCash(keyword, field_name, match_mode, value_name, currentPage, pageSize, allsize, issuedate, issuedate2, nil, company_id, account_id, sales_type_id, status_id, status_gl_id, user_id, report_grup, report_type, search_detail, field_name_top)
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

func (c *FinanceReportController) ReportPettyCashYearly() {
	var issuedate, issuedate2 *string
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	report_grup := 0
	report_type := 3
	keyword := strings.TrimSpace(c.GetString("keyword"))
	match_mode := strings.TrimSpace(c.GetString("match_mode"))
	value_name := strings.TrimSpace(c.GetString("value_name"))
	field_name := strings.TrimSpace(c.GetString("field_name"))
	field_name_top := strings.TrimSpace(c.GetString("field_name_top"))
	issue_date := strings.TrimSpace(c.GetString("issue_date"))
	issue_date2 := strings.TrimSpace(c.GetString("issue_date2"))
	status_id, _ := c.GetInt("status_id")
	status_gl_id, _ := c.GetInt("status_gl_id")
	company_id, _ := c.GetInt("company_id")
	account_id, _ := c.GetInt("account_id")
	sales_type_id, _ := c.GetInt("sales_type_id")
	search_detail, _ := c.GetInt("search_detail")
	currentPage, _ := c.GetInt("page")

	if currentPage == 0 {
		currentPage = 1
	}
	pageSize, _ := c.GetInt("pagesize")
	if pageSize == 0 {
		pageSize = 10
	}
	allsize, _ := c.GetInt("allsize")

	if issue_date == "" {
		currdate := utils.GetSvrDate().Format("2006-01-02")
		issuedate = &currdate
	} else {
		issuedate = &issue_date
	}
	if issue_date2 == "" {
		issuedate2 = issuedate
	} else {
		issuedate2 = &issue_date2
	}

	d, _, err := t_pettycashh.GetPettyCash(keyword, field_name, match_mode, value_name, currentPage, pageSize, allsize, issuedate, issuedate2, nil, company_id, account_id, sales_type_id, status_id, status_gl_id, user_id, report_grup, report_type, search_detail, field_name_top)
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

func (c *FinanceReportController) ReportVoucher() {
	var issuedate, issuedate2 *string
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	report_grup := 0
	report_type := 1
	keyword := strings.TrimSpace(c.GetString("keyword"))
	match_mode := strings.TrimSpace(c.GetString("match_mode"))
	value_name := strings.TrimSpace(c.GetString("value_name"))
	field_name := strings.TrimSpace(c.GetString("field_name"))
	field_name_top := strings.TrimSpace(c.GetString("field_name_top"))
	issue_date := strings.TrimSpace(c.GetString("issue_date"))
	issue_date2 := strings.TrimSpace(c.GetString("issue_date2"))
	status_id, _ := c.GetInt("status_id")
	status_gl_id, _ := c.GetInt("status_gl_id")
	company_id, _ := c.GetInt("company_id")
	account_id, _ := c.GetInt("account_id")
	sales_type_id, _ := c.GetInt("sales_type_id")
	search_detail, _ := c.GetInt("search_detail")
	currentPage, _ := c.GetInt("page")

	if currentPage == 0 {
		currentPage = 1
	}
	pageSize, _ := c.GetInt("pagesize")
	if pageSize == 0 {
		pageSize = 10
	}
	allsize, _ := c.GetInt("allsize")

	if issue_date == "" {
		currdate := utils.GetSvrDate().Format("2006-01-02")
		issuedate = &currdate
	} else {
		issuedate = &issue_date
	}
	if issue_date2 == "" {
		issuedate2 = issuedate
	} else {
		issuedate2 = &issue_date2
	}

	d, err := t_pettycashh.GetVoucher(keyword, field_name, match_mode, value_name, currentPage, pageSize, allsize, issuedate, issuedate2, id, company_id, account_id, sales_type_id, status_id, status_gl_id, user_id, report_grup, report_type, search_detail, field_name_top)
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
