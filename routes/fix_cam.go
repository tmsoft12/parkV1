package routes

import (
	camfix "park/controller/camFix"

	"github.com/gofiber/fiber/v2"
)

func FixRoute(app *fiber.App) {
	app.Post("/api/v1/addcam", camfix.AddCam)
	app.Get("/api/v1/cams", camfix.GetAllCams)
	app.Put("/api/v1/update-channel-ids", camfix.UpdateChannelIdsByChannelName)
	app.Put("/api/v1/updatemac", camfix.UpdateMacUser)
	app.Patch("/api/v1/type/:id", camfix.UpdateCamera)
	app.Delete("/api/v1/deletecam/:id", camfix.DeleteCam)
	app.Get("/api/v1/sync-camfix", camfix.SyncCamFixWithConfig)
}
