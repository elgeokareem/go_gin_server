package db

import (
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

	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " port=" + dbPort + " sslmode=" + dbSslMode
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// Uncomment this line to create the tables models
	// AutoGen(db)

	Service = &DBService{DB: db}
}
