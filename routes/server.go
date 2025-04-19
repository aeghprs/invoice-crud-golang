package routes

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
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

	r := mux.NewRouter()

	log.Printf("✅ Server started at port: %s \n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
