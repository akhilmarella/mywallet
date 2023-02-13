package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var mySigningKey = []byte("nooneneedtoknow")

func IsAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Token") != "" {
			token, err := jwt.Parse(c.GetHeader("Token"), func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return mySigningKey, nil
			})

			if err != nil {
				log.Error().Err(err).Any("token", token).
					Any("action:", "middleware_token.go_IsAuthorized").
					Msg("error in invalid ttoken")
				c.Abort()
			}

			if token.Valid {
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					authID_f, ok := claims["auth_id"].(float64)

					if !ok {
						log.Error().Any("user_id", authID_f).
					Any("action:", "middleware_token.go_IsAuthorized").
							Msg("user id not found")
						return
					}

					authID := strconv.FormatFloat(authID_f, 'g', 1, 64)

					if claims["role"] == "vendor" {
						c.Writer.Header().Add("auth_id", authID)
						c.Writer.Header().Add("role", "vendor")
					} else if claims["role"] == "customer" {
						c.Writer.Header().Add("auth_id", authID)
						c.Writer.Header().Add("role", "customer")
					}
				}
				c.Next()
			} else {
				c.AbortWithError(http.StatusNotAcceptable, fmt.Errorf(" token is not valid"))
			}
		} else {
			log.Error().Msg("empty in header")
			c.IndentedJSON(400, gin.H{"message": "missing header in token"})
			c.AbortWithError(400, fmt.Errorf("no token in header"))
		}
	}
}
