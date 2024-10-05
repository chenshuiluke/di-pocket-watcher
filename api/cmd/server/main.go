package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/chenshuiluke/di-pocket-watcher/api/internal/db"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	url := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	conn, err := pgxpool.New(ctx, url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		queries := db.New(conn)
		results, err := queries.ListUsers(ctx)
		if err != nil {
			log.Println(err)
		}
		return c.JSON(results)
	})

	app.Listen(":8080")
}
