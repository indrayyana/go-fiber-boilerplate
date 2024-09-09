package integration

import (
	"app/src/model"
	"app/src/response"
	"app/src/validation"
	"app/test"
	"app/test/fixture"
	"app/test/helper"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRoutes(t *testing.T) {
	t.Run("POST /v1/users", func(t *testing.T) {
		var newUser = validation.CreateUser{
			Name:     "Test",
			Email:    "test@gmail.com",
			Password: "password1",
			Role:     "user",
		}

		t.Run("should return 201 and successfully create new user if data is ok", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.Admin)

			adminAccessToken, err := fixture.AccessToken(fixture.Admin)
			assert.Nil(t, err)

			bodyJSON, err := json.Marshal(newUser)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")
			request.Header.Set("Authorization", "Bearer "+adminAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			bytes, err := io.ReadAll(apiResponse.Body)
			assert.Nil(t, err)

			responseBody := new(response.SuccessWithUser)

			err = json.Unmarshal(bytes, responseBody)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusCreated, apiResponse.StatusCode)
			assert.Equal(t, "success", responseBody.Status)
			assert.NotContains(t, string(bytes), "password")
			assert.NotNil(t, responseBody.User.ID)
			assert.Equal(t, newUser.Name, responseBody.User.Name)
			assert.Equal(t, newUser.Email, responseBody.User.Email)
			assert.Equal(t, "user", responseBody.User.Role)
			assert.Equal(t, false, responseBody.User.VerifiedEmail)

			user, err := helper.GetUserByID(test.DB, responseBody.User.ID.String())
			assert.Nil(t, err)

			assert.NotNil(t, user)
			assert.NotEqual(t, user.Password, newUser.Password)
			assert.Equal(t, user.Name, newUser.Name)
			assert.Equal(t, user.Email, newUser.Email)
			assert.Equal(t, user.Role, newUser.Role)
			assert.Equal(t, false, user.VerifiedEmail)
		})

		t.Run("should be able to create an admin as well", func(t *testing.T) {
			helper.ClearAll(test.DB)
			newUser.Role = "admin"
			helper.InsertUser(test.DB, fixture.Admin)

			adminAccessToken, err := fixture.AccessToken(fixture.Admin)
			assert.Nil(t, err)

			bodyJSON, err := json.Marshal(newUser)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")
			request.Header.Set("Authorization", "Bearer "+adminAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			bytes, err := io.ReadAll(apiResponse.Body)
			assert.Nil(t, err)

			responseBody := new(response.SuccessWithUser)

			err = json.Unmarshal(bytes, responseBody)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusCreated, apiResponse.StatusCode)
			assert.Equal(t, responseBody.User.Role, "admin")

			user, err := helper.GetUserByID(test.DB, responseBody.User.ID.String())
			assert.Nil(t, err)

			assert.Equal(t, user.Role, "admin")
		})

		t.Run("should return 401 error if access token is missing", func(t *testing.T) {
			helper.ClearAll(test.DB)

			bodyJSON, err := json.Marshal(newUser)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(string(bodyJSON)))

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusUnauthorized, apiResponse.StatusCode)
		})

		t.Run("should return 403 error if logged in user is not admin", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.UserOne)

			userOneAccessToken, err := fixture.AccessToken(fixture.UserOne)
			assert.Nil(t, err)

			bodyJSON, err := json.Marshal(newUser)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")
			request.Header.Set("Authorization", "Bearer "+userOneAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusForbidden, apiResponse.StatusCode)
		})

		t.Run("should return 400 error if email is invalid", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.Admin)
			newUser.Email = "invalidEmail"

			adminAccessToken, err := fixture.AccessToken(fixture.Admin)
			assert.Nil(t, err)

			bodyJSON, err := json.Marshal(newUser)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")
			request.Header.Set("Authorization", "Bearer "+adminAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusBadRequest, apiResponse.StatusCode)
		})

		t.Run("should return 409 error if email is already used", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.Admin, fixture.UserOne)
			newUser.Email = fixture.UserOne.Email

			adminAccessToken, err := fixture.AccessToken(fixture.Admin)
			assert.Nil(t, err)

			bodyJSON, err := json.Marshal(newUser)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")
			request.Header.Set("Authorization", "Bearer "+adminAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusConflict, apiResponse.StatusCode)
		})

		t.Run("should return 400 error if password length is less than 8 characters", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.Admin)
			newUser.Password = "passwo1"

			adminAccessToken, err := fixture.AccessToken(fixture.Admin)
			assert.Nil(t, err)

			bodyJSON, err := json.Marshal(newUser)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")
			request.Header.Set("Authorization", "Bearer "+adminAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusBadRequest, apiResponse.StatusCode)
		})

		t.Run("should return 400 error if password does not contain both letters and numbers", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.Admin)
			newUser.Password = "password"

			adminAccessToken, err := fixture.AccessToken(fixture.Admin)
			assert.Nil(t, err)

			bodyJSON, err := json.Marshal(newUser)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")
			request.Header.Set("Authorization", "Bearer "+adminAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusBadRequest, apiResponse.StatusCode)

			newUser.Password = "1111111"

			bodyJSON, err = json.Marshal(newUser)
			assert.Nil(t, err)

			request = httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")
			request.Header.Set("Authorization", "Bearer "+adminAccessToken)

			apiResponse, err = test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusBadRequest, apiResponse.StatusCode)
		})

		t.Run("should return 400 error if role is neither user nor admin", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.Admin)
			newUser.Role = "invalid"

			adminAccessToken, err := fixture.AccessToken(fixture.Admin)
			assert.Nil(t, err)

			bodyJSON, err := json.Marshal(newUser)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")
			request.Header.Set("Authorization", "Bearer "+adminAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusBadRequest, apiResponse.StatusCode)
		})

		t.Run("should return 400 error if role is neither user or admin", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.Admin)
			newUser.Role = "invalid"

			adminAccessToken, err := fixture.AccessToken(fixture.Admin)
			assert.Nil(t, err)

			bodyJSON, err := json.Marshal(newUser)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")
			request.Header.Set("Authorization", "Bearer "+adminAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusBadRequest, apiResponse.StatusCode)
		})
	})

	t.Run("GET /v1/users", func(t *testing.T) {
		t.Run("should return 200 and apply the default query options", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.UserOne, fixture.UserTwo, fixture.Admin)

			adminAccessToken, err := fixture.AccessToken(fixture.Admin)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
			request.Header.Set("Authorization", "Bearer "+adminAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			bytes, err := io.ReadAll(apiResponse.Body)
			assert.Nil(t, err)

			responseBody := new(response.SuccessWithPaginate[model.User])

			err = json.Unmarshal(bytes, responseBody)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusOK, apiResponse.StatusCode)
			assert.Equal(t, 1, responseBody.Page)
			assert.Equal(t, 10, responseBody.Limit)
			assert.Equal(t, int64(1), responseBody.TotalPages)
			assert.Equal(t, int64(3), responseBody.TotalResults)

			assert.Len(t, responseBody.Results, 3)
			assert.Equal(t, fixture.UserOne.ID, responseBody.Results[0].ID)
			assert.Equal(t, fixture.UserOne.Name, responseBody.Results[0].Name)
			assert.Equal(t, fixture.UserOne.Email, responseBody.Results[0].Email)
			assert.Equal(t, fixture.UserOne.Role, responseBody.Results[0].Role)
			assert.Equal(t, fixture.UserOne.VerifiedEmail, responseBody.Results[0].VerifiedEmail)
		})

		t.Run("should return 401 if access token is missing", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.UserOne, fixture.UserTwo, fixture.Admin)

			request := httptest.NewRequest(http.MethodGet, "/v1/users", nil)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusUnauthorized, apiResponse.StatusCode)
		})

		t.Run("should return 403 if a non-admin is trying to access all users", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.UserOne, fixture.UserTwo, fixture.Admin)

			userOneAccessToken, err := fixture.AccessToken(fixture.UserOne)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
			request.Header.Set("Authorization", "Bearer "+userOneAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusForbidden, apiResponse.StatusCode)
		})

		t.Run("should limit returned array if limit param is specified", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.UserOne, fixture.UserTwo, fixture.Admin)

			adminAccessToken, err := fixture.AccessToken(fixture.Admin)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodGet, "/v1/users?limit=2", nil)
			request.Header.Set("Authorization", "Bearer "+adminAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			bytes, err := io.ReadAll(apiResponse.Body)
			assert.Nil(t, err)

			responseBody := new(response.SuccessWithPaginate[model.User])

			err = json.Unmarshal(bytes, responseBody)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusOK, apiResponse.StatusCode)
			assert.Equal(t, 1, responseBody.Page)
			assert.Equal(t, 2, responseBody.Limit)
			assert.Equal(t, int64(2), responseBody.TotalPages)
			assert.Equal(t, int64(3), responseBody.TotalResults)

			assert.Len(t, responseBody.Results, 2)
			assert.Equal(t, fixture.UserOne.ID, responseBody.Results[0].ID)
			assert.Equal(t, fixture.UserTwo.ID, responseBody.Results[1].ID)
		})

		t.Run("should return the correct page if page and limit params are specified", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.UserOne, fixture.UserTwo, fixture.Admin)

			adminAccessToken, err := fixture.AccessToken(fixture.Admin)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodGet, "/v1/users?page=2&limit=2", nil)
			request.Header.Set("Authorization", "Bearer "+adminAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			bytes, err := io.ReadAll(apiResponse.Body)
			assert.Nil(t, err)

			responseBody := new(response.SuccessWithPaginate[model.User])

			err = json.Unmarshal(bytes, responseBody)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusOK, apiResponse.StatusCode)
			assert.Equal(t, 2, responseBody.Page)
			assert.Equal(t, 2, responseBody.Limit)
			assert.Equal(t, int64(2), responseBody.TotalPages)
			assert.Equal(t, int64(3), responseBody.TotalResults)

			assert.Len(t, responseBody.Results, 1)
			assert.Equal(t, fixture.Admin.ID, responseBody.Results[0].ID)
		})
	})

	t.Run("GET /v1/users/:userId", func(t *testing.T) {
		t.Run("should return 200 and the user object if data is ok", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.UserOne)

			userOneAccessToken, err := fixture.AccessToken(fixture.UserOne)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodGet, "/v1/users/"+fixture.UserOne.ID.String(), nil)
			request.Header.Set("Authorization", "Bearer "+userOneAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			bytes, err := io.ReadAll(apiResponse.Body)
			assert.Nil(t, err)

			responseBody := new(response.SuccessWithUser)
			err = json.Unmarshal(bytes, responseBody)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusOK, apiResponse.StatusCode)
			assert.NotContains(t, string(bytes), "password")
			assert.Equal(t, responseBody.User.ID, fixture.UserOne.ID)
			assert.Equal(t, responseBody.User.Email, fixture.UserOne.Email)
			assert.Equal(t, responseBody.User.Name, fixture.UserOne.Name)
			assert.Equal(t, responseBody.User.Role, fixture.UserOne.Role)
			assert.Equal(t, responseBody.User.VerifiedEmail, fixture.UserOne.VerifiedEmail)
		})

		t.Run("should return 401 error if access token is missing", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.UserOne)

			request := httptest.NewRequest(http.MethodGet, "/v1/users/"+fixture.UserOne.ID.String(), nil)
			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusUnauthorized, apiResponse.StatusCode)
		})

		t.Run("should return 403 error if user is trying to get another user", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.UserOne, fixture.UserTwo)

			userOneAccessToken, err := fixture.AccessToken(fixture.UserOne)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodGet, "/v1/users/"+fixture.UserTwo.ID.String(), nil)
			request.Header.Set("Authorization", "Bearer "+userOneAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusForbidden, apiResponse.StatusCode)
		})

		t.Run("should return 200 and the user object if admin is trying to get another user", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.UserOne, fixture.Admin)

			adminAccessToken, err := fixture.AccessToken(fixture.Admin)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodGet, "/v1/users/"+fixture.UserOne.ID.String(), nil)
			request.Header.Set("Authorization", "Bearer "+adminAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusOK, apiResponse.StatusCode)
		})

		t.Run("should return 400 error if userId is not a valid postgres id", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.Admin)

			adminAccessToken, err := fixture.AccessToken(fixture.Admin)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodGet, "/v1/users/invalidId", nil)
			request.Header.Set("Authorization", "Bearer "+adminAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusBadRequest, apiResponse.StatusCode)
		})

		t.Run("should return 404 error if user is not found", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.Admin)

			adminAccessToken, err := fixture.AccessToken(fixture.Admin)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodGet, "/v1/users/"+fixture.UserOne.ID.String(), nil)
			request.Header.Set("Authorization", "Bearer "+adminAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusNotFound, apiResponse.StatusCode)
		})
	})

	t.Run("DELETE /v1/users/:userId", func(t *testing.T) {
		t.Run("should return 200 if data is ok", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.UserOne)

			userOneAccessToken, err := fixture.AccessToken(fixture.UserOne)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodDelete, "/v1/users/"+fixture.UserOne.ID.String(), nil)
			request.Header.Set("Authorization", "Bearer "+userOneAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusOK, apiResponse.StatusCode)

			user, _ := helper.GetUserByID(test.DB, fixture.UserOne.ID.String())
			assert.Nil(t, user)
		})

		t.Run("should return 401 error if access token is missing", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.UserOne)

			request := httptest.NewRequest(http.MethodDelete, "/v1/users/"+fixture.UserOne.ID.String(), nil)
			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusUnauthorized, apiResponse.StatusCode)
		})

		t.Run("should return 403 error if user is trying to delete another user", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.UserOne, fixture.UserTwo)

			userOneAccessToken, err := fixture.AccessToken(fixture.UserOne)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodDelete, "/v1/users/"+fixture.UserTwo.ID.String(), nil)
			request.Header.Set("Authorization", "Bearer "+userOneAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusForbidden, apiResponse.StatusCode)
		})

		t.Run("should return 200 if admin is trying to delete another user", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.UserOne, fixture.Admin)

			adminAccessToken, err := fixture.AccessToken(fixture.Admin)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodDelete, "/v1/users/"+fixture.UserOne.ID.String(), nil)
			request.Header.Set("Authorization", "Bearer "+adminAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusOK, apiResponse.StatusCode)
		})

		t.Run("should return 400 error if userId is not a valid postgres id", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.Admin)

			adminAccessToken, err := fixture.AccessToken(fixture.Admin)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodDelete, "/v1/users/invalidId", nil)
			request.Header.Set("Authorization", "Bearer "+adminAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusBadRequest, apiResponse.StatusCode)
		})

		t.Run("should return 404 error if user already is not found", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.Admin)

			adminAccessToken, err := fixture.AccessToken(fixture.Admin)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodDelete, "/v1/users/"+fixture.UserOne.ID.String(), nil)
			request.Header.Set("Authorization", "Bearer "+adminAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusNotFound, apiResponse.StatusCode)
		})
	})

	t.Run("PATCH /v1/users/:userId", func(t *testing.T) {
		t.Run("should return 200 and successfully update user if data is ok", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.UserOne)
			updateBody := validation.UpdateUser{
				Name:     "Golang",
				Email:    "golang@gmail.com",
				Password: "newPassword1",
			}

			userOneAccessToken, err := fixture.AccessToken(fixture.UserOne)
			assert.Nil(t, err)

			bodyJSON, err := json.Marshal(updateBody)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPatch, "/v1/users/"+fixture.UserOne.ID.String(), strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")
			request.Header.Set("Authorization", "Bearer "+userOneAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			bytes, err := io.ReadAll(apiResponse.Body)
			assert.Nil(t, err)

			responseBody := new(response.SuccessWithUser)

			err = json.Unmarshal(bytes, responseBody)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusOK, apiResponse.StatusCode)
			assert.Equal(t, "success", responseBody.Status)
			assert.NotContains(t, string(bytes), "password")
			assert.Equal(t, fixture.UserOne.ID, responseBody.User.ID)
			assert.Equal(t, updateBody.Name, responseBody.User.Name)
			assert.Equal(t, updateBody.Email, responseBody.User.Email)
			assert.Equal(t, "user", responseBody.User.Role)
			assert.Equal(t, false, responseBody.User.VerifiedEmail)

			user, err := helper.GetUserByID(test.DB, responseBody.User.ID.String())
			assert.Nil(t, err)

			assert.NotNil(t, user)
			assert.NotEqual(t, user.Password, updateBody.Password)
			assert.Equal(t, user.Name, updateBody.Name)
			assert.Equal(t, user.Email, updateBody.Email)
			assert.Equal(t, user.Role, "user")
		})

		t.Run("should return 401 error if access token is missing", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.UserOne)
			updateBody := validation.UpdateUser{
				Name: "Golang",
			}

			bodyJSON, err := json.Marshal(updateBody)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPatch, "/v1/users/"+fixture.UserOne.ID.String(), strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusUnauthorized, apiResponse.StatusCode)
		})

		t.Run("should return 403 if user is updating another user", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.UserOne, fixture.UserTwo)
			updateBody := validation.UpdateUser{
				Name: "Golang",
			}

			userOneAccessToken, err := fixture.AccessToken(fixture.UserOne)
			assert.Nil(t, err)

			bodyJSON, err := json.Marshal(updateBody)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPatch, "/v1/users/"+fixture.UserTwo.ID.String(), strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")
			request.Header.Set("Authorization", "Bearer "+userOneAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusForbidden, apiResponse.StatusCode)
		})

		t.Run("should return 200 and successfully update user if admin is updating another user", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.UserOne, fixture.Admin)
			updateBody := validation.UpdateUser{
				Name: "Golang",
			}

			adminAccessToken, err := fixture.AccessToken(fixture.Admin)
			assert.Nil(t, err)

			bodyJSON, err := json.Marshal(updateBody)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPatch, "/v1/users/"+fixture.UserOne.ID.String(), strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")
			request.Header.Set("Authorization", "Bearer "+adminAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusOK, apiResponse.StatusCode)
		})

		t.Run("should return 404 if admin is updating another user that is not found", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.Admin)
			updateBody := validation.UpdateUser{
				Name: "Golang",
			}

			adminAccessToken, err := fixture.AccessToken(fixture.Admin)
			assert.Nil(t, err)

			bodyJSON, err := json.Marshal(updateBody)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPatch, "/v1/users/"+fixture.UserOne.ID.String(), strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")
			request.Header.Set("Authorization", "Bearer "+adminAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusNotFound, apiResponse.StatusCode)
		})

		t.Run("should return 400 error if userId is not a valid postgres id", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.Admin)
			updateBody := validation.UpdateUser{
				Name: "Golang",
			}

			adminAccessToken, err := fixture.AccessToken(fixture.Admin)
			assert.Nil(t, err)

			bodyJSON, err := json.Marshal(updateBody)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPatch, "/v1/users/invalidId", strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")
			request.Header.Set("Authorization", "Bearer "+adminAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusBadRequest, apiResponse.StatusCode)
		})

		t.Run("should return 400 if email is invalid", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.UserOne)
			updateBody := validation.UpdateUser{
				Email: "invalidEmail",
			}

			userOneAccessToken, err := fixture.AccessToken(fixture.UserOne)
			assert.Nil(t, err)

			bodyJSON, err := json.Marshal(updateBody)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPatch, "/v1/users/"+fixture.UserOne.ID.String(), strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")
			request.Header.Set("Authorization", "Bearer "+userOneAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusBadRequest, apiResponse.StatusCode)
		})

		t.Run("should return 409 if email is already taken", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.UserOne, fixture.UserTwo)
			updateBody := validation.UpdateUser{
				Email: fixture.UserTwo.Email,
			}

			userOneAccessToken, err := fixture.AccessToken(fixture.UserOne)
			assert.Nil(t, err)

			bodyJSON, err := json.Marshal(updateBody)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPatch, "/v1/users/"+fixture.UserOne.ID.String(), strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")
			request.Header.Set("Authorization", "Bearer "+userOneAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusConflict, apiResponse.StatusCode)
		})

		t.Run("should not return 400 if email is my email", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.UserOne)
			updateBody := validation.UpdateUser{
				Email: fixture.UserOne.Email,
			}

			userOneAccessToken, err := fixture.AccessToken(fixture.UserOne)
			assert.Nil(t, err)

			bodyJSON, err := json.Marshal(updateBody)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPatch, "/v1/users/"+fixture.UserOne.ID.String(), strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")
			request.Header.Set("Authorization", "Bearer "+userOneAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusOK, apiResponse.StatusCode)
		})

		t.Run("should return 400 if password length is less than 8 characters", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.UserOne)
			updateBody := validation.UpdateUser{
				Password: "passwo1",
			}

			userOneAccessToken, err := fixture.AccessToken(fixture.UserOne)
			assert.Nil(t, err)

			bodyJSON, err := json.Marshal(updateBody)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPatch, "/v1/users/"+fixture.UserOne.ID.String(), strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")
			request.Header.Set("Authorization", "Bearer "+userOneAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusBadRequest, apiResponse.StatusCode)
		})

		t.Run("should return 400 if password does not contain both letters and numbers", func(t *testing.T) {
			helper.ClearAll(test.DB)
			helper.InsertUser(test.DB, fixture.UserOne)
			updateBody := validation.UpdateUser{
				Password: "password",
			}

			userOneAccessToken, err := fixture.AccessToken(fixture.UserOne)
			assert.Nil(t, err)

			bodyJSON, err := json.Marshal(updateBody)
			assert.Nil(t, err)

			request := httptest.NewRequest(http.MethodPatch, "/v1/users/"+fixture.UserOne.ID.String(), strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")
			request.Header.Set("Authorization", "Bearer "+userOneAccessToken)

			apiResponse, err := test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusBadRequest, apiResponse.StatusCode)

			updateBody.Password = "11111111"

			bodyJSON, err = json.Marshal(updateBody)
			assert.Nil(t, err)

			request = httptest.NewRequest(http.MethodPatch, "/v1/users/"+fixture.UserOne.ID.String(), strings.NewReader(string(bodyJSON)))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/json")
			request.Header.Set("Authorization", "Bearer "+userOneAccessToken)

			apiResponse, err = test.App.Test(request)
			assert.Nil(t, err)

			assert.Equal(t, http.StatusBadRequest, apiResponse.StatusCode)
		})
	})
}
