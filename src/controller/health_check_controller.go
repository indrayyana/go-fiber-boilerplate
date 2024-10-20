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
func (h *HealthCheckController) addServiceStatus(serviceList *[]response.HealthCheck, name string, isUp bool, message *string) {
	status := "Up"
	if !isUp {
		status = "Down"
	}
	*serviceList = append(*serviceList, response.HealthCheck{
		Name:   name,
		Status: status,
		IsUp:   isUp,
		Message: message,
	})
}

// Check handles the health check endpoint
// @Summary Health Check
// @Description Check the status of services and database connections
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} response.HealthCheckResponse
// @Failure 500 {object} response.HealthCheckResponse
// @Router /health [get]
func (h *HealthCheckController) Check(c *fiber.Ctx) error {
	var isHealthy bool = true
	var serviceList []response.HealthCheck

	// Check the database connection
	if err := h.HealthCheckService.GormCheck(); err != nil {
		isHealthy = false
		errMsg := err.Error()
		h.addServiceStatus(&serviceList, "Postgre", false, &errMsg)
	} else {
		h.addServiceStatus(&serviceList, "Postgre", true, nil)
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
