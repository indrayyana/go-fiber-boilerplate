package integration

import (
	"app/src/config"
	"app/src/model"
	"app/src/response"
	"app/src/utils"
	"app/src/validation"
	"app/test"
	"app/test/fixture"
	"app/test/helper"
	_ "embed"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAuthRoutes(t *testing.T) {
	t.Run("POST /v1/auth/register", func(t *testing.T) {
		var requestBody = validation.Register{
			Name:     "Test",
			Email:    "test@gmail.com",
			Password: "password1",
		}

		t.Run("should return 201 and successfully register user if request data is ok", func(t *testing.T) {
			helper.ClearAll(test.DB)
			bodyJSON, err := json.Marshal(requestBody)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPost, "/v1/auth/register", strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")

			apiResponse, err := test.App.Test(request)

			assert.Nil(t, err)

			bytes, err := io.ReadAll(apiResponse.Body)
			assert.Nil(t, err)

			responseBody := new(response.SuccessWithTokens[model.User])

			err = json.Unmarshal(bytes, responseBody)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusCreated, apiResponse.StatusCode)
			assert.Equal(t, "success", responseBody.Status)
			assert.NotContains(t, string(bytes), "password")
			assert.NotNil(t, responseBody.Data.ID)
			assert.Equal(t, requestBody.Name, responseBody.Data.Name)
			assert.Equal(t, requestBody.Email, responseBody.Data.Email)
			assert.Equal(t, "user", responseBody.Data.Role)
			assert.Equal(t, false, responseBody.Data.VerifiedEmail)
			assert.NotNil(t, responseBody.Tokens.Access.Token)
			assert.NotNil(t, responseBody.Tokens.Refresh.Token)

			user, err := helper.GetUserByID(test.DB, responseBody.Data.ID.String())
			assert.Nil(t, err)

			assert.NotNil(t, user)
			assert.NotEqual(t, user.Password, requestBody.Password)
			assert.Equal(t, user.Name, requestBody.Name)
			assert.Equal(t, user.Email, requestBody.Email)
			assert.Equal(t, user.Role, "user")
			assert.Equal(t, user.VerifiedEmail, false)
		})

		t.Run("should return 400 error if email is invalid", func(t *testing.T) {
			helper.ClearAll(test.DB)
			requestBody.Email = "invalidEmail"

			bodyJSON, err := json.Marshal(requestBody)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPost, "/v1/auth/register", strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusBadRequest, apiResponse.StatusCode)
		})

		t.Run("should return 409 error if email is already used", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.CreateUser(test.DB, "test@gmail.com", "test1234", "Test")
			requestBody.Email = "test@gmail.com"

			bodyJSON, err := json.Marshal(requestBody)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPost, "/v1/auth/register", strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusConflict, apiResponse.StatusCode)
		})

		t.Run("should return 400 error if password length is less than 8 characters", func(t *testing.T) {
			helper.ClearAll(test.DB)
			requestBody.Password = "passwo1"

			bodyJSON, err := json.Marshal(requestBody)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPost, "/v1/auth/register", strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusBadRequest, apiResponse.StatusCode)
		})

		t.Run("should return 400 error if password does not contain both letters and numbers", func(t *testing.T) {
			helper.ClearAll(test.DB)
			requestBody.Password = "password"

			bodyJSON, err := json.Marshal(requestBody)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPost, "/v1/auth/register", strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusBadRequest, apiResponse.StatusCode)

			requestBody.Password = "11111111"

			bodyJSON, err = json.Marshal(requestBody)
			assert.Nil(t, err)

			request = httptest.NewRequest(http.MethodPost, "/v1/auth/register", strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")

			apiResponse, err = test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusBadRequest, apiResponse.StatusCode)
		})
	})
	t.Run("POST /v1/auth/login", func(t *testing.T) {
		t.Run("should return 200 and login user if email and password match", func(t *testing.T) {
			helper.CreateUser(test.DB, "test@gmail.com", "test1234", "Test User")
			loginCredentials := &validation.Login{
				Email:    "test@gmail.com",
				Password: "test1234",
			}

			bodyJSON, err := json.Marshal(loginCredentials)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			bytes, err := io.ReadAll(apiResponse.Body)
			assert.Nil(t, err)

			responseBody := new(response.SuccessWithTokens[model.User])

			err = json.Unmarshal(bytes, responseBody)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusOK, apiResponse.StatusCode)
			assert.Equal(t, "success", responseBody.Status)
			assert.NotNil(t, responseBody.Data.ID)
			assert.Equal(t, "Test User", responseBody.Data.Name)
			assert.Equal(t, "test@gmail.com", responseBody.Data.Email)
			assert.Equal(t, "user", responseBody.Data.Role)
			assert.Equal(t, false, responseBody.Data.VerifiedEmail)
			assert.NotNil(t, responseBody.Tokens.Access.Token)
			assert.NotNil(t, responseBody.Tokens.Refresh.Token)
		})

		t.Run("should return 401 error if there are no users with that email", func(t *testing.T) {
			helper.ClearAll(test.DB)
			loginCredentials := &validation.Login{
				Email:    "nonexistent@gmail.com",
				Password: "test1234",
			}

			bodyJSON, err := json.Marshal(loginCredentials)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			bytes, err := io.ReadAll(apiResponse.Body)
			assert.Nil(t, err)

			responseBody := make(map[string]interface{})

			err = json.Unmarshal(bytes, &responseBody)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusUnauthorized, apiResponse.StatusCode)
			assert.Equal(t, "error", responseBody["status"])
			assert.Equal(t, "Incorrect email or password", responseBody["message"])
		})

		t.Run("should return 401 error if password is wrong", func(t *testing.T) {
			helper.CreateUser(test.DB, "test@gmail.com", "test1234", "Test User")
			loginCredentials := &validation.Login{
				Email:    "test@gmail.com",
				Password: "wrongPassword1",
			}

			bodyJSON, err := json.Marshal(loginCredentials)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			bytes, err := io.ReadAll(apiResponse.Body)
			assert.Nil(t, err)

			responseBody := make(map[string]interface{})

			err = json.Unmarshal(bytes, &responseBody)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusUnauthorized, apiResponse.StatusCode)
			assert.Equal(t, "error", responseBody["status"])
			assert.Equal(t, "Incorrect email or password", responseBody["message"])
		})
	})
	t.Run("POST /v1/auth/logout", func(t *testing.T) {
		t.Run("should return 200 if refresh token is valid", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.UserOne)

			expires := time.Now().Add(time.Hour * 24 * time.Duration(config.JWTRefreshExp))
			refreshToken, err := helper.GenerateToken(fixture.UserOne.ID.String(), fixture.UserOne.Role, expires, config.TokenTypeRefresh)
			assert.Nil(t, err)

			err = helper.SaveToken(test.DB, refreshToken, fixture.UserOne.ID.String(), expires)
			assert.Nil(t, err)

			bodyJSON, err := json.Marshal(validation.RefreshToken{Token: refreshToken})
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPost, "/v1/auth/logout", strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")

			apiresponse, err := test.App.Test(request)
			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, apiresponse.StatusCode)

			dbRefreshTokenDoc, _ := helper.GetTokenByUserID(test.DB, refreshToken)
			assert.Nil(t, dbRefreshTokenDoc)
		})

		t.Run("should return 400 error if refresh token is missing from request body", func(t *testing.T) {
			helper.ClearAll(test.DB)

			request := httptest.NewRequest(http.MethodPost, "/v1/auth/logout", nil)
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")

			apiresponse, err := test.App.Test(request)
			assert.Nil(t, err)
			assert.Equal(t, http.StatusBadRequest, apiresponse.StatusCode)
		})

		t.Run("should return 404 error if refresh token is not found in the database", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.UserOne)

			expires := time.Now().Add(time.Hour * 24 * time.Duration(config.JWTRefreshExp))
			refreshToken, err := helper.GenerateToken(fixture.UserOne.ID.String(), fixture.UserOne.Role, expires, config.TokenTypeRefresh)
			assert.Nil(t, err)

			bodyJSON, err := json.Marshal(validation.RefreshToken{Token: refreshToken})
			assert.Nil(t, err)
			request := httptest.NewRequest(http.MethodPost, "/v1/auth/logout", strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")

			apiresponse, err := test.App.Test(request)
			assert.Nil(t, err)
			assert.Equal(t, http.StatusNotFound, apiresponse.StatusCode)
		})
	})
}

func TestAuthMiddleware(t *testing.T) {
	t.Run("should call next with no errors if access token is valid", func(t *testing.T) {
		helper.ClearAll(test.DB)
		helper.InsertUser(test.DB, fixture.UserOne)

		expires := time.Now().Add(time.Minute * time.Duration(config.JWTAccessExp))
		token, err := helper.GenerateToken(fixture.UserOne.ID.String(), fixture.UserOne.Role, expires, config.TokenTypeAccess)
		assert.Nil(t, err)

		request := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
		request.Header.Set("Authorization", "Bearer "+token)

		userID, err := utils.VerifyToken(token, config.JWTSecret, config.TokenTypeAccess)
		assert.Nil(t, err)

		assert.Equal(t, fixture.UserOne.ID.String(), userID)
	})

	t.Run("should call next with unauthorized error if access token is not found in header", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
		apiResponse, err := test.App.Test(request)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnauthorized, apiResponse.StatusCode)
	})

	t.Run("should call next with unauthorized error if access token is not a valid jwt token", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
		request.Header.Set("Authorization", "Bearer randomToken")

		apiResponse, err := test.App.Test(request)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnauthorized, apiResponse.StatusCode)
	})

	t.Run("should call next with unauthorized error if the token is not an access token", func(t *testing.T) {
		helper.ClearAll(test.DB)
		helper.InsertUser(test.DB, fixture.UserOne)

		expires := time.Now().Add(time.Minute * time.Duration(config.JWTAccessExp))
		refreshToken, err := helper.GenerateToken(fixture.UserOne.ID.String(), fixture.UserOne.Role, expires, config.TokenTypeRefresh)
		assert.Nil(t, err)

		request := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
		request.Header.Set("Authorization", "Bearer "+refreshToken)

		apiResponse, err := test.App.Test(request)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnauthorized, apiResponse.StatusCode)
	})

	t.Run("should call next with unauthorized error if access token is generated with an invalid secret", func(t *testing.T) {
		helper.ClearAll(test.DB)
		helper.InsertUser(test.DB, fixture.UserOne)

		expires := time.Now().Add(time.Minute * time.Duration(config.JWTAccessExp))
		accessToken, err := helper.GenerateInvalidToken(fixture.UserOne.ID.String(), fixture.UserOne.Role, expires, config.TokenTypeAccess)
		assert.Nil(t, err)

		request := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
		request.Header.Set("Authorization", "Bearer "+accessToken)

		apiResponse, err := test.App.Test(request)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnauthorized, apiResponse.StatusCode)
	})

	t.Run("should call next with unauthorized error if access token is expired", func(t *testing.T) {
		helper.ClearAll(test.DB)
		helper.InsertUser(test.DB, fixture.UserOne)

		expires := time.Now().Add(time.Second * 1)
		accessToken, err := helper.GenerateToken(fixture.UserOne.ID.String(), fixture.UserOne.Role, expires, config.TokenTypeAccess)
		assert.Nil(t, err)

		time.Sleep(2 * time.Second)

		request := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
		request.Header.Set("Authorization", "Bearer "+accessToken)

		apiResponse, err := test.App.Test(request)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnauthorized, apiResponse.StatusCode)
	})

	t.Run("should call next with unauthorized error if user is not found", func(t *testing.T) {
		helper.ClearAll(test.DB)

		expires := time.Now().Add(time.Minute * time.Duration(config.JWTAccessExp))
		accessToken, err := helper.GenerateToken(fixture.UserOne.ID.String(), fixture.UserOne.Role, expires, config.TokenTypeAccess)
		assert.Nil(t, err)

		request := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
		request.Header.Set("Authorization", "Bearer "+accessToken)

		apiResponse, err := test.App.Test(request)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusUnauthorized, apiResponse.StatusCode)
	})

	t.Run("should call next with forbidden error if user does not have required rights and userId is not in params", func(t *testing.T) {
		helper.ClearAll(test.DB)
		helper.InsertUser(test.DB, fixture.UserOne)

		expires := time.Now().Add(time.Minute * time.Duration(config.JWTAccessExp))
		accessToken, err := helper.GenerateToken(fixture.UserOne.ID.String(), fixture.UserOne.Role, expires, config.TokenTypeAccess)
		assert.Nil(t, err)

		request := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
		request.Header.Set("Authorization", "Bearer "+accessToken)

		apiResponse, err := test.App.Test(request)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusForbidden, apiResponse.StatusCode)
	})

	t.Run("should call next with no errors if user does not have required rights but userId is in params", func(t *testing.T) {
		helper.ClearAll(test.DB)
		helper.InsertUser(test.DB, fixture.UserOne)

		expires := time.Now().Add(time.Minute * time.Duration(config.JWTAccessExp))
		accessToken, err := helper.GenerateToken(fixture.UserOne.ID.String(), fixture.UserOne.Role, expires, config.TokenTypeAccess)
		assert.Nil(t, err)

		request := httptest.NewRequest(http.MethodGet, "/v1/users/"+fixture.UserOne.ID.String(), nil)
		request.Header.Set("Authorization", "Bearer "+accessToken)

		apiResponse, err := test.App.Test(request)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, apiResponse.StatusCode)
	})

	t.Run("should call next with no errors if user has required rights", func(t *testing.T) {
		helper.ClearAll(test.DB)
		helper.InsertUser(test.DB, fixture.UserOne, fixture.Admin)

		expires := time.Now().Add(time.Minute * time.Duration(config.JWTAccessExp))
		accessToken, err := helper.GenerateToken(fixture.Admin.ID.String(), fixture.Admin.Role, expires, config.TokenTypeAccess)
		assert.Nil(t, err)

		request := httptest.NewRequest(http.MethodGet, "/v1/users/"+fixture.UserOne.ID.String(), nil)
		request.Header.Set("Authorization", "Bearer "+accessToken)

		apiResponse, err := test.App.Test(request)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, apiResponse.StatusCode)
	})
}
