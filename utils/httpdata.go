package utils

import (
	"encoding/json"

	beego "github.com/beego/beego/v2/server/web"
)

type HTTPData struct {
	ErrNo  int         `json:"errno"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

func ReturnHTTPSuccessWithMessage(this *beego.Controller, errno int, errmsg string, val interface{}) {

	rtndata := HTTPData{
		ErrNo:  errno,
		ErrMsg: errmsg,
		Data:   val,
	}

	data, err := json.Marshal(rtndata)
	if err != nil {
		this.Data["json"] = err
	} else {
		this.Data["json"] = json.RawMessage(string(data))
	}
}

func ReturnHTTPError(this *beego.Controller, errno int, errmsg string) {

	rtndata := HTTPData{
		ErrNo:  errno,
		ErrMsg: errmsg,
		Data:   nil,
	}

	data, err := json.Marshal(rtndata)
	if err != nil {
		this.Data["json"] = err
	} else {
		this.Data["json"] = json.RawMessage(string(data))
	}
}
