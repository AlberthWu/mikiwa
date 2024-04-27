package controllers

import (
	"fmt"
	"mikiwa/utils"
	"net/http"
	"strings"
	"time"

	"mikiwa/models"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
)

type AuthController struct {
	beego.Controller
}

func (c *AuthController) Login() {
	var t_user_log models.UserLog
	sess := c.GetSession("profile")
	sid := c.CruSession.SessionID(c.Ctx.Request.Context())
	cookie := http.Cookie{Name: SessionName(), Value: sid, Path: "/", SameSite: http.SameSiteNoneMode, Secure: cookie_secure, HttpOnly: cookie_http, Expires: time.Now().Add(24 * time.Hour)}
	http.SetCookie(c.Ctx.ResponseWriter, &cookie)
	logs.Info("do login", " sessionId :", sess)

	email := strings.TrimSpace(c.GetString("email"))
	password := strings.TrimSpace(c.GetString("password"))
	valid := validation.Validation{}
	valid.Required(email, "email").Message("Email is required")
	valid.Required(password, "password").Message("Password is required")

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

	access_token_duration := AccessTokenDuration()
	access_token_private_key := AccessTokenPrivateKey()

	td, errTd := time.ParseDuration(access_token_duration)
	errcodeTd, errmessageTd := DecodeErr(errTd)
	if errTd != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcodeTd)
		utils.ReturnHTTPError(&c.Controller, errcodeTd, errmessageTd)
		c.ServeJSON()
		return
	}

	logs.Info("Token duration :", td, access_token_duration)

	var users models.Users
	err := models.Userss().Filter("email", email).One(&users)
	if err == orm.ErrNoRows {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Incorrect email or password")
		c.ServeJSON()
		return
	}

	var deletedat = users.DeletedAt.Format("2006-01-02")
	if deletedat != "0001-01-01" {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprintf("user_id '%s', has been been DELETED", users.Username))
		c.ServeJSON()
		return
	}

	if users.Status == 0 {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, fmt.Sprintf("user_id '%s', has been set as INACTIVE", users.Username))
		c.ServeJSON()
		return
	}

	if errPassword := utils.ComparePassword(users.Password, password); errPassword != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, "Incorrect email or password")
		c.ServeJSON()
		return
	}

	access_token, exp, err := models.CreateToken(td, users.Id, access_token_private_key)
	errcode, errmessage := DecodeErr(err)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(errcode)
		utils.ReturnHTTPError(&c.Controller, errcode, errmessage)
		c.ServeJSON()
		return
	}

	profile := make(map[string]interface{})
	profile["id"] = users.Id
	profile["username"] = users.Username
	profile["email"] = users.Email
	profile["refresh_token"] = access_token
	profile["session_id"] = sid
	profile["logged_in"] = true
	profile["expired_at"] = exp

	logs.Info("Check exp:", exp, utils.GetSvrDate(), "..")

	thedate, errdate := time.Parse("2006-01-02 15:04:05 -0700", exp+" +0700")
	if errdate != nil {
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPError(&c.Controller, 401, errdate.Error())
		c.ServeJSON()
		return
	}

	sess = c.SetSession("profile", profile)
	sess_profile := c.GetSession("profile")
	sessionID := sess_profile.(map[string]interface{})["session_id"].(string)

	s := strings.Split(c.Ctx.Request.RemoteAddr, ":")

	models.Userss().Filter("id", users.Id).Update(orm.Params{"token": access_token, "refresh_token": sessionID, "updated_at": utils.GetSvrDate()})
	t_user_log = models.UserLog{
		SessionId:    sessionID,
		UserId:       users.Id,
		Username:     users.Username,
		RefreshToken: access_token,
		CreatedAt:    utils.GetSvrDate(),
		ExpiredAt:    thedate,
		ClientIp:     s[0],
	}

	t_user_log.Insert(t_user_log)
	c.Ctx.ResponseWriter.WriteHeader(200)
	utils.ReturnHTTPSuccessWithMessage(&c.Controller, 200, "Success", map[string]interface{}{"message": "Login Success", "id": users.Id, "user": users.Username, "email": users.Email, "token": sessionID})
	c.ServeJSON()
}

func (c *AuthController) Logout() {
	logs.Info("POST logout")
	sess := c.GetSession("profile")
	if sess == nil {
		c.DestroySession()
		c.Ctx.SetCookie(SessionName(), "", -1)
		logs.Info("Ghost is logged in")
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 401, "Forbidden", map[string]interface{}{"message": "Ghost is logged in"})
		c.ServeJSON()
		return
	}

	c.DestroySession()
	c.Ctx.SetCookie(SessionName(), "", -1)

	user_id := sess.(map[string]interface{})["id"].(int)

	var users models.Users
	models.Userss().Filter("id", user_id).One(&users)
	models.Userss().Filter("id", users.Id).Update(orm.Params{"token": nil, "refresh_token": nil})
	logs.Info("Logout success")
	c.Ctx.ResponseWriter.WriteHeader(200)
	utils.ReturnHTTPSuccessWithMessage(&c.Controller, 200, "Success", map[string]interface{}{"message": "Logout success"})
	c.ServeJSON()
}

func (c *AuthController) Forgot() {

}

func (c *AuthController) GetMe() {

	sess := c.GetSession("profile")
	if sess != nil {
		sessUserId := sess.(map[string]interface{})["id"].(int)
		sessUserName := sess.(map[string]interface{})["username"].(string)
		c.Ctx.ResponseWriter.WriteHeader(200)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 200, "Success", map[string]interface{}{"message": "Who am i", "id": sessUserId, "user": sessUserName})

	} else {
		c.DestroySession()
		c.Ctx.SetCookie(SessionName(), "", -1)
		logs.Info("Ghost is logged in")
		c.Ctx.ResponseWriter.WriteHeader(401)
		utils.ReturnHTTPSuccessWithMessage(&c.Controller, 401, "Success", map[string]interface{}{"message": "Ghost is logged in"})
	}

	c.ServeJSON()
}
