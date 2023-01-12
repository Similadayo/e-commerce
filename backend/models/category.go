package models

import "github.com/jinzhu/gorm"

type Category struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);not null"`
	Description string `gorm:"type:text;not null"`
}
