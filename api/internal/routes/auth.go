package routes

import (
	"github.com/chenshuiluke/di-pocket-watcher/api/internal/controllers"
	"github.com/chenshuiluke/di-pocket-watcher/api/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(group fiber.Router, ac controllers.AuthController) {
	group.Get("/", middlewares.Jwt, ac.GetCurrentUser)
	group.Get("/google/login", ac.HandleGoogleLogin)
	group.Get("/google/callback", ac.HandleGoogleCallback)
}
