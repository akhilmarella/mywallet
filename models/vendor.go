package models

import "time"

type Vendor struct {
	ID          uint64 `gorm:"primary_key;auto_increment;not_null"`
	CompanyName string
	Name        string
	Email       string
	PhoneNumber string
	AddressID   uint64
	CreatedAt   time.Time
	LastUpdated time.Time
}

// type Address struct {
// 	ID          string
// 	UserType    string
// 	PersonID    string
// 	StreetNo    string
// 	Area        string
// 	Place       string
// 	District    string
// 	State       string
// 	PinCode     string
// 	CreatedAt   time.Time
// 	LastUpdated time.Time
// }

type AuthDetails struct {
	ID          uint64 `gorm:"primary_key;auto_increment;not_null"`
	UserType    string
	AccountID   uint64
	Email       string
	Password    string
	CreatedAt   time.Time
	LastUpdated time.Time
}
