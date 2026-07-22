package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/PeterTakahashi/gin-openapi/openapiui"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/db"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/di"
)

// @title           GoREST
// @version         1.0
// @description     Go exercise implementing a simple REST server.

// @contact.name   Ryan Suhartanto
// @contact.url    https://github.com/ryansuhartanto/koda-b8-backend2
// @contact.email  suhartanto@kekkon.nexus

// @license.name  MIT

// @servers.url  http://localhost:8080

// @securitydefinitions.bearerauth BearerAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
}

var ctx context.Context
var pool *pgxpool.Pool

func init() {
	var err error
	ctx = context.Background()
	pool, err = pgxpool.New(ctx, "")
	if err != nil {
		panic(err)
	}

	m, err := db.Migrate(pool)
	if err != nil {
		panic(err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}
}

func main() {
	r := gin.Default()
	container := di.NewContainer(pool, ctx, []byte(os.Getenv("JWT_SECRET")))

	r.Any("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/docs")
	})

	r.GET("/docs/*any", openapiui.WrapHandler(openapiui.Config{
		SpecURL:      "/docs/openapi.json",
		SpecFilePath: "./docs/swagger.json",
		Title:        "GoREST",
	}))

	container.Handle(r)

	r.Run("0.0.0.0:8080")
}
