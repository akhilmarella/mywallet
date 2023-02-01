package utils

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/rs/zerolog/log"
	"github.com/twinj/uuid"
)

var accessSecretKey = []byte("nooneneedtoknow")

var refreshSecretKey = []byte("nothing")

type TokenDetails struct {
	AccessToken   string
	AccessID      string
	RefreshToken  string
	RefreshID     string
	AccessExpiry  int64
	RefreshExpiry int64
}

func CreateToken(email, role string, id uint64) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AccessID = uuid.NewV4().String()
	td.AccessExpiry = time.Now().Add(time.Hour * 1).Unix()

	td.RefreshID = td.AccessID + "_" + strconv.Itoa(int(id))
	td.RefreshExpiry = time.Now().Add(time.Hour * 24 * 24).Unix()

	// creating access token
	accessClaims := jwt.MapClaims{}
	accessClaims["authorized"] = true
	accessClaims["email"] = email
	accessClaims["id"] = id
	accessClaims["role"] = role
	accessClaims["exp"] = td.AccessExpiry
	accessClaims["access_id"] = td.AccessID

	var err error
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	td.AccessToken, err = access.SignedString(accessSecretKey)
	if err != nil {
		log.Error().Err(err).Msg("error in access token")
		return nil, err
	}

	// creating refresh token
	refreshClaims := jwt.MapClaims{}
	refreshClaims["authorized"] = true
	refreshClaims["email"] = email
	refreshClaims["id"] = id
	refreshClaims["role"] = role
	refreshClaims["exp"] = td.RefreshExpiry
	refreshClaims["refresh_id"] = td.RefreshID

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	td.RefreshToken, err = refresh.SignedString(refreshSecretKey)
	if err != nil {
		log.Error().Err(err).Msg("error in refresh token")
		return nil, err
	}
	return td, nil
}
