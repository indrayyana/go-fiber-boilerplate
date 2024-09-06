package fixture

import (
	"app/src/config"
	"app/src/model"
	"app/test/helper"
	"time"
)

func AccessToken(user *model.User) (string, error) {
	expires := time.Now().Add(time.Minute * time.Duration(config.JWTAccessExp))
	accessToken, err := helper.GenerateToken(user.ID.String(), user.Role, expires, config.TokenTypeAccess)
	if err != nil {
		return accessToken, err
	}
	return accessToken, nil
}
