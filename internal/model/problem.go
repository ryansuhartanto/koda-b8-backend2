package model

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Problem struct {
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail,omitempty"`
}

func AbortProblem(ctx *gin.Context, status int, detail string) {
	ctx.Header("Content-Type", "application/problem+json")
	ctx.AbortWithStatusJSON(status, &Problem{
		Title:  http.StatusText(status),
		Status: status,
		Detail: detail,
	})
}
