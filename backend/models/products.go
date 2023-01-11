package models

type Product struct {
	ID       uint   `gorm:"primary_key"`
	Name     string `gorm:"not null"`
	Price    uint   `gorm:"not null"`
	Quantity uint   `gorm:"not null"`
}
