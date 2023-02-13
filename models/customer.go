package models

import "time"

type Customer struct {
	ID          int64 `gorm:"primary_key;auto_increment;not_null"`
	FirstName   string
	LastName    string
	UserName    string
	Email       string
	PhoneNumber string
	DOB         string
	AddressID   int64
	CreatedAt   time.Time
	LastUpdated time.Time
}
