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

	router.GET("/getAllCustomers", controllers.GetAllCustomers)
	router.POST("/createCustomers", controllers.CreateCustomers)
	router.PUT("/updateCustomers/:id", controllers.UpdateCustomerStatus)

	router.Run(":" + port)
}
