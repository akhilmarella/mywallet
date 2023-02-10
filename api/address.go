package api

type AddressRegisterRequest struct {
	StreetNo string `json:"street_no"`
	Area     string `json:"area"`
	Place    string `json:"place"`
	District string `json:"district"`
	State    string `json:"state"`
	PinCode  int    `json:"pin_code"`
}
