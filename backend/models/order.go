package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Order struct {
	gorm.Model
	UserID    uint      `gorm:"not null"`
	ProductID uint      `gorm:"not null"`
	Quantity  int       `gorm:"type:integer;not null"`
	Address   string    `gorm:"type:varchar(255);not null"`
	TotalCost float32   `gorm:"type:decimal(10,2);not null"`
	Status    string    `gorm:"type:varchar(100);not null"`
	Date      time.Time `gorm:"not null"`
}
