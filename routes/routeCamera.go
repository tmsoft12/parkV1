package routes

import (
	"os"
	"park/controller/getdata"
	"park/controller/operator"
	"park/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func CameraRoutes(app *fiber.App) {

	app.Get("/ws/notification", websocket.New(operator.Ws))

	plate := os.Getenv("IMAGE_URL")
	app.Static("/plate", plate)

	camera := app.Group("/api/v1/camera")
	camera.Post("/getdata", getdata.CreateCarEntry)
	camera.Put("/getdata", getdata.CreateCarExit)
	camera.Put("/updatecar/:plate", middleware.Auth, operator.UpdateCar)
}
