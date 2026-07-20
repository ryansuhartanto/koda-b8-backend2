package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryansuhartanto/koda-b8-backend1/model"
	"github.com/ryansuhartanto/koda-b8-backend1/repo"
)

func main() {
	r := gin.Default()

	repoUser := repo.NewRepoUser([]model.User{})

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, struct {
			Text string `json:"text"`
		}{"Hello, World!"})
	})

	r.GET("/users", func(ctx *gin.Context) {
		users, err := repoUser.FindAll()

		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusOK, users)
	})

	r.POST("/users", func(ctx *gin.Context) {
		var new model.User
		if err := ctx.ShouldBindJSON(&new); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		if err := repoUser.Create(new); err != nil {
			ctx.JSON(http.StatusConflict, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, new)
	})

	r.Run("0.0.0.0:8080")
}
