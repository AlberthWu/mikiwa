package controllers

import (
	"encoding/json"
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

type PettyCashController struct {
	base.BaseController
}

func (c *PettyCashController) Prepare() {
	c.Ctx.Request.Header.Set("token", "No Aut")
	c.BaseController.Prepare()
}

func (c *PettyCashController) Post() {
	o := orm.NewOrm()
	type InputBody struct {
		VoucherId   int     `json:"voucher_id"`
		AccountId   int     `json:"account_id"`
		Debet       float64 `json:"debet"`
		Credit      float64 `json:"credit"`
		Pic         string  `json:"pic"`
		Memo        string  `json:"memo"`
		ReceivingNo string  `json:"receiving_no"`
	}
	var input []models.PettyCash
	var ob []InputBody
	var err error
	var i int = 1
	var deletedat string
	var user_name string
	var user_id int = 0

	sess := c.GetSession("profile")
	if sess != nil {
		user_name = sess.(map[string]interface{})["username"].(string)
		user_id = sess.(map[string]interface{})["id"].(int)
	}

	fmt.Print("Check :", user_id, "..")
	body := c.Ctx.Input.RequestBody
	json.Unmarshal(body, &ob)
	valid := validation.Validation{}
	for _, v := range ob {

		var querydata models.PettyCashHeader
		err = models.PettyCashHeaders().Filter("id", v.VoucherId).One(&querydata)
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

		valid.Required(v.VoucherId, "voucher_id").Message("Is required")
		if v.Debet+v.Credit == 0 {
			if querydata.TransactionType == "In" {
				valid.Required(v.Debet, "debet").Message("Is required")
			} else if querydata.TransactionType == "Out" {
				valid.Required(v.Credit, "credit").Message("Is required")
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

		var coa models.CharOfAccount
		err = models.ChartOfAccounts().Filter("id", v.AccountId).One(&coa)
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

		issue_date := querydata.IssueDate.Format("2006-01-02")

		thedate, errdate := time.Parse("2006-01-02", issue_date)
		if errdate != nil {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, errdate.Error())
			c.ServeJSON()
			return
		}

		input = append(input, models.PettyCash{
			IssueDate:         thedate,
			VoucherId:         v.VoucherId,
			CompanyId:         querydata.CompanyId,
			CompanyCode:       querydata.CompanyCode,
			CompanyName:       querydata.CompanyName,
			AccountIdHeader:   querydata.AccountId,
			AccountCodeHeader: querydata.AccountCode,
			AccountNameHeader: querydata.AccountName,
			AccountId:         v.AccountId,
			AccountCode:       coa.CodeCoa,
			AccountName:       coa.NameCoa,
			VoucherSeqNo:      querydata.VoucherSeqNo,
			VoucherCode:       querydata.VoucherCode,
			VoucherNo:         querydata.VoucherNo,
			Debet:             v.Debet,
			Credit:            v.Credit,
			Pic:               v.Pic,
			Memo:              v.Memo,
			ReceivingNo:       v.ReceivingNo,
			TransactionType:   querydata.TransactionType,
			Period:            querydata.Period,
			BatchNo:           querydata.BatchNo,
			CreatedBy:         user_name,
			UpdatedBy:         user_name,
		})

		i += 1
	}
	o.InsertMulti(i, input)
	c.Data["json"] = input
	c.ServeJSON()
}

func (c *PettyCashController) Put() {
	var err error
	var user_name string
	var user_id int = 0
	var details models.PettyCashHeader
	sess := c.GetSession("profile")
	if sess != nil {
		user_name = sess.(map[string]interface{})["username"].(string)
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	fmt.Print("Check :", user_id, "..")

	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	err = models.PettyCashHeaders().Filter("id", id).One(&details)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Invoice id unregistered/Illegal data")
		c.ServeJSON()
		return
	}
	var deletedat = details.DeletedAt.Format("2006-01-02")
	if deletedat != "0001-01-01" {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("'%v' has been deleted", details.VoucherNo))
		c.ServeJSON()
		return
	}

	if details.StatusId > 0 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("'%v' has been CLOSED", details.VoucherNo))
		c.ServeJSON()
		return
	}

	if details.StatusGlId > 0 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("'%v' has been POSTED", details.VoucherNo))
		c.ServeJSON()
		return
	}

	issue_date := details.IssueDate.Format("2006-01-02")

	thedate, errdate := time.Parse("2006-01-02", issue_date)
	if errdate != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, errdate.Error())
		c.ServeJSON()
		return
	}

	type InputBody struct {
		Id          int     `json:"id"`
		VoucherId   int     `json:"voucher_id"`
		AccountId   int     `json:"account_id"`
		Debet       float64 `json:"debet"`
		Credit      float64 `json:"credit"`
		Pic         string  `json:"pic"`
		Memo        string  `json:"memo"`
		ReceivingNo string  `json:"receiving_no"`
	}

	var ob []InputBody
	body := c.Ctx.Input.RequestBody
	json.Unmarshal(body, &ob)
	valid := validation.Validation{}
	for _, v := range ob {

		valid.Required(v.VoucherId, "voucher_id").Message("Is required")
		if v.Debet+v.Credit == 0 {
			if details.TransactionType == "In" {
				valid.Required(v.Debet, "debet").Message("Is required")
			} else if details.TransactionType == "Out" {
				valid.Required(v.Credit, "credit").Message("Is required")

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

		var querydata models.PettyCash
		models.PettyCashs().Filter("deleted_at__isnull", true).Filter("id", v.Id).Filter("voucher_id", id).One(&querydata)
		if err == orm.ErrNoRows {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, "Detail id unregistered/Illegal data")
			c.ServeJSON()
			return
		}
		var deletedat = details.DeletedAt.Format("2006-01-02")
		if deletedat != "0001-01-01" {
			c.Ctx.ResponseWriter.WriteHeader(402)
			utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Id '%v' has been deleted", details.Id))
			c.ServeJSON()
			return
		}

		var coa models.CharOfAccount
		err = models.ChartOfAccounts().Filter("id", v.AccountId).One(&coa)
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

		t_pettycash.Id = v.Id
		t_pettycash.IssueDate = thedate
		t_pettycash.VoucherId = querydata.VoucherId
		t_pettycash.CompanyId = querydata.CompanyId
		t_pettycash.CompanyCode = querydata.CompanyCode
		t_pettycash.CompanyName = querydata.CompanyName
		t_pettycash.AccountIdHeader = querydata.AccountIdHeader
		t_pettycash.AccountCodeHeader = querydata.AccountCodeHeader
		t_pettycash.AccountNameHeader = querydata.AccountNameHeader
		t_pettycash.AccountId = v.AccountId
		t_pettycash.AccountCode = coa.CodeCoa
		t_pettycash.AccountName = coa.NameCoa
		t_pettycash.VoucherCode = querydata.VoucherCode
		t_pettycash.VoucherNo = querydata.VoucherNo
		t_pettycash.VoucherSeqNo = querydata.VoucherSeqNo
		t_pettycash.Debet = v.Debet
		t_pettycash.Credit = v.Credit
		t_pettycash.ArId = querydata.ArId
		t_pettycash.ArReferenceNo = querydata.ArReferenceNo
		t_pettycash.ApId = querydata.ApId
		t_pettycash.ApReferenceNo = querydata.ApReferenceNo
		t_pettycash.Pic = v.Pic
		t_pettycash.Memo = v.Memo
		t_pettycash.ReceivingId = querydata.ReceivingId
		t_pettycash.ReceivingNo = v.ReceivingNo
		t_pettycash.LoanId = querydata.LoanId
		t_pettycash.LoanReferenceNo = querydata.LoanReferenceNo
		t_pettycash.TransactionType = details.TransactionType
		t_pettycash.StatusId = details.StatusId
		t_pettycash.StatusGlId = details.StatusGlId
		t_pettycash.Period = details.Period
		t_pettycash.BatchNo = details.BatchNo
		t_pettycash.CreatedBy = details.CreatedBy
		t_pettycash.UpdatedBy = user_name

		t_pettycash.Update()

	}
	c.Data["json"] = "Success"
	c.ServeJSON()
}

func (c *PettyCashController) Delete() {
	var user_name string
	var err error
	sess := c.GetSession("profile")
	if sess != nil {
		user_name = sess.(map[string]interface{})["username"].(string)
	}
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	var querydata models.PettyCash
	models.PettyCashs().Filter("voucher_id", id).One(&querydata)
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

	var deletedat = querydata.DeletedAt.Format("2006-01-02")
	if deletedat != "0001-01-01" {
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

	var ub DeleteBody
	body := c.Ctx.Input.RequestBody
	json.Unmarshal(body, &ub)
	var joinId []string
	ids := strings.Split(ub.Id, ",")
	for _, st := range ids {
		joinId = append(joinId, st)
	}

	num, err := models.PettyCashs().Filter("deleted_at__isnull", true).Filter("voucher_id", id).Filter("id__in", joinId).Update(orm.Params{"deleted_at": utils.GetSvrDate(), "deleted_by": user_name})
	code, message := base.DecodeErr(err)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(code)
		utils.ReturnHTTPError(&c.Controller, code, message)
	} else {
		c.Ctx.ResponseWriter.WriteHeader(200)
		utils.ReturnHTTPError(&c.Controller, 200, fmt.Sprintf("Total '%v''s data has been deleted", num))
	}

	c.ServeJSON()
}

func (c *PettyCashController) CheckDelete() {

	var user_id, form_id int
	var querydata models.PettyCash
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	form_id = base.FormName("petty_cash")
	delete_aut := models.CheckPrivileges(user_id, form_id, 12)
	if !delete_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Delete detail not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}

	err := models.PettyCashs().Filter("deleted_at__isnull", true).Filter("id", id).One(&querydata)

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

	var deletedat = querydata.DeletedAt.Format("2006-01-02")
	if deletedat != "0001-01-01" {
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

	c.Ctx.ResponseWriter.WriteHeader(200)
	utils.ReturnHTTPError(&c.Controller, 200, "Success, delete allowed")
	c.ServeJSON()

}
