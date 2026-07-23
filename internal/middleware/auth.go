package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ryansuhartanto/koda-b8-backend2/internal/model"
)

var ErrUnexpectedSigningMethod = errors.New("middleware: unexpected signing method")

func AuthMiddleware(jwtKey []byte) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var www strings.Builder
		www.WriteString("Bearer realm=")
		www.WriteRune('"')
		www.WriteString(ctx.FullPath())
		www.WriteRune('"')

		auth := ctx.GetHeader("Authorization")
		tokenString, found := strings.CutPrefix(auth, "Bearer ")
		if !found {
			ctx.Header("WWW-Authenticate", www.String())
			model.AbortProblem(ctx, http.StatusUnauthorized, "missing bearer token")
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
			www.WriteString(`,error="invalid_token"`)
			ctx.Header("WWW-Authenticate", www.String())
			model.AbortProblem(ctx, http.StatusUnauthorized, err.Error())
			return
		}

		userId, err := claims.GetId()
		if err != nil {
			www.WriteString(`,error="invalid_token"`)
			ctx.Header("WWW-Authenticate", www.String())
			model.AbortProblem(ctx, http.StatusUnauthorized, err.Error())
			return
		}

		ctx.Set("user.id", userId)
		ctx.Next()
	}
}
