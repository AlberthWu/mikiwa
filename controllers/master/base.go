package controllers

import (
	base "mikiwa/controllers"
	"mikiwa/models"
)

type BaseController struct {
	base.BaseController
}

var form_product = "product"
var form_product_division = "product_division"
var form_product_type = "product_type"
var form_uom = "uom"

var t_product models.Product
var t_product_type models.ProductType
var t_product_division models.ProductDivision
var t_uom models.Uom

type DeleteBody struct {
	Id string `json:"id"`
}
