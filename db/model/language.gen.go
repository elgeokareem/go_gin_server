// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameLanguage = "language"

// Language mapped from table <language>
type Language struct {
	LanguageID int32     `gorm:"column:language_id;primaryKey;autoIncrement:true" json:"language_id"`
	Name       string    `gorm:"column:name;not null" json:"name"`
	LastUpdate time.Time `gorm:"column:last_update;not null;default:now()" json:"last_update"`
}

// TableName Language's table name
func (*Language) TableName() string {
	return TableNameLanguage
}
