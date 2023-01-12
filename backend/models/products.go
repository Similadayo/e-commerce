package models

import "github.com/jinzhu/gorm"

type Product struct {
	gorm.Model
	Name        string   `gorm:"type:varchar(100);not null"`
	Description string   `gorm:"type:text;not null"`
	Price       float32  `gorm:"type:decimal(10,2);not null"`
	Quantity    int      `gorm:"not null"`
	Images      []string `gorm:"type:varchar(255)"`
	Category    Category `gorm:"ForeignKey:CategoryID;AssociationForeignKey:ID"`
	CategoryID  uint     `gorm:"not null"`
	Sizes       []string `gorm:"type:varchar(255)"`
	Colors      []string `gorm:"type:varchar(255)"`
	Brand       *string  `gorm:"type:varchar(255)"`
}
