package router

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
)

func MetricsRoutes(app *fiber.App) {

	prometheus := fiberprometheus.New("go-fiber-boilerplate")
	app.Use(prometheus.Middleware)
	prometheus.RegisterAt(app, "/metrics")

}
