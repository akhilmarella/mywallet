package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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
				log.Printf("invalid token")
				c.Abort()
			}
			
			if token.Valid {
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					if claims["role"] == "vendor" {
						c.Writer.Header().Add("role", "vendor")
					} else if claims["role"] == "customer" {
						c.Writer.Header().Add("role", "customer")
					}
				}
				c.Next()
			} else {
				c.AbortWithError(http.StatusNotAcceptable, fmt.Errorf(" token is not valid"))
			}
		} else {
			log.Println("No token")
			c.IndentedJSON(400, gin.H{"message": "missing header in token"})
			c.AbortWithError(400, fmt.Errorf("no token in header"))
		}
	}
}
