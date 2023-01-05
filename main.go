package main

import (
	"goGinServer/db"
	"goGinServer/routes"
)

func main() {
	db.DbConnection()

	routes.Routes()
}
