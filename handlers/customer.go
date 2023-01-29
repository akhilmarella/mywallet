package handlers

import (
	"mywallet/db"
	"mywallet/models"
	"mywallet/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

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

func CustomerRegister(c *gin.Context) {
	var req CustomerRegisterReq

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

	var new models.Customer
	new.FirstName = req.FirstName
	new.LastName = req.LastName
	new.UserName = req.UserName
	new.Email = req.Email
	new.PhoneNumber = req.PhoneNumber
	new.DOB = req.DOB

	customerID, err := db.AddCustomer(new)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error in adding customer"})
		return
	}

	pass, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Error().Err(err).Any("password", req.Password).Msg("error in hashing password")
		return
	}

	var auth models.AuthDetails
	auth.UserType = "customer"
	auth.AccountID = customerID
	auth.Email = req.Email
	auth.Password = pass

	if error := db.AddAuth(auth); error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error in adding author"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "created"})
}
