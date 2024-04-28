package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

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

func String2Int(val string) int {

	goodsId_int, err := strconv.Atoi(val)
	if err != nil {
		return -1
	} else {
		return goodsId_int
	}
}

func Int2String(val int) string {
	return strconv.Itoa(val)
}

func Float2String(val float64) string {
	// return strconv.FormatFloat(val, 'E', 5, 64)
	return fmt.Sprintf("%f", val)
}

func Pointer2String(val *string) string {
	if val != nil {
		return *val
	}

	return ""
}

func ToUpper(s string) string {
	a := []rune(s)
	for i, c := range a {
		if unicode.IsLower(c) {
			a[i] = unicode.ToUpper(c)
		}
	}
	return string(a)
}
