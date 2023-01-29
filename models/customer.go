package models

import "time"

type Customer struct {
	ID          uint64 `gorm:"primary_key;auto_increment;not_null"`
	FirstName   string
	LastName    string
	UserName    string
	Email       string
	PhoneNumber string
	DOB         string
	AddressID   uint64
	CreatedAt   time.Time
	LastUpdated time.Time
}
