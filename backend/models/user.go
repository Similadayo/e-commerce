package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username    string     `json:"username" gorm:"type:varchar(100);not null"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Email       string     `json:"email" gorm:"type:varchar(100);not null"`
	Password    string     `json:"password" gorm:"type:varchar(100);not null"`
	Address     string     `json:"address" gorm:"type:varchar(100);not null"`
	PhoneNumber string     `json:"phone_number" gorm:"type:varchar(20);not null"`
	Role        string     `json:"role" gorm:"type:varchar(100);not null"`
	Suspension  Suspension `json:"suspension" gorm:"foreignkey:UserID"`
}
