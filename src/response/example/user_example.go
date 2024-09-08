package example

import "github.com/google/uuid"

type User struct {
	ID            uuid.UUID `json:"id" example:"e088d183-9eea-4a11-8d5d-74d7ec91bdf5"`
	Name          string    `json:"name" example:"fake name"`
	Email         string    `json:"email" example:"fake@example.com"`
	Role          string    `json:"role" example:"user"`
	VerifiedEmail bool      `json:"verified_email" example:"false"`
}

type GoogleUser struct {
	ID            uuid.UUID `json:"id" example:"e088d183-9eea-4a11-8d5d-74d7ec91bdf5"`
	Name          string    `json:"name" example:"fake name"`
	Email         string    `json:"email" example:"fake@example.com"`
	Role          string    `json:"role" example:"user"`
	VerifiedEmail bool      `json:"verified_email" example:"true"`
}
