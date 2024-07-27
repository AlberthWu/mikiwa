package controllers

import (
	base "mikiwa/controllers"
	"mikiwa/models"
	"sync"
)

type BaseController struct {
	base.BaseController
}

var errcode int
var errmessage string
var wg *sync.WaitGroup

var form_petty_cash = "petty_cash"

var t_pettycashh models.PettyCashHeader

type DeleteBody struct {
	Id string `json:"id"`
}
