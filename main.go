package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/di"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	// PostgreSQL
	ctx := context.Background()

	conn, err := pgxpool.New(ctx, "")
	if err != nil {
		log.Fatal("Unable to connect to database", err)
	}
	defer conn.Close()

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
