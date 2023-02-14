package api

import "time"

type AddressRegisterRequest struct {
	StreetNo string `json:"street_no"`
	Area     string `json:"area"`
	Place    string `json:"place"`
	District string `json:"district"`
	State    string `json:"state"`
	PinCode  int    `json:"pin_code"`
}

type Address struct {
	UserType    string    `json:"user_type"`
	StreetNo    string    `json:"street_no"`
	Area        string    `json:"area"`
	Place       string    `json:"place"`
	District    string    `json:"district"`
	State       string    `json:"state"`
	PinCode     int       `json:"pin_code"`
	CreatedAt   time.Time `json:"created_at"`
	LastUpdated time.Time `json:"last_updated"`
}
