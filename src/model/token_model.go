package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Token struct {
	ID        uuid.UUID `gorm:"primaryKey;not null"`
	Token     string    `gorm:"not null"`
	UserID    uuid.UUID `gorm:"not null"`
	Type      string    `gorm:"not null"`
	Expires   time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli"`
	UpdatedAt time.Time `gorm:"autoCreateTime:milli;autoUpdateTime:milli"`
	User      *User     `gorm:"foreignKey:user_id;references:id"`
}

func (token *Token) BeforeCreate(_ *gorm.DB) error {
	token.ID = uuid.New()
	return nil
}
