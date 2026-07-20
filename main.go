package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, struct {
			Text string `json:"text"`
		}{"Hello, World!"})
	})

	r.Run("0.0.0.0:8080")
}
