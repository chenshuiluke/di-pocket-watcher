package routes

import (
	"github.com/chenshuiluke/di-pocket-watcher/api/internal/controllers"
	"github.com/chenshuiluke/di-pocket-watcher/api/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

func TransactionAnalysisRoutes(group fiber.Router, tac controllers.TransactionAnalysisController) {
	group.Get("/email", middlewares.Jwt(), tac.GetTransactionEmails)
}
