package integration

import (
	"app/src/response"
	"app/test"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckRoutes(t *testing.T) {
	t.Run("GET /v1/health-check", func(t *testing.T) {
		t.Run("should return 200 and success response if request is ok", func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/v1/health-check", nil)

			msTimeout := 2000
			apiResponse, err := test.App.Test(request, msTimeout)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusOK, apiResponse.StatusCode)

			bytes, err := io.ReadAll(apiResponse.Body)
			assert.Nil(t, err)

			responseBody := new(response.HealthCheckResponse)

			err = json.Unmarshal(bytes, responseBody)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusOK, apiResponse.StatusCode)
			assert.Equal(t, http.StatusOK, responseBody.Code)
			assert.Equal(t, "success", responseBody.Status)
			assert.Equal(t, "Health check completed", responseBody.Message)
			assert.Equal(t, true, responseBody.IsHealthy)
			assert.Equal(t, []response.HealthCheck{
				{
					Name:   "Postgre",
					Status: "Up",
					IsUp:   true,
				},
				{
					Name:   "Memory",
					Status: "Up",
					IsUp:   true,
				},
			}, responseBody.Result)
		})

		// t.Run("should return 500 and error response if request failed", func(t *testing.T) {
		// 	request := httptest.NewRequest(http.MethodGet, "/v1/health-check", nil)

		// 	msTimeout := 2000
		// 	apiResponse, err := test.App.Test(request, msTimeout)
		// 	assert.Nil(t, err)

		// 	assert.Equal(t, http.StatusInternalServerError, apiResponse.StatusCode)

		// 	bytes, err := io.ReadAll(apiResponse.Body)
		// 	assert.Nil(t, err)

		// 	responseBody := new(response.HealthCheckResponse)

		// 	err = json.Unmarshal(bytes, responseBody)
		// 	assert.Nil(t, err)

		// 	assert.Equal(t, http.StatusInternalServerError, apiResponse.StatusCode)
		// 	assert.Equal(t, http.StatusInternalServerError, responseBody.Code)
		// 	assert.Equal(t, "error", responseBody.Status)
		// 	assert.Equal(t, "Health check completed", responseBody.Message)
		// 	assert.Equal(t, false, responseBody.IsHealthy)
		// 	assert.Equal(t, []response.HealthCheck{
		// 		{
		// 			Name:   "Postgre",
		// 			Status: "Down",
		// 			IsUp:   false,
		// 		},
		// 	}, responseBody.Result)
		// })
	})
}
