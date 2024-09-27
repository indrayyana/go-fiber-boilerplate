package model_test

import (
	"app/src/model"
	"app/src/validation"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var validate = validation.Validator()

func TestUserModel(t *testing.T) {
	t.Run("Create user validation", func(t *testing.T) {
		var newUser = validation.CreateUser{
			Name:     "John Doe",
			Email:    "johndoe@gmail.com",
			Password: "password1",
			Role:     "user",
		}

		t.Run("should correctly validate a valid user", func(t *testing.T) {
			err := validate.Struct(newUser)
			assert.NoError(t, err)
		})

		t.Run("should throw a validation error if email is invalid", func(t *testing.T) {
			newUser.Email = "invalidEmail"
			err := validate.Struct(newUser)
			assert.Error(t, err)
		})

		t.Run("should throw a validation error if password length is less than 8 characters", func(t *testing.T) {
			newUser.Password = "passwo1"
			err := validate.Struct(newUser)
			assert.Error(t, err)
		})

		t.Run("should throw a validation error if password does not contain numbers", func(t *testing.T) {
			newUser.Password = "password"
			err := validate.Struct(newUser)
			assert.Error(t, err)
		})

		t.Run("should throw a validation error if password does not contain letters", func(t *testing.T) {
			newUser.Password = "11111111"
			err := validate.Struct(newUser)
			assert.Error(t, err)
		})

		t.Run("should throw a validation error if role is unknown", func(t *testing.T) {
			newUser.Role = "invalid"
			err := validate.Struct(newUser)
			assert.Error(t, err)
		})
	})

	t.Run("Update user validation", func(t *testing.T) {
		var updateUser = validation.UpdateUser{
			Name:     "John Doe",
			Email:    "johndoe@gmail.com",
			Password: "password1",
		}

		t.Run("should correctly validate a valid user", func(t *testing.T) {
			err := validate.Struct(updateUser)
			assert.NoError(t, err)
		})

		t.Run("should throw a validation error if email is invalid", func(t *testing.T) {
			updateUser.Email = "invalidEmail"
			err := validate.Struct(updateUser)
			assert.Error(t, err)
		})

		t.Run("should throw a validation error if password length is less than 8 characters", func(t *testing.T) {
			updateUser.Password = "passwo1"
			err := validate.Struct(updateUser)
			assert.Error(t, err)
		})

		t.Run("should throw a validation error if password does not contain numbers", func(t *testing.T) {
			updateUser.Password = "password"
			err := validate.Struct(updateUser)
			assert.Error(t, err)
		})

		t.Run("should throw a validation error if password does not contain letters", func(t *testing.T) {
			updateUser.Password = "11111111"
			err := validate.Struct(updateUser)
			assert.Error(t, err)
		})
	})

	t.Run("Update user password validation", func(t *testing.T) {
		var newPassword = validation.UpdatePassOrVerify{
			Password: "password1",
		}

		t.Run("should correctly validate a valid user password", func(t *testing.T) {
			err := validate.Struct(newPassword)
			assert.NoError(t, err)
		})

		t.Run("should throw a validation error if password length is less than 8 characters", func(t *testing.T) {
			newPassword.Password = "passwo1"
			err := validate.Struct(newPassword)
			assert.Error(t, err)
		})

		t.Run("should throw a validation error if password does not contain numbers", func(t *testing.T) {
			newPassword.Password = "password"
			err := validate.Struct(newPassword)
			assert.Error(t, err)
		})

		t.Run("should throw a validation error if password does not contain letters", func(t *testing.T) {
			newPassword.Password = "11111111"
			err := validate.Struct(newPassword)
			assert.Error(t, err)
		})
	})

	t.Run("User toJSON()", func(t *testing.T) {
		t.Run("should not return user password when toJSON is called", func(t *testing.T) {
			user := &model.User{
				Name:     "John Doe",
				Email:    "johndoe@gmail.com",
				Password: "password1",
				Role:     "user",
			}

			bytes, _ := json.Marshal(user)
			assert.NotContains(t, string(bytes), "password")
		})
	})
}
