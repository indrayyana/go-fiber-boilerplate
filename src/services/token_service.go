package services

import (
	"app/src/config"
	"app/src/model"
	res "app/src/response"
	"app/src/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TokenService interface {
	GenerateToken(userID string, expires time.Time, tokenType string) (string, error)
	SaveToken(c *fiber.Ctx, token string, userID string, expires time.Time) error
	DeleteToken(c *fiber.Ctx, userID string) error
	GetTokenByUserID(c *fiber.Ctx, tokenStr string) (*model.Token, error)
	GenerateAuthTokens(c *fiber.Ctx, user *model.User) (*res.Tokens, error)
}

type tokenService struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewTokenService(db *gorm.DB) TokenService {
	return &tokenService{
		Log: utils.Log,
		DB:  db,
	}
}

func (s *tokenService) GenerateToken(userID string, expires time.Time, tokenType string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"iat":  time.Now().Unix(),
		"exp":  expires.Unix(),
		"type": tokenType,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.JWTSecret))
}

func (s *tokenService) SaveToken(c *fiber.Ctx, token string, userID string, expires time.Time) error {
	if err := s.DeleteToken(c, userID); err != nil {
		return err
	}

	tokenDoc := &model.Token{
		Token:   token,
		UserID:  uuid.MustParse(userID),
		Expires: expires,
	}

	result := s.DB.WithContext(c.Context()).Create(tokenDoc)

	if result.Error != nil {
		s.Log.Errorf("Failed save token: %+v", result.Error)
	}

	return result.Error
}

func (s *tokenService) DeleteToken(c *fiber.Ctx, userID string) error {
	tokenDoc := new(model.Token)

	result := s.DB.WithContext(c.Context()).Where("user_id = ?", userID).Delete(tokenDoc)

	if result.Error != nil {
		s.Log.Errorf("Failed to delete token: %+v", result.Error)
	}

	return result.Error
}

func (s *tokenService) GetTokenByUserID(c *fiber.Ctx, tokenStr string) (*model.Token, error) {
	userID, err := utils.VerifyToken(tokenStr, config.JWTSecret, config.TokenTypeRefresh)
	if err != nil {
		return nil, err
	}

	tokenDoc := new(model.Token)

	result := s.DB.WithContext(c.Context()).
		Where("token = ? AND user_id = ?", tokenStr, userID).
		First(tokenDoc)

	if result.Error != nil {
		s.Log.Errorf("Failed get token by user id: %+v", err)
		return nil, result.Error
	}

	return tokenDoc, nil
}

func (s *tokenService) GenerateAuthTokens(c *fiber.Ctx, user *model.User) (*res.Tokens, error) {
	accessTokenExpires := time.Now().UTC().Add(time.Minute * time.Duration(config.JWTAccessExp))
	accessToken, err := s.GenerateToken(user.ID.String(), accessTokenExpires, "access")
	if err != nil {
		s.Log.Errorf("Failed generate token: %+v", err)
		return nil, err
	}

	refreshTokenExpires := time.Now().UTC().Add(time.Hour * 24 * time.Duration(config.JWTRefreshExp))
	refreshToken, err := s.GenerateToken(user.ID.String(), refreshTokenExpires, "refresh")
	if err != nil {
		s.Log.Errorf("Failed generate token: %+v", err)
		return nil, err
	}

	err = s.SaveToken(c, refreshToken, user.ID.String(), refreshTokenExpires)
	if err != nil {
		s.Log.Errorf("Failed save token: %+v", err)
		return nil, err
	}

	return &res.Tokens{
		Access: res.TokenExpires{
			Token:   accessToken,
			Expires: accessTokenExpires,
		},
		Refresh: res.TokenExpires{
			Token:   refreshToken,
			Expires: refreshTokenExpires,
		},
	}, nil
}
