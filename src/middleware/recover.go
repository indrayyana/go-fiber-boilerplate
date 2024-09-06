package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func RecoverConfig() fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
	})
}
