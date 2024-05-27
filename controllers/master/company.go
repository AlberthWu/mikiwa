package controllers

import (
	"fmt"
	base "mikiwa/controllers"
	"mikiwa/models"
	"mikiwa/utils"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
)

type CompanyController struct {
	base.BaseController
}

func (c *CompanyController) Prepare() {
	c.Ctx.Request.Header.Set("token", "No Aut")
	c.BaseController.Prepare()
}

func (c *CompanyController) Post() {
	var user_id, form_id int
	var err error

	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}

	form_id = base.FormName(form_customer)

	write_aut := models.CheckPrivileges(user_id, form_id, base.Write)
	write_aut = true
	if !write_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Post not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}

	parentId, _ := c.GetInt("parent_id")
	code := strings.TrimSpace(c.GetString("code"))
	name := strings.TrimSpace(c.GetString("name"))
	phone := strings.TrimSpace(c.GetString("phone"))
	fax := strings.TrimSpace(c.GetString("fax"))
	npwp := strings.TrimSpace(c.GetString("npwp"))
	npwpaddress := strings.TrimSpace(c.GetString("npwp_address"))
	npwpname := strings.TrimSpace(c.GetString("npwp_name"))
	email := strings.TrimSpace(c.GetString("email"))
	terms, _ := c.GetInt("terms")
	credit, _ := c.GetFloat("credit")
	address := strings.TrimSpace(c.GetString("address"))
	cityId, _ := c.GetInt("city_id")
	zip := strings.TrimSpace(c.GetString("zip"))
	ispo, _ := c.GetInt8("is_po")
	istax, _ := c.GetInt8("is_tax")
	isReceipt, _ := c.GetInt8("is_receipt")
	iscash, _ := c.GetInt8("is_cash")
	status, _ := c.GetInt8("status")
	company_type := strings.TrimSpace(c.GetString("company_type"))
	business_type := strings.TrimSpace(c.GetString("business_type"))
	bank_id, _ := c.GetInt("bank_id")
	bank_no := strings.TrimSpace(c.GetString("bank_no"))
	bank_account_name := strings.TrimSpace(c.GetString("bank_account_name"))
	bank_branch := strings.TrimSpace(c.GetString("bank_branch"))

	valid := validation.Validation{}
	valid.Required(code, "code").Message("Code is required")
	valid.MinSize(code, 2, "code").Message("Code min char is 2")
	valid.MaxSize(code, 15, "code").Message("Code max char is 15")
	valid.Required(name, "name").Message("Name is required")
	valid.MinSize(name, 3, "name").Message("Name min char is 3")
	valid.MaxSize(name, 250, "name").Message("Name max char is 250")
	if len(email) > 0 {
		valid.Email(email, "email").Message("Invalid email format")
	}
	valid.Required(bank_id, "bank_id").Message("is required")
	valid.Required(cityId, "city_id").Message("is required")
	valid.Required(company_type, "company_type").Message("is required")
	valid.Required(business_type, "business_type").Message("is required")

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

	if t_company.CheckCode(0, code) {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("code : '%v' has been REGISTERED", code))
		c.ServeJSON()
		return
	}

	var cities models.City
	models.Cities().Filter("id", cityId).One(&cities)

	var states models.City
	models.Cities().Filter("id", cities.ParentId).One(&states)

	var districts models.City
	models.Cities().Filter("id", states.ParentId).One(&districts)

	var ctype models.CompanyTypes
	ctypearray := strings.Split(company_type, ",")
	for _, ch := range ctypearray {
		err = models.CompanyTypess().Filter("id", utils.String2Int(strings.TrimSpace(ch))).One(&ctype)
		if err == orm.ErrNoRows {
			//utils.ReturnHTTPSuccessWithMessage(&c.Controller, 401, "Invalid company type", map[string]string{"company_type": fmt.Sprintf("Company type is not valid")})
			utils.ReturnHTTPError(&c.Controller, 401, "Company type unregistered/Illegal data")
			c.ServeJSON()
			return
		}
	}

	var businessUnit models.BusinessUnit
	buIdArray := strings.Split(business_type, ",")
	for _, bu := range buIdArray {
		err = models.BusinessUnits().Filter("id", utils.String2Int(strings.TrimSpace(bu))).One(&businessUnit)
		if err == orm.ErrNoRows {
			// utils.ReturnHTTPSuccessWithMessage(&c.Controller, 401, "Invalid data", map[string]string{"business_unit": fmt.Sprintf("Business Unit is not valid")})
			utils.ReturnHTTPError(&c.Controller, 401, "Business unit unregistered/Illegal data")
			c.ServeJSON()
			return
		}
	}

	var banks models.Bank
	err = models.Banks().Filter("id", bank_id).One(&banks)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Bank unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	t_company = models.Company{
		Parent:          parentId,
		Code:            code,
		Name:            name,
		Phone:           phone,
		Fax:             fax,
		Npwp:            npwp,
		NpwpAddress:     npwpaddress,
		NpwpName:        npwpname,
		Email:           email,
		Terms:           terms,
		Credit:          credit,
		Address:         address,
		CityId:          cityId,
		StateId:         states.Id,
		DistrictId:      districts.Id,
		Zip:             zip,
		BankId:          bank_id,
		BankName:        banks.Name,
		BankNo:          bank_no,
		BankAccountName: bank_account_name,
		BankBranch:      bank_branch,
		IsPo:            ispo,
		IsCash:          iscash,
		IsTax:           istax,
		IsReceipt:       isReceipt,
		PriceMethod:     1,
		Status:          status,
	}

	d, err_ := t_company.Insert(t_company)
	errcode, errmessage = base.DecodeErr(err_)
	if err_ != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
	} else {
		_, err_ := models.InsertCType(d.Id, company_type)
		errcode, errmessage = base.DecodeErr(err_)
		if err_ != nil {
			utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
			c.ServeJSON()
			return
		}

		_, err_ = t_company_business_unit.InsertM2M(d.Id, business_type)
		errcode, errmessage = base.DecodeErr(err_)
		if err_ != nil {
			utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
			c.ServeJSON()
			return
		}

		v, err_ := t_company.GetById(d.Id, user_id)
		errcode, errmessage = base.DecodeErr(err_)
		if err_ != nil {
			c.Ctx.ResponseWriter.WriteHeader(errcode)
			utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
		} else {
			utils.ReturnHTTPSuccessWithMessage(&c.Controller, errcode, errmessage, v)
		}
	}
	c.ServeJSON()
}

func (c *CompanyController) Put() {
	var user_id, form_id int
	var err error
	var deletedAt string
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}

	form_id = base.FormName(form_customer)

	put_aut := models.CheckPrivileges(user_id, form_id, base.Write)
	put_aut = true
	if !put_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Post not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	var querydata models.Company
	err = models.Companies().Filter("id", id).One(&querydata)
	if err == orm.ErrNoRows {
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

	deletedAt = querydata.DeletedAt.Format("2006-01-02")
	if deletedAt != "0001-01-01" {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been DELETED", querydata.Code))
		c.ServeJSON()
		return
	}

	parentId, _ := c.GetInt("parent_id")
	code := strings.TrimSpace(c.GetString("code"))
	name := strings.TrimSpace(c.GetString("name"))
	phone := strings.TrimSpace(c.GetString("phone"))
	fax := strings.TrimSpace(c.GetString("fax"))
	npwp := strings.TrimSpace(c.GetString("npwp"))
	npwpaddress := strings.TrimSpace(c.GetString("npwp_address"))
	npwpname := strings.TrimSpace(c.GetString("npwp_name"))
	email := strings.TrimSpace(c.GetString("email"))
	terms, _ := c.GetInt("terms")
	credit, _ := c.GetFloat("credit")
	address := strings.TrimSpace(c.GetString("address"))
	cityId, _ := c.GetInt("city_id")
	zip := strings.TrimSpace(c.GetString("zip"))
	ispo, _ := c.GetInt8("is_po")
	istax, _ := c.GetInt8("is_tax")
	isReceipt, _ := c.GetInt8("is_receipt")
	iscash, _ := c.GetInt8("is_cash")
	status, _ := c.GetInt8("status")
	company_type := strings.TrimSpace(c.GetString("company_type"))
	business_type := strings.TrimSpace(c.GetString("business_type"))
	bank_id, _ := c.GetInt("bank_id")
	bank_no := strings.TrimSpace(c.GetString("bank_no"))
	bank_account_name := strings.TrimSpace(c.GetString("bank_account_name"))
	bank_branch := strings.TrimSpace(c.GetString("bank_branch"))

	valid := validation.Validation{}
	valid.Required(code, "code").Message("Code is required")
	valid.MinSize(code, 2, "code").Message("Code min char is 2")
	valid.MaxSize(code, 15, "code").Message("Code max char is 15")
	valid.Required(name, "name").Message("Name is required")
	valid.MinSize(name, 3, "name").Message("Name min char is 3")
	valid.MaxSize(name, 250, "name").Message("Name max char is 250")
	if len(email) > 0 {
		valid.Email(email, "email").Message("Invalid email format")
	}
	valid.Required(bank_id, "bank_id").Message("is required")
	valid.Required(cityId, "city_id").Message("is required")
	valid.Required(company_type, "company_type").Message("is required")
	valid.Required(business_type, "business_type").Message("is required")

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

	if t_company.CheckCode(id, code) {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("code : '%v' has been REGISTERED", code))
		c.ServeJSON()
		return
	}

	var cities models.City
	models.Cities().Filter("id", cityId).One(&cities)

	var states models.City
	models.Cities().Filter("id", cities.ParentId).One(&states)

	var districts models.City
	models.Cities().Filter("id", states.ParentId).One(&districts)

	var ctype models.CompanyTypes
	ctypearray := strings.Split(company_type, ",")
	for _, ch := range ctypearray {
		err = models.CompanyTypess().Filter("id", utils.String2Int(strings.TrimSpace(ch))).One(&ctype)
		if err == orm.ErrNoRows {
			//utils.ReturnHTTPSuccessWithMessage(&c.Controller, 401, "Invalid company type", map[string]string{"company_type": fmt.Sprintf("Company type is not valid")})
			utils.ReturnHTTPError(&c.Controller, 401, "Company type unregistered/Illegal data")
			c.ServeJSON()
			return
		}
	}

	var businessUnit models.BusinessUnit
	buIdArray := strings.Split(business_type, ",")
	for _, bu := range buIdArray {
		err = models.BusinessUnits().Filter("id", utils.String2Int(strings.TrimSpace(bu))).One(&businessUnit)
		if err == orm.ErrNoRows {
			utils.ReturnHTTPSuccessWithMessage(&c.Controller, 401, "Invalid data", map[string]string{"business_unit": "Business Unit is not valid"})
			//utils.ReturnHTTPError(&c.Controller, 401, "Company type unregistered/Illegal data")
			c.ServeJSON()
			return
		}
	}

	var banks models.Bank
	err = models.Banks().Filter("id", bank_id).One(&banks)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Bank unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}
	t_company.Id = id
	t_company.Parent = parentId
	t_company.Code = code
	t_company.Name = name
	t_company.Phone = phone
	t_company.Fax = fax
	t_company.Npwp = npwp
	t_company.NpwpAddress = npwpaddress
	t_company.NpwpName = npwpname
	t_company.Email = email
	t_company.Terms = terms
	t_company.Credit = credit
	t_company.Address = address
	t_company.CityId = cityId
	t_company.StateId = states.Id
	t_company.DistrictId = districts.Id
	t_company.Zip = zip
	t_company.BankId = bank_id
	t_company.BankName = banks.Name
	t_company.BankNo = bank_no
	t_company.BankAccountName = bank_account_name
	t_company.BankBranch = bank_branch
	t_company.IsPo = ispo
	t_company.IsCash = iscash
	t_company.IsTax = istax
	t_company.IsReceipt = isReceipt
	t_company.PriceMethod = 1
	t_company.Status = status
	err_ := t_company.Update()

	errcode, errmessage = base.DecodeErr(err_)
	if err_ != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
	} else {
		_, err_ := models.InsertCType(id, company_type)
		errcode, errmessage = base.DecodeErr(err_)
		if err_ != nil {
			utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
			c.ServeJSON()
			return
		}

		_, err_ = t_company_business_unit.InsertM2M(id, business_type)
		errcode, errmessage = base.DecodeErr(err_)
		if err_ != nil {
			utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
			c.ServeJSON()
			return
		}
		v, err_ := t_company.GetById(id, user_id)
		errcode, errmessage = base.DecodeErr(err_)
		if err_ != nil {
			c.Ctx.ResponseWriter.WriteHeader(errcode)
			utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
		} else {
			utils.ReturnHTTPSuccessWithMessage(&c.Controller, errcode, errmessage, v)
		}
	}
	c.ServeJSON()
}

func (c *CompanyController) Delete() {
	var user_id, form_id int
	var err error
	var deletedAt string
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	form_id = base.FormName(form_product)
	delete_aut := models.CheckPrivileges(user_id, form_id, base.Delete)
	if !delete_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Delete not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	var querydata models.Company
	err = models.Companies().Filter("id", id).One(&querydata)
	if err == orm.ErrNoRows {
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

	deletedAt = querydata.DeletedAt.Format("2006-01-02")
	if deletedAt != "0001-01-01" {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been DELETED", querydata.Code))
		c.ServeJSON()
		return
	}
	models.Companies().Filter("id", id).Filter("deleted_at__isnull", true).Update(orm.Params{"deleted_at": utils.GetSvrDate()})

	utils.ReturnHTTPError(&c.Controller, 200, "soft delete success")
	c.ServeJSON()
}

func (c *CompanyController) GetOne() {
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	v, err := t_company.GetById(id, user_id)
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

func (c *CompanyController) GetAllInternal() {
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	var updatedat *string

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

	status_ids := strings.TrimSpace(c.GetString("status_ids"))
	updated_at := strings.TrimSpace(c.GetString("updated_at"))
	company_type_id := 1

	if updated_at == "" {
		updatedat = nil

	} else {
		updatedat = &updated_at
	}

	d, err := t_company.GetAll(keyword, field_name, match_mode, value_name, currentPage, pageSize, allsize, user_id, company_type_id, status_ids, updatedat)
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

func (c *CompanyController) GetAllCustomer() {
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	var updatedat *string

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

	status_ids := strings.TrimSpace(c.GetString("status_ids"))
	updated_at := strings.TrimSpace(c.GetString("updated_at"))
	company_type_id := 2

	if updated_at == "" {
		updatedat = nil

	} else {
		updatedat = &updated_at
	}

	d, err := t_company.GetAll(keyword, field_name, match_mode, value_name, currentPage, pageSize, allsize, user_id, company_type_id, status_ids, updatedat)
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

func (c *CompanyController) GetAllCustOthers() {
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	var updatedat *string

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

	status_ids := strings.TrimSpace(c.GetString("status_ids"))
	updated_at := strings.TrimSpace(c.GetString("updated_at"))
	company_type_id := 3

	if updated_at == "" {
		updatedat = nil

	} else {
		updatedat = &updated_at
	}

	d, err := t_company.GetAll(keyword, field_name, match_mode, value_name, currentPage, pageSize, allsize, user_id, company_type_id, status_ids, updatedat)
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

func (c *CompanyController) GetAllWarehouse() {
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	var updatedat *string

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

	status_ids := strings.TrimSpace(c.GetString("status_ids"))
	updated_at := strings.TrimSpace(c.GetString("updated_at"))
	company_type_id := 4

	if updated_at == "" {
		updatedat = nil

	} else {
		updatedat = &updated_at
	}

	d, err := t_company.GetAll(keyword, field_name, match_mode, value_name, currentPage, pageSize, allsize, user_id, company_type_id, status_ids, updatedat)
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

func (c *CompanyController) GetAllSparepart() {
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	var updatedat *string

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

	status_ids := strings.TrimSpace(c.GetString("status_ids"))
	updated_at := strings.TrimSpace(c.GetString("updated_at"))
	company_type_id := 5

	if updated_at == "" {
		updatedat = nil

	} else {
		updatedat = &updated_at
	}

	d, err := t_company.GetAll(keyword, field_name, match_mode, value_name, currentPage, pageSize, allsize, user_id, company_type_id, status_ids, updatedat)
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

func (c *CompanyController) GetAllTransporter() {
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	var updatedat *string

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

	status_ids := strings.TrimSpace(c.GetString("status_ids"))
	updated_at := strings.TrimSpace(c.GetString("updated_at"))
	company_type_id := 6

	if updated_at == "" {
		updatedat = nil

	} else {
		updatedat = &updated_at
	}

	d, err := t_company.GetAll(keyword, field_name, match_mode, value_name, currentPage, pageSize, allsize, user_id, company_type_id, status_ids, updatedat)
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

func (c *CompanyController) GetAllGoods() {
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	var updatedat *string

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

	status_ids := strings.TrimSpace(c.GetString("status_ids"))
	updated_at := strings.TrimSpace(c.GetString("updated_at"))
	company_type_id := 7

	if updated_at == "" {
		updatedat = nil

	} else {
		updatedat = &updated_at
	}

	d, err := t_company.GetAll(keyword, field_name, match_mode, value_name, currentPage, pageSize, allsize, user_id, company_type_id, status_ids, updatedat)
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

func (c *CompanyController) GetAllSuppOthers() {
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	var updatedat *string

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

	status_ids := strings.TrimSpace(c.GetString("status_ids"))
	updated_at := strings.TrimSpace(c.GetString("updated_at"))
	company_type_id := 7

	if updated_at == "" {
		updatedat = nil

	} else {
		updatedat = &updated_at
	}

	d, err := t_company.GetAll(keyword, field_name, match_mode, value_name, currentPage, pageSize, allsize, user_id, company_type_id, status_ids, updatedat)
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

func (c *CompanyController) GetAllPartner() {
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	var updatedat *string

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

	status_ids := strings.TrimSpace(c.GetString("status_ids"))
	updated_at := strings.TrimSpace(c.GetString("updated_at"))
	company_type_id := 9

	if updated_at == "" {
		updatedat = nil

	} else {
		updatedat = &updated_at
	}

	d, err := t_company.GetAll(keyword, field_name, match_mode, value_name, currentPage, pageSize, allsize, user_id, company_type_id, status_ids, updatedat)
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

func (c *CompanyController) GetAllInsurance() {
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	var updatedat *string

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

	status_ids := strings.TrimSpace(c.GetString("status_ids"))
	updated_at := strings.TrimSpace(c.GetString("updated_at"))
	company_type_id := 10

	if updated_at == "" {
		updatedat = nil

	} else {
		updatedat = &updated_at
	}

	d, err := t_company.GetAll(keyword, field_name, match_mode, value_name, currentPage, pageSize, allsize, user_id, company_type_id, status_ids, updatedat)
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

func (c *CompanyController) GetAllListInternal() {
	keyword := strings.TrimSpace(c.GetString("keyword"))
	d, err := t_company.GetAllList(keyword, Internal)
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

func (c *CompanyController) GetAllListCustomer() {
	keyword := strings.TrimSpace(c.GetString("keyword"))
	d, err := t_company.GetAllList(keyword, Customer)
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

func (c *CompanyController) GetAllListCustOthers() {
	keyword := strings.TrimSpace(c.GetString("keyword"))
	d, err := t_company.GetAllList(keyword, CustomerOthers)
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

func (c *CompanyController) GetAllListWarehouse() {
	keyword := strings.TrimSpace(c.GetString("keyword"))
	d, err := t_company.GetAllList(keyword, Warehouse)
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

func (c *CompanyController) GetAllListSparepart() {
	keyword := strings.TrimSpace(c.GetString("keyword"))
	d, err := t_company.GetAllList(keyword, Sparepart)
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

func (c *CompanyController) GetAllListTransporter() {
	keyword := strings.TrimSpace(c.GetString("keyword"))
	d, err := t_company.GetAllList(keyword, Transporter)
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

func (c *CompanyController) GetAllListGoods() {
	keyword := strings.TrimSpace(c.GetString("keyword"))
	d, err := t_company.GetAllList(keyword, Goods)
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

func (c *CompanyController) GetAllListSuppOthers() {
	keyword := strings.TrimSpace(c.GetString("keyword"))
	d, err := t_company.GetAllList(keyword, Others)
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

func (c *CompanyController) GetAllListPartner() {
	keyword := strings.TrimSpace(c.GetString("keyword"))
	d, err := t_company.GetAllList(keyword, Partner)
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

func (c *CompanyController) GetAllListInsurance() {
	keyword := strings.TrimSpace(c.GetString("keyword"))
	d, err := t_company.GetAllList(keyword, Insurance)
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
