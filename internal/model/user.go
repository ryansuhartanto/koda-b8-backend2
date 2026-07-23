package model

import "time"

type Email string

type Password string

func (Password) IsZero() bool { return true }

type Timestamps struct {
	CreatedAt time.Time `db:"created_at" uri:"created_at" form:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" uri:"updated_at" form:"updated_at" json:"updated_at"`
}

type ProfileTimestamps struct {
	ProfileUpdatedAt time.Time `db:"profile_updated_at" uri:"profile_updated_at" form:"profile_updated_at" json:"profile_updated_at"`
}

type UserIdentified struct {
	Id

	Timestamps
	ProfileTimestamps

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
	Email    Email    `db:"email" form:"email" json:"email" binding:"required,email"`
	Password Password `db:"password" form:"password" json:"password,omitzero" binding:"required,base64"`
}
