package models

import "time"

type Wallet struct {
	Id          int64 `gorm:"primary_key;auto_increment;not_null"`
	UserID      int64
	UserType    string
	TotalMoney  float32
	CreatedAt   time.Time
	LastUpdated time.Time
}
