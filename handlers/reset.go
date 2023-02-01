package handlers

import (
	"mywallet/db"
	"mywallet/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ResetAuthPassword struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func ResetPassword(c *gin.Context) {
	var req ResetAuthPassword
	if err := c.BindJSON(&req); err != nil {
		log.Error().Err(err).Any("req", req).
			Msg("error in unmarshaling")
		c.JSON(http.StatusBadRequest, gin.H{"message": "error found in unmarshaling"})
		return
	}

	role := c.Writer.Header().Get("role")
	if role == "" {
		return
	}

	if len(req.Password) <= 8 || len(req.Password) >= 14 {
		log.Error().Any("password", req.Password).Msg("password must contain between 8 to 14 characters")
		c.JSON(http.StatusBadRequest,
			gin.H{"message": "password is less than the 8 characters or greatwr than the 14 characters "})
		return
	}

	if req.ConfirmPassword != req.Password {
		log.Error().Any("confirm_password", req.ConfirmPassword).Any("password", req.Password).
			Msg("password not equal to confirm password")
		c.JSON(http.StatusBadRequest, gin.H{"message": "password is not equal to confirm password"})
		return
	}

	valid := utils.ValidatePassword(req.Password)
	if !valid {
		log.Error().Any("password", req.Password).
			Msg("password is not valid")
		c.JSON(http.StatusBadRequest, gin.H{"message": "password is not valid"})
		return
	}

	if err := db.ChangePassword(req.Email, req.Password, role); err != nil {
		log.Error().Err(err).Any("email", req.Email).Any("password", req.Password).
			Any("confirm_passwors", req.ConfirmPassword).Msg("error in changig password")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error in changing password"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "password is changed"})
}
