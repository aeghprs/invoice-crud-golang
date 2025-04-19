package main

import (
	config "main/dbConfig"
	routes "main/routes"
)

func main() {
	config.InitDB()

	routes.InitServer()
}
