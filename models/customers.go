package models

import "gorm.io/gorm"

type Customers struct {
	gorm.Model
	ID       uint `gorm:"primaryKey"`
	Name     string
	Email    string
	Invoices []Invoices `gorm:"foreignKey:CustomersID;"`
}

func (Customers) TableName() string {
	return "public.customers"
}
