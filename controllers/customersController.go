package controllers

import (
	config "main/dbConfig"
	"main/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type StatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func GetAllCustomers(c *gin.Context) {
	var users []models.Customers

	config.DB.Unscoped().Find(&users)

	c.JSON(200, gin.H{
		"message": "Customers fetched successfully",
		"result":  users,
	})
}

func CreateCustomers(c *gin.Context) {
	var user models.Customers

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	result := config.DB.Create(&user)

	if result.Error != nil {
		if pgErr, ok := result.Error.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				c.JSON(409, gin.H{
					"message": "Email already exists",
				})
				return
			}
		} else {
			c.JSON(400, gin.H{
				"message": "Unexpected error occurred",
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"message": "Customer created successfully",
		"result":  user,
	})
}

func UpdateCustomerStatus(c *gin.Context) {
	// 1. Parse & validate customer ID
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// 2. Parse & validate incoming status
	var req StatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if req.Status != "active" && req.Status != "inactive" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status must be 'active' or 'inactive'"})
		return
	}

	// 3. Build a GORM model stub for operations
	customer := models.Customers{Model: gorm.Model{ID: uint(id)}}

	if req.Status == "active" {
		// ---- RESTORE CUSTOMER & INVOICES ----

		// a) Restore customer (clear DeletedAt)
		resCust := config.DB.Unscoped().
			Model(&customer).
			Update("deleted_at", nil)
		if resCust.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to restore customer"})
			return
		}
		if resCust.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "Customer not found"})
			return
		}

		// b) Restore all invoices for this customer
		resInv := config.DB.Unscoped().
			Model(&models.Invoices{}).
			Where("customers_id = ?", customer.ID).
			Update("deleted_at", nil)
		if resInv.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Customer restored but failed to restore invoices"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Customer and invoices restored successfully"})
		return
	}

	// ---- SOFT-DELETE CUSTOMER & INVOICES ----

	// a) Soft-delete all invoices first
	delInv := config.DB.
		Where("customers_id = ?", customer.ID).
		Delete(&models.Invoices{})
	if delInv.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete invoices"})
		return
	}

	// b) Soft-delete the customer
	delCust := config.DB.
		Delete(&customer)
	if delCust.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete customer"})
		return
	}
	if delCust.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer and invoices deleted successfully"})
}
