package routes

import (
	"park/controller/accountant"

	"github.com/gofiber/fiber/v2"
)

func AccountantRoutes(app *fiber.App) {
	act := app.Group("/api/v1/accountant")
	act.Get("/calculateMoney", accountant.CalculateMoney)
}
