package models

type (
	Employee struct {
		Id int `json:"id"  orm:"column(id);auto;pk"`
	}
)
