package controllers

import (
	"encoding/json"
	"fmt"
	base "mikiwa/controllers"
	"mikiwa/models"
	"mikiwa/utils"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
)

type PriceController struct {
	base.BaseController
}

func (c *PriceController) Prepare() {
	c.Ctx.Request.Header.Set("token", "No Aut")
	c.BaseController.Prepare()
}

type (
	InputHeaderPrice struct {
		EffectiveDate string             `json:"effective_date"`
		ExpiredDate   string             `json:"expired_date"`
		CompanyId     int                `json:"company_id"`
		ProductId     int                `json:"product_id"`
		SureName      string             `json:"sure_name"`
		StatusId      int8               `json:"status_id"`
		Price         []InputDetailPrice `json:"price"`
	}

	InputDetailPrice struct {
		Id        int     `json:"id"`
		UomId     int     `json:"uom_id"`
		Ratio     float64 `json:"ratio"`
		DiscOne   float64 `json:"disc_one"`
		DiscTwo   float64 `json:"disc_two"`
		DiscTpr   float64 `json:"disc_tpr"`
		IsDefault int8    `json:"is_default"`
	}
)

func (c *PriceController) Post() {
	o := orm.NewOrm()
	var user_id, form_id int
	var user_name string
	var err error

	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
		user_name = sess.(map[string]interface{})["username"].(string)
	}

	form_id = base.FormName(form_price)

	write_aut := models.CheckPrivileges(user_id, form_id, base.Write)
	write_aut = true
	if !write_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Post not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}

	var i int = 0
	var ob InputHeaderPrice
	var inputDetail []models.PriceProductUom

	price_type := "sales"
	var uom_id int
	var uom_code string
	var disc_one, disc_two, disc_tpr float64
	var normal_price, price, ratio float64
	body := c.Ctx.Input.RequestBody
	err = json.Unmarshal(body, &ob)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}
	valid := validation.Validation{}
	valid.Required(strings.TrimSpace(ob.EffectiveDate), "effective_date").Message("Is required")
	valid.Required(ob.ProductId, "product_id").Message("Is required")
	valid.Required(ob.CompanyId, "company_id").Message("Is required")

	if len(ob.Price) == 0 {
		valid.AddError("price", "0 not allowed")
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

	var company models.Company
	err = models.Companies().Filter("id", ob.CompanyId).Filter("CompanyTypes__TypeId__Id", base.Customer).Filter("deleted_at__isnull", true).One(&company)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Customer unregistered/Illegal data")
		c.ServeJSON()
		return
	}
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	if company.Status == 0 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been set as INACTIVE", company.Code))
		c.ServeJSON()
		return
	}

	var products models.Product
	err = models.Products().Filter("id", ob.ProductId).Filter("deleted_at__isnull", true).One(&products)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Product unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	if products.StatusId == 0 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been set as INACTIVE", products.ProductCode))
		c.ServeJSON()
		return
	}

	var uom models.Uom
	for _, v := range ob.Price {
		err = models.Uoms().Filter("id", v.UomId).One(&uom)
		if err == orm.ErrNoRows {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, "Uom unregistered/Illegal data")
			c.ServeJSON()
			return
		}

		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, err.Error())
			c.ServeJSON()
			return
		}

		if uom.StatusId == 0 {
			c.Ctx.ResponseWriter.WriteHeader(402)
			utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been set as INACTIVE", uom.UomCode))
			c.ServeJSON()
			return
		}

		if v.Ratio == 0 {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 402, "0 ratio not allowed")
			c.ServeJSON()
			return
		}

		if v.IsDefault == 1 {
			uom_id = v.UomId
			uom_code = uom.UomCode
			disc_one = v.DiscOne
			disc_two = v.DiscTwo
			disc_tpr = v.DiscTpr
			ratio = v.Ratio
			i += 1
		}
	}

	if i > 1 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Multiple default uom", map[string]interface{}{"is_default": "'Only allowed 1 uom as a default"})
		c.ServeJSON()
		return
	}

	thedate, errDate := time.Parse("2006-01-02", ob.EffectiveDate)
	if errDate != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprint("effective_date: ", errDate.Error()))
		c.ServeJSON()
		return
	}

	var expireddate *time.Time
	if ob.ExpiredDate == "" {
		expireddate = nil
	} else {
		expiredthedate, err_date := time.Parse("2006-01-02", ob.ExpiredDate)
		if err_date != nil {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprint("expired_date: ", err_date.Error()))
			c.ServeJSON()
			return
		}
		expireddate = &expiredthedate
	}

	t_price = models.Price{
		EffectiveDate: thedate,
		ExpiredDate:   expireddate,
		CompanyId:     ob.CompanyId,
		CompanyCode:   company.Code,
		ProductId:     ob.ProductId,
		ProductCode:   products.ProductCode,
		UomId:         uom_id,
		UomCode:       uom_code,
		DiscOne:       disc_one * -1,
		DiscTwo:       disc_two * -1,
		DiscTpr:       disc_tpr * -1,
		Ratio:         ratio,
		NormalPrice:   normal_price,
		Price:         price,
		SureName:      ob.SureName,
		PriceType:     price_type,
		StatusId:      ob.StatusId,
		CreatedBy:     user_name,
		UpdatedBy:     user_name,
	}

	d, err_ := t_price.Insert(t_price)
	errcode, errmessage = base.DecodeErr(err_)

	if err_ != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
	} else {
		i = 0
		for k, v := range ob.Price {
			inputDetail = append(inputDetail, models.PriceProductUom{
				ItemNo:    k + 1,
				PriceId:   d.Id,
				ProductId: ob.ProductId,
				UomId:     v.UomId,
				Ratio:     v.Ratio,
				DiscOne:   v.DiscOne,
				DiscTwo:   v.DiscTwo,
				DiscTpr:   v.DiscTpr,
				IsDefault: v.IsDefault,
				CreatedBy: user_name,
				UpdatedBy: user_name,
			})
			i += 1
		}
		o.InsertMulti(i, inputDetail)
		o.Raw("call sp_CalcPriceProductUom(" + utils.Int2String(ob.ProductId) + "," + utils.Int2String(d.Id) + "," + utils.Int2String(user_id) + ")").Exec()
		// conversion := t_product.GetConversion(1, ob.ProductId, ob.ProductId, user_id)
		// disc_one := (ob.Price * ob.DiscOne) / 100
		// disc_two := (ob.Price - disc_one) * ob.DiscTwo / 100
		// price := ob.Price - disc_one - disc_two - ob.DiscTpr

		_, err_ = t_price_company.InsertM2M(d.Id, utils.Int2String(ob.CompanyId))
		errcode, errmessage = base.DecodeErr(err_)
		if err_ != nil {
			utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
			c.ServeJSON()
			return
		}
		// o.Raw("call sp_GeneratePriceProduct()").Exec()
		// v, err := t_price.GetById(d.Id)
		// errcode, errmessage := base.DecodeErr(err)
		// if err != nil {
		// 	c.Ctx.ResponseWriter.WriteHeader(errcode)
		// 	utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
		// } else {
		// 	utils.ReturnHTTPSuccessWithMessage(&c.Controller, errcode, errmessage, v)
		// }
	}
	c.ServeJSON()
}
func (c *PriceController) Put() {
	o := orm.NewOrm()
	var user_id, form_id int
	var user_name string
	var err error

	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
		user_name = sess.(map[string]interface{})["username"].(string)
	}

	form_id = base.FormName(form_price)

	put_aut := models.CheckPrivileges(user_id, form_id, base.Write)
	put_aut = true
	if !put_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Put not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}

	var i int = 0
	var ob InputHeaderPrice
	var inputDetail []models.PriceProductUom

	price_type := "sales"
	var uom_id int
	var uom_code string
	var disc_one, disc_two, disc_tpr float64
	var normal_price, price, ratio float64
	body := c.Ctx.Input.RequestBody
	err = json.Unmarshal(body, &ob)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	var querydata models.Price
	err = models.Prices().Filter("id", id).Filter("deleted_at__isnull", true).Filter("price_type", price_type).One(&querydata)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, "Price id unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	valid := validation.Validation{}
	valid.Required(strings.TrimSpace(ob.EffectiveDate), "effective_date").Message("Is required")
	valid.Required(ob.ProductId, "product_id").Message("Is required")

	if len(ob.Price) == 0 {
		valid.AddError("price", "0 not allowed")
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

	var company models.Company
	err = models.Companies().Filter("id", ob.CompanyId).Filter("CompanyTypes__TypeId__Id", base.Customer).Filter("deleted_at__isnull", true).One(&company)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Customer unregistered/Illegal data")
		c.ServeJSON()
		return
	}
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	if company.Status == 0 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been set as INACTIVE", company.Code))
		c.ServeJSON()
		return
	}

	var products models.Product
	err = models.Products().Filter("id", ob.ProductId).Filter("deleted_at__isnull", true).One(&products)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Product unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	if products.StatusId == 0 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been set as INACTIVE", products.ProductCode))
		c.ServeJSON()
		return
	}

	var deleteIds []string
	var joinId string

	var uom models.Uom
	for _, v := range ob.Price {
		err = models.Uoms().Filter("id", v.UomId).One(&uom)
		if err == orm.ErrNoRows {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, "Uom unregistered/Illegal data")
			c.ServeJSON()
			return
		}

		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, err.Error())
			c.ServeJSON()
			return
		}

		if uom.StatusId == 0 {
			c.Ctx.ResponseWriter.WriteHeader(402)
			utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been set as INACTIVE", uom.UomCode))
			c.ServeJSON()
			return
		}

		if v.Ratio == 0 {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 402, "0 ratio not allowed")
			c.ServeJSON()
			return
		}

		if v.IsDefault == 1 {
			uom_id = v.UomId
			uom_code = uom.UomCode
			disc_one = v.DiscOne
			disc_two = v.DiscTwo
			disc_tpr = v.DiscTpr
			ratio = v.Ratio
			i += 1
		}

		if v.Id != 0 {
			deleteIds = append(deleteIds, utils.Int2String(v.Id))
		}
	}

	if i > 1 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Multiple default uom", map[string]interface{}{"is_default": "'Only allowed 1 uom as a default"})
		c.ServeJSON()
		return
	}

	thedate, errDate := time.Parse("2006-01-02", ob.EffectiveDate)
	if errDate != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprint("effective_date: ", errDate.Error()))
		c.ServeJSON()
		return
	}

	var expireddate *time.Time
	if ob.ExpiredDate == "" {
		expireddate = nil
	} else {
		expiredthedate, err_date := time.Parse("2006-01-02", ob.ExpiredDate)
		if err_date != nil {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprint("expired_date: ", err_date.Error()))
			c.ServeJSON()
			return
		}
		expireddate = &expiredthedate
	}

	if len(deleteIds) == 0 {
		joinId = "0"
	} else {
		joinId = strings.Join(deleteIds, ",")
	}
	o.Raw("update price_product_uom set deleted_at = now(), deleted_by = '" + user_name + "' where deleted_at is null and price_id = " + utils.Int2String(id) + " and id not in (" + joinId + ")").Exec()

	t_price.Id = id
	t_price.EffectiveDate = thedate
	t_price.ExpiredDate = expireddate
	t_price.CompanyId = ob.CompanyId
	t_price.CompanyCode = company.Code
	t_price.ProductId = ob.ProductId
	t_price.ProductCode = products.ProductCode
	t_price.UomId = uom_id
	t_price.UomCode = uom_code
	t_price.DiscOne = disc_one * -1
	t_price.DiscTwo = disc_two * -1
	t_price.DiscTpr = disc_tpr * -1
	t_price.Ratio = ratio
	t_price.NormalPrice = normal_price
	t_price.Price = price
	t_price.SureName = ob.SureName
	t_price.PriceType = price_type
	t_price.StatusId = ob.StatusId
	t_price.CreatedBy = querydata.CreatedBy
	t_price.UpdatedBy = user_name
	err_ := t_price.Update()

	errcode, errmessage = base.DecodeErr(err_)
	if err_ != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
	} else {
		i = 0
		for k, v := range ob.Price {
			if v.Id == 0 {
				inputDetail = append(inputDetail, models.PriceProductUom{
					ItemNo:    k + 1,
					PriceId:   id,
					ProductId: ob.ProductId,
					UomId:     v.UomId,
					Ratio:     v.Ratio,
					DiscOne:   v.DiscOne,
					DiscTwo:   v.DiscTwo,
					DiscTpr:   v.DiscTpr,
					IsDefault: v.IsDefault,
					CreatedBy: user_name,
					UpdatedBy: user_name,
				})
				i += 1
			} else {
				t_price_product_uom.Id = v.Id
				t_price_product_uom.ItemNo = k + 1
				t_price_product_uom.PriceId = id
				t_price_product_uom.ProductId = ob.ProductId
				t_price_product_uom.UomId = v.UomId
				t_price_product_uom.Ratio = v.Ratio
				t_price_product_uom.DiscOne = v.DiscOne
				t_price_product_uom.DiscTwo = v.DiscTwo
				t_price_product_uom.DiscTpr = v.DiscTpr
				t_price_product_uom.IsDefault = v.IsDefault
				t_price_product_uom.CreatedBy = querydata.CreatedBy
				t_price_product_uom.UpdatedBy = user_name
				t_price_product_uom.Update()
			}
		}
		o.InsertMulti(i, inputDetail)
		o.Raw("call sp_CalcPriceProductUom(" + utils.Int2String(ob.ProductId) + "," + utils.Int2String(id) + "," + utils.Int2String(user_id) + ")").Exec()
		// conversion := t_product.GetConversion(1, ob.ProductId, ob.ProductId, user_id)
		// disc_one := (ob.Price * ob.DiscOne) / 100
		// disc_two := (ob.Price - disc_one) * ob.DiscTwo / 100
		// price := ob.Price - disc_one - disc_two - ob.DiscTpr

		_, err_ = t_price_company.InsertM2M(id, utils.Int2String(ob.CompanyId))
		errcode, errmessage = base.DecodeErr(err_)
		if err_ != nil {
			utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
			c.ServeJSON()
			return
		}
		// o.Raw("call sp_GeneratePriceProduct()").Exec()
		// v, err := t_price.GetById(d.Id)
		// errcode, errmessage := base.DecodeErr(err)
		// if err != nil {
		// 	c.Ctx.ResponseWriter.WriteHeader(errcode)
		// 	utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
		// } else {
		// 	utils.ReturnHTTPSuccessWithMessage(&c.Controller, errcode, errmessage, v)
		// }
	}
	c.ServeJSON()
}
func (c *PriceController) Delete() {
	var user_id, form_id int
	var err error
	var user_name string
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
		user_name = sess.(map[string]interface{})["username"].(string)
	}
	form_id = base.FormName(form_price)
	delete_aut := models.CheckPrivileges(user_id, form_id, base.Delete)
	if !delete_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Delete not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	price_type := "sales"

	var querydata models.Price
	err = models.Prices().Filter("id", id).Filter("deleted_at__isnull", true).Filter("price_type", price_type).One(&querydata)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, "Price id unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	models.Prices().Filter("id", id).Filter("deleted_at__isnull", true).Update(orm.Params{"deleted_at": utils.GetSvrDate(), "deleted_by": user_name})
	models.PriceProductUoms().Filter("price_id", id).Filter("deleted_at__isnull", true).Update(orm.Params{"deleted_at": utils.GetSvrDate(), "deleted_by": user_name})

	utils.ReturnHTTPError(&c.Controller, 200, "soft delete success")
	c.ServeJSON()
}
func (c *PriceController) GetOne() {
	// var user_id int
	// sess := c.GetSession("profile")
	// if sess != nil {
	// 	user_id = sess.(map[string]interface{})["id"].(int)
	// }
	// id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	// v, err := t_price.GetById(id, user_id)
	// code, message := base.DecodeErr(err)
	// if err == orm.ErrNoRows {
	// 	code = 200
	// 	c.Ctx.ResponseWriter.WriteHeader(code)
	// 	utils.ReturnHTTPError(&c.Controller, code, "No data")
	// } else if err != nil {
	// 	c.Ctx.ResponseWriter.WriteHeader(code)
	// 	utils.ReturnHTTPError(&c.Controller, code, message)
	// } else {

	// 	utils.ReturnHTTPSuccessWithMessage(&c.Controller, code, message, v)
	// }
	// c.ServeJSON()
}
func (c *PriceController) GetAll() {}
func (c *PriceController) CalcPrice() {
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	issue_date := strings.TrimSpace(c.GetString("issue_date"))
	company_id, _ := c.GetInt("company_id")
	product_id, _ := c.GetInt("product_id")
	uom_id, _ := c.GetInt("uom_id")
	disc_one, _ := c.GetFloat("disc_one")
	disc_two, _ := c.GetFloat("disc_two")
	disc_tpr, _ := c.GetFloat("disc_tpr")

	if issue_date == "" {
		issue_date = utils.GetSvrDate().Format("2006-01-02")
	}
	price := t_product.GetPrice(issue_date, company_id, product_id, uom_id, user_id, 0, 1, disc_one, disc_two, disc_tpr)
	utils.ReturnHTTPSuccessWithMessage(&c.Controller, 200, "Success", map[string]interface{}{"price": fmt.Sprintf("%.2f", price)})
	c.Ctx.ResponseWriter.WriteHeader(200)
	c.ServeJSON()
}
