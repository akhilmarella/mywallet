package api

type VendorRegisterRequest struct {
	CompanyName     string `json:"company_name"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	PhoneNumber     string `json:"phone_number"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}
