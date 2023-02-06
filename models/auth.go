package models

import "time"

type AuthDetails struct {
	ID          uint64 `gorm:"primary_key;auto_increment;not_null"`
	UserType    string
	AccountID   uint64
	Email       string
	Password    string
	CreatedAt   time.Time
	LastUpdated time.Time
}
