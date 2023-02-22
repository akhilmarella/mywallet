package handlers

import (
	"mywallet/api"
	"mywallet/db"
	"mywallet/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func AddAddress(c *gin.Context) {
	var req api.AddressRegisterRequest
	if err := c.BindJSON(&req); err != nil {
		log.Error().Err(err).Any("req", req).
			Any("action:", "handlers_addres.go_AddAddress").
			Msg("error in unmarshal in address")
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in unmarshalling in address "})
		return
	}

	userType := c.Writer.Header().Get("role")
	if userType == "" {
		log.Error().Any("user_type", userType).
			Any("action:", "handlers_addres.go_AddAddress").
			Msg("empty header in usertype")
		c.JSON(http.StatusBadRequest, gin.H{"message": "empty found in header usertype"})
		return
	}

	authID_s := c.Writer.Header().Get("auth_id")
	if authID_s == "" {
		log.Error().Any("user_id", authID_s).
			Any("action:", "handlers_addres.go_AddAddress").
			Msg("empty header in user id")
		c.JSON(http.StatusBadRequest, gin.H{"message": "empty found in header userid"})
		return
	}

	authID, err := strconv.ParseInt(authID_s, 10, 64)
	if err != nil {
		log.Error().Err(err).Any("auth_id", authID).
			Any("action:", "handlers_addres.go_AddAddress").
			Msg("error in converting  userid")
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in converting userid"})
		return
	}

	var new models.Address
	new.StreetNo = req.StreetNo
	new.Area = req.Area
	new.Place = req.Place
	new.District = req.District
	new.State = req.State
	new.PinCode = req.PinCode
	new.UserType = userType
	new.AuthID = authID

	if err := db.AddAddress(new); err != nil {
		log.Error().Err(err).Any("address_details", new).
			Any("action:", "handlers_addres.go_AddAddress").
			Msg("error in adding address details")
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in adding new address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "created"})
}

func UpdateAddress(c *gin.Context) {
	userType := c.Writer.Header().Get("role")
	if userType == "" {
		log.Error().Any("user_type", userType).
			Any("action", "handlers_address.go_UpdateAddress").Msg("empty header in usertype")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "empty header found in usertype"})
		return
	}

	authID_s := c.Writer.Header().Get("auth_id")
	if authID_s == "" {
		log.Error().Any("auth_id", authID_s).Any("action", "handlers_address.go_UpdateAddress").
			Msg("empty header  in authID")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "empty header in authID"})
		return
	}

	authID, err := strconv.ParseInt(authID_s, 10, 64)
	if err != nil {
		log.Error().Err(err).Any("authID", authID).Any("action", "handlers_address.go_UpdateAddress").
			Msg("error in converting authID")
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in converting authID from string to int"})
		return
	}

	id := c.Param("id")
	if id == "" {
		log.Error().Any("id", id).Any("action", "handlers_address.go_UpdateAddress").
			Msg("id not found")
		c.JSON(http.StatusBadRequest, gin.H{"message": "id not found"})
		return
	}

	addressID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Error().Err(err).Any("address_id", id).Any("action", "handlers_address.go_UpdateAddress").
			Msg("error in converting addressID")
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in converting addressID from string to int"})
		return
	}

	var change api.AddressUpdate
	if Err := c.BindJSON(&change); Err != nil {
		log.Error().Err(Err).Any("details", change).Any("action", "handlers_address.go_UpdateAddress").
			Msg("error in unmarshalling")
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in unmarshalling"})
		return
	}

	var address models.Address
	address.ID = addressID
	address.UserType = userType
	address.StreetNo = change.StreetNo
	address.AuthID = authID
	address.Area = change.Area
	address.Place = change.Place
	address.District = change.District
	address.State = change.State
	address.PinCode = change.PinCode

	res, err := db.UpdateAddress(address)
	if err != nil {
		log.Error().Err(err).Any("address", address).Any("action", "handlers_address.go_UpdateAddress").
			Msg("error in updating address")
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in updating address"})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": res})
}
