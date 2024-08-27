package models

import "time"

type BinanceSpotWallet struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"foreignKey:ID;references:ID"`
	Asset     string
	Free      float64
	Locked    float64
	Freeze    float64
	CreatedAt time.Time
	UpdatedAt time.Time
	User      User
}

type BinanceFundWallet struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"foreignKey:ID;references:ID"`
	Asset     string
	Free      float64
	Locked    float64
	CreatedAt time.Time
	UpdatedAt time.Time
	User      User
}
