package routes

import (
	"park/controller/accountant"
	tarifcontrol "park/controller/tarifControl"
	"park/middleware"

	"github.com/gofiber/fiber/v2"
)

func AccountantRoutes(app *fiber.App) {
	act := app.Group("/api/v1/accountant")
	act.Get("/calculateMoney", middleware.Auth, accountant.CalculateMoney)
	act.Get("/operators", middleware.Auth, accountant.GetOperators)
	act.Post("/tarif", tarifcontrol.CreateTarif)
	act.Delete("/tarif/:id", tarifcontrol.DeleteTarif)
	act.Get("/tarif", tarifcontrol.GetAllTarif)
}
