package model

type Password string

func (Password) IsZero() bool { return true }

type User struct {
	Name     string   `db:"name" json:"name" binding:"required"`
	Email    string   `db:"email" json:"email" binding:"required"`
	Password Password `db:"password" json:"password,omitzero" binding:"required"`
}
