package models

type User struct {
	ID       uint   `gorm:"primary_key"`
	Email    string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
	Role     string `gorm:"not null"`
}
