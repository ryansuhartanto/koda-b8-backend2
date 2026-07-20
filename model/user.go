package model

type Password string

func (Password) IsZero() bool { return true }

type User struct {
	Name     string   `json:"name" binding:"required"`
	Email    string   `json:"email" binding:"required"`
	Password Password `json:"password,omitzero" binding:"required"`
}
