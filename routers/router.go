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
	master "mikiwa/controllers/master"
	sys_manager "mikiwa/controllers/sys_manager"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// utility
	beego.Router("/v1/getsvrdate", &controllers.BaseController{}, "get:GetSvrDate")

	// business_unit
	beego.Router("/v1/business_unit/list", &master.BuController{}, "get:GetAllList")

	// aut
	beego.Router("/v1/users/login", &controllers.AuthController{}, "post:Login")
	beego.Router("/v1/users/logout", &controllers.AuthController{}, "post:Logout")
	beego.Router("/v1/users/forgot", &controllers.AuthController{}, "post:Forgot")
	beego.Router("/v1/users/whoami", &controllers.AuthController{}, "get:GetMe")

	// users
	beego.Router("/v1/users", &controllers.UsersControllers{}, "post:Post;get:GetAll")
	beego.Router("/v1/users/:id", &controllers.UsersControllers{}, "put:Put;get:GetOne")
	beego.Router("/v1/users/menu/:id", &controllers.UsersControllers{}, "get:GetMenu")

	// privileges
	beego.Router("/v1/sys_manager/menu/list/:id", &sys_manager.MenuController{}, "get:GetAll")

	// gudang
	beego.Router("/v1/pool", &master.PoolController{}, "post:Post;get:GetAll")
	beego.Router("/v1/pool/list", &master.PoolController{}, "get:GetAllList")
	beego.Router("/v1/pool/:id", &master.PoolController{}, "put:Put;delete:Delete;get:GetOne")

	// product
	// division
	beego.Router("/v1/product/division", &master.ProductDivisionController{}, "post:Post;get:GetAll")
	beego.Router("/v1/product/division/:id", &master.ProductDivisionController{}, "put:Put;get:GetOne")
	beego.Router("/v1/product/division/list", &master.ProductDivisionController{}, "get:GetAllList")
	// type
	beego.Router("/v1/product/type", &master.ProductTypeController{}, "post:Post;get:GetAll")
	beego.Router("/v1/product/type/:id", &master.ProductTypeController{}, "put:Put;get:GetOne")
	beego.Router("/v1/product/type/list", &master.ProductTypeController{}, "get:GetAllList")
	// product
	beego.Router("/v1/product", &master.ProductController{}, "post:Post;get:GetAll")
	beego.Router("/v1/product/detail", &master.ProductController{}, "get:GetDetail")
	beego.Router("/v1/product/:id", &master.ProductController{}, "put:Put;get:GetOne;delete:Delete")
	beego.Router("/v1/product/list/raw", &master.ProductController{}, "get:GetAllListRaw")
	beego.Router("/v1/product/list/wip", &master.ProductController{}, "get:GetAllListWip")
	beego.Router("/v1/product/list/finishing", &master.ProductController{}, "get:GetAllListFinishing")
	beego.Router("/v1/product/list/acc", &master.ProductController{}, "get:GetAllListAcc")
	beego.Router("/v1/product/list/others", &master.ProductController{}, "get:GetAllListOthers")
	beego.Router("/v1/product/list/recycle", &master.ProductController{}, "get:GetAllListRecycle")
	// image
	beego.Router("/v1/product/document/:id", &master.ProductController{}, "post:PostDocument;get:GetDocument")
	// uom
	beego.Router("/v1/product/uom", &master.UomController{}, "post:Post;get:GetAll")
	beego.Router("/v1/product/uom/:id", &master.UomController{}, "put:Put;get:GetOne")
	beego.Router("/v1/product/uom/list", &master.UomController{}, "get:GetAllList")

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

	// report
	beego.Router("/v1/finance/report/pettycash/summary/1", &finance.FinanceReportController{}, "get:ReportPettyCashSummaryDaily")
	beego.Router("/v1/finance/report/pettycash/summary/2", &finance.FinanceReportController{}, "get:ReportPettyCashSummaryMonthly")
	beego.Router("/v1/finance/report/pettycash/summary/3", &finance.FinanceReportController{}, "get:ReportPettyCashSummaryYearly")
	beego.Router("/v1/finance/report/pettycash/detail/1", &finance.FinanceReportController{}, "get:ReportPettyCashDaily")
	beego.Router("/v1/finance/report/pettycash/detail/2", &finance.FinanceReportController{}, "get:ReportPettyCashMonthly")
	beego.Router("/v1/finance/report/pettycash/detail/3", &finance.FinanceReportController{}, "get:ReportPettyCashYearly")
	beego.Router("/v1/finance/report/pettycash/voucher/:id", &finance.FinanceReportController{}, "get:ReportVoucher")

	// companies
	// list
	beego.Router("/v1/internal/list", &master.CompanyController{}, "get:GetAllListInternal")
	beego.Router("/v1/customer/list", &master.CompanyController{}, "get:GetAllListCustomer")
	beego.Router("/v1/customerothers/list", &master.CompanyController{}, "get:GetAllListCustOthers")
	beego.Router("/v1/warehouse/list", &master.CompanyController{}, "get:GetAllListWarehouse")
	beego.Router("/v1/sparepart/list", &master.CompanyController{}, "get:GetAllListSparepart")
	beego.Router("/v1/transporter/list", &master.CompanyController{}, "get:GetAllListTransporter")
	beego.Router("/v1/goods/list", &master.CompanyController{}, "get:GetAllListGoods")
	beego.Router("/v1/supplierothers/list", &master.CompanyController{}, "get:GetAllListSuppOthers")
	beego.Router("/v1/partner/list", &master.CompanyController{}, "get:GetAllListPartner")
	beego.Router("/v1/insurance/list", &master.CompanyController{}, "get:GetAllListInsurance")
	// getAll
	beego.Router("/v1/internal", &master.CompanyController{}, "get:GetAllInternal")
	beego.Router("/v1/customer", &master.CompanyController{}, "get:GetAllCustomer")
	beego.Router("/v1/customerothers", &master.CompanyController{}, "get:GetAllCustOthers")
	beego.Router("/v1/warehouse", &master.CompanyController{}, "get:GetAllWarehouse")
	beego.Router("/v1/sparepart", &master.CompanyController{}, "get:GetAllSparepart")
	beego.Router("/v1/transporter", &master.CompanyController{}, "get:GetAllTransporter")
	beego.Router("/v1/goods", &master.CompanyController{}, "get:GetAllGoods")
	beego.Router("/v1/supplierothers", &master.CompanyController{}, "get:GetAllSuppOthers")
	beego.Router("/v1/partner", &master.CompanyController{}, "get:GetAllPartner")
	beego.Router("/v1/insurance", &master.CompanyController{}, "get:GetAllInsurance")
	// crud
	beego.Router("/v1/company", &master.CompanyController{}, "post:Post")
	beego.Router("/v1/company/:id", &master.CompanyController{}, "put:Put;get:GetOne;delete:Delete")

	// plants
	beego.Router("/v1/plant", &master.PlantController{}, "post:Post")
	beego.Router("/v1/plant/:id", &master.PlantController{}, "put:Put;get:GetOne;delete:Delete")
	beego.Router("/v1/plant/list/:id", &master.PlantController{}, "get:GetAllList")

	// customer_types
	beego.Router("/v1/companytype/list", &master.CompanyTypeController{}, "get:GetAllList")

	// cities
	beego.Router("/v1/city/list", &master.CityController{}, "get:GetAllList")

	// banks
	beego.Router("/v1/bank/list", &master.BankController{}, "get:GetAllList")
}
