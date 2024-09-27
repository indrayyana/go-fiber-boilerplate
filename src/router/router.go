package router

import (
	"app/src/config"
	"app/src/services"
	"app/src/validation"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Routes(app *fiber.App, db *gorm.DB) {
	validate := validation.Validator()

	emailService := services.NewEmailService()
	userService := services.NewUserService(db, validate)
	tokenService := services.NewTokenService(db, validate, userService)
	authService := services.NewAuthService(db, validate, userService, tokenService)

	v1 := app.Group("/v1")

	AuthRoutes(v1, authService, userService, tokenService, emailService)
	UserRoutes(v1, userService, tokenService)
	// TODO: add another routes here...

	if !config.IsProd {
		DocsRoutes(v1)
	}
}
