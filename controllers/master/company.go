package controllers

import (
	base "mikiwa/controllers"
	"mikiwa/models"
	"mikiwa/utils"
	"strings"
)

type CompanyController struct {
	base.BaseController
}

func (c *CompanyController) Prepare() {
	c.Ctx.Request.Header.Set("token", "No Aut")
	c.BaseController.Prepare()
}

func (c *CompanyController) GetAllInternalList() {
	keyword := strings.TrimSpace(c.GetString("keyword"))
	d, err := models.GetAllInternalLimit(keyword)
	code, message := base.DecodeErr(err)

	if err != nil {
		utils.ReturnHTTPError(&c.Controller, code, message)
	} else {
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, code, message, d)
	}
	c.ServeJSON()
}
