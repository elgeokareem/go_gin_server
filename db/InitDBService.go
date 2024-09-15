package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBService struct {
	DB *gorm.DB
}

var Service *DBService

func InitDBService() {
	dbHost := os.Getenv("HOST_DB")
	dbUser := os.Getenv("USER_DB")
	dbPassword := os.Getenv("PASSWORD_DB")
	dbName := os.Getenv("DATABASE_DB")
	dbPort := os.Getenv("PORT_DB")
	dbSslMode := os.Getenv("SSL_MODE_DB")

	// Initial connection string without specifying the database
	dsnWithoutDB := "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " port=" + dbPort + " sslmode=" + dbSslMode
	fmt.Println(dsnWithoutDB)

	// Connection string with the specified database
	dsnWithDB := dsnWithoutDB + " dbname=" + dbName
	db, err := gorm.Open(postgres.Open(dsnWithDB), &gorm.Config{})

	if err != nil {
		fmt.Println("Failed to connect to database")

		temporalDb, err := gorm.Open(postgres.Open(dsnWithoutDB+" dbname=postgres"), &gorm.Config{})

		if err != nil {
			panic(err)
		}

		createDB := fmt.Sprintf("CREATE DATABASE %s;", dbName)

		err = temporalDb.Exec(createDB).Error
		if err != nil {
			panic(err)
		}

		fmt.Println("Database created successfully!")

		// Reconnect to the newly created database
		db, err = gorm.Open(postgres.Open(dsnWithDB), &gorm.Config{})
		if err != nil {
			panic("failed to connect to the newly created database")
		}

		Service = &DBService{DB: db}

		return
	}

	// Uncomment this line to create the tables models
	// AutoGen(db)

	Service = &DBService{DB: db}
}
