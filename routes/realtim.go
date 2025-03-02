package routes

import (
	"park/controller/realtime"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func InitRealtime(app *fiber.App) {
	app.Put("/api/v1/update/count", realtime.UpdateCount)
	app.Get("/api/v1/update/count", websocket.New(realtime.HandleWebSocketCount), realtime.GetAllCounts)
}
