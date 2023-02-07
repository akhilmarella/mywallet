package handlers

import (
	"mywallet/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func Logout(c *gin.Context) {
	accessToken := c.GetHeader("Access-Token")
	if accessToken == "" {
		log.Error().Any("Access_token", accessToken).Msg("couldn't get access token from header")
		c.JSON(http.StatusBadRequest, gin.H{"message": "couldn't get access token"})
		return
	}

	deleteToken, err := service.DeleteToken(accessToken)
	if err != nil {
		log.Error().Err(err).Any("delete_token", deleteToken).Any("access_token", accessToken).
			Msg("error in deleting token ")
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in deleting token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "token successfully deleted"})
}
