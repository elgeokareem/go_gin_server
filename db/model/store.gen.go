// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameStore = "store"

// Store mapped from table <store>
type Store struct {
	StoreID        int32     `gorm:"column:store_id;primaryKey;autoIncrement:true" json:"store_id"`
	ManagerStaffID int16     `gorm:"column:manager_staff_id;not null" json:"manager_staff_id"`
	AddressID      int16     `gorm:"column:address_id;not null" json:"address_id"`
	LastUpdate     time.Time `gorm:"column:last_update;not null;default:now()" json:"last_update"`
}

// TableName Store's table name
func (*Store) TableName() string {
	return TableNameStore
}
