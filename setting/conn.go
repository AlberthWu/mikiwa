package setting

import (
	"github.com/beego/beego/v2/core/logs"
)

func init() {
	logs.Info("Initialized connection")
	InitDB()
}
