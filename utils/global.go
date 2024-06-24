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

func GetDppPpnTotal(issue_date string, vat, pph_22, pph_23, pbb_kb_1, pbb_kb_2 int, dpp_amount float64) (float64, float64, float64, float64, float64, float64, float64) {
	o := orm.NewOrm()
	var dpp, pph22, pph23, pbbkb1, pbbkb2, ppn, total float64 = 0, 0, 0, 0, 0, 0, 0
	o.Raw("call sp_calcdppppntotal('"+issue_date+"',"+Int2String(vat)+","+Int2String(pph_22)+","+Int2String(pph_23)+","+Int2String(pbb_kb_1)+","+Int2String(pbb_kb_2)+","+Float2String(dpp_amount)+"); ").QueryRow(&dpp, &pph22, &pph23, &pbbkb1, &pbbkb2, &ppn, &total)
	return dpp, pph22, pph23, pbbkb1, pbbkb2, ppn, total
}
