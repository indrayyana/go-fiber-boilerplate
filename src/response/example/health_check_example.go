package example

type HealthCheck struct {
	Name   string `json:"name" example:"Postgre"`
	Status string `json:"status" example:"Up"`
	IsUp   bool   `json:"is_up" example:"true"`
}

type HealthCheckResponse struct {
	Code      int           `json:"code" example:"200"`
	Status    string        `json:"status" example:"success"`
	Message   string        `json:"message" example:"Health check completed"`
	IsHealthy bool          `json:"is_healthy" example:"true"`
	Result    []HealthCheck `json:"result"`
}

type HealthCheckError struct {
	Name    string  `json:"name" example:"Postgre"`
	Status  string  `json:"status" example:"Down"`
	IsUp    bool    `json:"is_up" example:"false"`
	Message *string `json:"message,omitempty" example:"failed to connect to 'host=localhost user=postgres database=wrongdb': server error (FATAL: database \"wrongdb\" does not exist (SQLSTATE 3D000))"`
}

type HealthCheckResponseError struct {
	Code      int                `json:"code" example:"500"`
	Status    string             `json:"status" example:"error"`
	Message   string             `json:"message" example:"Health check completed"`
	IsHealthy bool               `json:"is_healthy" example:"false"`
	Result    []HealthCheckError `json:"result"`
}
