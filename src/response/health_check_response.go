package response

type HealthCheck struct {
	Name    string  `json:"name"`
	Status  string  `json:"status"`
	IsUp    bool    `json:"is_up"`
	Message *string `json:"message,omitempty"`
}

type HealthCheckResponse struct {
	Code      int           `json:"code"`
	Status    string        `json:"status"`
	Message   string        `json:"message"`
	IsHealthy bool          `json:"is_healthy"`
	Result    []HealthCheck `json:"result"`
}
