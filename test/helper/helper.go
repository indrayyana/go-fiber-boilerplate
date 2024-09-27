package helper

import (
	"app/src/config"
	"app/src/model"
	"app/src/utils"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func ClearAll(db *gorm.DB) {
	ClearToken(db)
	ClearUsers(db)
}

func ClearUsers(db *gorm.DB) {
	err := db.Where("id is not null").Delete(&model.User{}).Error
	if err != nil {
		logrus.Fatalf("Failed clear user data : %+v", err)
	}
}

func ClearToken(db *gorm.DB) {
	err := db.Where("id is not null").Delete(&model.Token{}).Error
	if err != nil {
		logrus.Fatalf("Failed clear user token : %+v", err)
	}
}

func CreateUser(db *gorm.DB, email, password, name string) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		logrus.Errorf("Failed hashed password : %+v", err)
	}

	user := &model.User{
		Email:    email,
		Password: hashedPassword,
		Name:     name,
	}

	err = db.Create(user).Error
	if err != nil {
		logrus.Errorf("Failed create user : %+v", err)
	}
}

func InsertUser(db *gorm.DB, users ...*model.User) {
	now := time.Now()

	for i, user := range users {
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			logrus.Errorf("Failed to hash password: %+v", err)
			continue
		}
		user.Password = hashedPassword
		user.CreatedAt = now.Add(time.Duration(i) * time.Second)

		if errDB := db.Create(user).Error; errDB != nil {
			logrus.Errorf("Failed to create user: %+v", errDB)
		}
	}
}

func SaveToken(db *gorm.DB, token, userID, tokenType string, expires time.Time) error {
	if err := DeleteToken(db, tokenType, userID); err != nil {
		return err
	}

	tokenDoc := &model.Token{
		Token:   token,
		UserID:  uuid.MustParse(userID),
		Type:    tokenType,
		Expires: expires,
	}

	result := db.Create(tokenDoc)

	return result.Error
}

func DeleteToken(db *gorm.DB, tokenType, userID string) error {
	tokenDoc := new(model.Token)

	result := db.Where("type = ? AND user_id = ?", tokenType, userID).Delete(tokenDoc)

	return result.Error
}

func GenerateToken(
	userID string, expires time.Time, tokenType string,
) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"iat":  time.Now().Unix(),
		"exp":  expires.Unix(),
		"type": tokenType,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.JWTSecret))
}

func GenerateInvalidToken(
	userID string, expires time.Time, tokenType string,
) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"iat":  time.Now().Unix(),
		"exp":  expires.Unix(),
		"type": tokenType,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte("invalidSecret"))
}

func GetTokenByUserID(db *gorm.DB, tokenStr string) (*model.Token, error) {
	userID, err := utils.VerifyToken(tokenStr, config.JWTSecret, config.TokenTypeRefresh)
	if err != nil {
		return nil, err
	}

	tokenDoc := new(model.Token)
	result := db.Where("token = ? AND user_id = ?", tokenStr, userID).
		First(tokenDoc)

	if result.Error != nil {
		return nil, result.Error
	}

	return tokenDoc, nil
}

func GetTokenByType(db *gorm.DB, userID string, tokenType string) (*model.Token, error) {
	tokenDoc := new(model.Token)
	result := db.Where("type = ? AND user_id = ?", tokenType, userID).
		First(tokenDoc)

	if result.Error != nil {
		return nil, result.Error
	}

	return tokenDoc, nil
}

func GetUserByID(db *gorm.DB, id string) (*model.User, error) {
	user := new(model.User)

	result := db.First(user, "id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	if result.Error != nil {
		logrus.Errorf("Failed get user by id: %+v", result.Error)
	}

	return user, result.Error
}
