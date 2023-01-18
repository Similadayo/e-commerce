package models

import "github.com/jinzhu/gorm"

type Color struct {
	gorm.Model
	Name   string `json:"name" gorm:"type:varchar(100);not null"`
	Hex    string `json:"hex" gorm:"type:varchar(100);not null"`
	Active bool   `json:"active" gorm:"default:true"`
}
