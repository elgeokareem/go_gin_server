package models

import "time"

type BinanceWallet struct {
	ID         uint   `gorm:"primaryKey"`
	UserID     uint   `gorm:"foreignKey:ID;references:ID"`
	WalletType string `gorm:"type:varchar(16);not null;check:wallet_type IN ('spot','fund')"`
	Asset      string
	Free       float64
	Locked     float64
	Freeze     float64
	CreatedAt  time.Time
	UpdatedAt  time.Time
	User       User
}
