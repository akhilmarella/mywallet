package api

import "time"

type VendorRegisterRequest struct {
	CompanyName     string `json:"company_name"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	PhoneNumber     string `json:"phone_number"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type VendorList struct {
	ID          int64     `json:"id"`
	CompanyName string    `json:"company_name"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	AddressID   Address   `json:"address_id"`
	CreatedAt   time.Time `json:"created_at"`
	LastUpdated time.Time `json:"last_updated"`
}
