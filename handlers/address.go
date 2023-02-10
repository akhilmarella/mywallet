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

func UpdateAddress(c *gin.Context) {
	var req api.AddressRegisterRequest
	if err := c.BindJSON(&req); err != nil {
		log.Error().Err(err).Any("req", req).
			Msg("error in unmarshal in address")
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in unmarshalling in address "})
		return
	}

	userType := c.Writer.Header().Get("role")
	if userType == "" {
		log.Error().Any("user_type", userType).
			Msg("empty header in usertype")
		c.JSON(http.StatusBadRequest, gin.H{"message": "empty found in header usertype"})
		return
	}

	userID_f := c.Writer.Header().Get("user_id")
	if userID_f == "" {
		log.Error().Any("user_id", userID_f).Msg("empty header in user id")
		c.JSON(http.StatusBadRequest, gin.H{"message": "empty found in header userid"})
		return
	}

	userID, err := strconv.ParseUint(userID_f, 10, 64)
	if err != nil {
		log.Error().Err(err).Any("user_id", userID).
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
	new.UserID = userID

	if err := db.AddAddress(new); err != nil {
		log.Error().Err(err).Any("address_details", new).
			Msg("error in adding address details")
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in adding new address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "created"})

}
