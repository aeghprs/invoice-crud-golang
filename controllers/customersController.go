package controllers

import (
	config "main/dbConfig"
	"main/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
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
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req StatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Status != "active" && req.Status != "inactive" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid status value. It must be 'active' or 'inactive'.",
		})
		return
	}

	customer := models.Customers{}
	customer.ID = uint(id)

	if req.Status == "active" {
		result := config.DB.Unscoped().Model(&customer).Update("deleted_at", nil)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to restore customer"})
			return
		}
		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "Customer not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Customer restored successfully"})
		return
	}

	result := config.DB.Delete(&customer)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unexpected error occurred"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer terminated successfully"})
}
