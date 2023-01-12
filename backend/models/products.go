package models

import "github.com/jinzhu/gorm"

type Product struct {
    gorm.Model
    Name string `gorm:"type:varchar(100);not null"`
    Description string `gorm:"type:text;not null"`
    Price float32 `gorm:"type:decimal(10,2);not null"`
    Quantity int `gorm:"type:integer;not null"`
    Category string `gorm:"type:varchar(100);not null"`
    ImageURL string `gorm:"type:varchar(255);not null"`
}
