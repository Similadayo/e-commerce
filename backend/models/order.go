package models

type Order struct {
	ID        uint `gorm:"primary_key"`
	ProductID uint `gorm:"not null"`
	UserID    uint `gorm:"not null"`
	Quantity  uint `gorm:"not null"`
}
