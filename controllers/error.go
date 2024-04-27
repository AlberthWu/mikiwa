package controllers

import (
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
)

type ErrorController struct {
	beego.Controller
}

type ErrResponse struct {
	ErrCode    int    `json:"err_code"`
	ErrMessage string `json:"err_message"`
}

type ErrNo struct {
	ErrCode    int
	ErrMessage string
}

type Err struct {
	ErrCode    int
	ErrMessage string
	Err        error
}

var (
	// Common errors
	SUCCESS             = &ErrNo{ErrCode: 0, ErrMessage: "Success"}
	InternalServerError = &ErrNo{ErrCode: 401, ErrMessage: "Internal service error"}
	ErrBind             = &ErrNo{ErrCode: 402, ErrMessage: "Parameter error"}

	ErrDB           = &ErrNo{ErrCode: 20001, ErrMessage: "Database error"}
	ErrToken        = &ErrNo{ErrCode: 20002, ErrMessage: "Token issue error"}
	ErrNoPermission = &ErrNo{ErrCode: 20003, ErrMessage: "No permission"}

	// user errors
	ErrUserNotFound = &ErrNo{ErrCode: 20101, ErrMessage: "User not registered"}
	ErrUserExist    = &ErrNo{ErrCode: 20102, ErrMessage: "User already exists"}
)

func (err ErrNo) Error() string {
	return err.ErrMessage
}

func (err *Err) Error() string {
	return fmt.Sprintf("Err - code: %d, message: %s, error: %s", err.ErrCode, err.ErrMessage, err.Err)
}

func DecodeErr(err error) (int, string) {
	if err == nil {
		return SUCCESS.ErrCode, SUCCESS.ErrMessage
	}

	switch typed := err.(type) {
	case *Err:
		return typed.ErrCode, typed.ErrMessage
	case *ErrNo:
		return typed.ErrCode, typed.ErrMessage
	default:
	}

	return InternalServerError.ErrCode, err.Error()
}
