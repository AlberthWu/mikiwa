package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type (
	Users struct {
		Id           int       `json:"id" orm:"column(id);auto;pk"`
		Email        string    `json:"email" orm:"column(email);size(100);null" `
		Password     string    `json:"-" orm:"column(password);size(255)"`
		Username     string    `json:"username" orm:"column(username)"`
		Status       int8      `json:"status" orm:"column(status);default(0)"`
		Token        string    `json:"token" orm:"column(token);type(text);size(255)"`
		RefreshToken string    `json:"refresh_token" orm:"column(refresh_token);type(text);size(255)"`
		CreatedAt    time.Time `json:"created_at" orm:"column(created_at);type(timestamp);auto_now_add"`
		UpdatedAt    time.Time `json:"updated_at" orm:"column(updated_at);type(timestamp);auto_now"`
		DeletedAt    time.Time `json:"deleted_at" orm:"column(deleted_at);type(timestamp);null"`
		// OrderTypes   []*OrderType `json:"-" orm:"reverse(many);rel_through(api.sampurna-group.com/models.UserOrderType)"`
		// SysMenu      []*SysMenu   `json:"-" orm:"reverse(many);rel_through(api.sampurna-group.com/models.UsersSysMenu)"`
		// SysRole      []*SysRole   `json:"-" orm:"reverse(many);rel_through(api.sampurna-group.com/models.SysUserRole)"`
	}

	UserLog struct {
		Id           int       `json:"id" orm:"column(id);auto;pk"`
		SessionId    string    `json:"session_id" orm:"column(session_id);type(text);size(100)"`
		UserId       int       `json:"user_id" orm:"column(user_id);default(0)"`
		Username     string    `json:"username" orm:"column(username)"`
		RefreshToken string    `json:"refresh_token" orm:"column(refresh_token);type(text);size(255)"`
		CreatedAt    time.Time `json:"created_at" orm:"column(created_at);type(datetime)"`
		ExpiredAt    time.Time `json:"expired_at" orm:"column(expired_at);type(datetime)"`
		ClientIp     string    `json:"client_ip" orm:"column(client_ip);type(text);size(45)"`
	}
)

func (t *Users) TableName() string {
	return "users"
}

func Userss() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(Users))
}

func (t *UserLog) TableName() string {
	return "user_log"
}

func UserLogs() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(UserLog))
}

func init() {
	orm.RegisterModel(new(Users), new(UserLog))
}

func (t *Users) Insert(m Users) (*Users, error) {
	o := orm.NewOrm()

	if _, err := o.Insert(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (t *Users) Update(fields ...string) error {
	o := orm.NewOrm()
	if _, err := o.Update(t, fields...); err != nil {
		return err
	}
	return nil
}

func (t *UserLog) Insert(m UserLog) (*UserLog, error) {
	o := orm.NewOrm()

	if _, err := o.Insert(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (t *UserLog) Update(fields ...string) error {
	o := orm.NewOrm()
	if _, err := o.Update(t, fields...); err != nil {
		return err
	}
	return nil
}
