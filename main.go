package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"

	"park/controller/imagetoplate"
	"park/controller/operator"
	"park/database"
	_ "park/docs"
	"park/routes"
	"park/util"
)

// @title Airline REST API
// @host 192.168.100.7:3000
// @BasePath /
func main() {
	database.ConnectDB()
	util.LoadVIPPlates()

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		// AllowOrigins: "http://172.16.4.204",
		AllowOrigins: "*",
		// AllowCredentials: true,
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))
	app.Get("/swagger/*", swagger.HandlerDefault)
	go operator.HandleMessages()
	go imagetoplate.WatchDirectory("image", database.DB)

	routes.AuthRoute(app)
	routes.InitAdminRoute(app)
	routes.CameraRoutes(app)
	routes.AccountantRoutes(app)
	routes.InitZreport(app)
	routes.InitRealtime(app)
	routes.Init(app)
	app.Listen(":3000")
}
