package controllers

import (
	base "mikiwa/controllers"
	"mikiwa/models"
)

type BaseController struct {
	base.BaseController
}

type DeleteBody struct {
	Id string `json:"id"`
}

var errcode int
var errmessage string

var form_sales_order = "sales_order"

var t_sales_order models.SalesOrder
