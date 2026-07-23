package model

import (
	"github.com/golang-jwt/jwt/v5"
)

type AuthResult struct {
	*UserIdentified
	Token string `json:"jwt" format:"jwt"`
}

func NewAuthResult(
	user *UserIdentified,
	key []byte,
) (*AuthResult, error) {
	claims := NewAuthClaim(user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	tokenString, err := token.SignedString(key)

	if err != nil {
		return nil, err
	}

	return &AuthResult{
		user,
		tokenString,
	}, nil
}
