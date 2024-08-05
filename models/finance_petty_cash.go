package models

import (
	"fmt"
	"mikiwa/utils"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type (
	PettyCashHeader struct {
		Id              int       `json:"id" orm:"column(id);auto;pk"`
		IssueDate       time.Time `json:"issue_date" orm:"column(issue_date);type(date)"`
		CompanyId       int       `json:"company_id"  orm:"column(company_id)"`
		CompanyCode     string    `json:"company_code" orm:"column(company_code)"`
		CompanyName     string    `json:"company_name" orm:"column(company_name)"`
		AccountId       int       `json:"account_id"  orm:"column(account_id)"`
		AccountCode     string    `json:"account_code" orm:"column(account_code)"`
		AccountName     string    `json:"account_name" orm:"column(account_name)"`
		VoucherSeqNo    int       `json:"voucher_seq_no"  orm:"column(voucher_seq_no)"`
		VoucherCode     string    `json:"voucher_code" orm:"column(voucher_code)"`
		VoucherNo       string    `json:"voucher_no" orm:"column(voucher_no)"`
		TransactionType string    `json:"transaction_type" orm:"column(transaction_type)"` // IN or OUT
		Debet           float64   `json:"debet" orm:"column(debet);digits(18);decimals(2);default(0)"`
		Credit          float64   `json:"credit" orm:"column(credit);digits(18);decimals(2);default(0)"`
		BatchNo         string    `json:"batch_no" orm:"column(batch_no)"`
		StatusId        int8      `json:"status_id" orm:"column(status_id)"`
		StatusGlId      int8      `json:"status_gl_id" orm:"column(status_gl_id)"`
		Period          int       `json:"period" orm:"column(period)"`
		Pic             string    `json:"pic" orm:"column(pic)"`
		Memo            string    `json:"memo" orm:"column(memo)"`
		ArId            int       `json:"ar_id"  orm:"column(ar_id)"`
		ArReferenceNo   string    `json:"ar_reference_no" orm:"column(ar_reference_no)"`
		ApId            int       `json:"ap_id"  orm:"column(ap_id)"`
		ApReferenceNo   string    `json:"ap_reference_no" orm:"column(ap_reference_no)"`
		LoanId          int       `json:"loan_id"  orm:"column(loan_id)"`
		LoanReferenceNo string    `json:"loan_reference_no" orm:"column(loan_reference_no)"`
		CreatedAt       time.Time `json:"created_at" orm:"column(created_at);type(timestamp);auto_now_add"`
		UpdatedAt       time.Time `json:"updated_at" orm:"column(updated_at);type(timestamp);auto_now"`
		DeletedAt       time.Time `json:"deleted_at" orm:"column(deleted_at);type(timestamp);null"`
		CreatedBy       string    `json:"created_by" orm:"column(created_by)"`
		UpdatedBy       string    `json:"updated_by" orm:"column(updated_by)"`
		DeletedBy       string    `json:"deleted_by" orm:"column(deleted_by)"`
	}

	PettyCash struct {
		Id                int       `json:"id" orm:"column(id);auto;pk"`
		VoucherId         int       `json:"voucher_id"  orm:"column(voucher_id)"`
		IssueDate         time.Time `json:"issue_date" orm:"column(issue_date);type(date)"`
		CompanyId         int       `json:"company_id"  orm:"column(company_id)"`
		CompanyCode       string    `json:"company_code" orm:"column(company_code)"`
		CompanyName       string    `json:"company_name" orm:"column(company_name)"`
		AccountIdHeader   int       `json:"account_id_header"  orm:"column(account_id_header)"`
		AccountCodeHeader string    `json:"account_code_header" orm:"column(account_code_header)"`
		AccountNameHeader string    `json:"account_name_header" orm:"column(account_name_header)"`
		ItemNo            int       `json:"item_no"  orm:"column(item_no)"`
		AccountId         int       `json:"account_id"  orm:"column(account_id)"`
		AccountCode       string    `json:"account_code" orm:"column(account_code)"`
		AccountName       string    `json:"account_name" orm:"column(account_name)"`
		VoucherSeqNo      int       `json:"voucher_seq_no"  orm:"column(voucher_seq_no)"`
		VoucherCode       string    `json:"voucher_code" orm:"column(voucher_code)"`
		VoucherNo         string    `json:"voucher_no" orm:"column(voucher_no)"`
		Debet             float64   `json:"debet" orm:"column(debet);digits(18);decimals(2);default(0)"`
		Credit            float64   `json:"credit" orm:"column(credit);digits(18);decimals(2);default(0)"`
		Pic               string    `json:"pic" orm:"column(pic)"`
		Memo              string    `json:"memo" orm:"column(memo)"`
		EmployeeId        int       `json:"employee_id"  orm:"column(employee_id)"`
		EmployeeName      string    `json:"employee_name" orm:"column(employee_name)"`
		ArId              int       `json:"ar_id"  orm:"column(ar_id)"`
		ArReferenceNo     string    `json:"ar_reference_no" orm:"column(ar_reference_no)"`
		ApId              int       `json:"ap_id"  orm:"column(ap_id)"`
		ApReferenceNo     string    `json:"ap_reference_no" orm:"column(ap_reference_no)"`
		ReceivingId       int       `json:"receiving_id"  orm:"column(receiving_id)"`
		ReceivingNo       string    `json:"receiving_no" orm:"column(receiving_no)"`
		LoanId            int       `json:"loan_id"  orm:"column(loan_id)"`
		LoanReferenceNo   string    `json:"loan_reference_no" orm:"column(loan_reference_no)"`
		TransactionType   string    `json:"transaction_type" orm:"column(transaction_type)"` // IN or OUT
		StatusId          int8      `json:"status_id" orm:"column(status_id)"`
		StatusGlId        int8      `json:"status_gl_id" orm:"column(status_gl_id)"`
		Period            int       `json:"period" orm:"column(period)"`
		BatchNo           string    `json:"batch_no" orm:"column(batch_no)"`
		CreatedAt         time.Time `json:"created_at" orm:"column(created_at);type(timestamp);auto_now_add"`
		UpdatedAt         time.Time `json:"updated_at" orm:"column(updated_at);type(timestamp);auto_now"`
		DeletedAt         time.Time `json:"deleted_at" orm:"column(deleted_at);type(timestamp);null"`
		CreatedBy         string    `json:"created_by" orm:"column(created_by)"`
		UpdatedBy         string    `json:"updated_by" orm:"column(updated_by)"`
		DeletedBy         string    `json:"deleted_by" orm:"column(deleted_by)"`
	}
)

func (t *PettyCashHeader) TableName() string {
	return "petty_cash_header"
}

func PettyCashHeaders() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(PettyCashHeader))
}

func (t *PettyCash) TableName() string {
	return "petty_cash"
}

func PettyCashs() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(PettyCash))
}

func init() {
	orm.RegisterModel(new(PettyCash), new(PettyCashHeader))
}

func (t *PettyCashHeader) Insert(m PettyCashHeader) (*PettyCashHeader, error) {
	o := orm.NewOrm()

	if _, err := o.Insert(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (t *PettyCashHeader) Update(fields ...string) error {
	o := orm.NewOrm()
	if _, err := o.Update(t, fields...); err != nil {
		return err
	}
	return nil
}

func (t *PettyCash) Insert(m PettyCash) (*PettyCash, error) {
	o := orm.NewOrm()

	if _, err := o.Insert(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (t *PettyCash) Update(fields ...string) error {
	o := orm.NewOrm()
	if _, err := o.Update(t, fields...); err != nil {
		return err
	}
	return nil
}

func (t *PettyCashHeader) InsertWithDetail(m PettyCashHeader, d []PettyCash) (data *PettyCashHeader, err error) {
	o := orm.NewOrm()

	if _, err = o.Insert(&m); err != nil {
		return nil, err
	}

	for i := range d {
		d[i].VoucherId = m.Id
	}

	_, err = o.InsertMulti(len(d), d)
	if err != nil {
		o.Raw("update petty_cash_header set deleted_at = now(),deleted_by = 'Failed' where id = " + utils.Int2String(m.Id)).Exec()
		return nil, err
	}

	return &m, nil
}

func (t *PettyCashHeader) UpdateWithDetail(m PettyCashHeader, data_post, data_put []PettyCash, user_name string) error {
	o := orm.NewOrm()
	tx, err := o.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	if _, err := tx.Update(&m); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update order: %v", err)
	}

	// Update existing PettyCashDetail (Details) and delete
	var deleteIds []string
	var joinId string
	for _, detail := range data_put {
		if detail.Id != 0 {
			if _, err := tx.Update(&detail); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update detail: %v", err)
			}
		}

		deleteIds = append(deleteIds, utils.Int2String(detail.Id))
	}

	if len(deleteIds) == 0 {
		joinId = "0"
	} else {
		joinId = strings.Join(deleteIds, ",")
	}
	_, err = o.Raw("update petty_cash set deleted_at = now(), deleted_by = '" + user_name + "' where deleted_at is null and voucher_id = " + utils.Int2String(m.Id) + " and id not in (" + joinId + ") ").Exec()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete existing details: %v", err)
	}

	// Insert new PettyCashDetail (Details)
	for i := range data_post {
		data_post[i].VoucherId = m.Id
	}

	if len(data_post) > 0 {
		_, err = o.InsertMulti(len(data_post), data_post)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert new details: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("failed to rollback transaction: %v after commit failed: %v", rbErr, err)
		}
		return fmt.Errorf("failed to commit transaction: %v", err)
	}
	return nil
}

// TODO TEMPORARAY

type (
	PettyCashRtn struct {
		Id               int                    `json:"id"`
		IssueDate        string                 `json:"issue_date"`
		CompanyId        int                    `json:"company_id"`
		CompanyCode      string                 `json:"company_code"`
		CompanyName      string                 `json:"company_name"`
		CompanyStatus    int                    `json:"company_status"`
		EffectiveDate    string                 `json:"effective_date"`
		Expired_Date     *string                `json:"expired_date"`
		AccountId        int                    `json:"account_id"`
		AccountCode      string                 `json:"account_code"`
		AccountName      string                 `json:"account_name"`
		CodeIn           string                 `json:"code_in"`
		CodeOut          string                 `json:"code_out"`
		VoucherSeqNo     int                    `json:"voucher_seq_no"`
		VoucherCode      string                 `json:"voucher_code"`
		VoucherNo        string                 `json:"voucher_no"`
		Debet            float64                `json:"debet"`
		Credit           float64                `json:"credit"`
		BatchNo          string                 `json:"batch_no"`
		Period           int                    `json:"period"`
		TransactionType  string                 `json:"transaction_type"`
		StatusId         int8                   `json:"status_id"`
		StatusGlId       int8                   `json:"status_gl_id"`
		Memo             string                 `json:"memo"`
		FieldKey         string                 `json:"field_key"`
		FieldLabel       string                 `json:"field_label"`
		FieldInt         string                 `json:"field_int"`
		FieldLevel       string                 `json:"field_level"`
		FieldFooter      string                 `json:"field_footer"`
		FieldExport      string                 `json:"field_export"`
		FieldExportLabel string                 `json:"field_export_label"`
		ApproveButton    int8                   `json:"approve_button"`
		RejectButton     int8                   `json:"reject_button"`
		DetailList       []PettyCashVoucherJson `json:"detail_list"`
	}

	PettyCashRtnHeader struct {
		Id          int            `json:"id"`
		CompanyId   int            `json:"company_id"`
		CompanyCode string         `json:"company_code"`
		CompanyName string         `json:"company_name"`
		AccountId   int            `json:"account_id"`
		AccountCode string         `json:"account_code"`
		AccountName string         `json:"account_name"`
		Opening     float64        `json:"opening"`
		Debet       float64        `json:"debet"`
		Credit      float64        `json:"credit"`
		Balance     float64        `json:"balance"`
		Detail      []PettyCashRtn `json:"detail"`
	}

	PettyCashRtnJson struct {
		Id               int                    `json:"id"`
		IssueDate        string                 `json:"issue_date"`
		CompanyId        CompanyListRtnJson     `json:"company_id"`
		AccountId        CoaRtnJson             `json:"account_id"`
		VoucherSeqNo     int                    `json:"voucher_seq_no"`
		VoucherCode      string                 `json:"voucher_code"`
		VoucherNo        string                 `json:"voucher_no"`
		TransactionType  string                 `json:"transaction_type"`
		Debet            float64                `json:"debet"`
		Credit           float64                `json:"credit"`
		Balance          float64                `json:"balance"`
		BatchNo          string                 `json:"batch_no"`
		Memo             string                 `json:"memo"`
		FieldKey         string                 `json:"field_key"`
		FieldLabel       string                 `json:"field_label"`
		FieldInt         string                 `json:"field_int"`
		FieldLevel       string                 `json:"field_level"`
		FieldFooter      string                 `json:"field_footer"`
		FieldExport      string                 `json:"field_export"`
		FieldExportLabel string                 `json:"field_export_label"`
		ApproveButton    int8                   `json:"approve_button"`
		RejectButton     int8                   `json:"reject_button"`
		DetailList       []PettyCashVoucherJson `json:"detail_list"`
		Document         []DocumentRtn          `json:"document"`
	}

	PettyCashVoucherJson struct {
		Id                int     `json:"id"`
		IssueDate         string  `json:"issue_date"`
		CompanyId         int     `json:"company_id"`
		CompanyCode       string  `json:"company_code"`
		CompanyName       string  `json:"company_name"`
		AccountIdHeader   int     `json:"account_id_header"`
		AccountCodeHeader string  `json:"account_code_header"`
		AccountNameHeader string  `json:"account_name_header"`
		AccountId         int     `json:"account_id"`
		AccountCode       string  `json:"account_code"`
		AccountName       string  `json:"account_name"`
		Debet             float64 `json:"debet"`
		Credit            float64 `json:"credit"`
		Balance           float64 `json:"balance"`
		Memo              string  `json:"memo"`
		Pic               string  `json:"pic"`
		ReceivingId       int     `json:"receiving_id"`
		ReceivingNo       string  `json:"receiving_no"`
		VoucherId         int     `json:"voucher_id"`
		VoucherSeqNo      int     `json:"voucher_seq_no"`
		VoucherNo         string  `json:"voucher_no"`
		BatchNo           string  `json:"batch_no"`
		StatusId          int8    `json:"status_id"`
		StatusGlId        int8    `json:"status_gl_id"`
	}

	PettyCashDailyJson struct {
		Id              int     `json:"id"`
		IssueDate       string  `json:"issue_date"`
		CompanyId       int     `json:"company_id"`
		CompanyCode     string  `json:"company_code"`
		CompanyName     string  `json:"company_name"`
		VoucherSeqNo    int     `json:"voucher_seq_no"`
		VoucherCode     string  `json:"voucher_code"`
		VoucherNo       string  `json:"voucher_no"`
		TransactionType string  `json:"transaction_type"`
		AccountId       int     `json:"account_id"`
		AccountCode     string  `json:"account_code"`
		AccountName     string  `json:"account_name"`
		Memo            string  `json:"memo"`
		Debet           float64 `json:"debet"`
		Credit          float64 `json:"credit"`
		Balance         float64 `json:"balance"`
		BatchNo         string  `json:"batch_no"`
		StatusId        int8    `json:"status_id"`
		StatusGlId      int8    `json:"status_gl_id"`
	}

	PettyCashReOrderJson struct {
		Id              int     `json:"id"`
		IssueDate       string  `json:"issue_date"`
		CompanyId       int     `json:"company_id"`
		CompanyCode     string  `json:"company_code"`
		CompanyName     string  `json:"company_name"`
		PreVoucherSeqNo int     `json:"pre_voucher_seq_no"`
		NewVoucherSeqNo int     `json:"new_voucher_seq_no"`
		PreVoucherCode  string  `json:"pre_voucher_code"`
		NewVoucherCode  string  `json:"new_voucher_code"`
		PreVoucherNo    string  `json:"pre_voucher_no"`
		NewVoucherNo    string  `json:"new_voucher_no"`
		AccountId       int     `json:"account_id"`
		AccountCode     string  `json:"account_code"`
		AccountName     string  `json:"account_name"`
		Memo            string  `json:"memo"`
		Debet           float64 `json:"debet"`
		Credit          float64 `json:"credit"`
		Balance         float64 `json:"balance"`
		BatchNo         string  `json:"batch_no"`
		StatusId        int8    `json:"status_id"`
		StatusGlId      int8    `json:"status_gl_id"`
		TransactionType string  `json:"transaction_type"`
	}

	VoucherSimpleRtnJson struct {
		Id          int    `json:"id"`
		VoucherNo   string `json:"voucher_no"`
		VoucherCode string `json:"voucher_code"`
	}

	vCount struct {
		Data0    int `json:"data0"`
		Open     int `json:"open"`
		Unposted int `json:"unposted"`
		Posted   int `json:"posted"`
		Total    int `json:"total"`
	}
)

func (t *PettyCashHeader) GetAll(keyword, field_name, match_mode, value_name string, p, size, allsize, user_id, search_detail, report_Type int, company_id, account_id, is_transaction int, status string, thedate, thedate2, updated_at *string) (u utils.PageDynamicAdd, err error) {
	o := orm.NewOrm()
	var m []orm.Params

	var dataCount vCount

	o.Raw("call sp_PettyCashCountV2(?,?,?,null,"+utils.Int2String(company_id)+","+utils.Int2String(account_id)+","+utils.Int2String(is_transaction)+",'"+status+"',"+utils.Int2String(user_id)+","+utils.Int2String(report_Type)+",'"+keyword+"',"+utils.Int2String(search_detail)+",'"+field_name+"','"+match_mode+"','"+value_name+"',null,null)", &thedate, &thedate2, &updated_at).QueryRow(&dataCount)

	if allsize == 1 && dataCount.Data0 > 0 {
		size = dataCount.Data0
	}
	_, err = o.Raw("call sp_PettyCashV2(?,?,?,null,"+utils.Int2String(company_id)+","+utils.Int2String(account_id)+","+utils.Int2String(is_transaction)+",'"+status+"',"+utils.Int2String(user_id)+","+utils.Int2String(report_Type)+",'"+keyword+"',"+utils.Int2String(search_detail)+",'"+field_name+"','"+match_mode+"','"+value_name+"',"+utils.Int2String(size)+", "+utils.Int2String((p-1)*size)+")", &thedate, &thedate2, &updated_at).Values(&m)
	if dataCount.Data0 == 0 && err == nil {
		err = orm.ErrNoRows
		return utils.PaginationDynamicAdd(dataCount.Data0, p, size, "", "", "", "", "", "", "", m, nil), err
	} else if err != nil {
		return utils.PaginationDynamicAdd(dataCount.Data0, p, size, "", "", "", "", "", "", "", m, nil), err
	}
	return utils.PaginationDynamicAdd(dataCount.Data0, p, size, fmt.Sprintf("%v", m[0]["field_key"]), fmt.Sprintf("%v", m[0]["field_label"]), fmt.Sprintf("%v", m[0]["field_int"]), fmt.Sprintf("%v", m[0]["field_level"]), fmt.Sprintf("%v", m[0]["field_export"]), fmt.Sprintf("%v", m[0]["field_export_label"]), fmt.Sprintf("%v", m[0]["field_footer"]), m, dataCount), err
}

func (t *PettyCashHeader) GetAllDetail(keyword, field_name, match_mode, value_name string, p, size, allsize, user_id, search_detail, report_Type int, company_id, account_id, is_transaction int, status string, id, thedate, thedate2, updated_at *string) (m []orm.Params, err error) {
	o := orm.NewOrm()
	var c int64
	c, err = o.Raw("call sp_PettyCashV2(?,?,?,?,"+utils.Int2String(company_id)+","+utils.Int2String(account_id)+","+utils.Int2String(is_transaction)+",'"+status+"',"+utils.Int2String(user_id)+","+utils.Int2String(report_Type)+",'"+keyword+"',"+utils.Int2String(search_detail)+",'"+field_name+"','"+match_mode+"','"+value_name+"',"+utils.Int2String(size)+", "+utils.Int2String((p-1)*size)+")", &thedate, &thedate2, &updated_at, &id).Values(&m)
	if c == 0 && err == nil {
		if _, err = o.Raw("call sp_PettyCashV2(null,null,null,0,null,null,null,null,null,null,null,null,0,null,null,null,null,null,null,null, null)").Values(&m); err == nil {
			err = orm.ErrNoRows
		}
	}
	return m, err
}

func (t *PettyCashHeader) GetById(id, user_id int) (m *PettyCashRtnJson, err error) {
	o := orm.NewOrm()
	var v PettyCashHeader

	var querydata PettyCashRtn
	err = o.Raw("call sp_PettyCashOne(" + utils.Int2String(id) + "," + utils.Int2String(user_id) + ")").QueryRow(&querydata)

	var companyrtn = CompanyListRtnJson{Id: querydata.CompanyId, Code: querydata.CompanyCode, Name: querydata.CompanyName, Status: querydata.CompanyStatus}
	var coartn = CoaRtnJson{Id: querydata.AccountId, EffectiveDate: querydata.EffectiveDate, ExpiredDate: querydata.Expired_Date, CodeCoa: querydata.AccountCode, NameCoa: querydata.AccountName, CodeIn: querydata.CodeIn, CodeOut: querydata.CodeOut}
	dlist, _ := v.GetByIdVoucher(querydata.Id)
	var td Document
	ilist := td.GetDocument(id, "petty_cash")

	m = &PettyCashRtnJson{
		Id:               querydata.Id,
		IssueDate:        querydata.IssueDate,
		CompanyId:        companyrtn,
		AccountId:        coartn,
		VoucherSeqNo:     querydata.VoucherSeqNo,
		VoucherCode:      querydata.VoucherCode,
		VoucherNo:        querydata.VoucherNo,
		TransactionType:  querydata.TransactionType,
		Debet:            querydata.Debet,
		Credit:           querydata.Credit,
		BatchNo:          querydata.BatchNo,
		Memo:             querydata.Memo,
		FieldKey:         querydata.FieldKey,
		FieldLabel:       querydata.FieldLabel,
		FieldInt:         querydata.FieldInt,
		FieldLevel:       querydata.FieldLevel,
		FieldFooter:      querydata.FieldFooter,
		FieldExport:      querydata.FieldExport,
		FieldExportLabel: querydata.FieldExportLabel,
		ApproveButton:    querydata.ApproveButton,
		RejectButton:     querydata.RejectButton,
		DetailList:       dlist,
		Document:         ilist,
	}

	return m, err
}

func (t *PettyCash) GetAll(keyword string, p, size int, thedate, status_gl_id string, company_id, account_id int) (u utils.Page, err error) {

	return utils.Pagination(int(0), p, size, nil), err
}

func (t *PettyCash) GetById(id int) (m *PettyCashVoucherJson, err error) {
	return nil, nil
}

func (t *PettyCashHeader) GetByIdVoucher(voucher_id int) (m []PettyCashVoucherJson, err error) {
	var querydata []PettyCash
	num, err := PettyCashs().Filter("voucher_id", voucher_id).Filter("deleted_at__isnull", true).All(&querydata)

	for _, list := range querydata {
		m = append(m, PettyCashVoucherJson{
			Id:                list.Id,
			IssueDate:         list.IssueDate.Format("2006-01-02"),
			AccountIdHeader:   list.AccountIdHeader,
			AccountCodeHeader: list.AccountCodeHeader,
			AccountNameHeader: list.AccountNameHeader,
			AccountId:         list.AccountId,
			AccountCode:       list.AccountCode,
			AccountName:       list.AccountName,
			Debet:             list.Debet,
			Credit:            list.Credit,
			Memo:              list.Memo,
			VoucherNo:         list.VoucherNo,
			BatchNo:           list.BatchNo,
			StatusId:          list.StatusId,
			StatusGlId:        list.StatusGlId,
			Pic:               list.Pic,
			ReceivingId:       list.ReceivingId,
			ReceivingNo:       list.ReceivingNo,
		})
	}

	if err == nil && num == 0 {
		err = orm.ErrNoRows
	}
	return m, err
}

func (t *PettyCashHeader) GetAllList(account_id int, issue_date, transaction_type, keyword string) (m []PettyCashVoucherJson, err error) {
	var querydata []PettyCashHeader
	var num int64
	cond := orm.NewCondition()
	cond1 := cond.And("issue_date", issue_date).And("deleted_at__isnull", true).And("account_id", account_id)
	qs := PettyCashHeaders().SetCond(cond1)
	cond2 := cond.AndCond(cond1)
	if keyword != "" {
		cond2.AndCond(cond.Or("voucher_no__icontains", keyword))
	}
	qs = qs.SetCond(cond2)
	// if transaction_type == "OUT" {
	num, err = qs.Filter("ap_id", 0).Filter("ar_id", 0).All(&querydata)
	// } else if transaction_type == "IN" {
	// 	num, err = qs.Filter("ar_id", 0).Filter("ar_id", 0).All(&querydata)
	// }

	for _, list := range querydata {
		m = append(m, PettyCashVoucherJson{
			Id:          list.Id,
			AccountId:   list.AccountId,
			AccountCode: list.AccountCode,
			AccountName: list.AccountName,
			Debet:       list.Debet,
			Credit:      list.Credit,
			Memo:        list.Memo,
			VoucherNo:   list.VoucherNo,
			BatchNo:     list.BatchNo,
			StatusId:    list.StatusId,
			StatusGlId:  list.StatusGlId,
			Pic:         list.Pic,
		})
	}

	if num == 0 && err == nil {
		err = orm.ErrNoRows
	}
	return m, err
}

func (t *PettyCashHeader) ReOrderNum(the_year, the_month, account_id, action_status int, ids, user_name string) (m []PettyCashReOrderJson, err error) {
	o := orm.NewOrm()
	var num int64
	num, err = o.Raw("call sp_ReorderNumber(" + utils.Int2String(the_year) + "," + utils.Int2String(the_month) + ",'" + ids + "'," + utils.Int2String(account_id) + "," + utils.Int2String(action_status) + ",'" + user_name + "')").QueryRows(&m)
	if num == 0 && err == nil {
		err = orm.ErrNoRows
	}
	return m, err
}

func (t *PettyCashHeader) ReOrderNumList(keyword, field_name, match_mode, value_name string, p, size, the_year, the_month, account_id, user_id int) (m []PettyCashReOrderJson, err error) {
	o := orm.NewOrm()
	var num int64
	num, err = o.Raw("call sp_ReorderList(" + utils.Int2String(the_year) + "," + utils.Int2String(the_month) + "," + utils.Int2String(account_id) + "," + utils.Int2String(user_id) + ",'" + keyword + "','" + field_name + "','" + match_mode + "','" + value_name + "',null,null)").QueryRows(&m)
	if num == 0 && err == nil {
		err = orm.ErrNoRows
	}
	return m, err
}

func (t *PettyCashHeader) GetPettyCashHeader(keyword, field_name, match_mode, value_name string, p, size, allsize int, issue_date, issue_date2 *string, voucher_id *int, company_id, account_id, sales_type_id, status_id, status_gl_id, user_id, report_grup, report_type, search_detail int, field_nameTop string) (u utils.Page, num int64, err error) {
	o := orm.NewOrm()
	var c int
	var querydata []PettyCashRtnHeader

	o.Raw("call sp_PettyCashReportCount(?,?,null,"+utils.Int2String(company_id)+","+utils.Int2String(account_id)+","+utils.Int2String(sales_type_id)+","+utils.Int2String(status_id)+","+utils.Int2String(status_gl_id)+","+utils.Int2String(user_id)+","+utils.Int2String(report_grup)+","+utils.Int2String(report_type)+",'"+keyword+"',"+utils.Int2String(search_detail)+",'"+field_nameTop+"','"+field_name+"','"+match_mode+"','"+value_name+"',null,null)", &issue_date, &issue_date2).QueryRow(&c)
	if allsize == 1 && c > 0 {
		size = c
	}
	num, err = o.Raw("call sp_PettyCashReport(?,?,null,"+utils.Int2String(company_id)+","+utils.Int2String(account_id)+","+utils.Int2String(sales_type_id)+","+utils.Int2String(status_id)+","+utils.Int2String(status_gl_id)+","+utils.Int2String(user_id)+","+utils.Int2String(report_grup)+","+utils.Int2String(report_type)+",'"+keyword+"',"+utils.Int2String(search_detail)+",'"+field_nameTop+"','"+field_name+"','"+match_mode+"','"+value_name+"', "+utils.Int2String(size)+", "+utils.Int2String((p-1)*size)+")", &issue_date, &issue_date2).QueryRows(&querydata)

	if num == 0 && err == nil {
		err = orm.ErrNoRows
	}

	return utils.Pagination(int(c), p, size, querydata), num, err
}

func (t *PettyCashHeader) GetPettyCash(keyword, field_name, match_mode, value_name string, p, size, allsize int, issue_date, issue_date2 *string, voucher_id *int, company_id, account_id, sales_type_id, status_id, status_gl_id, user_id, report_grup, report_type, search_detail int, field_nameTop string) (u utils.Page, num int64, err error) {
	o := orm.NewOrm()
	var c int
	var querydata []PettyCashDailyJson

	o.Raw("call sp_PettyCashReportCount(?,?,null,"+utils.Int2String(company_id)+","+utils.Int2String(account_id)+","+utils.Int2String(sales_type_id)+","+utils.Int2String(status_id)+","+utils.Int2String(status_gl_id)+","+utils.Int2String(user_id)+","+utils.Int2String(report_grup)+","+utils.Int2String(report_type)+",'"+keyword+"',"+utils.Int2String(search_detail)+",'"+field_nameTop+"','"+field_name+"','"+match_mode+"','"+value_name+"',null,null)", &issue_date, &issue_date2).QueryRow(&c)
	if allsize == 1 && c > 0 {
		size = c
	}
	num, err = o.Raw("call sp_PettyCashReport(?,?,null,"+utils.Int2String(company_id)+","+utils.Int2String(account_id)+","+utils.Int2String(sales_type_id)+","+utils.Int2String(status_id)+","+utils.Int2String(status_gl_id)+","+utils.Int2String(user_id)+","+utils.Int2String(report_grup)+","+utils.Int2String(report_type)+",'"+keyword+"',"+utils.Int2String(search_detail)+",'"+field_nameTop+"','"+field_name+"','"+match_mode+"','"+value_name+"', "+utils.Int2String(size)+", "+utils.Int2String((p-1)*size)+")", &issue_date, &issue_date2).QueryRows(&querydata)

	if num == 0 && err == nil {
		err = orm.ErrNoRows
	}

	return utils.Pagination(int(c), p, size, querydata), num, err
}

func (t *PettyCashHeader) GetVoucher(keyword, field_name, match_mode, value_name string, p, size, allsize int, issue_date, issue_date2 *string, voucher_id int, company_id, account_id, sales_type_id, status_id, status_gl_id, user_id, report_grup, report_type, search_detail int, field_nameTop string) (m []PettyCashVoucherJson, err error) {
	o := orm.NewOrm()
	var num int64

	num, err = o.Raw("call sp_PettyCashReport(?,?,"+utils.Int2String(voucher_id)+","+utils.Int2String(company_id)+","+utils.Int2String(account_id)+","+utils.Int2String(sales_type_id)+","+utils.Int2String(status_id)+","+utils.Int2String(status_gl_id)+","+utils.Int2String(user_id)+","+utils.Int2String(report_grup)+","+utils.Int2String(report_type)+",'"+keyword+"',"+utils.Int2String(search_detail)+",'"+field_nameTop+"','"+field_name+"','"+match_mode+"','"+value_name+"', null,null)", &issue_date, &issue_date2).QueryRows(&m)
	if num == 0 && err == nil {
		err = orm.ErrNoRows
	}
	return m, err
}
