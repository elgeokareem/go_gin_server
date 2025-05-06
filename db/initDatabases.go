package db

import "goGinServer/db/models"

func InitDatabases() {
	// AutoMigrate will ONLY create tables, missing columns and missing indexes,
	// and WON'T change existing column's type or delete unused columns to protect your data.
	Service.DB.AutoMigrate(models.User{}, models.BinanceWallet{})
}
