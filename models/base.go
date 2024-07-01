package models

import (
	"fmt"
	"mikiwa/utils"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

func GeneratePettyCashNumber(issue_date time.Time, company_id, account_id int, company_code, transaction_code, transaction_type string) (int, string) {

	o := orm.NewOrm()
	var seqno int
	qb := []string{"select max(voucher_seq_no)+1 seqno from petty_cash_header where voucher_seq_no REGEXP '^[0-9]+$' and   year(issue_date) = " + utils.Int2String(issue_date.Year()) + " and   month(issue_date) = ? and company_id = " + utils.Int2String(company_id) + " and account_id = " + utils.Int2String(account_id) + " and transaction_type ='" + transaction_type + "' and deleted_at is null"}
	sql := strings.Join(qb, "")
	o.Raw(sql, issue_date.Month()).QueryRow(&seqno)
	if seqno == 0 {
		seqno = 1
	}

	strnum := fmt.Sprintf("%04d", int(seqno))

	year := issue_date.Year()
	month := issue_date.Month()
	stryear := utils.Int2String(year % 1e2)
	strmonth := utils.Int2String(int(month))
	strmonth2 := fmt.Sprintf("%02s", strmonth)

	return seqno, transaction_code + stryear + strmonth2 + strnum
}

func GenerateNumber(issue_date time.Time, pool_id, customer_id int) (int, string) {
	o := orm.NewOrm()
	var seqno int
	var format string
	qb := []string{"call sp_GenerateNumber('" + issue_date.Format("2006-01-02") + "','SalesOrder'," + utils.Int2String(pool_id) + "," + utils.Int2String(customer_id) + ")"}

	sql := strings.Join(qb, "")

	o.Raw(sql).QueryRow(&seqno, &format)
	if seqno == 0 {
		seqno = 1
	}
	return seqno, format
}

func GenerateBatchNumber(issue_date time.Time, outlet_id, customer_id int) (int, string) {
	o := orm.NewOrm()
	var seqno int
	var format string
	qb := []string{"call sp_GenerateNumber('" + issue_date.Format("2006-01-02") + "','DeliveryOrder'," + utils.Int2String(outlet_id) + "," + utils.Int2String(customer_id) + ")"}

	sql := strings.Join(qb, "")

	o.Raw(sql).QueryRow(&seqno, &format)
	if seqno == 0 {
		seqno = 1
	}
	return seqno, format
}
