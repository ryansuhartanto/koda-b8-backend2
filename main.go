package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/db"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/di"
)

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

	m.Up()
}

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
