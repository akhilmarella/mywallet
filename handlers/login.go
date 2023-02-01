package handlers

import (
	"fmt"
	"mywallet/db"
	"mywallet/store"
	"mywallet/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var req LoginReq
	if err := c.BindJSON(&req); err != nil {
		log.Error().Err(err).Any("req", req).
			Msg("eror in unmarshaling")
		c.JSON(http.StatusBadRequest, gin.H{"msg": "error is found in unmarshaling"})
		return
	}

	role := c.GetHeader("User-Type")

	var msg string
	var Error error

	switch role {
	case "":
		msg = "empty value in User-Type header"
		Error = fmt.Errorf(msg)
	case "vendor", "customer":
		//	msg = ""
		//	Error = nil
		break
	default:
		msg = fmt.Sprintf("invalid User-Type: %v", role)
		Error = fmt.Errorf(msg)
	}
	if Error != nil {
		log.Error().Err(Error).Msg(msg)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid User-Type in header"})
		return
	}

	if req.Password == "" {
		log.Error().Any("password", req.Password).Msg("password is empty")
		c.JSON(http.StatusBadRequest, gin.H{"message": "password is empty"})
		return
	}

	if req.Email == "" {
		log.Error().Any("email", req.Email).Msg("email is empty")
		c.JSON(http.StatusBadRequest, gin.H{"message": "email is empty"})
		return
	}

	id, err := db.ReadUser(req.Email, req.Password, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in login"})
		return
	}
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"messsage": "wrong password"})
		return
	}

	tokenDetails, err := utils.CreateToken(req.Email, role, id)
	if err != nil {
		c.JSON(http.StatusFailedDependency, gin.H{"message": "error in token"})
		return
	}
	err = store.AddToken(id, tokenDetails)
	if err != nil {
		log.Error().Err(err).Any("id", id).Any("token_detais", tokenDetails).
			Msg("error in token details")
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in token details"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"token": tokenDetails})
}
