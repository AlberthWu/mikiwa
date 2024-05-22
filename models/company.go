package models

import (
	"errors"
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
)

func (t *Company) TableName() string {
	return "companies"
}

func Companies() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(Company))
}

func init() {
	orm.RegisterModel(new(Company))
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
)

func GetAllInternalLimit(keyword string) (m []CompanyListRtnJson, err error) {
	var companies []Company
	cond := orm.NewCondition()
	cond1 := cond.And("deleted_at__isnull", true).And("CompanyTypes__TypeId__Id", 1)
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
