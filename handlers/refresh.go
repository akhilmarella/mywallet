package handlers

import (
	"fmt"
	"mywallet/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func RefreshToken(c *gin.Context) {
	refreshToken := c.GetHeader("Refresh-Token")
	if refreshToken == "" {
		log.Error().Any("refresh_token", refreshToken).Msg("couldn't get token from header")
		c.JSON(http.StatusBadRequest, gin.H{"message": "couldn't get token"})
		return
	}

	newTokenDetails, err := service.RefreshToken(refreshToken)
	if err != nil {
		log.Error().Err(err).Any("token_details_new", newTokenDetails).
			Msg("error in new token details")
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in token details"})
		return
	}

	if newTokenDetails == nil {

		if newTokenDetails.AccessExpiry == 0 {
			fmt.Println("empty access expiry found", newTokenDetails.AccessExpiry)
			c.JSON(http.StatusBadRequest, gin.H{"message": "empty access expiry found"})
			return
		}
		if newTokenDetails.AccessID == "" {
			fmt.Println("empty access id found", newTokenDetails.AccessID)
			c.JSON(http.StatusBadRequest, gin.H{"message": "empty access id found"})
			return
		}
		if newTokenDetails.AccessToken == "" {
			fmt.Println("empty access token found", newTokenDetails.AccessToken)
			c.JSON(http.StatusBadRequest, gin.H{"message": "empty access token found"})
			return
		}
		if newTokenDetails.RefreshExpiry == 0 {
			fmt.Println("empty refresh expiry found", newTokenDetails.RefreshExpiry)
			c.JSON(http.StatusBadRequest, gin.H{"message": "empty refresh expiry found"})
			return
		}
		if newTokenDetails.RefreshID == "" {
			fmt.Println("empty refresh id found", newTokenDetails.RefreshID)
			c.JSON(http.StatusBadRequest, gin.H{"message": "empty refresh id found"})
			return
		}
		if newTokenDetails.RefreshToken == "" {
			fmt.Println("empty refresh token found", newTokenDetails.RefreshToken)
			c.JSON(http.StatusBadRequest, gin.H{"message": "empty refresh token found"})
			return
		}
		fmt.Println("error in new token details", newTokenDetails)
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in new token details"})
		return
	}
	response := map[string]string{
		"access_token":  newTokenDetails.AccessToken,
		"refresh_token": newTokenDetails.RefreshToken,
	}

	c.JSON(http.StatusCreated, response)
}
