package handlers

import (
	"mywallet/api"
	"mywallet/db"
	"mywallet/models"
	"mywallet/utils"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
)

func VendorRegister(c *gin.Context) {
	var req api.VendorRegisterRequest

	if err := c.BindJSON(&req); err != nil {
		log.Error().Err(err).Any("req", req).
			Any("action:", "handlers_vendor.go_VendorRegister").
			Msg("error in unmarshal")
		c.JSON(http.StatusBadRequest, gin.H{"message": "error found in unmarshaling"})
		return
	}

	if len(req.PhoneNumber) != 10 {
		log.Error().Any("phonenumber", req.PhoneNumber).
			Any("action:", "handlers_vendor.go_VendorRegister").
			Msg("phone number  must contain 10 numbers ")
		c.JSON(http.StatusBadRequest, gin.H{"message": "phone number is not valid"})
		return
	}

	ok := utils.ValidateEmail(req.Email)
	if !ok {
		log.Error().Any("email", req.Email).
			Any("action:", "handlers_vendor.go_VendorRegister").
			Msg("email is not valid")
		c.JSON(http.StatusBadRequest, gin.H{"message": "email is not valid "})
		return
	}

	if len(req.Password) <= 8 || len(req.Password) >= 14 {
		log.Error().Any("password", req.Password).
			Any("action:", "handlers_vendor.go_VendorRegister").
			Msg("password must contain between 8 to 14 characters")
		c.JSON(http.StatusBadRequest,
			gin.H{"message": "password is less than 8 characters (or) greater than 14 charaters"})
		return
	}

	if req.ConfirmPassword != req.Password {
		log.Error().Any("confirm_password", req.ConfirmPassword).Any("password", req.Password).
			Any("action:", "handlers_vendor.go_VendorRegister").
			Msg("password not equal to confirm password")
		c.JSON(http.StatusBadRequest, gin.H{"message": "password is not equal to confirm password"})
		return
	}

	valid := utils.ValidatePassword(req.Password)
	if !valid {
		log.Error().Any("password", req.Password).
			Any("action:", "handlers_vendor.go_VendorRegister").
			Msg("password is not valid")
		c.JSON(http.StatusBadRequest, gin.H{"message": "password is not valid"})
		return
	}

	var new models.Vendor
	new.CompanyName = req.CompanyName
	new.Name = req.Name
	new.Email = req.Email
	new.PhoneNumber = req.PhoneNumber
	vendorID, err := db.AddVendor(new)
	if err != nil {
		log.Error().Err(err).Any("vendor_details", new).
			Any("action:", "handlers_vendor.go_VendorRegister").
			Msg("error in adding vendor details")
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error in adding vendor details"})
		return
	}

	// after registered the vendor details,rewrite the vendorid in accountid
	var auth models.AuthDetails
	auth.UserType = "vendor"
	auth.AccountID = vendorID
	auth.Email = req.Email
	auth.Password = req.Password

	if error := db.AddAuth(auth); error != nil {
		log.Error().Err(error).
			Any("action:", "handlers_vendor.go_VendorRegister").
			Msg("error in adding authdetails for vendor")
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error in adding authdetails for vendor"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "created"})
}

func GetVendor(c *gin.Context) {
	userType := c.Writer.Header().Get("role")
	if userType != "vendor" {
		log.Error().Any("user_type", userType).Any("auth_id", c.Writer.Header().Get("auth_id")).
			Any("action", "handlers_vendor.go_GetVendor").Msg("unauthorized user")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized user"})
		return
	}

	id := c.Param("id")
	if id == "" {
		log.Error().Any("id", id).Any("action", "db_vendor.go_GetVendor").
			Msg("id not found")
		return
	}

	vendorID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Error().Err(err).Any("vendor_id", id).Any("action", "handlers_vendor.go_GetVendor").
			Msg("error in converting id string to int64 ")
		c.JSON(http.StatusNotFound, gin.H{"message": "error in converting id "})
		return
	}

	vendor, err := db.GetVendor(vendorID)
	if err != nil {
		log.Error().Err(err).Any("id", id).Any("action", "handlers_vendor.go_GetVendor").
			Msg("error  in id not found")
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	address, err := db.GetAddress(vendor.AddressID)
	if err != nil {
		log.Error().Err(err).Any("id", id).Any("action", "handlers_vendor.go_GetVendor").
			Msg("error  in id not found")
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	res := api.VendorList{
		ID:          vendorID,
		CompanyName: vendor.CompanyName,
		Name:        vendor.Name,
		Email:       vendor.Email,
		PhoneNumber: vendor.PhoneNumber,
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
		CreatedAt:   vendor.CreatedAt,
		LastUpdated: vendor.LastUpdated,
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": res})
}
