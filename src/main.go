package main

import (
	"app/src/config"
	"app/src/database"
	"app/src/middleware"
	"app/src/router"
	"app/src/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
)

// @title go-fiber-boilerplate API documentation
// @version 1.0.0
// @license.name MIT
// @license.url https://github.com/indrayyana/go-fiber-boilerplate/blob/main/LICENSE
// @host localhost:3000
// @BasePath /v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Example Value: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
func main() {
	app := fiber.New(config.FiberConfig())

	db := database.Connect(config.DBHost, config.DBName)
	database.Migrate(db)

	// limit repeated failed requests to auth endpoints
	app.Use("/v1/auth", middleware.LimiterConfig())

	// logging
	app.Use(middleware.LoggerConfig())

	// set security HTTP headers
	app.Use(helmet.New())

	// gzip compression
	app.Use(compress.New())

	// enable cors
	app.Use(cors.New())

	// recover panic handler
	app.Use(middleware.RecoverConfig())

	// API routes
	router.Routes(app, db)

	// send back a 404 error for any unknown api request
	app.Use(utils.NotFoundHandler)

	address := fmt.Sprintf("%s:%d", config.AppHost, config.AppPort)
	utils.Log.Fatal(app.Listen(address))
}
