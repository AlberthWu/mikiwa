package controllers

import (
	base "mikiwa/controllers"
	"mikiwa/utils"
	"strings"

	"github.com/beego/beego/v2/client/orm"
)

type CompanyTypeController struct {
	base.BaseController
}

func (c *CompanyTypeController) Prepare() {
	c.Ctx.Request.Header.Set("token", "No Aut")
	c.BaseController.Prepare()
}

func (c *CompanyTypeController) GetAllList() {
	keyword := strings.TrimSpace(c.GetString("keyword"))

	d, err := t_company_type.GetAllList(keyword)
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
