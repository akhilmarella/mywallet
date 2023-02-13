package models

import "time"

type Address struct {
	ID          int64 `gorm:"primary_key;auto_increment;not_null"`
	UserType    string
	StreetNo    string
	AuthID      int64
	Area        string
	Place       string
	District    string
	State       string
	PinCode     int
	CreatedAt   time.Time
	LastUpdated time.Time
}
