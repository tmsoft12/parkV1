package routes

import (
	"park/controller/operator"
	"park/middleware"

	"github.com/gofiber/fiber/v2"
)

func Init(app *fiber.App) {
	cars := app.Group("/api/v1")
	cars.Get("/getallcars", middleware.Auth, operator.GetCars)
	cars.Get("/getcar/:id", middleware.Auth, operator.GetCar)
	cars.Get("/searchcar", middleware.Auth, operator.SearchCar)

}
