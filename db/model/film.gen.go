// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameFilm = "film"

// Film mapped from table <film>
type Film struct {
	FilmID          int32     `gorm:"column:film_id;primaryKey;autoIncrement:true" json:"film_id"`
	Title           string    `gorm:"column:title;not null" json:"title"`
	Description     string    `gorm:"column:description" json:"description"`
	ReleaseYear     int32     `gorm:"column:release_year" json:"release_year"`
	LanguageID      int16     `gorm:"column:language_id;not null" json:"language_id"`
	RentalDuration  int16     `gorm:"column:rental_duration;not null;default:3" json:"rental_duration"`
	RentalRate      float64   `gorm:"column:rental_rate;not null;default:4.99" json:"rental_rate"`
	Length          int16     `gorm:"column:length" json:"length"`
	ReplacementCost float64   `gorm:"column:replacement_cost;not null;default:19.99" json:"replacement_cost"`
	Rating          string    `gorm:"column:rating;default:G" json:"rating"`
	LastUpdate      time.Time `gorm:"column:last_update;not null;default:now()" json:"last_update"`
	SpecialFeatures string    `gorm:"column:special_features" json:"special_features"`
	Fulltext        string    `gorm:"column:fulltext;not null" json:"fulltext"`
}

// TableName Film's table name
func (*Film) TableName() string {
	return TableNameFilm
}
