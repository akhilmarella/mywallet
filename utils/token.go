package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog/log"
)

var mySigningKey = []byte("nooneneedtoknow")

func GenJWT(email, role string, id uint64) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["email"] = email
	claims["id"] = id
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		log.Error().Err(err).Msg("something went wrong in token")
		return "", err
	}
	return tokenString, nil
}
