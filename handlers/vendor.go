package handlers

import (
	"mywallet/db"
	"mywallet/models"
	"mywallet/utils"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
)

type VendorRegisterRequest struct {
	CompanyName     string `json:"company_name"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	PhoneNumber     string `json:"phone_number"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func VendorRegister(c *gin.Context) {
	var req VendorRegisterRequest

	if err := c.BindJSON(&req); err != nil {
		log.Error().Err(err).Any("req", req).
			Msg("error in unmarshal")
		return
	}

	if len(req.PhoneNumber) != 10 {
		log.Error().Any("phonenumber", req.PhoneNumber).
			Msg("phone number  must contain 10 numbers ")
		return
	}

	ok := utils.ValidateEmail(req.Email)
	if !ok {
		log.Error().Any("email", req.Email).
			Msg("email is not valid")
		return
	}

	if len(req.Password) <= 8 && len(req.Password) >= 12 {
		log.Error().Any("password", req.Password).Msg("password must contain between 8 to 12 characters")
		return
	}

	if req.ConfirmPassword != req.Password {
		log.Error().Any("confirm_password", req.ConfirmPassword).Any("password", req.Password).
			Msg("password not equal to confirm password")
		return
	}

	valid := utils.ValidatePassword(req.Password)
	if !valid {
		log.Error().Any("password", req.Password).
			Msg("password is not valid")
		return
	}

	var new models.Vendor
	new.CompanyName = req.CompanyName
	new.Name = req.Name
	new.Email = req.Email
	new.PhoneNumber = req.PhoneNumber
	vendorID, err := db.AddVendor(new)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error in adding vendor"})
		return
	}

	

	// after registered the vendor details,rewrite the vendorid in accountid
	var auth models.AuthDetails
	auth.UserType = "vendor"
	auth.AccountID = vendorID
	auth.Email = req.Email
	auth.Password = req.Password

	if error := db.AddAuth(auth); error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error in adding author"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "created"})
}
