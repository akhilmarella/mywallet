package service

import (
	"fmt"
	"mywallet/api"
	"mywallet/store"

	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog/log"
)

var accessSecretKey = []byte("nooneneedtoknow")

func DeleteToken(token string) (*api.TokenDetails, error) {
	accessToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Error().Msg("error found in access token ")
			return nil, fmt.Errorf("There was an error ")
		}
		return accessSecretKey, nil
	})
	if err != nil {
		log.Error().Err(err).Msg("access token invald ")
		return nil, err
	}

	_, ok := accessToken.Claims.(jwt.Claims)
	if !ok && !accessToken.Valid {
		log.Error().Any("access_token", accessToken).
			Msg("access token expired")
		return nil, fmt.Errorf("error in access token expired")
	}

	claims, ok := accessToken.Claims.(jwt.MapClaims)
	if ok && accessToken.Valid {
		accessID, ok := claims["access_id"].(string)
		if !ok {
			log.Error().Any("access_id", accessID).
				Msg("access id not found")
			return nil, fmt.Errorf("access id not there")
		}

		userID, ok := claims["id"].(float64)
		if !ok {
			log.Error().Any("user_id", userID).
				Msg("user id not found")
			return nil, fmt.Errorf("user id not there")
		}

		deleteAccess, deleteRefresh, err := store.DeleteTokens(accessID, uint64(userID))
		if err != nil || deleteAccess == 0 || deleteRefresh == 0 {
			log.Error().Err(err).Any("delete_Access_id", deleteAccess).
				Any("delete_refresh_id", deleteRefresh).Any("id", userID).
				Any("access_id", accessID).
				Msg("error in deleting access id")
			return nil, err
		}
		return nil, err
	}
	return nil, fmt.Errorf("error in deleting access token")
}
