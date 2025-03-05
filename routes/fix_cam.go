package routes

import (
	camfix "park/controller/camFix"

	"github.com/gofiber/fiber/v2"
)

func FixRoute(app *fiber.App) {
	app.Post("/api/v1/addcam", camfix.AddCam)
}
