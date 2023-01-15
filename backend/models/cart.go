package models

import "github.com/jinzhu/gorm"

type Cart struct {
	gorm.Model
	UserID    uint    `json:"user_id" gorm:"not null"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	Quantity  int     `json:"quantity" gorm:"type:integer;not null"`
	TotalCost float32 `json:"total_cost" gorm:"type:decimal(10,2);not null"`
}
