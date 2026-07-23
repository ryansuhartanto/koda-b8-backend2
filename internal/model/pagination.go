package model

type Pagination struct {
	Limit  int    `db:"limit"  form:"limit"  json:"limit"  binding:"omitempty,min=1,max=100"`
	Offset int    `db:"offset" form:"offset" json:"offset" binding:"omitempty,min=0"`
	Query  string `db:"query"  form:"query"  json:"query"`
}
