package controllers

import (
	"mikiwa/utils"
	"strconv"

	base "mikiwa/controllers"

	"github.com/beego/beego/v2/client/orm"
)

type MenuController struct {
	base.BaseController
}

func (c *MenuController) Prepare() {
	c.Ctx.Request.Header.Set("token", "No Aut")
	c.BaseController.Prepare()
}

func (c *MenuController) GetAll() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	v, err := t_sys_menu.GetAllMenu(id)
	code, message := base.DecodeErr(err)
	if err == orm.ErrNoRows {
		code := 200
		c.Ctx.ResponseWriter.WriteHeader(code)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, code, "No data", nil)
	} else if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(code)
		utils.ReturnHTTPError(&c.Controller, code, message)
	} else {
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, code, message, v)
	}
	c.ServeJSON()
}
