// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"mikiwa/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/v1/users/login", &controllers.AuthController{}, "post:Login")
	beego.Router("/v1/users/logout", &controllers.AuthController{}, "post:Logout")
	beego.Router("/v1/users/forgot", &controllers.AuthController{}, "post:Forgot")
	beego.Router("/v1/users/whoami", &controllers.AuthController{}, "get:GetMe")
}
