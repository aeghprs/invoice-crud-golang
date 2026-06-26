package routes

import (
	"log"
	"main/controllers"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var port string

func initENVVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port = os.Getenv("SERVER_PORT")
}

func InitServer() {
	initENVVariables()

	router := gin.Default()

	log.Printf("✅ Server started at port: %s \n", port)
	{
		v1 := router.Group("/customers")
		v1.GET("/all", controllers.GetAllCustomers)
		v1.POST("/create", controllers.CreateCustomers)
		v1.PUT("/update/:id", controllers.UpdateCustomerStatus)
	}

	{
		v1 := router.Group("/invoice")
		v1.GET("/all", controllers.GetAllInvoices)
		v1.POST("/create", controllers.CreateInvoice)
	}

	router.Run(":" + port)
}
