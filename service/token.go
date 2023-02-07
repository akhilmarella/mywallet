package service

import (
	"fmt"
	"mywallet/api"
	"mywallet/store"
	"mywallet/utils"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog/log"
)

var refreshSecretKey = []byte("nothing")

func RefreshToken(token string) (*api.TokenDetails, error) {
	refreshToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Error().Msg("error found in refresh token ")
			return nil, fmt.Errorf("There was an error ")
		}
		return refreshSecretKey, nil
	})

	if err != nil {
		log.Error().Err(err).Msg("refresh token invald ")
		return nil, err
	}

	_, ok := refreshToken.Claims.(jwt.Claims)
	if !ok && !refreshToken.Valid {
		log.Error().Any("refresh_token", refreshToken).
			Msg("refresh token expired")
		return nil, fmt.Errorf("refresh token expired")
	}

	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if ok && refreshToken.Valid {
		refreshID, ok := claims["refresh_id"].(string)
		if !ok {
			log.Error().Any("refresh_id", refreshID).
				Msg("refresh id is not found")
			return nil, fmt.Errorf("refresh id is not there")
		}

		email, ok := claims["email"].(string)
		if !ok {
			log.Error().Any("email", email).
				Msg("error in email ")
			return nil, err
		}

		role, ok := claims["role"].(string)
		if !ok {
			log.Error().Any("role", role).
				Msg("error in role")
			return nil, err
		}

		splitID := strings.Split(refreshID, "_")
		id := splitID[1]
		newID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			log.Error().Any("id", newID).
				Msg("error in id")
			return nil, err
		}

		deleted, delErr := store.DeleteRefreshID(refreshID)
		if delErr != nil || deleted == 0 {
			log.Error().Err(delErr).Any("deleted_ID", deleted).Any("refresh_id", refreshID).
				Msg("error in deleting refresh ID")
			return nil, delErr
		}

		newTokenDetails, err := utils.CreateToken(email, role, newID)
		if err != nil {
			log.Error().Err(err).Any("new_token_details", newTokenDetails).
				Msg("eror in  creating new token details")
			return nil, err
		}

		err = store.AddToken(newID, newTokenDetails)
		if err != nil {
			log.Error().Err(err).Any("id", id).Any("new_token_details", newTokenDetails).
				Msg("error in  Adding new token details")
			return nil, err
		}
		return newTokenDetails, nil

	}
	return nil, fmt.Errorf("error in refresh token")
}
