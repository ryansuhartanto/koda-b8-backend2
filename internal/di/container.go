package di

import (
	"github.com/gin-gonic/gin"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/handler"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/model"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/repository"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/service"
)

type Container struct {
	*repository.UserRepository
	*service.UserService
	*handler.UserHandler
}

func NewContainer() *Container {
	UserRepository := repository.NewUserRepository([]model.User{})
	UserService := service.NewUserService(UserRepository)
	UserHandler := handler.NewUserHandler(UserService)

	return &Container{
		UserRepository,
		UserService,
		UserHandler,
	}
}

func (c *Container) Handle(r *gin.Engine) {
	r.POST("/users/register", c.UserHandler.HandleRegister)
	r.POST("/users/auth", c.UserHandler.HandleAuth)
	r.GET("/users", c.UserHandler.HandleList)
}
