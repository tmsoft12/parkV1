package routes

import (
	"park/controller/accountant"
	"park/middleware"

	"github.com/gofiber/fiber/v2"
)

func AccountantRoutes(app *fiber.App) {
	act := app.Group("/api/v1/accountant")
	act.Get("/calculateMoney", middleware.Auth, accountant.CalculateMoney)
	act.Get("/operators", middleware.Auth, accountant.GetOperators)
}
