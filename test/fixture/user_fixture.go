package fixture

import (
	"app/src/model"

	"github.com/google/uuid"
)

var UserOne = &model.User{
	ID:            uuid.New(),
	Name:          "Test1",
	Email:         "test1@gmail.com",
	Password:      "password1",
	Role:          "user",
	VerifiedEmail: false,
}

var UserTwo = &model.User{
	ID:            uuid.New(),
	Name:          "Test2",
	Email:         "test2@gmail.com",
	Password:      "password1",
	Role:          "user",
	VerifiedEmail: false,
}

var Admin = &model.User{
	ID:            uuid.New(),
	Name:          "Admin",
	Email:         "admin@gmail.com",
	Password:      "password1",
	Role:          "admin",
	VerifiedEmail: false,
}
