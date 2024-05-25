package controllers

import (
	"fmt"
	base "mikiwa/controllers"
	"mikiwa/models"
	"mikiwa/utils"
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
	ispo := strings.TrimSpace(c.GetString("is_po"))
	istax := strings.TrimSpace(c.GetString("is_tax"))
	isReceipt, _ := c.GetInt("is_receipt")
	iscash, _ := c.GetInt("is_cash")
	status, _ := c.GetInt("status")
	company_type := strings.TrimSpace(c.GetString("company_type"))
	bank_id, _ := c.GetInt("bank_id")
	bank_name := strings.TrimSpace(c.GetString("bank_name"))
	bank_no := strings.TrimSpace(c.GetString("bank_no"))
	bank_account_name := strings.TrimSpace(c.GetString("bank_account_name"))
	bank_branch := strings.TrimSpace(c.GetString("bank_branch"))

	if len(ispo) == 0 {
		ispo = "1"
	}

	if len(istax) == 0 {
		istax = "1"
	}

	var cities models.City
	models.Cities().Filter("id", cityId).One(&cities)

	var states models.City
	models.Cities().Filter("id", cities.ParentId).One(&states)

	var districts models.City
	models.Cities().Filter("id", states.ParentId).One(&districts)

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
	valid.Required(cityId, "city_id").Message("City is required")
	valid.Required(company_type, "company_type").Message("Company Type is required")

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
		BankName:        bank_name,
		BankNo:          bank_no,
		BankAccountName: bank_account_name,
		BankBranch:      bank_branch,
		IsPo:            int8(utils.String2Int(ispo)),
		IsCash:          int8(iscash),
		IsTax:           int8(utils.String2Int(istax)),
		IsReceipt:       int8(isReceipt),
		PriceMethod:     1,
		Status:          int8(status),
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
}

func (c *CompanyController) Dekete() {
}

func (c *CompanyController) GetOne() {
}

func (c *CompanyController) GetAllInternal() {
}

func (c *CompanyController) GetAllCustomer() {
}

func (c *CompanyController) GetAllCustOthers() {
}

func (c *CompanyController) GetAllWarehouse() {
}

func (c *CompanyController) GetAllSparepart() {
}

func (c *CompanyController) GetAllTransporter() {
}

func (c *CompanyController) GetAllGoods() {
}

func (c *CompanyController) GetAllSuppOthers() {
}

func (c *CompanyController) GetAllInsurance() {
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
