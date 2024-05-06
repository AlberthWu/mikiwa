package controllers

import (
	"fmt"
	"mikiwa/utils"
	"strconv"
	"strings"
	"time"

	base "mikiwa/controllers"
	"mikiwa/models"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
)

type PoolController struct {
	base.BaseController
}

func (c *PoolController) Prepare() {
	c.Ctx.Request.Header.Set("token", "No Aut")
	c.BaseController.Prepare()
}

func (c *PoolController) Post() {
	name := strings.TrimSpace(c.GetString("name"))
	status, _ := c.GetInt("status")

	valid := validation.Validation{}
	valid.Required(name, "name").Message("Name is required")
	valid.MinSize(name, 3, "name").Message("Name min char is 3")
	valid.MaxSize(name, 250, "name").Message("Name max char is 250")

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

	if models.CheckPoolName(name) {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprintf("'%s' has been registered", name))
		c.ServeJSON()
		return
	}

	pools := models.Pool{
		Name:   name,
		Status: int8(status),
	}

	d, err := models.InsertPool(pools)
	errcode, errmessage := base.DecodeErr(err)
	if err != nil {
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
	} else {
		v, _ := models.GetByIdPool(d.Id)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, errcode, errmessage, v)
	}
	c.ServeJSON()
}

func (c *PoolController) Put() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	name := strings.TrimSpace(c.GetString("name"))
	status, _ := c.GetInt("status")

	valid := validation.Validation{}
	valid.Required(name, "name").Message("Name is required")
	valid.MinSize(name, 3, "name").Message("Name min char is 3")
	valid.MaxSize(name, 250, "name").Message("Name max char is 250")

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

	if models.CheckPoolNamePut(id, name) {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprintf("'%s' has been registered", name))
		c.ServeJSON()
		return
	}

	pools := models.Pool{
		Id:     id,
		Name:   name,
		Status: int8(status),
	}

	err := models.UpdateByIdPool(&pools)
	errcode, errmessage := base.DecodeErr(err)
	if err == nil {
		d, _ := models.GetByIdPool(id)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, errcode, errmessage, d)
	} else {
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
	}
	c.ServeJSON()
}

func (c *PoolController) GetOne() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	v, err := models.GetByIdPool(id)
	code, message := base.DecodeErr(err)
	if err != nil {
		utils.ReturnHTTPError(&c.Controller, code, message)
	} else {

		utils.ReturnHTTPSuccessWithMessage(&c.Controller, code, message, v)
	}
	c.ServeJSON()
}

func (c *PoolController) GetAll() {
	currentPage, _ := c.GetInt("page")
	if currentPage == 0 {
		currentPage = 1
	}

	pageSize, _ := c.GetInt("pagesize")
	if pageSize == 0 {
		pageSize = 10
	}
	keyword := strings.TrimSpace(c.GetString("keyword"))
	d, err := models.GetAllPool(keyword, currentPage, pageSize)
	code, message := base.DecodeErr(err)

	if err != nil {
		utils.ReturnHTTPError(&c.Controller, code, message)
	} else {
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, code, message, d)
	}
	c.ServeJSON()
}

func (c *PoolController) GetAllList() {

	keyword := strings.TrimSpace(c.GetString("keyword"))
	d, err := models.GetAllPoolLimit(keyword)
	code, message := base.DecodeErr(err)

	if err != nil {
		utils.ReturnHTTPError(&c.Controller, code, message)
	} else {
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, code, message, d)
	}
	c.ServeJSON()
}

func (c *PoolController) Delete() {
	o := orm.NewOrm()
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	var pools models.Pool
	pool := new(models.Pool)
	err := o.QueryTable(pool).Filter("id", id).Filter("deleted_at__isnull", true).One(&pools)
	if err == orm.ErrNoRows {
		utils.ReturnHTTPError(&c.Controller, 400, "No data")
		c.ServeJSON()
	}
	pools.DeletedAt = time.Now()
	pools.Status = 0
	o.Update(&pools, "deleted_at", "status")
	utils.ReturnHTTPSuccessWithMessage(&c.Controller, 200, fmt.Sprintf("'%s' has been deleted", pools.Name), nil)
	c.ServeJSON()
}
