package handlers

import (
	"mywallet/api"
	"mywallet/db"
	"mywallet/models"
	"mywallet/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func CustomerRegister(c *gin.Context) {
	var req api.CustomerRegisterReq

	if err := c.BindJSON(&req); err != nil {
		log.Error().Err(err).Any("req", req).
			Msg("error in unmarshal")
		c.JSON(http.StatusBadRequest, gin.H{"message": "error found in unmarshaling"})
		return
	}

	if len(req.PhoneNumber) != 10 {
		log.Error().Any("phonenumber", req.PhoneNumber).
			Msg("phone number  must contain 10 numbers ")
		c.JSON(http.StatusBadRequest, gin.H{"message": "phone number not valid "})
		return
	}

	ok := utils.ValidateEmail(req.Email)
	if !ok {
		log.Error().Any("email", req.Email).
			Msg("email is not valid")
		c.JSON(http.StatusBadRequest, gin.H{"message": "email is not valid"})
		return
	}

	if len(req.Password) <= 8 || len(req.Password) >= 14 {
		log.Error().Any("password", req.Password).Msg("password must contain between  8 to 14 characters ")
		c.JSON(http.StatusBadRequest,
			gin.H{"message": "password is lessthan than the 8 characters or greater than the 14 characters"})
		return
	}

	if req.ConfirmPassword != req.Password {
		log.Error().Any("confirm_password", req.ConfirmPassword).Any("password", req.Password).
			Msg("password not equal to confirm password")
		c.JSON(http.StatusBadRequest, gin.H{"message": "password is not equal to confirm password"})
		return
	}

	valid := utils.ValidatePassword(req.Password)
	if !valid {
		log.Error().Any("password", req.Password).
			Msg("password is not valid")
		c.JSON(http.StatusBadRequest, gin.H{"message": "password is not valid"})
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
		log.Error().Err(err).Any("customerID", customerID).
			Msg("error in adding customer details")
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error in adding customer details"})
		return
	}

	var auth models.AuthDetails
	auth.UserType = "customer"
	auth.AccountID = customerID
	auth.Email = req.Email
	auth.Password = req.Password

	if error := db.AddAuth(auth); error != nil {
		log.Error().Err(error).Msg("error in adding authDetails")
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error in adding authDetails"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "created"})
}
