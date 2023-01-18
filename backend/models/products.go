package models

import "github.com/jinzhu/gorm"

type Product struct {
	gorm.Model
	Name        string   `json:"name" gorm:"type:varchar(100);not null"`
	Description string   `json:"description" gorm:"type:text;not null"`
	Price       float32  `json:"price" gorm:"type:decimal(10,2);not null"`
	Quantity    int      `json:"quantity" gorm:"not null"`
	Category    Category `json:"category" gorm:"many2many:product_categories;ForeignKey:CategoryID;AssociationForeignKey:ID"`
	CategoryID  uint     `json:"category_id" gorm:"not null"`
	Sizes       []Size   `json:"sizes" gorm:"many2many:product_sizes;ForeignKey:SizeID;AssociationForeignKey:ID"`
	Colors      []Color  `json:"colors" gorm:"many2many:product_colors;ForeignKey:ColorID;AssociationForeignKey:ID"`
	Brand       *string  `json:"brand" gorm:"type:varchar(255)"`
}
