package db

import (
	"fmt"
	"goGinServer/db/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DbConnection() *gorm.DB {
	dsn := "host=localhost user=postgres password=12345678 dbname=DVDRental port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// AutoGen(db)

	// make a query and print the result
	var actor model.Actor
	db.First(&actor, 1) // find product with integer primary key
	fmt.Println("------------------")
	fmt.Println(actor)

	return db
}
