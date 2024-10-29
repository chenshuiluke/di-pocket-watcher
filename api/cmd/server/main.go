package main

import (
	"log"

	connection_manager "github.com/chenshuiluke/di-pocket-watcher/api/internal"
	"github.com/chenshuiluke/di-pocket-watcher/api/internal/controllers"
	"github.com/chenshuiluke/di-pocket-watcher/api/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	app := fiber.New()
	app.Use(cors.New())
	app.Get("/", func(c *fiber.Ctx) error {

		results, err := connection_manager.Mgr.Queries.ListUsers(c.Context())
		if err != nil {
			log.Println(err)
		}
		return c.JSON(results)
	})
	apiGroup := app.Group("/api")

	authGroup := apiGroup.Group("/auth")
	authController := controllers.AuthController{}

	routes.AuthRoutes(authGroup, authController)

	//TODO: Make CORS more strict in future

	app.Listen(":8080")
}
