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

type DeleteBody struct {
	Id string `json:"id"`
}

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
