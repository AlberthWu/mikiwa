package main

import (
	"mikiwa/middleware"
	_ "mikiwa/routers"
	_ "mikiwa/setting"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/filter/cors"
)

func main() {
	dbType, _ := beego.AppConfig.String("db_type")

	dbSession, _ := beego.AppConfig.Bool("sessionon")
	dbSessionName, _ := beego.AppConfig.String("session_name")
	dbSessionConfigProvider, _ := beego.AppConfig.String(dbType + "::sessionproviderconfig")
	sessionprovider, _ := beego.AppConfig.String(dbType + "::sessionprovider")

	logs.SetLogFuncCall(true)
	logs.Info("AppPath:", beego.AppPath)
	logs.Info("Session:", dbSession)
	logs.Info("Session name:", dbSessionName)

	beego.BConfig.WebConfig.Session.SessionOn = dbSession
	beego.BConfig.WebConfig.Session.SessionName = dbSessionName
	beego.BConfig.WebConfig.Session.SessionAutoSetCookie = false
	beego.BConfig.WebConfig.Session.SessionProvider = sessionprovider
	beego.BConfig.WebConfig.Session.SessionProviderConfig = dbSessionConfigProvider
	EnableCors()

	dbtoken, _ := beego.AppConfig.Bool("jwt::tokenon")
	if dbtoken {
		EnableToken()
	}
	beego.Run()

}

func EnableCors() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "authorize", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "x-xsrf-token", "AxiosHeaders", "X-Requested-With", "X-CSRF-Token", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "authorize", "Set-Cookie", "Cookie"},
		AllowCredentials: true,
	}))
}

func EnableToken() {
	beego.InsertFilter("*", beego.BeforeRouter, func(c *context.Context) {
		middleware.Jwt(c)
	})
}
