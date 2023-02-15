package handlers

import (
	"mywallet/api"
	"mywallet/db"
	"mywallet/models"
	"mywallet/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func CustomerRegister(c *gin.Context) {
	var req api.CustomerRegisterReq

	if err := c.BindJSON(&req); err != nil {
		log.Error().Err(err).Any("req", req).
			Any("action:", "handlers_customer.go_CustomerRegister").
			Msg("error in unmarshal")
		c.JSON(http.StatusBadRequest, gin.H{"message": "error found in unmarshaling"})
		return
	}

	if len(req.PhoneNumber) != 10 {
		log.Error().Any("phonenumber", req.PhoneNumber).
			Any("action:", "handlers_customer.go_CustomerRegister").
			Msg("phone number  must contain 10 numbers ")
		c.JSON(http.StatusBadRequest, gin.H{"message": "phone number not valid "})
		return
	}

	ok := utils.ValidateEmail(req.Email)
	if !ok {
		log.Error().Any("email", req.Email).
			Any("action:", "handlers_customer.go_CustomerRegister").
			Msg("email is not valid")
		c.JSON(http.StatusBadRequest, gin.H{"message": "email is not valid"})
		return
	}

	if len(req.Password) <= 8 || len(req.Password) >= 14 {
		log.Error().Any("password", req.Password).
			Any("action:", "handlers_customer.go_CustomerRegister").
			Msg("password must contain between  8 to 14 characters ")
		c.JSON(http.StatusBadRequest,
			gin.H{"message": "password is lessthan than the 8 characters or greater than the 14 characters"})
		return
	}

	if req.ConfirmPassword != req.Password {
		log.Error().Any("confirm_password", req.ConfirmPassword).Any("password", req.Password).
			Any("action:", "handlers_customer.go_CustomerRegister").
			Msg("password not equal to confirm password")
		c.JSON(http.StatusBadRequest, gin.H{"message": "password is not equal to confirm password"})
		return
	}

	valid := utils.ValidatePassword(req.Password)
	if !valid {
		log.Error().Any("password", req.Password).
			Any("action:", "handlers_customer.go_CustomerRegister").
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
			Any("action:", "handlers_customer.go_CustomerRegister").
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
		log.Error().Err(error).
			Any("action:", "handlers_customer.go_CustomerRegister").
			Msg("error in adding authDetails")
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error in adding authDetails"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "created"})
}

func GetCustomer(c *gin.Context) {
	userType := c.Writer.Header().Get("role")
	if userType != "customer" {
		log.Error().Any("user_type", userType).Any("auth_id", c.Writer.Header().Get("auth_id")).
			Any("action", "handlers_customer.go_GetCustomer").Msg("unauthorized user")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized user"})
		return
	}

	id := c.Param("id")
	if id == "" {
		log.Error().Any("id", id).Any("action", "db_customer.go_GetCustomer").
			Msg("id not found")
		return
	}

	customerID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Error().Err(err).Any("customer_id", id).Any("action", "handlers_customer.go_GetCustomer").
			Msg("error in converting id string to int64 ")
		c.JSON(http.StatusNotFound, gin.H{"message": "error in converting id "})
		return
	}

	customer, err := db.GetCustomer(customerID)
	if err != nil {
		log.Error().Err(err).Any("id", id).Any("action", "handlers_customer.go_GetCustomer").
			Msg("error  in id not found")
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	address, err := db.GetAddress(customer.AddressID)
	if err != nil {
		log.Error().Err(err).Any("id", id).Any("action", "handlers_customer.go_GetCustomer").
			Msg("error  in id not found")
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	res := api.CustomerList{
		Id:          customerID,
		FirstName:   customer.FirstName,
		LastName:    customer.LastName,
		UserName:    customer.UserName,
		PhoneNumber: customer.PhoneNumber,
		DOB:         customer.DOB,
		Email:       customer.Email,
		AddressID: api.Address{
			UserType:    address.UserType,
			StreetNo:    address.StreetNo,
			Area:        address.Area,
			Place:       address.Place,
			District:    address.District,
			State:       address.State,
			PinCode:     address.PinCode,
			CreatedAt:   address.CreatedAt,
			LastUpdated: address.LastUpdated,
		},
		CreatedAt:   address.CreatedAt,
		LastUpdated: address.LastUpdated,
	}
	
	c.IndentedJSON(http.StatusCreated, gin.H{"message": res})
}
