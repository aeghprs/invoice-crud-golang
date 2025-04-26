package models

import "gorm.io/gorm"

type Customers struct {
	gorm.Model
	Name     string     `gorm:"not null" json:"name" binding:"required"`
	Email    string     `gorm:"not null;unique" json:"email" binding:"required,email"`
	Invoices []Invoices `gorm:"foreignKey:CustomersID;" json:"invoices"`
}

func (Customers) TableName() string {
	return "public.customers"
}
