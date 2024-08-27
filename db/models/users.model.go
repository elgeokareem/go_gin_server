package models

import "time"

type User struct {
	ID           uint `gorm:"primaryKey"`
	Email        string
	Password     string
	BinanceToken string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
