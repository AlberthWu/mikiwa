package controllers

import (
	base "mikiwa/controllers"
	"mikiwa/models"
)

type BaseController struct {
	base.BaseController
}

var errcode int
var errmessage string

var form_product = "product"
var form_product_division = "product_division"
var form_product_type = "product_type"
var form_uom = "uom"
var form_customer = "customer"
var form_plant = "plant"
var form_price = "sales_price"

var t_product models.Product
var t_product_type models.ProductType
var t_product_division models.ProductDivision
var t_uom models.Uom
var t_company models.Company
var t_business_unit models.BusinessUnit
var t_company_business_unit models.CompanyBusinessUnit
var t_plant models.Plant
var t_bank models.Bank
var t_city models.City
var t_company_type models.CompanyTypes
var t_price models.Price
var t_price_company models.PriceCompany
var t_price_product_uom models.PriceProductUom

type DeleteBody struct {
	Id string `json:"id"`
}
