package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID `gorm:"primaryKey;not null" json:"id"`
	Name          string    `gorm:"not null" json:"name"`
	Email         string    `gorm:"uniqueIndex;not null" json:"email"`
	Password      string    `gorm:"not null" json:"-"`
	Role          string    `gorm:"default:user;not null" json:"role"`
	VerifiedEmail bool      `gorm:"default:false;not null" json:"verified_email"`
	CreatedAt     time.Time `gorm:"autoCreateTime:milli" json:"-"`
	UpdatedAt     time.Time `gorm:"autoCreateTime:milli;autoUpdateTime:milli" json:"-"`
	Token         []Token   `gorm:"foreignKey:user_id;references:id" json:"-"`
}

func (user *User) BeforeCreate(_ *gorm.DB) error {
	user.ID = uuid.New() // Generate UUID before create
	return nil
}
