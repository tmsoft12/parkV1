package routes

import (
	testz "park/controller/testZ"

	"github.com/gofiber/fiber/v2"
)

func InitZreport(app *fiber.App) {
	app.Post("/zreport", testz.GetZdata)
}
