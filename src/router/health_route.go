package router

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func HealthRoutes(v1 fiber.Router, db *gorm.DB) {

	// Health check endpoint
	v1.Get("/health", func(c *fiber.Ctx) error {
		// Use GORM's DB() method to get the *sql.DB instance
		sqlDB, err := db.DB()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "DOWN",
				"error":  err.Error(),
			})
		}

		// Check if the database connection is alive
		if err := sqlDB.Ping(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "DOWN",
				"database": fiber.Map{
					"status": "DOWN",
					"error":  err.Error(),
				},
			})
		}

		// Everything is fine
		return c.JSON(fiber.Map{
			"status":   "UP",
			"database": "UP",
		})
	})
}
