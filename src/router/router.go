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

	userService := services.NewUserService(db, validate)
	tokenService := services.NewTokenService(db)
	authService := services.NewAuthService(db, validate, userService, tokenService)

	v1 := app.Group("/v1")

	AuthRoutes(v1, authService, userService, tokenService)
	UserRoutes(v1, userService, tokenService)
	// add another routes here...

	if !config.IsProd {
		DocsRoutes(v1)
	}
}
