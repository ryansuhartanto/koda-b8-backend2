package model

type Id struct {
	Id int64 `db:"id" uri:"id" form:"id" json:"id"`
}
