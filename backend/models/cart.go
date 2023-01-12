package models

import "github.com/jinzhu/gorm"

type Cart struct {
	gorm.Model
	UserID    uint    `gorm:"not null"`
	ProductID uint    `gorm:"not null"`
	Quantity  int     `gorm:"type:integer;not null"`
	TotalCost float32 `gorm:"type:decimal(10,2);not null"`
}
