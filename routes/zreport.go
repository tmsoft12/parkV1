package routes

import (
	zreport "park/controller/zreport"

	"github.com/gofiber/fiber/v2"
)

func InitZreport(app *fiber.App) {
	app.Post("/zreport", zreport.GetZdata)
}
