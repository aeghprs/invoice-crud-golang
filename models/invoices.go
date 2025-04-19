package models

import (
	"time"

	"gorm.io/gorm"
)

type Invoices struct {
	gorm.Model
	ID          uint `gorm:"primaryKey"`
	CustomersID uint `gorm:"index;not null"`
	Name        string
	Amount      float64
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Invoices) TableName() string {
	return "public.invoices"
}
