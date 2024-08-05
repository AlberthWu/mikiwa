package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type (
	Document struct {
		Id          int       `json:"id" orm:"column(id);auto;pk"`
		ReferenceId int       `json:"reference_id" orm:"column(reference_id);"`
		FileName    string    `json:"file_name" orm:"column(file_name);"`
		PathName    string    `json:"path_name" orm:"column(path_name)"`
		PathFile    string    `json:"path_file" orm:"column(path_file)"`
		FileType    string    `json:"file_type" orm:"column(file_type)"`
		FolderName  string    `json:"folder_name" orm:"column(folder_name)"`
		CreatedAt   time.Time `json:"created_at" orm:"column(created_at);type(timestamp);auto_now_add"`
		UpdatedAt   time.Time `json:"updated_at" orm:"column(updated_at);type(timestamp);auto_now"`
		DeletedAt   time.Time `json:"deleted_at" orm:"column(deleted_at);type(timestamp);null"`
		CreatedBy   string    `json:"created_by" orm:"column(created_by)"`
		UpdatedBy   string    `json:"updated_by" orm:"column(updated_by)"`
		DeletedBy   string    `json:"deleted_by" orm:"column(deleted_by)"`
	}
)

func (t *Document) TableName() string {
	return "documents"
}

func Documents() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(Document))
}

func init() {
	orm.RegisterModel(new(Document))
}

func (t *Document) Insert(m Document) (*Document, error) {
	o := orm.NewOrm()
	if _, err := o.Insert(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (t *Document) Update(fields ...string) error {
	o := orm.NewOrm()
	if _, err := o.Update(t, fields...); err != nil {
		return err
	}
	return nil
}

type (
	DocumentRtn struct {
		Id          int    `json:"id"`
		ReferenceId int    `json:"reference_id"`
		FileName    string `json:"file_name"`
		PathName    string `json:"path_name"`
		PathFile    string `json:"path_file"`
		FileType    string `json:"file_type"`
	}
)

func (t *Document) GetDocument(referenceId int, types string) []DocumentRtn {
	var data []Document

	cond1 := orm.NewCondition()
	cond2 := cond1.And("deleted_at__isnull", true).And("reference_id", referenceId).And("file_type__icontains", types)

	qs := Documents().SetCond(cond2)
	num, _ := qs.All(&data)

	var querydata []DocumentRtn

	for _, list := range data {
		querydata = append(querydata, DocumentRtn{
			Id:          list.Id,
			ReferenceId: list.ReferenceId,
			FileName:    list.FileName,
			PathName:    list.PathName,
			PathFile:    list.PathFile,
			FileType:    list.FileType,
		})
	}

	if num == 0 {
		return nil
	}
	return querydata
}
