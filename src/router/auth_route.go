package router

import (
	"app/src/config"
	"app/src/controllers"
	m "app/src/middleware"
	"app/src/services"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(
	v1 fiber.Router, a services.AuthService, u services.UserService,
	t services.TokenService, e services.EmailService,
) {
	authController := controllers.NewAuthController(a, u, t, e)
	config.GoogleConfig()

	auth := v1.Group("/auth")

	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Post("/logout", authController.Logout)
	auth.Post("/refresh-tokens", authController.RefreshTokens)
	auth.Post("/forgot-password", authController.ForgotPassword)
	auth.Post("/reset-password", authController.ResetPassword)
	auth.Post("/send-verification-email", m.Auth(u), authController.SendVerificationEmail)
	auth.Post("/verify-email", authController.VerifyEmail)
	auth.Get("/google", authController.GoogleLogin)
	auth.Get("/google-callback", authController.GoogleCallback)
}
