package api

import "time"

type CustomerRegisterReq struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	UserName        string `json:"user_name"`
	PhoneNumber     string `json:"phone_number"`
	DOB             string `json:"dob"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type CustomerList struct {
	Id          int64     `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	UserName    string    `json:"user_name"`
	PhoneNumber string    `json:"phone_number"`
	DOB         string    `json:"dob"`
	Email       string    `json:"email"`
	AddressID   Address   `json:"address"`
	CreatedAt   time.Time `json:"created_at"`
	LastUpdated time.Time `json:"last_updated"`
}
