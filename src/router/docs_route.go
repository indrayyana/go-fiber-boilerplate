package router

import (
	// initialize the Swagger documentation
	_ "app/src/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func DocsRoutes(v1 fiber.Router) {
	docs := v1.Group("/docs")

	docs.Get("/*", swagger.HandlerDefault)
}
