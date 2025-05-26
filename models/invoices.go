package models

import "gorm.io/gorm"

type Invoices struct {
	gorm.Model
	CustomersID uint      `gorm:"index;not null" binding:"required" json:"customers_id"`
	Customer    Customers `gorm:"foreignKey:CustomersID" json:"customer"`
	Name        string    `json:"name"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"`
}

func (Invoices) TableName() string {
	return "public.invoices"
}
