package models

import (
	"strconv"
	"time"

	"mikiwa/utils"

	"github.com/beego/beego/v2/client/orm"
)

type (
	GlAccountType struct {
		Id               int       `json:"id" orm:"column(id);auto;pk"`
		AccountCode      string    `json:"account_code" orm:"column(account_code)"`
		Name             string    `json:"name" orm:"column(name)"`
		JournalPosition  string    `json:"journal_position" orm:"column(journal_position)"`
		ComponentAccount string    `json:"component_account" orm:"column(component_account)"`
		CreatedAt        time.Time `json:"created_at" orm:"column(created_at);type(timestamp);auto_now_add"`
		UpdatedAt        time.Time `json:"updated_at" orm:"column(updated_at);type(timestamp);auto_now"`
		DeletedAt        time.Time `json:"deleted_at" orm:"column(deleted_at);type(timestamp);null"`
		CreatedBy        string    `json:"created_by" orm:"column(created_by)"`
		UpdatedBy        string    `json:"updated_by" orm:"column(updated_by)"`
		DeletedBy        string    `json:"deleted_by" orm:"column(deleted_by)"`
	}

	CharOfAccount struct {
		Id              int        `json:"id" orm:"column(id);auto;pk"`
		AccountTypeId   int        `json:"account_type_id" orm:"column(account_type_id)"`
		AccountTypeName string     `json:"account_type_name" orm:"column(account_type_name)"`
		CompanyId       int        `json:"company_id" orm:"column(company_id)"`
		CompanyCode     string     `json:"company_code" orm:"column(company_code)"`
		CompanyName     string     `json:"company_name" orm:"column(company_name)"`
		LevelNo         int        `json:"level_no" orm:"column(level_no)"`
		ParentId        int        `json:"parent_id" orm:"column(parent_id)"`
		ParentCode      string     `json:"parent_code" orm:"column(parent_code)"`
		CodeCoa         string     `json:"code_coa" orm:"column(code_coa)"`
		NameCoa         string     `json:"name_coa" orm:"column(name_coa)"`
		CodeOut         string     `json:"code_out"  orm:"column(code_out)"`
		CodeIn          string     `json:"code_in" orm:"column(code_in)"`
		EffectiveDate   time.Time  `json:"effective_date" orm:"column(effective_date);type(date)"`
		ExpiredDate     *time.Time `json:"expired_date" orm:"column(expired_date);type(date);null"`
		StatusId        int8       `json:"status_id" orm:"column(status_id);default(0)"`
		JournalPosition string     `json:"journal_position" orm:"column(journal_position)"`
		IsHeader        int8       `json:"is_header" orm:"column(is_header);default(0)"`
		CreatedAt       time.Time  `json:"created_at" orm:"column(created_at);type(timestamp);auto_now_add"`
		UpdatedAt       time.Time  `json:"updated_at" orm:"column(updated_at);type(timestamp);auto_now"`
		DeletedAt       time.Time  `json:"deleted_at" orm:"column(deleted_at);type(timestamp);null"`
		CreatedBy       string     `json:"created_by" orm:"column(created_by)"`
		UpdatedBy       string     `json:"updated_by" orm:"column(updated_by)"`
		DeletedBy       string     `json:"deleted_by" orm:"column(deleted_by)"`
	}
)

func (t *GlAccountType) TableName() string {
	return "gl_account_type"
}

func GlAccountTypes() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(GlAccountType))
}

func (t *CharOfAccount) TableName() string {
	return "chart_of_accounts"
}

func ChartOfAccounts() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(CharOfAccount))
}

func init() {
	orm.RegisterModel(new(GlAccountType), new(CharOfAccount))
}

func (t *GlAccountType) Insert(m GlAccountType) (*GlAccountType, error) {
	o := orm.NewOrm()

	if _, err := o.Insert(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (t *GlAccountType) Update(fields ...string) error {
	o := orm.NewOrm()
	if _, err := o.Update(t, fields...); err != nil {
		return err
	}
	return nil
}

func (t *CharOfAccount) Insert(m CharOfAccount) (*CharOfAccount, error) {
	o := orm.NewOrm()

	if _, err := o.Insert(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (t *CharOfAccount) Update(fields ...string) error {
	o := orm.NewOrm()
	if _, err := o.Update(t, fields...); err != nil {
		return err
	}
	return nil
}

type (
	CoaRtnJson struct {
		Id              int     `json:"id"`
		EffectiveDate   string  `json:"effective_date"`
		ExpiredDate     *string `json:"expired_date"`
		AccountTypeId   int     `json:"account_type_id"`
		AccountTypeName string  `json:"account_type_name"`
		CompanyId       int     `json:"company_id"`
		CompanyCode     string  `json:"company_code"`
		CompanyName     string  `json:"company_name"`
		SalesTypeId     int     `json:"sales_type_id"`
		SalesTypeName   string  `json:"sales_type_name"`
		LevelNo         int     `json:"level_no"`
		ParentId        int     `json:"parent_id"`
		ParentCode      string  `json:"parent_code"`
		ParentName      string  `json:"parent_name"`
		CodeCoa         string  `json:"code_coa"`
		NameCoa         string  `json:"name_coa"`
		CodeOut         string  `json:"code_out"`
		CodeIn          string  `json:"code_in"`
		StatusId        int8    `json:"status_id"`
		JournalPosition string  `json:"journal_position"`
		StatusData      string  `json:"status_data"`
		IsHeader        int8    `json:"is_header"`
	}

	CoaListRtnJson struct {
		Id      int    `json:"id"`
		CodeCoa string `json:"code_coa"`
		NameCoa string `json:"name_coa"`
	}

	CoaRtn struct {
		Id              int                `json:"id"`
		EffectiveDate   string             `json:"effective_date"`
		ExpiredDate     *string            `json:"expired_date"`
		AccountTypeId   GlAccountType      `json:"account_type_id"`
		CompanyId       CompanyListRtnJson `json:"company_id"`
		LevelNo         int                `json:"level_no"`
		ParentId        CoaListRtnJson     `json:"parent_id"`
		CodeCoa         string             `json:"code_coa"`
		NameCoa         string             `json:"name_coa"`
		CodeOut         string             `json:"code_out"`
		CodeIn          string             `json:"code_in"`
		StatusId        int8               `json:"status_id"`
		JournalPosition string             `json:"journal_position"`
		StatusData      string             `json:"status_data"`
		IsHeader        int8               `json:"is_header"`
	}
)

func (t *CharOfAccount) CheckCode(id, company_id int, code string) bool {
	exist := ChartOfAccounts().Exclude("id", id).Filter("company_id", company_id).Filter("code_coa", code).Filter("deleted_at__isnull", true).Exist()
	return exist
}

func (t *GlAccountType) GetById(id int) (m *GlAccountType, err error) {
	m = &GlAccountType{}
	cond := orm.NewCondition()
	cond1 := cond.And("id", id)
	qs := GlAccountTypes().SetCond(cond1)

	if err = qs.One(m); err != nil {
		return nil, err
	}

	return m, err
}

func (t *GlAccountType) GetAll(keyword string, p, size int) (u utils.Page, err error) {

	var details []GlAccountType
	var d int64
	cond := orm.NewCondition()
	cond1 := cond.And("name__icontains", keyword).And("deleted_at__isnull", true)
	qs := GlAccountTypes().SetCond(cond1)

	d, err = qs.Limit(size).Offset((p - 1) * size).OrderBy("-id").All(&details)
	count, _ := qs.Limit(-1).Count()
	c, _ := strconv.Atoi(strconv.FormatInt(count, 10))

	if err == nil && d == 0 {
		err = orm.ErrNoRows
	}
	return utils.Pagination(c, p, size, details), err
}

func (t *CharOfAccount) GetById(id int) (m *CoaRtn, err error) {
	o := orm.NewOrm()
	var detail CharOfAccount
	cond := orm.NewCondition()
	cond1 := cond.And("deleted_at__isnull", true).And("id", id)
	qs := ChartOfAccounts().SetCond(cond1)
	err = qs.One(&detail)

	var companyrtn CompanyListRtnJson
	o.Raw("select id,code,name,status from companies where id  = ?", detail.CompanyId).QueryRow(&companyrtn)

	var parentrtn CoaListRtnJson
	o.Raw("select id,code_coa,name_coa from chart_of_accounts where id = " + utils.Int2String(detail.ParentId)).QueryRow(&parentrtn)

	var accounttypertn GlAccountType
	GlAccountTypes().Filter("id", detail.AccountTypeId).One(&accounttypertn)

	var expiredDate *string
	if detail.ExpiredDate == nil {
		expiredDate = nil
	} else {
		thedate := detail.ExpiredDate.Format("2006-01-02")
		expiredDate = &thedate
	}

	m = &CoaRtn{
		Id:              detail.Id,
		EffectiveDate:   detail.EffectiveDate.Format("2006-01-02"),
		ExpiredDate:     expiredDate,
		AccountTypeId:   accounttypertn,
		CompanyId:       companyrtn,
		LevelNo:         detail.LevelNo,
		ParentId:        parentrtn,
		CodeCoa:         detail.CodeCoa,
		NameCoa:         detail.NameCoa,
		CodeOut:         detail.CodeOut,
		CodeIn:          detail.CodeIn,
		StatusId:        detail.StatusId,
		JournalPosition: detail.JournalPosition,
		IsHeader:        detail.IsHeader,
	}

	return m, err
}

func (t *CharOfAccount) GetAll(keyword, field_name, match_mode, value_name string, p, size, parent_level_no, level_no, account_type_id, company_id, sales_type_id, user_id int, issue_date, updated_at *string) (u utils.Page, err error) {

	o := orm.NewOrm()
	var c int
	var querydata []CoaRtnJson

	o.Raw("call sp_ChartOfAccountCount(?,?,"+utils.Int2String(parent_level_no)+","+utils.Int2String(level_no)+","+utils.Int2String(account_type_id)+","+utils.Int2String(company_id)+","+utils.Int2String(sales_type_id)+","+utils.Int2String(user_id)+",'"+keyword+"','"+field_name+"','"+match_mode+"','"+value_name+"',null,null)", &issue_date, &updated_at).QueryRow(&c)
	d, err := o.Raw("call sp_ChartOfAccount(?,?,"+utils.Int2String(parent_level_no)+","+utils.Int2String(level_no)+","+utils.Int2String(account_type_id)+","+utils.Int2String(company_id)+","+utils.Int2String(sales_type_id)+","+utils.Int2String(user_id)+",'"+keyword+"','"+field_name+"','"+match_mode+"','"+value_name+"',"+utils.Int2String(size)+","+utils.Int2String((p-1)*size)+")", &issue_date, &updated_at).QueryRows(&querydata)
	var detailrtn []CoaRtn
	for _, val := range querydata {
		var companyrtn CompanyListRtnJson
		o.Raw("select id,code,name,status from companies where id  = ?", val.CompanyId).QueryRow(&companyrtn)

		var parentrtn CoaListRtnJson
		o.Raw("select id,code_coa,name_coa from chart_of_accounts where id = " + utils.Int2String(val.ParentId)).QueryRow(&parentrtn)

		var accounttypertn GlAccountType
		GlAccountTypes().Filter("id", val.AccountTypeId).One(&accounttypertn)

		detailrtn = append(detailrtn, CoaRtn{
			Id:              val.Id,
			EffectiveDate:   val.EffectiveDate,
			ExpiredDate:     val.ExpiredDate,
			AccountTypeId:   accounttypertn,
			CompanyId:       companyrtn,
			LevelNo:         val.LevelNo,
			ParentId:        parentrtn,
			CodeCoa:         val.CodeCoa,
			NameCoa:         val.NameCoa,
			CodeOut:         val.CodeOut,
			CodeIn:          val.CodeIn,
			StatusId:        val.StatusId,
			JournalPosition: val.JournalPosition,
			StatusData:      val.StatusData,
			IsHeader:        val.IsHeader,
		})
	}

	if d == 0 && err == nil {
		err = orm.ErrNoRows
	}
	return utils.Pagination(int(c), p, size, detailrtn), err
}

func (t *CharOfAccount) GetAllLimit(keyword string, parent_level_no, level_no, account_type_id, company_id, sales_type_id, user_id int) (m []CoaRtnJson, err error) {
	o := orm.NewOrm()
	d, err := o.Raw("call sp_ChartOfAccount(now(),null," + utils.Int2String(parent_level_no) + "," + utils.Int2String(level_no) + "," + utils.Int2String(account_type_id) + "," + utils.Int2String(company_id) + "," + utils.Int2String(sales_type_id) + "," + utils.Int2String(user_id) + ",'" + keyword + "',null,null,null,null,null)").QueryRows(&m)

	if d == 0 && err == nil {
		err = orm.ErrNoRows
	}
	return m, err
}

func (t *CharOfAccount) GetAllLimitChild(issue_date, keyword, component string, company_id, account_type_id, sales_type_id, user_id int) (m []CoaRtnJson, err error) {
	o := orm.NewOrm()
	var querydata []CoaRtnJson
	d, err := o.Raw("call sp_ChartOfAccountChild('" + issue_date + "'," + utils.Int2String(company_id) + "," + utils.Int2String(account_type_id) + ",0," + utils.Int2String(sales_type_id) + ",'" + component + "'," + utils.Int2String(user_id) + ",'" + keyword + "',null,null,null,null,null)").QueryRows(&querydata)
	if d == 0 && err == nil {
		err = orm.ErrNoRows
	}
	m = querydata
	return m, err
}

func (t *GlAccountType) GetAllList(keyword, component_account string) (m []GlAccountType, err error) {
	var num int64
	cond := orm.NewCondition()
	if component_account == "" {
		cond = cond.And("name__icontains", keyword).And("deleted_at__isnull", true)
	} else {
		cond = cond.And("component_account", component_account).AndCond(cond.Or("name__icontains", keyword))
	}

	qs := GlAccountTypes().SetCond(cond).OrderBy("id")
	num, err = qs.Limit(100).Offset(0).All(&m)

	if num == 0 && err == nil {
		err = orm.ErrNoRows
	}
	return m, err
}

func (t *GlAccountType) GetAllListAsset(keyword string) (m []GlAccountType, err error) {
	var num int64
	cond := orm.NewCondition()
	cond1 := cond.And("component_account", "Assets").And("deleted_at__isnull", true)
	qs := GlAccountTypes().SetCond(cond1)
	cond2 := cond.AndCond(cond1).AndCond(cond.Or("name__icontains", keyword))
	qs = qs.SetCond(cond2).OrderBy("id")
	num, err = qs.Limit(100).Offset(0).All(&m)

	if num == 0 && err == nil {
		err = orm.ErrNoRows
	}
	return m, err
}

func (t *GlAccountType) GetAllListExpenses(keyword string) (m []GlAccountType, err error) {
	var num int64
	cond := orm.NewCondition()
	cond1 := cond.And("component_account", "Expenses").And("deleted_at__isnull", true)
	qs := GlAccountTypes().SetCond(cond1)
	cond2 := cond.AndCond(cond1).AndCond(cond.Or("name__icontains", keyword))
	qs = qs.SetCond(cond2).OrderBy("id")
	num, err = qs.Limit(100).Offset(0).All(&m)

	if num == 0 && err == nil {
		err = orm.ErrNoRows
	}
	return m, err
}

func (t *GlAccountType) GetAllListLiability(keyword string) (m []GlAccountType, err error) {
	var num int64
	cond := orm.NewCondition()
	cond1 := cond.And("component_account", "Liability").And("deleted_at__isnull", true)
	qs := GlAccountTypes().SetCond(cond1)
	cond2 := cond.AndCond(cond1).AndCond(cond.Or("name__icontains", keyword))
	qs = qs.SetCond(cond2).OrderBy("id")
	num, err = qs.Limit(100).Offset(0).All(&m)

	if num == 0 && err == nil {
		err = orm.ErrNoRows
	}
	return m, err
}

func (t *GlAccountType) GetAllListEquity(keyword string) (m []GlAccountType, err error) {
	var num int64
	cond := orm.NewCondition()
	cond1 := cond.And("component_account", "Equity").And("deleted_at__isnull", true)
	qs := GlAccountTypes().SetCond(cond1)
	cond2 := cond.AndCond(cond1).AndCond(cond.Or("name__icontains", keyword))
	qs = qs.SetCond(cond2).OrderBy("id")
	num, err = qs.Limit(100).Offset(0).All(&m)

	if num == 0 && err == nil {
		err = orm.ErrNoRows
	}
	return m, err
}

func (t *GlAccountType) GetAllListRevenue(keyword string) (m []GlAccountType, err error) {
	var num int64
	cond := orm.NewCondition()
	cond1 := cond.And("component_account", "Revenue").And("deleted_at__isnull", true)
	qs := GlAccountTypes().SetCond(cond1)
	cond2 := cond.AndCond(cond1).AndCond(cond.Or("name__icontains", keyword))
	qs = qs.SetCond(cond2).OrderBy("id")
	num, err = qs.Limit(100).Offset(0).All(&m)

	if num == 0 && err == nil {
		err = orm.ErrNoRows
	}
	return m, err
}

func (t *GlAccountType) GetAllListCogs(keyword string) (m []GlAccountType, err error) {
	var num int64
	cond := orm.NewCondition()
	cond1 := cond.And("component_account", "Cogs").And("deleted_at__isnull", true)
	qs := GlAccountTypes().SetCond(cond1)
	cond2 := cond.AndCond(cond1).AndCond(cond.Or("name__icontains", keyword))
	qs = qs.SetCond(cond2).OrderBy("id")
	num, err = qs.Limit(100).Offset(0).All(&m)

	if num == 0 && err == nil {
		err = orm.ErrNoRows
	}
	return m, err
}
