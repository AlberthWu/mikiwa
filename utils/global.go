package utils

import (
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

func GetSvrDate() time.Time {
	o := orm.NewOrm()
	var thetime time.Time
	qb := []string{"select now() thetime"}
	sql := strings.Join(qb, "")

	o.Raw(sql).QueryRow(&thetime)

	return thetime
}

func Int2String(val int) string {
	return strconv.Itoa(val)
}
