package controller

import (
	"app/src/response"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

type HealthCheckController struct {
	HealthCheckService service.HealthCheckService
}

func NewHealthCheckController(healthCheckService service.HealthCheckService) *HealthCheckController {
	return &HealthCheckController{
		HealthCheckService: healthCheckService,
	}
}

func (h *HealthCheckController) addServiceStatus(
	serviceList *[]response.HealthCheck, name string, isUp bool, message *string,
) {
	status := "Up"
	if !isUp {
		status = "Down"
	}
	*serviceList = append(*serviceList, response.HealthCheck{
		Name:    name,
		Status:  status,
		IsUp:    isUp,
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
	isHealthy := true
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
	status := "success"

	if !isHealthy {
		statusCode = fiber.StatusInternalServerError
		status = "error"
	}

	return c.Status(statusCode).JSON(response.HealthCheckResponse{
		Status:    status,
		Message:   "Health check completed",
		Code:      statusCode,
		IsHealthy: isHealthy,
		Result:    serviceList,
	})
}
