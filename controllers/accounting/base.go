package controllers

import (
	base "mikiwa/controllers"
	"mikiwa/models"
)

type BaseController struct {
	base.BaseController
}

var form_coa = "chart_of_accounts"
var form_account_type = "account_type"

var t_coa models.CharOfAccount
var t_account_type models.GlAccountType
