package controllers

import (
	"encoding/json"
	"fmt"
	base "mikiwa/controllers"
	"mikiwa/models"
	"mikiwa/utils"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
)

type PettyCashV2Controller struct {
	base.BaseController
}

func (c *PettyCashV2Controller) Prepare() {
	c.Ctx.Request.Header.Set("token", "No Aut")
	c.BaseController.Prepare()
}

type (
	InputHeaderPettyCash struct {
		IssueDate       string                 `json:"issue_date"`
		CompanyId       int                    `json:"company_id"`
		AccountId       int                    `json:"account_id"`
		TransactionType string                 `json:"transaction_type"`
		VoucherNo       string                 `json:"voucher_no"`
		Pic             string                 `json:"pic"`
		Memo            string                 `json:"memo"`
		StatusIds       string                 `json:"status_ids"`
		UploadFile      models.DocumentList    `json:"upload_file"`
		Detail          []InputDetailPettyCash `json:"detail"`
	}

	InputDetailPettyCash struct {
		Id          int     `json:"id"`
		AccountId   int     `json:"account_id"`
		AccountCode string  `json:"account_code"`
		AccountName string  `json:"account_name"`
		Debet       float64 `json:"debet"`
		Credit      float64 `json:"credit"`
		Pic         string  `json:"pic"`
		Memo        string  `json:"memo"`
		CreatedBy   string  `json:"created_by"`
	}

	coaUserRtn struct {
		Id          int    `json:"id"`
		CodeCoa     string `json:"code_coa"`
		NameCoa     string `json:"name_coa"`
		CodeIn      string `json:"code_in"`
		CodeOut     string `json:"code_out"`
		StatusId    int8   `json:"status_id"`
		CompanyId   int    `json:"company_id"`
		CompanyCode string `json:"company_code"`
		UserId      int    `json:"user_id"`
	}
)

func (c *PettyCashV2Controller) Post() {
	o := orm.NewOrm()
	var user_id, form_id int
	var user_name string
	var err error
	var folderName string = "petty_cash"
	var num, voucher_seq_no int = 0, 0
	var status_id, status_gl_id int8 = 0, 0
	var reference_no, voucher_code, period, batch_no string
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
		user_name = sess.(map[string]interface{})["username"].(string)
	}

	form_id = base.FormName(form_petty_cash)
	write_aut := models.CheckPrivileges(user_id, form_id, base.Write)
	aut_aut := models.CheckPrivileges(user_id, form_id, base.Author)
	app_aut := models.CheckPrivileges(user_id, form_id, base.Pending)
	close_aut := models.CheckPrivileges(user_id, form_id, base.Approval)
	if !write_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Post not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}

	var i int = 0
	var ob InputHeaderPettyCash
	var inputDetail []models.PettyCash

	body := c.Ctx.Input.RequestBody
	err = json.Unmarshal(body, &ob)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	ob.CompanyId = 1

	valid := validation.Validation{}
	valid.Required(strings.TrimSpace(ob.IssueDate), "issue_date").Message("Is required")
	valid.Required(ob.AccountId, "account_id").Message("Is required")
	valid.Required(ob.TransactionType, "transaction_type").Message("Is required")

	if len(ob.Detail) == 0 {
		valid.AddError("detail", "Detail list is required")
	}

	if utils.ToUpper(ob.TransactionType) != "MASUK" && utils.ToUpper(ob.TransactionType) != "KELUAR" && utils.ToUpper(ob.TransactionType) != "IN" && utils.ToUpper(ob.TransactionType) != "OUT" {
		valid.AddError("transaction_type", "Invalid transaction type")
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
	if err = models.Companies().Filter("id", ob.CompanyId).Filter("CompanyTypes__TypeId__Id", base.Internal).Filter("deleted_at__isnull", true).One(&companies); err == orm.ErrNoRows {
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

	var coa coaUserRtn
	if !aut_aut {
		err = o.Raw("select id,code_coa,name_coa,code_in,code_out,status_id,company_id,company_code,user_id from chart_of_accounts t0 left join (select user_id,account_id from sys_user_account where  user_id = " + utils.Int2String(user_id) + " ) t1 on t1.account_id = t0.id where deleted_at is null  and id not in (select parent_id from chart_of_accounts where deleted_at is null)  and is_header = 1 and t0.id = " + utils.Int2String(ob.AccountId)).QueryRow(&coa)
	} else {
		err = o.Raw("select id,code_coa,name_coa,code_in,code_out,status_id,company_id,company_code," + utils.Int2String(user_id) + " user_id from chart_of_accounts where deleted_at is null  and id not in (select parent_id from chart_of_accounts where deleted_at is null)  and is_header = 1 and id = " + utils.Int2String(ob.AccountId)).QueryRow(&coa)
	}
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Chart of accounts unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	if coa.UserId != user_id {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprintf("Error '%v' is not your account", coa.CodeCoa))
		c.ServeJSON()
		return
	}

	if coa.StatusId == 0 {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprintf("Error account '%v' has been set as INACTIVE", coa.CodeCoa))
		c.ServeJSON()
		return
	}

	if coa.CompanyId != ob.CompanyId {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprintf("Error account '%v' is not under company '%s'", coa.CodeCoa, companies.Code))
		c.ServeJSON()
		return
	}

	var coaDetail coaUserRtn
	var debet, credit float64
	wg = new(sync.WaitGroup)
	var mutex sync.Mutex
	resultChan := make(chan utils.ResultChan, len(ob.Detail))
	var queryResults []utils.ResultChan
	wg.Add(len(ob.Detail))
	for _, v := range ob.Detail {
		go func(v InputDetailPettyCash) {
			defer wg.Done()
			mutex.Lock()
			defer mutex.Unlock()
			if !aut_aut {
				err = o.Raw("select id,code_coa,name_coa,code_in,code_out,status_id,user_id from chart_of_accounts  t0 left join (select user_id,account_id from sys_user_account where  user_id = " + utils.Int2String(user_id) + " ) t1 on t1.account_id = t0.id where deleted_at is null  and id not in (select parent_id from chart_of_accounts where deleted_at is null) and  t0.id = " + utils.Int2String(v.AccountId)).QueryRow(&coaDetail)
			} else {
				err = o.Raw("select id,code_coa,name_coa,code_in,code_out,status_id," + utils.Int2String(user_id) + " user_id from chart_of_accounts where deleted_at is null  and id not in (select parent_id from chart_of_accounts where deleted_at is null)  and id = " + utils.Int2String(v.AccountId)).QueryRow(&coaDetail)
			}
			if err == orm.ErrNoRows {
				resultChan <- utils.ResultChan{Id: v.AccountId, Data: "Invalid detail id", Message: "chart of account unregistered/Illegal data"}
				return
			}

			if err != nil {
				resultChan <- utils.ResultChan{Id: v.AccountId, Data: coaDetail.CodeCoa, Message: err.Error()}
				return
			}

			if coaDetail.UserId != user_id {
				resultChan <- utils.ResultChan{Id: v.AccountId, Data: coaDetail.CodeCoa, Message: fmt.Sprintf("'%v' is not your account", coaDetail.CodeCoa)}
				return
			}

			if coaDetail.StatusId == 0 {
				resultChan <- utils.ResultChan{Id: v.AccountId, Data: coaDetail.CodeCoa, Message: fmt.Sprintf("'%v' has been set as INACTIVE", coaDetail.CodeCoa)}
				return
			}

			if coa.CompanyId != ob.CompanyId {
				resultChan <- utils.ResultChan{Id: v.AccountId, Data: coaDetail.CodeCoa, Message: fmt.Sprintf("Error account '%v' is not under company '%s'", coa.CodeCoa, companies.Code)}
				return
			}
		}(v)
		debet += v.Debet
		credit += v.Credit
	}

	// Use goroutine to wait until all goroutines are finished
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		if result.Message != "" {
			queryResults = append(queryResults, utils.ResultChan{
				Id:      result.Id,
				Data:    result.Data,
				Message: result.Message,
			})
		}
	}

	if len(queryResults) != 0 {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 401, "Error", map[string]interface{}{"Invalid detail field": queryResults})
		c.ServeJSON()
		return
	}

	if utils.ToUpper(ob.TransactionType) == "IN" || utils.ToUpper(ob.TransactionType) == "MASUK" {
		voucher_code = coa.CodeIn
	} else if utils.ToUpper(ob.TransactionType) == "OUT" || utils.ToUpper(ob.TransactionType) == "KELUAR" {
		voucher_code = coa.CodeOut
	}

	thedate, errDate := time.Parse("2006-01-02", ob.IssueDate)
	if errDate != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprint("issue_date: ", errDate.Error()))
		c.ServeJSON()
		return
	}

	num, reference_no = models.GeneratePettyCashNumber(thedate, ob.CompanyId, ob.AccountId, companies.Code, voucher_code, ob.TransactionType)
	period = string(thedate.Format("20060102"))
	batch_no = string(thedate.Format("200601")) + reference_no

	if ob.VoucherNo == "" {
		ob.VoucherNo = reference_no
		voucher_seq_no = num
	}

	if ob.StatusIds == "1" {
		if close_aut && app_aut {
			status_id = 1
			status_gl_id = 1
		} else if close_aut && !app_aut {
			status_id = 0
			status_gl_id = 0
		} else if app_aut {
			status_id = 1
			status_gl_id = 0
		} else {
			status_gl_id = 0
			status_id = 0
		}
	}

	t_pettycashh = models.PettyCashHeader{
		IssueDate:       thedate,
		CompanyId:       ob.CompanyId,
		CompanyCode:     companies.Code,
		CompanyName:     companies.Name,
		AccountId:       ob.AccountId,
		AccountCode:     coa.CodeCoa,
		AccountName:     coa.NameCoa,
		VoucherSeqNo:    voucher_seq_no,
		VoucherCode:     voucher_code,
		VoucherNo:       ob.VoucherNo,
		Debet:           debet,
		Credit:          credit,
		BatchNo:         batch_no,
		TransactionType: ob.TransactionType,
		Period:          utils.String2Int(period),
		Pic:             ob.Pic,
		Memo:            ob.Memo,
		StatusId:        status_id,
		StatusGlId:      status_gl_id,
		CreatedBy:       user_name,
		UpdatedBy:       user_name,
	}

	for k, v := range ob.Detail {
		i = 0
		wg.Add(1)
		go func(k int, v InputDetailPettyCash) {
			defer wg.Done()
			mutex.Lock()
			if v.Id == 0 {
				inputDetail = append(inputDetail, models.PettyCash{
					IssueDate:         thedate,
					VoucherId:         t_pettycashh.Id,
					CompanyId:         ob.CompanyId,
					CompanyCode:       companies.Code,
					CompanyName:       companies.Name,
					AccountIdHeader:   ob.AccountId,
					AccountCodeHeader: coa.CodeCoa,
					AccountNameHeader: coa.NameCoa,
					ItemNo:            k + 1,
					AccountId:         v.AccountId,
					AccountCode:       v.AccountCode,
					AccountName:       v.AccountName,
					VoucherSeqNo:      voucher_seq_no,
					VoucherCode:       voucher_code,
					VoucherNo:         ob.VoucherNo,
					Debet:             v.Debet,
					Credit:            v.Credit,
					Pic:               v.Pic,
					Memo:              v.Memo,
					TransactionType:   ob.TransactionType,
					Period:            utils.String2Int(period),
					StatusId:          status_id,
					StatusGlId:        status_gl_id,
					CreatedBy:         user_name,
					UpdatedBy:         user_name,
				})
				i += 1
			}
			mutex.Unlock()
		}(k, v)
	}
	wg.Wait()

	d, err_ := t_pettycashh.InsertWithDetail(t_pettycashh, inputDetail)
	errcode, errmessage = base.DecodeErr(err_)
	if err_ != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
		c.ServeJSON()
		return
	} else {
		if err = base.PostFirebaseRaw(ob.UploadFile, user_name, d.Id, folderName+"/"+utils.Int2String(d.Id), folderName+"/"+utils.Int2String(d.Id), folderName); err != nil {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprint("Error processing data and uploading to Firebase: ", err.Error()))
		} else {
			v, err := t_pettycashh.GetById(d.Id, user_id)
			errcode, errmessage = base.DecodeErr(err)
			if err != nil {
				c.Ctx.ResponseWriter.WriteHeader(errcode)
				utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
			} else {
				utils.ReturnHTTPSuccessWithMessage(&c.Controller, errcode, errmessage, v)
			}
		}
	}

	c.ServeJSON()
}

func (c *PettyCashV2Controller) Put() {
	o := orm.NewOrm()
	var user_id, form_id int
	var user_name string
	var err error
	var folderName string = "petty_cash"
	var num int = 0
	var status_id, status_gl_id int8 = 0, 0
	var reference_no, voucher_code, period, deletedatData, batch_no string
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
		user_name = sess.(map[string]interface{})["username"].(string)
	}

	form_id = base.FormName(form_petty_cash)
	put_aut := models.CheckPrivileges(user_id, form_id, base.Update)
	aut_aut := models.CheckPrivileges(user_id, form_id, base.Author)
	app_aut := models.CheckPrivileges(user_id, form_id, base.Pending)
	close_aut := models.CheckPrivileges(user_id, form_id, base.Approval)

	if !put_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Put not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}

	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	var i int = 0
	var ob InputHeaderPettyCash
	var inputDetail []models.PettyCash
	var putDetail []models.PettyCash

	body := c.Ctx.Input.RequestBody
	err = json.Unmarshal(body, &ob)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	var querydata models.PettyCashHeader
	err = models.PettyCashHeaders().Filter("id", id).One(&querydata)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Voucher id unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	deletedatData = querydata.DeletedAt.Format("2006-01-02")
	if deletedatData != "0001-01-01" {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprintf("'%v' has been deleted", querydata.VoucherNo))
		c.ServeJSON()
		return
	}

	thedate, errDate := time.Parse("2006-01-02", ob.IssueDate)
	if errDate != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, errDate.Error())
		c.ServeJSON()
		return

	}

	if thedate.Month() != querydata.IssueDate.Month() || thedate.Year() != querydata.IssueDate.Year() {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Allowed changes date part only")
		c.ServeJSON()
		return
	}

	if querydata.StatusId == 0 && !app_aut && close_aut {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprintf("'%v'  has not been VERIFIED", querydata.VoucherNo))
		c.ServeJSON()
		return
	}

	if querydata.StatusGlId == 1 && !close_aut {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprintf("'%v'  has been POSTED", querydata.VoucherNo))
		c.ServeJSON()
		return
	}

	if querydata.StatusId == 1 && !app_aut && !close_aut {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprintf("'%v' has been APPROVED", querydata.VoucherNo))
		c.ServeJSON()
		return
	}

	if querydata.LoanId > 0 {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 401, "Unable to edit", map[string]interface{}{"id": "'" + querydata.VoucherNo + "' has been CLAIM by loan '" + querydata.LoanReferenceNo + "'"})
		c.ServeJSON()
		return
	}

	if querydata.ArId > 0 {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 401, "Unable to edit", map[string]interface{}{"id": "'" + querydata.VoucherNo + "' has been CLAIM by ar '" + querydata.ArReferenceNo + "'"})
		c.ServeJSON()
		return
	}

	if querydata.ApId > 0 {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 401, "Unable to edit", map[string]interface{}{"id": "'" + querydata.VoucherNo + "' has been CLAIM by ap '" + querydata.ApReferenceNo + "'"})
		c.ServeJSON()
		return
	}

	ob.CompanyId = querydata.CompanyId
	ob.TransactionType = querydata.TransactionType
	ob.AccountId = querydata.AccountId
	if querydata.VoucherSeqNo != 0 {
		ob.VoucherNo = querydata.VoucherNo
	}
	valid := validation.Validation{}
	valid.Required(strings.TrimSpace(ob.IssueDate), "issue_date").Message("Is required")
	valid.Required(ob.AccountId, "account_id").Message("Is required")
	valid.Required(ob.TransactionType, "transaction_type").Message("Is required")

	if len(ob.Detail) == 0 {
		valid.AddError("detail", "Detail list is required")
	}

	if utils.ToUpper(ob.TransactionType) != "MASUK" && utils.ToUpper(ob.TransactionType) != "KELUAR" && utils.ToUpper(ob.TransactionType) != "IN" && utils.ToUpper(ob.TransactionType) != "OUT" {
		valid.AddError("transaction_type", "Invalid transaction type")
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
	if err = models.Companies().Filter("id", ob.CompanyId).Filter("CompanyTypes__TypeId__Id", base.Internal).Filter("deleted_at__isnull", true).One(&companies); err == orm.ErrNoRows {
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

	var coa coaUserRtn
	if !aut_aut {
		err = o.Raw("select id,code_coa,name_coa,code_in,code_out,status_id,company_id,company_code,user_id from chart_of_accounts t0 left join (select user_id,account_id from sys_user_account where  user_id = " + utils.Int2String(user_id) + " ) t1 on t1.account_id = t0.id where deleted_at is null  and id not in (select parent_id from chart_of_accounts where deleted_at is null)  and is_header = 1 and t0.id = " + utils.Int2String(ob.AccountId)).QueryRow(&coa)
	} else {
		err = o.Raw("select id,code_coa,name_coa,code_in,code_out,status_id,company_id,company_code," + utils.Int2String(user_id) + " user_id from chart_of_accounts where deleted_at is null  and id not in (select parent_id from chart_of_accounts where deleted_at is null)  and is_header = 1 and id = " + utils.Int2String(ob.AccountId)).QueryRow(&coa)
	}
	// err = models.ChartOfAccounts().Filter("deleted_at__isnull", true).Filter("id", ob.AccountId).One(coa)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Chart of accounts unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	if coa.UserId != user_id {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprintf("Error '%v' is not your account", coa.CodeCoa))
		c.ServeJSON()
		return
	}

	if coa.StatusId == 0 {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprintf("Error '%v' has been set as INACTIVE", coa.CodeCoa))
		c.ServeJSON()
		return
	}

	if coa.CompanyId != ob.CompanyId {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprintf("Error account '%v' is not under company '%s'", coa.CodeCoa, companies.Code))
		c.ServeJSON()
		return
	}

	var coaDetail coaUserRtn
	var querydetail models.PettyCash
	var debet, credit float64
	wg = new(sync.WaitGroup)
	var mutex sync.Mutex
	resultChan := make(chan utils.ResultChan, len(ob.Detail))
	var queryResults []utils.ResultChan
	wg.Add(len(ob.Detail))
	for _, k := range ob.Detail {
		go func(v InputDetailPettyCash) {
			defer wg.Done()
			mutex.Lock()
			defer mutex.Unlock()
			if v.Id != 0 {
				if err = models.PettyCashs().Filter("deleted_at__isnull", true).Filter("voucher_id", id).Filter("id", v.Id).One(&querydetail); err == orm.ErrNoRows {
					resultChan <- utils.ResultChan{Id: v.Id, Data: "Invalid detail id", Message: "detail unregistered/Illegal data"}
					return
				}

				if err != nil {
					resultChan <- utils.ResultChan{Id: v.Id, Data: querydetail.AccountCode, Message: err.Error()}
					return
				}
			}
			if !aut_aut {
				err = o.Raw("select id,code_coa,name_coa,code_in,code_out,status_id,user_id from chart_of_accounts  t0 left join (select user_id,account_id from sys_user_account where  user_id = " + utils.Int2String(user_id) + " ) t1 on t1.account_id = t0.id where deleted_at is null  and id not in (select parent_id from chart_of_accounts where deleted_at is null) and  t0.id = " + utils.Int2String(v.AccountId)).QueryRow(&coaDetail)
			} else {
				err = o.Raw("select id,code_coa,name_coa,code_in,code_out,status_id," + utils.Int2String(user_id) + " user_id from chart_of_accounts where deleted_at is null  and id not in (select parent_id from chart_of_accounts where deleted_at is null)  and id = " + utils.Int2String(ob.AccountId)).QueryRow(&coaDetail)
			}
			if err == orm.ErrNoRows {
				resultChan <- utils.ResultChan{Id: v.AccountId, Data: "Invalid detail id", Message: "chart of account unregistered/Illegal data"}
				return
			}

			if err != nil {
				resultChan <- utils.ResultChan{Id: v.AccountId, Data: coaDetail.CodeCoa, Message: err.Error()}
				return
			}

			if coaDetail.UserId != user_id {
				resultChan <- utils.ResultChan{Id: v.AccountId, Data: coaDetail.CodeCoa, Message: fmt.Sprintf("'%v' is not your account", coaDetail.CodeCoa)}
				return
			}

			if coaDetail.StatusId == 0 {
				resultChan <- utils.ResultChan{Id: v.AccountId, Data: coaDetail.CodeCoa, Message: fmt.Sprintf("'%v' has been set as INACTIVE", coaDetail.CodeCoa)}
				return
			}

			if coa.CompanyId != ob.CompanyId {
				resultChan <- utils.ResultChan{Id: v.AccountId, Data: coaDetail.CodeCoa, Message: fmt.Sprintf("Error account '%v' is not under company '%s'", coa.CodeCoa, companies.Code)}
				return
			}
		}(k)
		debet += k.Debet
		credit += k.Credit
	}

	// Use goroutine to wait until all goroutines are finished
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		if result.Message != "" {
			queryResults = append(queryResults, utils.ResultChan{
				Id:      result.Id,
				Data:    result.Data,
				Message: result.Message,
			})
		}
	}

	if len(queryResults) != 0 {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 401, "Error", map[string]interface{}{"Invalid detail field": queryResults})
		c.ServeJSON()
		return
	}

	voucher_code = querydata.VoucherCode
	num = querydata.VoucherSeqNo
	period = string(thedate.Format("20060102"))
	reference_no = querydata.VoucherNo
	batch_no = string(thedate.Format("200601")) + reference_no

	if ob.StatusIds == "1" {
		if close_aut && app_aut {
			status_id = 1
			status_gl_id = 1
		} else if close_aut {
			status_id = querydata.StatusId
			status_gl_id = 1
		} else if app_aut {
			status_id = 1
			status_gl_id = querydata.StatusGlId
		} else {
			status_gl_id = querydata.StatusGlId
			status_id = querydata.StatusId
		}
	} else if ob.StatusIds == "0" {
		if close_aut && app_aut {
			status_id = 0
			status_gl_id = 0
		} else if close_aut {
			status_id = querydata.StatusId
			status_gl_id = 0
		} else if app_aut {
			status_id = 0
			status_gl_id = querydata.StatusGlId
		} else {
			status_gl_id = querydata.StatusGlId
			status_id = querydata.StatusId
		}
	} else {
		status_id = querydata.StatusId
		status_gl_id = querydata.StatusGlId
	}

	t_pettycashh.Id = id
	t_pettycashh.IssueDate = thedate
	t_pettycashh.CompanyId = querydata.CompanyId
	t_pettycashh.CompanyCode = querydata.CompanyCode
	t_pettycashh.CompanyName = querydata.CompanyName
	t_pettycashh.AccountId = ob.AccountId
	t_pettycashh.AccountCode = coa.CodeCoa
	t_pettycashh.AccountName = coa.NameCoa
	t_pettycashh.VoucherSeqNo = num
	t_pettycashh.VoucherCode = voucher_code
	t_pettycashh.VoucherNo = reference_no
	t_pettycashh.Debet = debet
	t_pettycashh.Credit = credit
	t_pettycashh.BatchNo = batch_no
	t_pettycashh.TransactionType = ob.TransactionType
	t_pettycashh.Pic = ob.Pic
	t_pettycashh.Memo = ob.Memo
	t_pettycashh.ArId = querydata.ArId
	t_pettycashh.ArReferenceNo = querydata.ArReferenceNo
	t_pettycashh.ApId = querydata.ApId
	t_pettycashh.ApReferenceNo = querydata.ApReferenceNo
	t_pettycashh.LoanId = querydata.LoanId
	t_pettycashh.LoanReferenceNo = querydata.LoanReferenceNo
	t_pettycashh.Period = utils.String2Int(period)
	t_pettycashh.StatusId = status_id
	t_pettycashh.StatusGlId = status_gl_id
	t_pettycashh.CreatedBy = querydata.CreatedBy
	t_pettycashh.UpdatedBy = user_name

	for k, v := range ob.Detail {
		i = 0
		wg.Add(1)
		go func(k int, v InputDetailPettyCash) {
			defer wg.Done()
			mutex.Lock()
			if v.Id == 0 {
				inputDetail = append(inputDetail, models.PettyCash{
					IssueDate:         thedate,
					VoucherId:         id,
					CompanyId:         ob.CompanyId,
					CompanyCode:       companies.Code,
					CompanyName:       companies.Name,
					AccountIdHeader:   ob.AccountId,
					AccountCodeHeader: coa.CodeCoa,
					AccountNameHeader: coa.NameCoa,
					ItemNo:            k + 1,
					AccountId:         v.AccountId,
					AccountCode:       v.AccountCode,
					AccountName:       v.AccountName,
					VoucherSeqNo:      num,
					VoucherCode:       voucher_code,
					VoucherNo:         reference_no,
					Debet:             v.Debet,
					Credit:            v.Credit,
					Pic:               v.Pic,
					Memo:              v.Memo,
					TransactionType:   ob.TransactionType,
					Period:            utils.String2Int(period),
					StatusId:          status_id,
					StatusGlId:        status_gl_id,
					BatchNo:           batch_no,
					CreatedBy:         user_name,
					UpdatedBy:         user_name,
				})
				i += 1
			} else {
				putDetail = append(putDetail, models.PettyCash{
					Id:                v.Id,
					IssueDate:         thedate,
					VoucherId:         id,
					CompanyId:         ob.CompanyId,
					CompanyCode:       companies.Code,
					CompanyName:       companies.Name,
					AccountIdHeader:   ob.AccountId,
					AccountCodeHeader: coa.CodeCoa,
					AccountNameHeader: coa.NameCoa,
					ItemNo:            k + 1,
					AccountId:         v.AccountId,
					AccountCode:       v.AccountCode,
					AccountName:       v.AccountName,
					VoucherSeqNo:      num,
					VoucherCode:       voucher_code,
					VoucherNo:         reference_no,
					Debet:             v.Debet,
					Credit:            v.Credit,
					Pic:               v.Pic,
					Memo:              v.Memo,
					ArId:              querydata.ArId,
					ArReferenceNo:     querydata.ArReferenceNo,
					ApId:              querydata.ApId,
					ApReferenceNo:     querydata.ApReferenceNo,
					LoanId:            querydata.LoanId,
					LoanReferenceNo:   querydata.LoanReferenceNo,
					TransactionType:   ob.TransactionType,
					Period:            utils.String2Int(period),
					StatusId:          status_id,
					StatusGlId:        status_gl_id,
					BatchNo:           batch_no,
					CreatedBy:         v.CreatedBy,
					UpdatedBy:         user_name,
				})
			}
			mutex.Unlock()
		}(k, v)
	}
	wg.Wait()

	err_ := t_pettycashh.UpdateWithDetail(t_pettycashh, inputDetail, putDetail, user_name)
	errcode, errmessage = base.DecodeErr(err_)
	if err_ != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
		c.ServeJSON()
		return
	} else {
		if err = base.PutFirebaseRaw(ob.UploadFile, user_name, id, folderName+"/"+utils.Int2String(id), folderName+"/"+utils.Int2String(id), folderName); err != nil {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprint("Error processing data and uploading to Firebase: ", err.Error()))
		} else {
			v, err := t_pettycashh.GetById(id, user_id)
			errcode, errmessage = base.DecodeErr(err)
			if err != nil {
				c.Ctx.ResponseWriter.WriteHeader(errcode)
				utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
			} else {
				utils.ReturnHTTPSuccessWithMessage(&c.Controller, errcode, errmessage, v)
			}
		}
	}

	c.ServeJSON()
}

func (c *PettyCashV2Controller) Delete() {
	var user_id, form_id int
	var err error
	var user_name string
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
		user_name = sess.(map[string]interface{})["username"].(string)
	}
	form_id = base.FormName(form_petty_cash)
	delete_aut := models.CheckPrivileges(user_id, form_id, base.Delete)
	if !delete_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Delete not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}

	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	var querydata models.PettyCashHeader
	err = models.PettyCashHeaders().Filter("id", id).One(&querydata)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, "Voucher id unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	if querydata.StatusGlId == 1 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("'%v'  has been POSTED", querydata.VoucherNo))
		c.ServeJSON()
		return
	}

	if querydata.StatusId == 1 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("'%v' has been APPROVED", querydata.VoucherNo))
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

	models.PettyCashHeaders().Filter("id", id).Filter("deleted_at__isnull", true).Update(orm.Params{"deleted_at": utils.GetSvrDate(), "deleted_by": user_name})
	models.PettyCashs().Filter("voucher_id", id).Filter("deleted_at__isnull", true).Update(orm.Params{"deleted_at": utils.GetSvrDate(), "deleted_by": user_name})

	utils.ReturnHTTPError(&c.Controller, 200, "soft delete success")
	c.ServeJSON()
}

func (c *PettyCashV2Controller) GetOne() {
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

func (c *PettyCashV2Controller) GetAll() {
	var issueDate, issueDate2, updatedat *string
	var user_id int

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
	match_mode := strings.TrimSpace(c.GetString("match_mode"))
	value_name := strings.TrimSpace(c.GetString("value_name"))
	field_name := strings.TrimSpace(c.GetString("field_name"))
	allsize, _ := c.GetInt("allsize")
	search_detail, _ := c.GetInt("search_detail")

	issue_date := strings.TrimSpace(c.GetString("issue_date"))
	issue_date2 := strings.TrimSpace(c.GetString("issue_date2"))
	updated_at := strings.TrimSpace(c.GetString("updated_at"))
	status := strings.TrimSpace(c.GetString("status"))
	account_id, _ := c.GetInt("account_id")
	is_transaction, _ := c.GetInt("is_transaction")
	company_id, _ := c.GetInt("company_id")
	report_type := 1

	if issue_date == "" {
		issueDate = nil

	} else {
		issueDate = &issue_date
	}

	if issue_date2 == "" {
		issueDate2 = nil

	} else {
		issueDate2 = &issue_date2
	}

	if updated_at == "" {
		updatedat = nil
	} else {
		updatedat = &updated_at
	}
	d, err := t_pettycashh.GetAll(keyword, field_name, match_mode, value_name, currentPage, pageSize, allsize, user_id, search_detail, report_type, company_id, account_id, is_transaction, status, issueDate, issueDate2, updatedat)
	code, message := base.DecodeErr(err)
	if err == orm.ErrNoRows {
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

func (c *PettyCashV2Controller) GetAllChild() {
	var issueDate, issueDate2, updatedat *string
	var user_id int

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
	match_mode := strings.TrimSpace(c.GetString("match_mode"))
	value_name := strings.TrimSpace(c.GetString("value_name"))
	field_name := strings.TrimSpace(c.GetString("field_name"))
	allsize, _ := c.GetInt("allsize")
	search_detail, _ := c.GetInt("search_detail")
	account_id, _ := c.GetInt("account_id")
	issue_date := strings.TrimSpace(c.GetString("issue_date"))
	issue_date2 := strings.TrimSpace(c.GetString("issue_date2"))
	updated_at := strings.TrimSpace(c.GetString("updated_at"))
	status := strings.TrimSpace(c.GetString("status"))
	is_transaction, _ := c.GetInt("is_transaction")
	company_id, _ := c.GetInt("company_id")
	report_type := 2

	if issue_date == "" {
		issueDate = nil

	} else {
		issueDate = &issue_date
	}

	if issue_date2 == "" {
		issueDate2 = nil

	} else {
		issueDate2 = &issue_date2
	}

	if updated_at == "" {
		updatedat = nil
	} else {
		updatedat = &updated_at
	}
	d, err := t_pettycashh.GetAll(keyword, field_name, match_mode, value_name, currentPage, pageSize, allsize, user_id, search_detail, report_type, company_id, account_id, is_transaction, status, issueDate, issueDate2, updatedat)
	code, message := base.DecodeErr(err)
	if err == orm.ErrNoRows {
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

func (c *PettyCashV2Controller) GetAllDetail() {
	var issueDate, issueDate2, updatedat, voucherId *string
	var user_id int
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
	match_mode := strings.TrimSpace(c.GetString("match_mode"))
	value_name := strings.TrimSpace(c.GetString("value_name"))
	field_name := strings.TrimSpace(c.GetString("field_name"))
	allsize, _ := c.GetInt("allsize")
	search_detail, _ := c.GetInt("search_detail")

	issue_date := strings.TrimSpace(c.GetString("issue_date"))
	issue_date2 := strings.TrimSpace(c.GetString("issue_date2"))
	updated_at := strings.TrimSpace(c.GetString("updated_at"))
	voucher_ids := strings.TrimSpace(c.GetString("voucher_ids"))
	status := strings.TrimSpace(c.GetString("status"))
	account_id, _ := c.GetInt("account_id")
	is_transaction, _ := c.GetInt("is_transaction")
	company_id, _ := c.GetInt("company_id")

	if issue_date == "" {
		issueDate = nil

	} else {
		issueDate = &issue_date
	}

	if issue_date2 == "" {
		issueDate2 = nil

	} else {
		issueDate2 = &issue_date2
	}

	if updated_at == "" {
		updatedat = nil
	} else {
		updatedat = &updated_at
	}

	if voucher_ids == "" {
		voucherId = nil
	} else {
		voucherId = &voucher_ids
	}

	report_type := 0
	d, err := t_pettycashh.GetAllDetail(keyword, field_name, match_mode, value_name, currentPage, pageSize, allsize, user_id, search_detail, report_type, company_id, account_id, is_transaction, status, voucherId, issueDate, issueDate2, updatedat)
	code, message := base.DecodeErr(err)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(code)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, code, "No data", d)
	} else if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(code)
		utils.ReturnHTTPError(&c.Controller, code, message)
	} else {
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, code, message, map[string]interface{}{
			"field_key":          d[0]["field_key"],
			"field_label":        d[0]["field_label"],
			"field_int":          d[0]["field_int"],
			"field_level":        d[0]["field_level"],
			"field_export":       d[0]["field_export"],
			"field_export_label": d[0]["field_export_label"],
			"field_footer":       d[0]["field_footer"],
			"list":               d,
		})
	}
	c.ServeJSON()
}

func (c *PettyCashV2Controller) GetAllList() {
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

func (c *PettyCashV2Controller) GetDocument() {
	o := orm.NewOrm()
	var user_id int
	var folder_name string = "petty_cash"
	var err error
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}

	form_id := base.FormName(form_petty_cash)
	aut_aut := models.CheckPrivileges(user_id, form_id, base.Author)
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	var coa coaUserRtn
	var document models.Document
	if !aut_aut {
		err = o.Raw("select id,code_coa,name_coa,code_in,code_out,status_id,company_id,company_code,user_id from chart_of_accounts t0 left join (select user_id,account_id from sys_user_account where  user_id = " + utils.Int2String(user_id) + " ) t1 on t1.account_id = t0.id where deleted_at is null  and id not in (select parent_id from chart_of_accounts where deleted_at is null)  and is_header = 1 and t0.id in (select account_id from petty_cash_header where id = " + utils.Int2String(id) + ")").QueryRow(&coa)
	} else {
		err = o.Raw("select id,code_coa,name_coa,code_in,code_out,status_id,company_id,company_code," + utils.Int2String(user_id) + " user_id from chart_of_accounts where deleted_at is null  and id not in (select parent_id from chart_of_accounts where deleted_at is null)  and is_header = 1 and id in (select account_id from petty_cash_header where id = " + utils.Int2String(id) + ")").QueryRow(&coa)
	}
	code, message := base.DecodeErr(err)
	if err == orm.ErrNoRows {
		code = 200
		c.Ctx.ResponseWriter.WriteHeader(code)
		utils.ReturnHTTPError(&c.Controller, code, "No data")
	} else if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(code)
		utils.ReturnHTTPError(&c.Controller, code, message)
	} else {

		d := document.GetDocument(id, folder_name)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 200, "Success", d)
	}
	c.ServeJSON()
}

func (c *PettyCashV2Controller) ReOrderNum() {
	var user_name string
	var user_id, form_id int
	fmt.Print("Check :", user_id, form_id, user_name, "..")
	sess := c.GetSession("profile")
	if sess != nil {
		user_name = sess.(map[string]interface{})["username"].(string)
		user_id = sess.(map[string]interface{})["id"].(int)
	}

	form_id = base.FormName("reorder_petty_cash")
	write_aut := models.CheckPrivileges(user_id, form_id, base.Write)
	if !write_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Reorder not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
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

func (c *PettyCashV2Controller) ReOrderNumList() {
	var user_name string
	var user_id, form_id int
	fmt.Print("Check :", user_id, form_id, user_name, "..")
	sess := c.GetSession("profile")
	if sess != nil {
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
