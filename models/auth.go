package models

import "time"

type AuthDetails struct {
	ID          int64 `gorm:"primary_key;auto_increment;not_null"`
	UserType    string
	AccountID   int64
	Email       string
	Password    string
	CreatedAt   time.Time
	LastUpdated time.Time
}
