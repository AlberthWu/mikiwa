package models

import (
	"errors"
	"mikiwa/utils"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type (
	Company struct {
		Id              int             `json:"id"  orm:"column(id);auto;pk"`
		Parent          int             `json:"parent" orm:"column(parent_id);default(0)"`
		Code            string          `json:"code" orm:"column(code);size(5);unique" Valid:"Required;MixSize(3);MaxSize(5)"`
		Name            string          `json:"name" orm:"column(name);size(200)" Valid:"Required"`
		Phone           string          `json:"phone" orm:"column(phone);size(20);null"`
		Fax             string          `json:"fax" orm:"column(fax);size(20);null"`
		Npwp            string          `json:"npwp" orm:"column(npwp);size(100);null"`
		NpwpName        string          `json:"npwp_name" orm:"column(npwp_name);size(200);null"`
		NpwpAddress     string          `json:"npwp_address" orm:"column(npwp_address);null"`
		Email           string          `json:"email" orm:"column(email);size(100);null" `
		Terms           int             `json:"terms" orm:"column(terms)"`
		Credit          float64         `json:"credit" orm:"column(credit);digits(18);decimals(2)"`
		Address         string          `json:"address" orm:"column(address);null"`
		CityId          int             `json:"city_id" orm:"column(city_id)"`
		StateId         int             `json:"state_id" orm:"column(state_id)"`
		DistrictId      int             `json:"district_id" orm:"column(district_id)"`
		Zip             string          `json:"zip" orm:"column(zip);size(5)"`
		IsPo            int8            `json:"is_po" orm:"column(is_po);default(1)"`
		IsTax           int8            `json:"is_tax" orm:"column(is_tax);default(1)"`
		IsCash          int8            `json:"is_cash" orm:"column(is_cash);default(0)"`
		IsReceipt       int8            `json:"is_receipt" orm:"column(is_receipt);default(0)"`
		PriceMethod     int8            `json:"price_method" orm:"column(price_method);default(1)"`
		Status          int8            `json:"status" orm:"column(status);default(0)"`
		Position        int8            `json:"position" orm:"column(position)"`
		BankId          int             `json:"bank_id" orm:"column(bank_id)"`
		BankName        string          `json:"bank_name" orm:"column(bank_name)"`
		BankNo          string          `json:"bank_no" orm:"column(bank_no)"`
		BankAccountName string          `json:"bank_account_name" orm:"column(bank_account_name)"`
		BankBranch      string          `json:"bank_branch" orm:"column(bank_branch)"`
		CreatedAt       time.Time       `json:"created_at" orm:"column(created_at);type(timestamp);auto_now_add"`
		UpdatedAt       time.Time       `json:"updated_at" orm:"column(updated_at);type(timestamp);auto_now"`
		DeletedAt       time.Time       `json:"deleted_at" orm:"column(deleted_at);type(timestamp);null"`
		CompanyTypes    []*CompanyTypes `json:"-" orm:"reverse(many);rel_through(mikiwa/models.CompanyCompanyType)"`
	}

	Plant struct {
		Id          int       `json:"id"  orm:"column(id);auto;pk"`
		CompanyId   int       `json:"company_id" orm:"column(company_id)"`
		Name        string    `json:"name" orm:"column(name);size(200)"`
		Pic         string    `json:"pic" orm:"column(pic);size(200)"`
		Phone       string    `json:"phone" orm:"column(phone);size(200)"`
		Fax         string    `json:"fax" orm:"column(fax);size(200)"`
		Address     string    `json:"address" orm:"column(address);type(text)"`
		IsDo        int8      `json:"is_do" orm:"column(is_do)"`
		IsPo        int8      `json:"is_po" orm:"column(is_po)"`
		IsSchedule  int8      `json:"is_schedule" orm:"column(is_schedule)"`
		IsReceipt   int8      `json:"is_receipt" orm:"column(is_receipt);default(0)"`
		PriceMethod int8      `json:"price_method" orm:"column(price_method);default(1)"`
		Status      int8      `json:"status" orm:"column(status)"`
		CreatedAt   time.Time `json:"created_at" orm:"column(created_at);type(timestamp);auto_now_add"`
		UpdatedAt   time.Time `json:"updated_at" orm:"column(updated_at);type(timestamp);auto_now"`
		DeletedAt   time.Time `json:"deleted_at" orm:"column(deleted_at);type(timestamp);null"`
	}
)

func (t *Company) TableName() string {
	return "companies"
}

func Companies() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(Company))
}

func (t *Plant) TableName() string {
	return "plants"
}

func Plants() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(Plant))
}

func init() {
	orm.RegisterModel(new(Company), new(Plant))
}

func (t *Company) Insert(m Company) (*Company, error) {
	o := orm.NewOrm()

	if _, err := o.Insert(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (t *Company) Update(fields ...string) error {
	o := orm.NewOrm()
	if _, err := o.Update(t, fields...); err != nil {
		return err
	}
	return nil
}

func (t *Plant) Insert(m Plant) (*Plant, error) {
	o := orm.NewOrm()

	if _, err := o.Insert(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (t *Plant) Update(fields ...string) error {
	o := orm.NewOrm()
	if _, err := o.Update(t, fields...); err != nil {
		return err
	}
	return nil
}

func (t *Company) CheckCode(id int, code string) bool {
	exist := Companies().Exclude("id", id).Filter("code", code).Exist()
	return exist
}

type (
	CompanyListRtnJson struct {
		Id              int    `json:"id"`
		Code            string `json:"code"`
		Name            string `json:"name"`
		Status          int    `json:"status"`
		BankId          int    `json:"bank_id"`
		BankName        string `json:"bank_name"`
		BankNo          string `json:"bank_no"`
		BankAccountName string `json:"bank_account_name"`
		BankBranch      string `json:"bank_branch"`
		IsTax           int8   `json:"is_tax"`
		IsPo            int8   `json:"is_po"`
		TypeId          int8   `json:"type_id"`
		Terms           int    `json:"terms"`
	}

	CompanyDetailRtnJson struct {
		Id              int            `json:"id"`
		Parent          int            `json:"parent_id"`
		Code            string         `json:"code"`
		Name            string         `json:"name"`
		Email           string         `json:"email"`
		Phone           string         `json:"phone"`
		Fax             string         `json:"fax"`
		Npwp            string         `json:"npwp"`
		NpwpName        string         `json:"npwp_name"`
		NpwpAddress     string         `json:"npwp_address"`
		Terms           int            `json:"terms"`
		Credit          float64        `json:"credit"`
		IsCash          int8           `json:"is_cash"`
		IsPo            int8           `json:"is_po"`
		IsTax           int8           `json:"is_tax"`
		IsReceipt       int8           `json:"is_receipt"`
		Status          int8           `json:"status"`
		Address         string         `json:"address"`
		Teritory        string         `json:"teritory"`
		Zip             string         `json:"zip"`
		CityId          int            `json:"city_id"`
		CityName        string         `json:"city_name"`
		BankId          int            `json:"bank_id"`
		BankName        string         `json:"bank_name"`
		BankNo          string         `json:"bank_no"`
		BankAccountName string         `json:"bank_account_name"`
		BankBranch      string         `json:"bank_branch"`
		Plant           []PlantRtnJson `json:"plant"`
		CompanyType     []CompanyTy    `json:"company_type"`
	}

	CompanyDetailReturn struct {
		Id              int            `json:"id"`
		Parent          int            `json:"parent_id"`
		Code            string         `json:"code"`
		Name            string         `json:"name"`
		Email           string         `json:"email"`
		Phone           string         `json:"phone"`
		Fax             string         `json:"fax"`
		Npwp            string         `json:"npwp"`
		NpwpName        string         `json:"npwp_name"`
		NpwpAddress     string         `json:"npwp_address"`
		Terms           int            `json:"terms"`
		Credit          float64        `json:"credit"`
		IsCash          int8           `json:"is_cash"`
		IsPo            int8           `json:"is_po"`
		IsTax           int8           `json:"is_tax"`
		IsReceipt       int8           `json:"is_receipt"`
		Status          int8           `json:"status"`
		Address         string         `json:"address"`
		Teritory        string         `json:"teritory"`
		Zip             string         `json:"zip"`
		CityId          CityRtnJson    `json:"city_id"`
		StateId         CityRtnJson    `json:"state_id"`
		DistrictId      CityRtnJson    `json:"district_id"`
		BankId          BankReturn     `json:"bank_id"`
		BankNo          string         `json:"bank_no"`
		BankAccountName string         `json:"bank_account_name"`
		BankBranch      string         `json:"bank_branch"`
		Plant           []PlantRtnJson `json:"plant"`
		CompanyType     []CompanyTy    `json:"company_type"`
	}

	PlantRtnJson struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		Pic         string `json:"pic"`
		Address     string `json:"address"`
		IsDo        int8   `json:"is_do"`
		IsPo        int8   `json:"is_po"`
		IsSchedule  int8   `json:"is_schedule"`
		IsReceipt   int8   `json:"is_receipt"`
		PriceMethod int8   `json:"price_method"`
		Status      int8   `json:"status"`
	}
)

func (t *Company) GetById(id, user_id int) (m *CompanyDetailReturn, err error) {
	o := orm.NewOrm()
	var companies Company
	cond := orm.NewCondition()
	cond1 := cond.And("deleted_at__isnull", true).And("id", id)
	qs := Companies().SetCond(cond1)
	err = qs.One(&companies)

	var bankrtn BankReturn
	o.Raw("select id,name,swift from banks where id  =" + utils.Int2String(companies.BankId) + " ").QueryRow(&bankrtn)

	var cityrtn CityRtnJson
	o.Raw("select c1.id,concat(c1.name,', ',c2.name,', ', c3.name) name  from cities c1 Left Join(select * from cities) c2 On c2.id=c1.parent_id	Left Join(select * from cities) c3 On c3.id=c2.parent_id where c1.id  =" + utils.Int2String(companies.CityId) + " ").QueryRow(&cityrtn)

	var statertn CityRtnJson
	o.Raw("select id,name  from cities where id  =" + utils.Int2String(companies.StateId) + " ").QueryRow(&statertn)

	var districtrtn CityRtnJson
	o.Raw("select id,name  from cities where id  =" + utils.Int2String(companies.DistrictId) + " ").QueryRow(&districtrtn)

	var plantrtn []PlantRtnJson
	var plants []Plant
	Plants().Filter("company_id", companies.Id).All(&plants)
	for _, plantval := range plants {
		plantrtn = append(plantrtn, PlantRtnJson{
			Id:         plantval.Id,
			Name:       plantval.Name,
			Address:    plantval.Address,
			Status:     plantval.Status,
			IsDo:       plantval.IsDo,
			IsPo:       plantval.IsPo,
			IsSchedule: plantval.IsSchedule,
			IsReceipt:  plantval.IsReceipt,
		})
	}

	qb, _ := orm.NewQueryBuilder("mysql")
	var companyty []CompanyTy
	qb.Select("type_id", "name").
		From("company_type ct").
		LeftJoin("company_types cts").On("ct.type_id=cts.id").
		Where("ct.company_id = " + utils.Int2String(companies.Id) + "").
		OrderBy("cts.position")
	sql := qb.String()
	o.Raw(sql).QueryRows(&companyty)

	m = &CompanyDetailReturn{
		Id:              companies.Id,
		Parent:          companies.Parent,
		Code:            companies.Code,
		Name:            companies.Name,
		Email:           companies.Email,
		Phone:           companies.Phone,
		Fax:             companies.Fax,
		Address:         companies.Address,
		Npwp:            companies.Npwp,
		NpwpName:        companies.NpwpAddress,
		NpwpAddress:     companies.NpwpAddress,
		Terms:           companies.Terms,
		Credit:          companies.Credit,
		IsCash:          companies.IsCash,
		IsPo:            companies.IsPo,
		IsTax:           companies.IsTax,
		IsReceipt:       companies.IsReceipt,
		Status:          companies.Status,
		Teritory:        cityrtn.Name,
		DistrictId:      districtrtn,
		CityId:          cityrtn,
		StateId:         statertn,
		Zip:             companies.Zip,
		BankId:          bankrtn,
		BankNo:          companies.BankNo,
		BankAccountName: companies.BankAccountName,
		BankBranch:      companies.BankBranch,
		Plant:           plantrtn,
		CompanyType:     companyty}

	return m, err
}

func (t *Company) GetAll(keyword, field_name, match_mode, value_name string, p, size, allsize, user_id, company_type int, updated_at *string) (u utils.Page, err error) {
	o := orm.NewOrm()
	var querydata []CompanyDetailRtnJson

	var m []CompanyDetailRtnJson
	var c int

	o.Raw("call sp_CompanyCount(?,"+utils.Int2String(company_type)+","+utils.Int2String(user_id)+",'"+keyword+"','"+field_name+"','"+match_mode+"','"+value_name+"',null,null)", &updated_at).QueryRow(&c)

	if allsize == 1 && c > 0 {
		size = c
	}

	_, err = o.Raw("call sp_Company(?,"+utils.Int2String(company_type)+","+utils.Int2String(user_id)+",'"+keyword+"','"+field_name+"','"+match_mode+"','"+value_name+"',"+utils.Int2String(size)+", "+utils.Int2String((p-1)*size)+")", &updated_at).QueryRows(&m)

	for _, val := range m {
		clist := t.GetDetailCt(val.Id)
		olist := t.GetDetail(val.Id)
		querydata = append(querydata, CompanyDetailRtnJson{
			Id:              val.Id,
			Parent:          val.Parent,
			Code:            val.Code,
			Name:            val.Name,
			Email:           val.Email,
			Phone:           val.Phone,
			Fax:             val.Fax,
			Npwp:            val.Npwp,
			NpwpName:        val.NpwpName,
			NpwpAddress:     val.NpwpAddress,
			Terms:           val.Terms,
			Credit:          val.Credit,
			IsCash:          val.IsCash,
			IsPo:            val.IsPo,
			IsTax:           val.IsTax,
			IsReceipt:       val.IsReceipt,
			Status:          val.Status,
			Address:         val.Address,
			Teritory:        val.Teritory,
			Zip:             val.Zip,
			CityId:          val.CityId,
			CityName:        val.CityName,
			BankId:          val.BankId,
			BankName:        val.BankName,
			BankNo:          val.BankNo,
			BankAccountName: val.BankAccountName,
			BankBranch:      val.BankBranch,
			Plant:           olist,
			CompanyType:     clist,
		})
	}

	if c == 0 && err == nil {
		err = orm.ErrNoRows
	}

	return utils.Pagination(int(c), p, size, querydata), err
}

func (c *Company) GetDetail(id int) (m []PlantRtnJson) {
	o := orm.NewOrm()
	o.Raw("select id,name,pic,address,is_do,is_po,is_schedule,is_receipt,price_method,status from plants where deleted_at is null and company_id = " + utils.Int2String(id) + ")").QueryRows(&m)
	return m
}

func (c *Company) GetDetailCt(id int) (m []CompanyTy) {
	o := orm.NewOrm()
	o.Raw("select id,name from company_types where id in (select type_id from company_type where company_id = " + utils.Int2String(id) + ")").QueryRows(&m)
	return m
}

func (t *Company) GetAllList(keyword string, company_type int) (m []CompanyListRtnJson, err error) {
	var companies []Company
	cond := orm.NewCondition()
	cond1 := cond.And("deleted_at__isnull", true).And("CompanyTypes__TypeId__Id", company_type).And("status", 1)
	qs := Companies().SetCond(cond1)
	cond2 := cond.AndCond(cond1).AndCond(cond.Or("name__icontains", keyword).Or("code__icontains", keyword))
	qs = qs.SetCond(cond2).RelatedSel()
	_, err = qs.Limit(100).Offset(0).All(&companies)

	var companylist []CompanyListRtnJson
	for _, val := range companies {

		companylist = append(companylist, CompanyListRtnJson{
			Id:              val.Id,
			Code:            val.Code,
			Name:            val.Name,
			Status:          int(val.Status),
			BankId:          val.BankId,
			BankName:        val.BankName,
			BankNo:          val.BankNo,
			BankAccountName: val.BankAccountName,
			BankBranch:      val.BankBranch,
			IsTax:           val.IsTax,
		})
	}

	if len(companylist) == 0 {
		return companylist, errors.New("no data")
	}
	return companylist, err
}
