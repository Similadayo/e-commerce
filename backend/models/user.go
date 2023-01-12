package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username    string `gorm:"type:varchar(100);not null"`
	Email       string `gorm:"type:varchar(100);not null"`
	Password    []byte `gorm:"type:varchar(100);not null"`
	Address     string `gorm:"type:varchar(100);not null"`
	PhoneNumber uint64 `gorm:"not null"`
	Role        Role   `gorm:"foreignkey:RoleID"`
	RoleID      uint
}

func (u *User) BeforeSave() (err error) {
	u.Password, err = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	return
}
