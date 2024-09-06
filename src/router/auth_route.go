package router

import (
	"app/src/config"
	"app/src/controllers"
	"app/src/services"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(v1 fiber.Router, a services.AuthService, u services.UserService, t services.TokenService) {
	authController := controllers.NewAuthController(a, u, t)
	config.GoogleConfig()

	auth := v1.Group("/auth")

	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Post("/logout", authController.Logout)
	auth.Post("/refresh-tokens", authController.RefreshToken)
	auth.Get("/google", authController.GoogleLogin)
	auth.Get("/google-callback", authController.GoogleCallback)
}
