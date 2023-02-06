package api

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