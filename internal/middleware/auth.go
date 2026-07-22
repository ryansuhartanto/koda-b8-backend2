package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/model"
)

var ErrUnexpectedSigningMethod = errors.New("middleware: unexpected signing method")

func AuthMiddleware(jwtKey []byte) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.GetHeader("Authorization")
		tokenString, found := strings.CutPrefix(auth, "Bearer ")
		if !found {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var claims model.AuthClaims
		_, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrUnexpectedSigningMethod
			}
			return jwtKey, nil
		})
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}

		userId, err := claims.GetId()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}

		ctx.Set("user.id", userId)
		ctx.Next()
	}
}
