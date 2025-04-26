package models

import (
	"gorm.io/gorm"
)

type Invoices struct {
	gorm.Model
	CustomersID uint `gorm:"index;not null"`
	Name        string
	Amount      float64
	Status      string
}

func (Invoices) TableName() string {
	return "public.invoices"
}
