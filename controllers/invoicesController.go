package controllers

import (
	config "main/dbConfig"
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateInvoiceInput struct {
	CustomersID uint    `json:"customers_id" binding:"required"`
	Name        string  `json:"name" binding:"required,min=3"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Status      string  `json:"status" binding:"required,oneof=paid unpaid pending"`
}

func GetAllInvoices(c *gin.Context) {
	var invoices []models.Invoices

	// Preload the associated Customers record
	config.DB.
		// Unscoped().
		Preload("Customer").
		Find(&invoices)

	c.JSON(200, gin.H{
		"message": "Invoices fetched successfully",
		"result":  invoices,
	})
}

func CreateInvoice(c *gin.Context) {
	var input CreateInvoiceInput

	// Validate input only against input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	// Check if customer exists
	var customer models.Customers
	if err := config.DB.First(&customer, input.CustomersID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"error":   "Customer not found",
		})
		return
	}

	// Create the invoice
	invoice := models.Invoices{
		CustomersID: input.CustomersID,
		Name:        input.Name,
		Amount:      input.Amount,
		Status:      input.Status,
	}

	if err := config.DB.Create(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create invoice",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Invoice created successfully",
		"invoice": invoice,
	})
}
