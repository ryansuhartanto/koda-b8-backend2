package model

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthClaims struct {
	jwt.RegisteredClaims
}

func NewAuthClaim(user *UserIdentified) *AuthClaims {
	return &AuthClaims{
		jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(user.Id.Id, 10),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
}

func (c *AuthClaims) GetId() (Id, error) {
	id, err := strconv.ParseInt(c.Subject, 10, 64)
	if err != nil {
		return Id{}, err
	}

	return Id{id}, nil
}
