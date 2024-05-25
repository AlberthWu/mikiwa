package middleware

import (
	"encoding/json"
	"mikiwa/utils"
	"strings"

	beego "github.com/beego/beego/v2/server/web"

	"github.com/beego/beego/v2/server/web/context"
)

type ResponseRtn struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Jwt(ctx *context.Context) {

	var access_token_public_key, _ = beego.AppConfig.String("jwt::access_token_public_key")

	var access_token string

	if ctx.Input.Header("authorize") == "" {
		ctx.Output.SetStatus(403)
		resBody, err := json.Marshal(ResponseRtn{403, "You are not logged in"})
		ctx.Output.Body(resBody)
		if err != nil {
			panic(err)
		}
		ctx.Redirect(302, "/v1/users/login")
		return
	}

	authorizeHeader := ctx.Request.Header.Get("authorize")
	fields := strings.Fields(authorizeHeader)

	var tokenString string = ctx.Input.Header("authorize")
	authorize := strings.Split(tokenString, "Bearer ")

	if len(fields) != 0 && fields[0] == "Bearer" {
		access_token = fields[1]
	} else {
		access_token = authorize[1]
	}

	if access_token == "" {
		ctx.Output.SetStatus(403)
		resBody, err := json.Marshal(ResponseRtn{403, "You are not logged in"})
		ctx.Output.Body(resBody)
		if err != nil {
			panic(err)
		}
		ctx.Redirect(302, "/v1/users/login")
		return
	}

	_, err := utils.VerifyToken(access_token, access_token_public_key)
	if err != nil {
		ctx.Output.SetStatus(401)
		resBody, err := json.Marshal(ResponseRtn{401, err.Error()})
		ctx.Output.Body(resBody)
		if err != nil {
			panic(err)
		}
		ctx.Redirect(302, "/v1/users/login")
		return
	}

}
