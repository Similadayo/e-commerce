package models

import (
	"github.com/jinzhu/gorm"
)

type Order struct {
	gorm.Model
	UserId          uint        `gorm:"not null"`
	OrderNumber     string      `gorm:"type:varchar(100);not null"`
	ShippingAddress string      `gorm:"type:varchar(100);not null"`
	BillingAddress  string      `gorm:"type:varchar(100);not null"`
	Items           []OrderItem `gorm:"foreignkey:OrderId"`
	TotalCost       float32     `gorm:"not null"`
	PaymentMethod   string      `gorm:"type:varchar(100);not null"`
	PaymentStatus   bool        `gorm:"not null"`
	DeliveryStatus  bool        `gorm:"not null"`
}

type OrderItem struct {
	gorm.Model
	OrderId   uint    `gorm:"not null"`
	ProductId uint    `gorm:"not null"`
	Quantity  uint    `gorm:"not null"`
	SubTotal  float32 `gorm:"not null"`
}
