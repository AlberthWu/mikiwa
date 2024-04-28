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
	accounting "mikiwa/controllers/accounting"
	finance "mikiwa/controllers/finance"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// users
	beego.Router("/v1/users/login", &controllers.AuthController{}, "post:Login")
	beego.Router("/v1/users/logout", &controllers.AuthController{}, "post:Logout")
	beego.Router("/v1/users/forgot", &controllers.AuthController{}, "post:Forgot")
	beego.Router("/v1/users/whoami", &controllers.AuthController{}, "get:GetMe")

	// accounting
	beego.Router("/v1/accounting/account_type/list", &accounting.AccountTypeController{}, "get:GetAllList")
	beego.Router("/v1/accounting/account_type/list/assets", &accounting.AccountTypeController{}, "get:GetAllListAssets")
	beego.Router("/v1/accounting/account_type/list/expenses", &accounting.AccountTypeController{}, "get:GetAllListExpenses")
	beego.Router("/v1/accounting/account_type/list/liability", &accounting.AccountTypeController{}, "get:GetAllListLiability")
	beego.Router("/v1/accounting/account_type/list/equity", &accounting.AccountTypeController{}, "get:GetAllListEquity")
	beego.Router("/v1/accounting/account_type/list/cogs", &accounting.AccountTypeController{}, "get:GetAllListCogs")
	beego.Router("/v1/accounting/account_type", &accounting.AccountTypeController{}, "post:Post;get:GetAll")
	beego.Router("/v1/accounting/account_type/:id", &accounting.AccountTypeController{}, "put:Put;get:GetOne;delete:Delete")

	beego.Router("/v1/accounting/coa/list", &accounting.CoaController{}, "get:GetAllLimit")
	beego.Router("/v1/accounting/coa/list/:id", &accounting.CoaController{}, "get:GetAllLimiChildByCompany")
	beego.Router("/v1/accounting/coa/list/assets/:id", &accounting.CoaController{}, "get:GetAllLimiChildByCompanyAssets")
	beego.Router("/v1/accounting/coa", &accounting.CoaController{}, "post:Post;get:GetAll")
	beego.Router("/v1/accounting/coa/:id", &accounting.CoaController{}, "put:Put;get:GetOne;delete:Delete")

	// finance
	// petty cash
	beego.Router("/v1/finance/pettycash", &finance.PettyCashHController{}, "post:Post;get:GetAll")
	beego.Router("/v1/finance/pettycash/:id", &finance.PettyCashHController{}, "put:Put;get:GetOne;delete:Delete")
	beego.Router("/v1/finance/pettycash/list/:id", &finance.PettyCashHController{}, "get:GetAllList")
	beego.Router("/v1/finance/pettycash/reorder", &finance.PettyCashHController{}, "get:ReOrderNumList;post:ReOrderNum")

	beego.Router("/v1/finance/pettycash/detail/checkdelete/:id", &finance.PettyCashController{}, "get:CheckDelete")
	beego.Router("/v1/finance/pettycash/detail", &finance.PettyCashController{}, "post:Post")
	beego.Router("/v1/finance/pettycash/detail/:id", &finance.PettyCashController{}, "put:Put;delete:Delete")

}
