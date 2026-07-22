package di

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/db"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/handler"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/middleware"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/repository"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/service"
)

type Container struct {
	*repository.UserRepository
	*service.UserService
	*handler.UserHandler

	jwtKey []byte
}

func NewContainer(
	querier db.Querier,
	ctx context.Context,

	jwtKey []byte,
) *Container {
	UserRepository := repository.NewUserRepository(querier)
	UserService := service.NewUserService(UserRepository, ctx, jwtKey)
	UserHandler := handler.NewUserHandler(UserService)

	return &Container{
		UserRepository,
		UserService,
		UserHandler,

		jwtKey,
	}
}

func (c *Container) Handle(r *gin.Engine) {
	r.Use(middleware.CorsMiddleware())

	r.Static("/uploads", "uploads")

	{
		auth := r.Group("/auth")

		auth.POST("/register", c.UserHandler.HandleRegister)
		auth.POST("/login", c.UserHandler.HandleLogin)
	}

	{
		users := r.Group("/users", middleware.AuthMiddleware(c.jwtKey))

		users.GET("/", c.UserHandler.HandleList)
		users.PATCH("/:id", c.UserHandler.HandlePatch)
		users.PUT("/:id/picture", c.UserHandler.HandlePutPicture)
		users.DELETE("/:id", c.UserHandler.HandleDelete)
	}
}
