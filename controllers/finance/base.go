package controllers

import (
	base "mikiwa/controllers"
	"mikiwa/models"
)

type BaseController struct {
	base.BaseController
}

var form_petty_cash = "petty_cash"

var t_pettycashh models.PettyCashHeader
var t_pettycash models.PettyCash

type DeleteBody struct {
	Id string `json:"id"`
}
