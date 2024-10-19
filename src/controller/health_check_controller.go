package controller

import (
	"app/src/response"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

type HealthCheckController struct {
	HealthCheckService service.HealthCheckService
}

// NewHealthCheckController creates a new instance of the HealthCheckController
func NewHealthCheckController(
	healthCheckService service.HealthCheckService,
) *HealthCheckController {
	return &HealthCheckController{
		HealthCheckService: healthCheckService,
	}
}

// addServiceStatus adds the service status to the health check response list
func (h *HealthCheckController) addServiceStatus(serviceList *[]response.HealthCheck, name string, isUp bool) {
	status := "Up"
	if !isUp {
		status = "Down"
	}
	*serviceList = append(*serviceList, response.HealthCheck{
		Name:   name,
		Status: status,
		IsUp:   isUp,
	})
}

// Check handles the health check endpoint
func (h *HealthCheckController) Check(c *fiber.Ctx) error {
	var isHealthy bool = true
	var serviceList []response.HealthCheck

	// Check the database connection
	if err := h.HealthCheckService.GormCheck(); err != nil {
		isHealthy = false
		h.addServiceStatus(&serviceList, "Postgre", false)
	} else {
		h.addServiceStatus(&serviceList, "Postgre", true)
	}

	// Return the response based on health check result
	statusCode := fiber.StatusOK
	if !isHealthy {
		statusCode = fiber.StatusInternalServerError
	}

	return c.Status(statusCode).JSON(response.HealthCheckResponse{
		Status:    "success",
		Message:   "Health check successful",
		Code: statusCode,
		IsHealthy: isHealthy,
		Result: serviceList,
	})
}
