package model

type id int64

func (id) IsZero() bool { return true }

type Id struct {
	Id id `db:"id" form:"id" json:"id,omitzero"`
}
