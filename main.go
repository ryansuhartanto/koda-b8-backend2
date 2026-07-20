package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/di"
)

func main() {
	r := gin.Default()
	container := di.NewContainer()

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, struct {
			Text string `json:"text"`
		}{"Hello, World!"})
	})

	container.Handle(r)

	r.Run("0.0.0.0:8080")
}
