package main

import (
	"goGinServer/db"
	"goGinServer/routes"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.DbConnection()

	routes.Routes()
}
