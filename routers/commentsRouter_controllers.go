package routers

// import (
// 	beego "github.com/beego/beego/v2/server/web"
// 	"github.com/beego/beego/v2/server/web/context/param"
// )

// func init() {
// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers/nideshop:GoodsController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers/nideshop:GoodsConGoodsControllertroller"],
// 		beego.ControllerComments{
// 			Method:           "Goods_Count",
// 			Router:           "/goods/count",
// 			AllowHTTPMethods: []string{"get"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})
// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers/nideshop:GoodsController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers/nideshop:GoodsController"],
// 		beego.ControllerComments{
// 			Method:           "Goods_List",
// 			Router:           "/goods/list",
// 			AllowHTTPMethods: []string{"get"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})
// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers/nideshop:GoodsController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers/nideshop:GoodsController"],
// 		beego.ControllerComments{
// 			Method:           "Goods_Category",
// 			Router:           "/goods/category",
// 			AllowHTTPMethods: []string{"get"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})
// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers/nideshop:GoodsController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers/nideshop:GoodsController"],
// 		beego.ControllerComments{
// 			Method:           "Goods_Detail",
// 			Router:           "/goods/detail",
// 			AllowHTTPMethods: []string{"get"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})
// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers:CompanyTypeController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers:CompanyTypeController"],
// 		beego.ControllerComments{
// 			Method:           "GetAll",
// 			Router:           "/",
// 			AllowHTTPMethods: []string{"get"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})

// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers:BankController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers:BankController"],
// 		beego.ControllerComments{
// 			Method:           "Post",
// 			Router:           "/",
// 			AllowHTTPMethods: []string{"post"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})
// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers:BankController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers:BankController"],
// 		beego.ControllerComments{
// 			Method:           "Put",
// 			Router:           "/:id",
// 			AllowHTTPMethods: []string{"put"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})

// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers:BankController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers:BankController"],
// 		beego.ControllerComments{
// 			Method:           "GetOne",
// 			Router:           "/:id",
// 			AllowHTTPMethods: []string{"get"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})

// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers:BankController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers:BankController"],
// 		beego.ControllerComments{
// 			Method:           "GetAll",
// 			Router:           "/",
// 			AllowHTTPMethods: []string{"get"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})
// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers:BankController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers:BankController"],
// 		beego.ControllerComments{
// 			Method:           "GetAllPage",
// 			Router:           "/All",
// 			AllowHTTPMethods: []string{"get"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})
// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers:CityController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers:CityController"],
// 		beego.ControllerComments{
// 			Method:           "GetAllDistrict",
// 			Router:           "/district",
// 			AllowHTTPMethods: []string{"get"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})
// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers:CityController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers:CityController"],
// 		beego.ControllerComments{
// 			Method:           "GetAllCity",
// 			Router:           "/city",
// 			AllowHTTPMethods: []string{"get"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})
// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers:CityController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers:CityController"],
// 		beego.ControllerComments{
// 			Method:           "GetAllState",
// 			Router:           "/state",
// 			AllowHTTPMethods: []string{"get"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})
// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers:CityController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers:CityController"],
// 		beego.ControllerComments{
// 			Method:           "GetAll",
// 			Router:           "/",
// 			AllowHTTPMethods: []string{"get"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})

// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers:CityController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers:CityController"],
// 		beego.ControllerComments{
// 			Method:           "Post",
// 			Router:           "/",
// 			AllowHTTPMethods: []string{"post"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})

// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers:CityController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers:CityController"],
// 		beego.ControllerComments{
// 			Method:           "GetOne",
// 			Router:           "/:id",
// 			AllowHTTPMethods: []string{"get"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})

// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers:CustomerController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers:CustomerController"],
// 		beego.ControllerComments{
// 			Method:           "Post",
// 			Router:           "/",
// 			AllowHTTPMethods: []string{"post"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})

// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers:CustomerController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers:CustomerController"],
// 		beego.ControllerComments{
// 			Method:           "GetOne",
// 			Router:           "/:id",
// 			AllowHTTPMethods: []string{"get"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})
// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers:CustomerController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers:CustomerController"],
// 		beego.ControllerComments{
// 			Method:           "GetAll",
// 			Router:           "/",
// 			AllowHTTPMethods: []string{"get"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})

// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers/sample:ObjectController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers/sample:ObjectController"],
// 		beego.ControllerComments{
// 			Method:           "Post",
// 			Router:           "/",
// 			AllowHTTPMethods: []string{"post"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})

// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers/sample:ObjectController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers/sample:ObjectController"],
// 		beego.ControllerComments{
// 			Method:           "GetAll",
// 			Router:           "/",
// 			AllowHTTPMethods: []string{"get"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})

// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers/sample:ObjectController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers/sample:ObjectController"],
// 		beego.ControllerComments{
// 			Method:           "Get",
// 			Router:           "/:objectId",
// 			AllowHTTPMethods: []string{"get"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})

// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers/sample:ObjectController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers/sample:ObjectController"],
// 		beego.ControllerComments{
// 			Method:           "Put",
// 			Router:           "/:objectId",
// 			AllowHTTPMethods: []string{"put"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})

// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers/sample:ObjectController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers/sample:ObjectController"],
// 		beego.ControllerComments{
// 			Method:           "Delete",
// 			Router:           "/:objectId",
// 			AllowHTTPMethods: []string{"delete"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})

// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers/sample:UserController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers/sample:UserController"],
// 		beego.ControllerComments{
// 			Method:           "Post",
// 			Router:           "/",
// 			AllowHTTPMethods: []string{"post"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})

// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers/sample:UserController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers/sample:UserController"],
// 		beego.ControllerComments{
// 			Method:           "GetAll",
// 			Router:           "/",
// 			AllowHTTPMethods: []string{"get"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})

// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers/sample:UserController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers/sample:UserController"],
// 		beego.ControllerComments{
// 			Method:           "Get",
// 			Router:           "/:uid",
// 			AllowHTTPMethods: []string{"get"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})

// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers/sample:UserController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers:UserController"],
// 		beego.ControllerComments{
// 			Method:           "Put",
// 			Router:           "/:uid",
// 			AllowHTTPMethods: []string{"put"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})

// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers/sample:UserController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers/sample:UserController"],
// 		beego.ControllerComments{
// 			Method:           "Delete",
// 			Router:           "/:uid",
// 			AllowHTTPMethods: []string{"delete"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})

// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers/sample:UserController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers/sample:UserController"],
// 		beego.ControllerComments{
// 			Method:           "Login",
// 			Router:           "/login",
// 			AllowHTTPMethods: []string{"get"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})

// 	beego.GlobalControllerRouter["api.sampurna-group.com/controllers/sample:UserController"] = append(beego.GlobalControllerRouter["api.sampurna-group.com/controllers/sample:UserController"],
// 		beego.ControllerComments{
// 			Method:           "Logout",
// 			Router:           "/logout",
// 			AllowHTTPMethods: []string{"get"},
// 			MethodParams:     param.Make(),
// 			Filters:          nil,
// 			Params:           nil})

// }
