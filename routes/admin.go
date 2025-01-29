package routes

import (
	admincontrol "park/controller/adminControl"
	pdfGenerator "park/controller/pdf"

	"github.com/gofiber/fiber/v2"
)

func InitAdminRoute(app *fiber.App) {
	user := app.Group("/api/v1")
	user.Post("/users", admincontrol.CreateUser)
	user.Get("/users", admincontrol.GetAllUsers)
	user.Get("/user/operators", admincontrol.GetOperator)
	user.Get("/users/:id", admincontrol.UserGetByID)
	user.Put("/users/:id", admincontrol.UserUpdate)
	user.Delete("/users/:id", admincontrol.UserDelete)
	user.Get("/userCount", admincontrol.UsersCount)
	user.Post("/pdf", pdfGenerator.CreatePDF)

	camera := app.Group("/api/v1/")
	camera.Post("/cameras", admincontrol.CreateCamera)
	camera.Put("/cameras/:id", admincontrol.UpdateCamera)
	camera.Get("/cameras/:id", admincontrol.GetCameraByID)
	camera.Delete("/cameras/:id", admincontrol.DeleteCamera)
	camera.Get("/cameras/", admincontrol.GetCameras)
}
