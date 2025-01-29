package routes

import (
	usercontrol "park/controller/authConrol"
	"park/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuthRoute(app *fiber.App) {
	auth := app.Group("/api/v1/auth")
	auth.Post("/register", usercontrol.Register)
	auth.Post("/login", usercontrol.Login)
	auth.Post("/logout", middleware.Auth, usercontrol.Logout)
	auth.Get("/me", middleware.Auth, usercontrol.Me)

}
