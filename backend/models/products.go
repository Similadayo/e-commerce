package models

import "github.com/jinzhu/gorm"

type Product struct {
	gorm.Model
	Name        string   `json:"name" gorm:"type:varchar(100);not null"`
	Description string   `json:"description" gorm:"type:text;not null"`
	Price       float32  `json:"price" gorm:"type:decimal(10,2);not null"`
	Quantity    int      `json:"quantity" gorm:"not null"`
	Images      []string `json:"images" gorm:"type:varchar(255)"`
	Category    Category `json:"category" gorm:"ForeignKey:CategoryID;AssociationForeignKey:ID"`
	CategoryID  uint     `json:"category_id" gorm:"not null"`
	Sizes       []string `json:"sizes" gorm:"type:varchar(255)"`
	Colors      []string `json:"colors" gorm:"type:varchar(255)"`
	Brand       *string  `json:"brand" gorm:"type:varchar(255)"`
}
