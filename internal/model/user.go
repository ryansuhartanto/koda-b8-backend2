package model

type Email string

type Password string

func (Password) IsZero() bool { return true }

type Auth struct {
	Email    `db:"email" form:"email" json:"email" binding:"required,email"`
	Password `db:"password" form:"password" json:"password,omitzero" binding:"required"`
}

type User struct {
	Id

	Name string `db:"name" form:"name" json:"name" binding:"required"`
	Auth
}
