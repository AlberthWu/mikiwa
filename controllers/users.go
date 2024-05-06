package controllers

import (
	"fmt"
	"mikiwa/utils"
	"strconv"
	"strings"

	"mikiwa/models"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
)

type UsersControllers struct {
	BaseController
}

func (c *UsersControllers) Prepare() {
	c.Ctx.Request.Header.Set("token", "No Aut")
	c.BaseController.Prepare()
}

var t_user models.Users

func (c *UsersControllers) Post() {
	// var data map[string]string
	username := strings.TrimSpace(c.GetString("username"))
	email := strings.TrimSpace(c.GetString("email"))

	password := strings.TrimSpace(c.GetString("password"))
	retype_password := strings.TrimSpace(c.GetString("retype_password"))

	status, _ := c.GetInt("status")

	valid := validation.Validation{}
	valid.Required(username, "username").Message("User name is required")
	valid.Required(password, "password").Message("Password is required")
	valid.Required(retype_password, "retype_password").Message("Retype Password is required")
	valid.Required(email, "email").Message("Email is required")
	valid.Email(email, "email").Message("Invalid email format")

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

	hashedPassword := utils.HashPassword(password)

	if password != retype_password {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Password not match")
		c.ServeJSON()
		return
	}

	if models.ChecKUserName(username) {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprintf("'%s' has been registered", username))
		c.ServeJSON()
		return
	}

	if models.CheckEmail(email) {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprintf("'%s' has been registered", email))
		c.ServeJSON()
		return
	}

	t_user := models.Users{
		Email:    email,
		Password: hashedPassword,
		Username: username,
		Status:   int8(status),
	}

	d, err_ := t_user.Insert(t_user)
	errcode, errmessage := DecodeErr(err_)
	if err_ != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
	} else {
		// v, _ := models.GetByIdFleetFormation(d.Id)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, errcode, errmessage, map[string]string{
			"id":       utils.Int2String(d.Id),
			"email":    d.Email,
			"username": d.Username,
			"status":   utils.Int2String(status),
		})
	}

	c.ServeJSON()
}

func (c *UsersControllers) Put() {
	// var data map[string]string
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	username := strings.TrimSpace(c.GetString("username"))
	email := strings.TrimSpace(c.GetString("email"))

	password := strings.TrimSpace(c.GetString("password"))
	retype_password := strings.TrimSpace(c.GetString("retype_password"))

	status, _ := c.GetInt("status")

	valid := validation.Validation{}
	valid.Required(username, "username").Message("User name is required")
	valid.Required(password, "password").Message("Password is required")
	valid.Required(retype_password, "retype_password").Message("Retype Password is required")
	valid.Required(email, "email").Message("Email is required")
	valid.Email(email, "email").Message("Invalid email format")

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

	hashedPassword := utils.HashPassword(password)

	if password != retype_password {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Password not match")
		c.ServeJSON()
		return
	}

	var users models.Users
	err := models.Userss().Filter("id", id).One(&users)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "User unregistered/Illegal data")
		c.ServeJSON()
		return
	}

	var deletedat = users.DeletedAt.Format("2006-01-02")
	if deletedat != "0001-01-01" {
		c.Ctx.ResponseWriter.WriteHeader(402)
		utils.ReturnHTTPError(&c.Controller, 402, fmt.Sprintf("user_id :'%v' has been deleted", users.Username))
		c.ServeJSON()
		return
	}

	if users.Status == 0 {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprintf("user_id '%s', has been set as INACTIVE", users.Username))
		c.ServeJSON()
		return
	}

	if models.ChecKUserNamePut(id, username) {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprintf("'%s' has been registered", username))
		c.ServeJSON()
		return
	}

	if models.CheckEmailPut(id, email) {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprintf("'%s' has been registered", email))
		c.ServeJSON()
		return
	}

	t_user.Id = id
	t_user.Email = email
	t_user.Password = hashedPassword
	t_user.Username = username
	t_user.Status = int8(status)

	err_ := t_user.Update()
	errcode, errmessage := DecodeErr(err_)
	if err_ != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
	} else {
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, errcode, errmessage, map[string]string{
			"id":       utils.Int2String(id),
			"email":    t_user.Email,
			"username": t_user.Username,
			"status":   utils.Int2String(int(t_user.Status)),
		})
	}
	c.ServeJSON()
}

func (c *UsersControllers) Delete() {}

func (c *UsersControllers) GetOne() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	username := c.GetSession("username")
	password := c.GetSession("password")

	if username != "" {
		c.Ctx.WriteString("Username:" + username.(string) + " Password:" + password.(string))
	}

	c.Data["json"] = id
	c.ServeJSON()
}

func (c *UsersControllers) GetAll() {}

func (c *UsersControllers) GetMenu() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	v, err := models.GetAllMenu(id)
	code, message := DecodeErr(err)
	if err == orm.ErrNoRows {
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
