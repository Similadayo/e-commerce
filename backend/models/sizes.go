package models

import "github.com/jinzhu/gorm"

type Size struct {
	gorm.Model
	Name   string `json:"name" gorm:"type:varchar(100);not null"`
	Active bool   `json:"active" gorm:"default:true"`
}
