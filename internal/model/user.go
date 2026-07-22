package model

type Email string

type Password string

func (Password) IsZero() bool { return true }

type UserIdentified struct {
	Id
	User
}

type User struct {
	Profile
	Credentials
}

type Profile struct {
	Name       string  `db:"name" form:"name" json:"name" binding:"required"`
	PictureUrl *string `db:"picture_url" form:"picture_url" json:"picture_url,omitzero"`
}

type Credentials struct {
	Email    `db:"email" form:"email" json:"email" binding:"required,email"`
	Password `db:"password" form:"password" json:"password,omitzero" binding:"required,base64"`
}
