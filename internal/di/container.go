package di

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/db"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/handler"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/repository"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/service"
)

type Container struct {
	*repository.UserRepository
	*service.UserService
	*handler.UserHandler
}

func NewContainer(querier db.Querier, ctx context.Context) *Container {
	UserRepository := repository.NewUserRepository(querier)
	UserService := service.NewUserService(UserRepository, ctx)
	UserHandler := handler.NewUserHandler(UserService)

	return &Container{
		UserRepository,
		UserService,
		UserHandler,
	}
}

func (c *Container) Handle(r *gin.Engine) {
	r.POST("/users/register", c.UserHandler.HandleRegister)
	r.POST("/users/login", c.UserHandler.HandleAuth)
	r.GET("/users", c.UserHandler.HandleList)
}
