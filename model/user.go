package model

import "encoding/json"

type User struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		Password string `json:"password,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(&u),
	})
}
