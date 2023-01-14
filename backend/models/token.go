package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type BlacklistToken struct {
	gorm.Model
	Token     string    `json:"token" gorm:"type:varchar(1000);not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
}
