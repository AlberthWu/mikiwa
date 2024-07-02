package controllers

import (
	"encoding/json"
	"fmt"
	base "mikiwa/controllers"
	"mikiwa/models"
	"mikiwa/utils"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
)

type ProductController struct {
	base.BaseController
}

func (c *ProductController) Prepare() {
	c.Ctx.Request.Header.Set("token", "No Aut")
	c.BaseController.Prepare()
}

type (
	InputHeaderProduct struct {
		ProductCode      string              `json:"product_code"`
		ProductName      string              `json:"product_name"`
		ProductTypeId    int                 `json:"product_type_id"`
		ProductDivisonId int                 `json:"product_division_id"`
		SerialNumber     string              `json:"serial_number"`
		ArtNo            string              `json:"art_no"`
		Barcode          string              `json:"barcode"`
		LeadTime         int                 `json:"lead_time"`
		StatusId         int8                `json:"status_id"`
		UploadFile       models.DocumentList `json:"upload_file"`
		Uom              []InputBodyUom      `json:"uom"`
	}

	InputBodyUom struct {
		UomId     int     `json:"uom_id"`
		Ratio     float64 `json:"ratio"`
		IsDefault int8    `json:"is_default"`
		Price     float64 `json:"price"`
	}
)

func (c *ProductController) Post() {
	o := orm.NewOrm()
	var user_id, form_id int
	var user_name string
	var err error
	var folderName string = "product"
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
		user_name = sess.(map[string]interface{})["username"].(string)
	}

	form_id = base.FormName(form_product)
	write_aut := models.CheckPrivileges(user_id, form_id, base.Write)
	write_aut = true
	if !write_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Post not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}

	var i int = 0
	var ob InputHeaderProduct
	var inputDetail []models.ProductUom

	body := c.Ctx.Input.RequestBody
	err = json.Unmarshal(body, &ob)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}
	valid := validation.Validation{}
	valid.Required(strings.TrimSpace(ob.ProductCode), "product_code").Message("Is required")
	valid.Required(strings.TrimSpace(ob.ProductName), "product_name").Message("Is required")
	valid.Required(ob.ProductTypeId, "product_type_id").Message("Is required")
	valid.Required(ob.ProductDivisonId, "product_division_id").Message("Is required")

	if len(ob.Uom) == 0 {
		valid.AddError("uom", "Uom list is required")
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

	if t_product.CheckCode(0, ob.ProductCode) {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("product_code : '%v' has been REGISTERED", ob.ProductCode))
		c.ServeJSON()
		return
	}

	var division models.ProductDivision
	err = models.ProductDivisions().Filter("id", ob.ProductDivisonId).One(&division)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Product division unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	if division.StatusId == 0 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been set as INACTIVE", division.DivisionName))
		c.ServeJSON()
		return
	}

	var types models.ProductType
	err = models.ProductTypes().Filter("id", ob.ProductTypeId).One(&types)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Product type unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	var uom_id int
	var uom_code string
	var price float64

	var uom models.Uom
	for _, v := range ob.Uom {
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
			utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Invalid data", map[string]interface{}{"uom_id": "'" + uom.UomCode + "' has been set as INACTIVE"})
			c.ServeJSON()
			return
		}

		if v.IsDefault == 1 {
			uom_id = v.UomId
			uom_code = uom.UomCode
			price = v.Price
			i += 1
		}
	}

	if i > 1 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Multiple default uom", map[string]interface{}{"is_default": "'Only allowed 1 uom as a default"})
		c.ServeJSON()
		return
	}

	var errcode int
	var errmessage string
	var ratio float64 = 1

	t_product = models.Product{
		ProductCode:         ob.ProductCode,
		ProductName:         ob.ProductName,
		ProductTypeId:       ob.ProductTypeId,
		ProductTypeName:     types.TypeName,
		ProductDivisionId:   ob.ProductDivisonId,
		ProductDivisionCode: division.DivisionCode,
		ProductDivisionName: division.DivisionName,
		SerialNumber:        ob.SerialNumber,
		Barcode:             ob.Barcode,
		ArtNo:               ob.ArtNo,
		Price:               price,
		UomId:               uom_id,
		UomCode:             uom_code,
		LeadTime:            ob.LeadTime,
		StatusId:            ob.StatusId,
		CreatedBy:           user_name,
	}

	d, err_ := t_product.Insert(t_product)
	errcode, errmessage = base.DecodeErr(err_)
	if err_ != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
	} else {
		i = 0
		for k, v := range ob.Uom {
			if v.Ratio == 0 {
				ratio = 1
			} else {
				ratio = v.Ratio
			}
			inputDetail = append(inputDetail, models.ProductUom{
				ProductId: d.Id,
				ItemNo:    k + 1,
				UomId:     v.UomId,
				Ratio:     ratio,
				IsDefault: v.IsDefault,
				Price:     v.Price,
			})
			i += 1
		}

		o.InsertMulti(i, inputDetail)
		if err := base.PostFirebaseRaw(ob.UploadFile, user_name, d.Id, folderName+"/"+utils.Int2String(d.Id), folderName+"/"+utils.Int2String(d.Id)); err != nil {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprint("Error processing data and uploading to Firebase: ", err.Error()))
			c.ServeJSON()
			return
		}
		o.Raw("call sp_CalcProductUom(" + utils.Int2String(d.Id) + "," + utils.Int2String(user_id) + ")").Exec()
		v, err := t_product.GetById(d.Id, user_id)
		errcode, errmessage = base.DecodeErr(err)
		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(errcode)
			utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
		} else {
			utils.ReturnHTTPSuccessWithMessage(&c.Controller, errcode, errmessage, v)
		}
	}
	c.ServeJSON()
}

func (c *ProductController) Put() {
	o := orm.NewOrm()
	var user_id, form_id int
	var user_name string
	var err error
	var folderName string = "product"
	var deletedat string
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
		user_name = sess.(map[string]interface{})["username"].(string)
	}

	form_id = base.FormName(form_product)
	put_aut := models.CheckPrivileges(user_id, form_id, base.Write)
	put_aut = true
	if !put_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Put not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}

	var i int = 0
	var ob InputHeaderProduct
	var inputDetail []models.ProductUom

	body := c.Ctx.Input.RequestBody
	err = json.Unmarshal(body, &ob)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	var querydata models.Product
	err = models.Products().Filter("id", id).One(&querydata)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, "Product id unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	deletedat = querydata.DeletedAt.Format("2006-01-02")
	if deletedat != "0001-01-01" {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been DELETED", querydata.ProductCode))
		c.ServeJSON()
		return
	}

	valid := validation.Validation{}
	valid.Required(strings.TrimSpace(ob.ProductCode), "product_code").Message("Is required")
	valid.Required(strings.TrimSpace(ob.ProductName), "product_name").Message("Is required")
	valid.Required(ob.ProductTypeId, "product_type_id").Message("Is required")
	valid.Required(ob.ProductDivisonId, "product_division_id").Message("Is required")

	if len(ob.Uom) == 0 {
		valid.AddError("uom", "Uom list is required")
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

	if t_product.CheckCode(id, ob.ProductCode) {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("product_code : '%v' has been REGISTERED", ob.ProductCode))
		c.ServeJSON()
		return
	}

	var division models.ProductDivision
	err = models.ProductDivisions().Filter("id", ob.ProductDivisonId).One(&division)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Product division unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	if division.StatusId == 0 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been set as INACTIVE", division.DivisionName))
		c.ServeJSON()
		return
	}

	var types models.ProductType
	err = models.ProductTypes().Filter("id", ob.ProductTypeId).One(&types)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Product type unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	var uom_id int
	var uom_code string
	var price float64

	var uom models.Uom
	for _, v := range ob.Uom {
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
			utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Invalid data", map[string]interface{}{"uom_id": "'" + uom.UomCode + "' has been set as INACTIVE"})
			c.ServeJSON()
			return
		}

		if v.IsDefault == 1 {
			uom_id = v.UomId
			uom_code = uom.UomCode
			price = v.Price
			i += 1
		}
	}

	if i > 1 {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Multiple default uom", map[string]interface{}{"is_default": "'Only allowed 1 uom as a default"})
		c.ServeJSON()
		return
	}

	var errcode int
	var errmessage string
	var ratio float64

	t_product.Id = id
	t_product.ProductCode = ob.ProductCode
	t_product.ProductName = ob.ProductName
	t_product.ProductTypeId = ob.ProductTypeId
	t_product.ProductTypeName = types.TypeName
	t_product.ProductDivisionId = ob.ProductDivisonId
	t_product.ProductDivisionCode = division.DivisionCode
	t_product.ProductDivisionName = division.DivisionName
	t_product.SerialNumber = ob.SerialNumber
	t_product.ArtNo = ob.ArtNo
	t_product.Barcode = ob.Barcode
	t_product.UomId = uom_id
	t_product.UomCode = uom_code
	t_product.Price = price
	t_product.LeadTime = ob.LeadTime
	t_product.StatusId = ob.StatusId
	t_product.CreatedBy = querydata.CreatedBy
	t_product.UpdatedBy = user_name
	err_ := t_product.Update()

	errcode, errmessage = base.DecodeErr(err_)
	if err_ != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
	} else {
		o.Raw("delete from product_uom where product_id = " + utils.Int2String(id) + " ").Exec()
		i = 0
		for k, v := range ob.Uom {
			if v.Ratio == 0 {
				ratio = 1
			} else {
				ratio = v.Ratio
			}
			inputDetail = append(inputDetail, models.ProductUom{
				ProductId: id,
				ItemNo:    k + 1,
				UomId:     v.UomId,
				Ratio:     ratio,
				IsDefault: v.IsDefault,
				Price:     v.Price,
			})
			i += 1
		}

		o.InsertMulti(i, inputDetail)
		if err := base.PutFirebaseRaw(ob.UploadFile, user_name, id, folderName+"/"+utils.Int2String(id), folderName+"/"+utils.Int2String(id)); err != nil {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprint("Error processing data and uploading to Firebase: ", err.Error()))
			c.ServeJSON()
			return
		}
		o.Raw("call sp_CalcProductUom(" + utils.Int2String(id) + "," + utils.Int2String(user_id) + ")").Exec()
		v, err := t_product.GetById(id, user_id)
		errcode, errmessage = base.DecodeErr(err)
		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(errcode)
			utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
		} else {
			utils.ReturnHTTPSuccessWithMessage(&c.Controller, errcode, errmessage, v)
		}
	}
	c.ServeJSON()
}

func (c *ProductController) Delete() {
	var user_id, form_id int
	var err error
	var deletedat, user_name string
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
		user_name = sess.(map[string]interface{})["username"].(string)
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

	var querydata models.Product
	err = models.Products().Filter("id", id).One(&querydata)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, "Product id unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, err.Error())
		c.ServeJSON()
		return
	}

	deletedat = querydata.DeletedAt.Format("2006-01-02")
	if deletedat != "0001-01-01" {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("Error '%v' has been DELETED", querydata.ProductCode))
		c.ServeJSON()
		return
	}

	models.Products().Filter("id", id).Filter("deleted_at__isnull", true).Update(orm.Params{"deleted_at": utils.GetSvrDate(), "deleted_by": user_name})

	utils.ReturnHTTPError(&c.Controller, 200, "soft delete success")
	c.ServeJSON()
}

func (c *ProductController) GetOne() {
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	v, err := t_product.GetById(id, user_id)
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

func (c *ProductController) GetAll() {
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
	purchase_ids := strings.TrimSpace(c.GetString("purchase_ids"))
	sales_ids := strings.TrimSpace(c.GetString("sales_ids"))
	production_ids := strings.TrimSpace(c.GetString("production_ids"))
	division_ids := strings.TrimSpace(c.GetString("division_ids"))
	type_ids := strings.TrimSpace(c.GetString("type_ids"))

	if updated_at == "" {
		updatedat = nil

	} else {
		updatedat = &updated_at
	}

	d, err := t_product.GetAll(keyword, field_name, match_mode, value_name, currentPage, pageSize, allsize, user_id, division_ids, type_ids, production_ids, purchase_ids, sales_ids, status_ids, updatedat)
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

func (c *ProductController) GetDetail() {
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	var updatedat *string

	keyword := strings.TrimSpace(c.GetString("keyword"))
	match_mode := strings.TrimSpace(c.GetString("match_mode"))
	value_name := strings.TrimSpace(c.GetString("value_name"))
	field_name := strings.TrimSpace(c.GetString("field_name"))

	status_ids := strings.TrimSpace(c.GetString("status_ids"))
	updated_at := strings.TrimSpace(c.GetString("updated_at"))
	purchase_ids := strings.TrimSpace(c.GetString("purchase_ids"))
	sales_ids := strings.TrimSpace(c.GetString("sales_ids"))
	production_ids := strings.TrimSpace(c.GetString("production_ids"))
	division_ids := strings.TrimSpace(c.GetString("division_ids"))
	type_ids := strings.TrimSpace(c.GetString("type_ids"))

	if updated_at == "" {
		updatedat = nil

	} else {
		updatedat = &updated_at
	}

	d, err := t_product.GetAllDetail(keyword, field_name, match_mode, value_name, user_id, division_ids, type_ids, production_ids, purchase_ids, sales_ids, status_ids, updatedat)
	code, message := base.DecodeErr(err)
	if err == orm.ErrNoRows {
		code = 200
		c.Ctx.ResponseWriter.WriteHeader(code)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, code, "No Data", map[string]interface{}{
			"field_key":          nil,
			"field_label":        nil,
			"field_int":          nil,
			"field_level":        nil,
			"field_export":       nil,
			"field_export_label": nil,
			"field_footer":       nil,
			"list":               nil})
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

func (c *ProductController) PostDocument() {
	var user_id, form_id int
	var user_name string
	var folderName string = "product"
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
		user_name = sess.(map[string]interface{})["username"].(string)
	}
	write_aut := models.CheckPrivileges(user_id, form_id, base.Write)
	write_aut = true
	if !write_aut {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 402, "Post not authorize", map[string]interface{}{"message": "Please contact administrator"})
		c.ServeJSON()
		return
	}

	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	upload_file, err := c.GetFiles("upload_file")
	code, message := base.DecodeErr(err)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(code)
		utils.ReturnHTTPError(&c.Controller, code, message)
	} else {
		if err := base.PutFilesFirebase(upload_file, user_name, id, folderName+"/"+utils.Int2String(id), folderName+"/"+utils.Int2String(id)); err != nil {
			c.Ctx.ResponseWriter.WriteHeader(401)
			utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprintf("Error while posting files to Firebase: %s", err.Error()))
			c.ServeJSON()
			return
		}
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 200, "Sucess", "File uploaded")
	}
	c.ServeJSON()
}

func (c *ProductController) GetDocument() {
	var user_id int
	var folderName string = "product"
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}

	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	v := t_product.Document(id, user_id, folderName)

	utils.ReturnHTTPSuccessWithMessage(&c.Controller, 200, "Sucess", v)

	c.ServeJSON()
}

func (c *ProductController) GetAllListRaw() {
	keyword := strings.TrimSpace(c.GetString("keyword"))

	d, err := t_product.GetAllListRaw(keyword)
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

func (c *ProductController) GetAllListWip() {
	keyword := strings.TrimSpace(c.GetString("keyword"))

	d, err := t_product.GetAllListWip(keyword)
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

func (c *ProductController) GetAllListFinishing() {
	keyword := strings.TrimSpace(c.GetString("keyword"))
	product_division_id, _ := c.GetInt("product_division_id")

	d, err := t_product.GetAllListFinishing(keyword, product_division_id)
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

func (c *ProductController) GetAllListAcc() {
	keyword := strings.TrimSpace(c.GetString("keyword"))

	d, err := t_product.GetAllListAcc(keyword)
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

func (c *ProductController) GetAllListOthers() {
	keyword := strings.TrimSpace(c.GetString("keyword"))

	d, err := t_product.GetAllListOthers(keyword)
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

func (c *ProductController) GetAllListRecycle() {
	keyword := strings.TrimSpace(c.GetString("keyword"))

	d, err := t_product.GetAllListRecycle(keyword)
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

func (c *ProductController) GetAllListSales() {
	keyword := strings.TrimSpace(c.GetString("keyword"))

	d, err := t_product.GetAllListSales(keyword)
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

func (c *ProductController) GetAllListPurchase() {
	keyword := strings.TrimSpace(c.GetString("keyword"))

	d, err := t_product.GetAllListPurchase(keyword)
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

func (c *ProductController) GetAllListProduction() {
	keyword := strings.TrimSpace(c.GetString("keyword"))

	d, err := t_product.GetAllListProduction(keyword)
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

func (c *ProductController) GetProductUom() {
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	uom_id, _ := c.GetInt("uom_id")
	d := t_product.GetProductUom(id, uom_id, user_id)
	utils.ReturnHTTPSuccessWithMessage(&c.Controller, 200, "Success", d)
	c.ServeJSON()
}

func (c *ProductController) GetConversion() {
	var user_id int
	sess := c.GetSession("profile")
	if sess != nil {
		user_id = sess.(map[string]interface{})["id"].(int)
	}
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	uom_id, _ := c.GetInt("uom_id")
	customer_id, _ := c.GetInt("customer_id")
	qty, _ := c.GetFloat("qty")
	issue_date := strings.TrimSpace(c.GetString("issue_date"))
	if issue_date == "" {
		issue_date = utils.GetSvrDate().Format("2006-01-02")
	}
	d := t_product.GetConversion(issue_date, qty, customer_id, id, uom_id, user_id)
	utils.ReturnHTTPSuccessWithMessage(&c.Controller, 200, "Success", d)
	c.ServeJSON()
}
