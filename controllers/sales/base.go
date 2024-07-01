package controllers

import (
	base "mikiwa/controllers"
	"mikiwa/models"
	"sync"
)

type BaseController struct {
	base.BaseController
}

type DeleteBody struct {
	Id string `json:"id"`
}

var errcode int
var errmessage string
var wg *sync.WaitGroup

var form_sales_order = "sales_order"
var form_do = "deliver_order"
var form_do_confirm = "deliver_order_confirm"

var t_sales_order models.SalesOrder
var t_delivery_order models.Do
var t_delivery_order_confirm models.DoConfirm
