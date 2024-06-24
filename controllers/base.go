package controllers

import (
	"fmt"
	"mikiwa/models"
	"mikiwa/utils"
	"strings"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type BaseController struct {
	beego.Controller
}

var cookie_secure, _ = beego.AppConfig.Bool("cookie_secure")
var cookie_http, _ = beego.AppConfig.Bool("cookie_http")

const (
	Read         = 1
	Write        = 2
	Update       = 3
	Delete       = 4
	Approval     = 5
	Author       = 6
	Browser      = 7
	Pending      = 8
	Creator      = 9
	WriteDetail  = 10
	UpdateDetail = 11
	DeleteDetail = 12
)
const (
	Internal       = 1
	Customer       = 2
	CustomerOthers = 3
	Warehouse      = 4
	Sparepart      = 5
	Transporter    = 8
	Goods          = 9
	Others         = 10
	Partner        = 11
	Insurance      = 12
)

const (
	OpenSo     = 1
	ConfirmSo  = 2
	ProgressSo = 3
	DoneSo     = 4
	RejectSo   = 5
	CloseSo    = 6
	VoidSo     = 99
)

func (c *BaseController) Prepare() {
	o := orm.NewOrm()
	sess := c.GetSession("profile")
	logs.Info("Session :", sess)
	tokenString := c.Ctx.Request.Header.Get("token")
	var tokenDb string
	err := o.Raw("select refresh_token from sessions where refresh_token = '" + tokenString + "' and is_blocked = 0 and expired_at >= now()").QueryRow(&tokenDb)
	if err == orm.ErrNoRows {
		if cookie_secure {
			logs.Info("Init check JWT")
			if sess == nil {
				c.DestroySession()
				c.Ctx.SetCookie(SessionName(), "", -1)
				c.Ctx.ResponseWriter.WriteHeader(498)
				utils.ReturnHTTPSuccessWithMessage(&c.Controller, 498, "Forbidden", map[string]interface{}{"message": "Ghost is logged in"})
				c.ServeJSON()
				return
			}

			refresh_token := sess.(map[string]interface{})["refresh_token"].(string)

			sub, err := models.VerifyToken(refresh_token, AccessTokenPublicKey())
			if err != nil {
				c.DestroySession()
				c.Ctx.SetCookie(SessionName(), "", -1)
				c.Ctx.ResponseWriter.WriteHeader(498)
				utils.ReturnHTTPSuccessWithMessage(&c.Controller, 498, "Error", err.Error())
				c.ServeJSON()
				return
			}

			var users models.Users
			errUser := models.Userss().Filter("id", sub).One(&users)
			if errUser == orm.ErrNoRows {
				c.Ctx.ResponseWriter.WriteHeader(400)
				utils.ReturnHTTPError(&c.Controller, 400, "Incorrect email or password")
				c.ServeJSON()
				return
			}

			if users.Token != refresh_token {
				c.DestroySession()
				c.Ctx.SetCookie(SessionName(), "", -1)
				c.Ctx.ResponseWriter.WriteHeader(498)
				utils.ReturnHTTPError(&c.Controller, 498, "You've been kicked out'")
				c.ServeJSON()
				return
			}
			// utils.ReturnHTTPSuccessWithMessage(&c.Controller, 200, "Success", map[string]interface{}{"message": "Session profile", "session": sessSessionId, "id": sessUserIdBase, "user": sessUserNameBase})
			// c.ServeJSON()

		}
	}
}

func AccessTokenPrivateKey() string {
	access_token_private_key, _ := beego.AppConfig.String("jwt::access_token_private_key")
	return access_token_private_key
}

func AccessTokenPublicKey() string {
	access_token_private_key, _ := beego.AppConfig.String("jwt::access_token_public_key")
	return access_token_private_key
}

func AccessTokenDuration() string {
	access_token_duration, _ := beego.AppConfig.String("jwt::access_token_duration")
	return access_token_duration
}

func RefreshTokenPrivateKey() string {
	refresh_token_private_key, _ := beego.AppConfig.String("jwt::refresh_token_private_key")
	return refresh_token_private_key
}

func RefreshTokenPublicKey() string {
	refresh_token_public_key, _ := beego.AppConfig.String("jwt::refresh_token_public_key")
	return refresh_token_public_key
}

func SessionName() string {
	session_name, _ := beego.AppConfig.String("session_name")
	return session_name
}

func FormName(form_name string) int {
	o := orm.NewOrm()
	var id int
	o.Raw("select id from sys_menus where form_name = ?", form_name).QueryRow(&id)
	return id
}

func (c *BaseController) GetSvrDate() {
	c.Data["json"] = utils.GetSvrDate()
	c.ServeJSON()
}

func (c *BaseController) GetDppPpnTotal() {
	issue_date := strings.TrimSpace(c.GetString("issue_date"))
	pph_22, _ := c.GetInt("pph_22")
	pph_23, _ := c.GetInt("pph_23")
	pbb_kb_1, _ := c.GetInt("pbb_kb_1")
	pbb_kb_2, _ := c.GetInt("pbb_kb_2")
	vat, _ := c.GetInt("vat")
	dpp, _ := c.GetFloat("dpp")
	fmt.Printf("Check : %f\n", dpp)
	dpp_amount, pph_22_amount, pph_23_amount, pbbkb1_amount, pbbkb2_amount, ppn, total := utils.GetDppPpnTotal(issue_date, vat, pph_22, pph_23, pbb_kb_1, pbb_kb_2, dpp)
	// dpp_ := fmt.Sprintf("%.2f", math.Trunc(dpp_amount))
	// ppn_ := math.Trunc(ppn)
	// total_ := math.Trunc(total)
	// utils.ReturnHTTPSuccessWithMessage(&c.Controller, 200, "Success", map[string]interface{}{"dpp": dpp_amount, "ppn": ppn, "total": total})
	utils.ReturnHTTPSuccessWithMessage(&c.Controller, 200, "Success", map[string]interface{}{"dpp": fmt.Sprintf("%.2f", dpp_amount), "pph_22": fmt.Sprintf("%.2f", pph_22_amount), "pph_23": fmt.Sprintf("%.2f", pph_23_amount), "pbb_kb_1": fmt.Sprintf("%.2f", pbbkb1_amount), "pbb_kb_2": fmt.Sprintf("%.2f", pbbkb2_amount), "ppn": fmt.Sprintf("%.2f", ppn), "total": fmt.Sprintf("%.2f", total)})
	c.Ctx.ResponseWriter.WriteHeader(200)
	c.ServeJSON()
}
