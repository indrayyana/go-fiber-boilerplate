package fixture

import (
	"app/src/config"
	"app/src/model"
	"app/test/helper"
	"time"
)

var ExpiresAccessToken = time.Now().UTC().Add(time.Minute * time.Duration(config.JWTAccessExp))
var ExpiresRefreshToken = time.Now().UTC().Add(time.Hour * 24 * time.Duration(config.JWTRefreshExp))

func AccessToken(user *model.User) (string, error) {
	accessToken, err := helper.GenerateToken(user.ID.String(), ExpiresAccessToken, config.TokenTypeAccess)
	if err != nil {
		return accessToken, err
	}
	return accessToken, nil
}

func RefreshToken(user *model.User) (string, error) {
	refreshToken, err := helper.GenerateToken(user.ID.String(), ExpiresRefreshToken, config.TokenTypeRefresh)
	if err != nil {
		return refreshToken, err
	}
	return refreshToken, nil
}
