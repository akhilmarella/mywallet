package models

import "time"

type Vendor struct {
	ID          int64 `gorm:"primary_key;auto_increment;not_null"`
	CompanyName string
	Name        string
	Email       string
	PhoneNumber string
	AddressID   int64
	CreatedAt   time.Time
	LastUpdated time.Time
}



type Check struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}
