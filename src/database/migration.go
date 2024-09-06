package database

import (
	"app/src/model"
	"app/src/utils"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	if err := db.AutoMigrate(
		&model.User{},
		&model.Token{},
	); err != nil {
		utils.Log.Errorf("Error migrating database: %+v", err)
	}
}
