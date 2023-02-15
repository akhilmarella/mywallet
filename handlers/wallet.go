package handlers

import (
	"fmt"
	"mywallet/api"
	"mywallet/db"
	"mywallet/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func AddWalletRegister(c *gin.Context) {
	userType := c.Writer.Header().Get("role")
	if userType == "" {
		log.Error().Any("user_type", userType).Any("action", "handlers_wallet.go_AddWalletRegister").
			Msg("usertype is nil")
		c.JSON(http.StatusBadRequest, gin.H{"message": "usertype is nil"})
		return
	}

	authID := c.Writer.Header().Get("auth_id")
	fmt.Println("authid", authID)
	if authID == "" {
		log.Error().Any("auth_id", authID).Any("action", "handlers_wallet.go_AddWalletRegister").
			Msg("auth id is nil")
		c.JSON(http.StatusBadRequest, gin.H{"message": "authid is nil"})
		return
	}

	authid, err := strconv.ParseInt(authID, 10, 64)
	if err != nil {
		log.Error().Err(err).Any("id", authid).Any("action", "handlers_wallet.go_AddWalletRegister").
			Msg("error in converting id ")
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in converting authid "})
		return
	}

	var req api.WalletRequest
	if err := c.BindJSON(&req); err != nil {
		log.Error().Err(err).Any("details", req).Any("action", "handlers_wallet.go_AddWalletRegister").
			Msg("error in unmarshling")
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in unmarshaling"})
		return
	}
	//convert authid into id(if customer then id is customer.id)

	userid, err := db.GetAccountID(authid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in fetching accountID from auth"})
		return
	}

	if userid == 0 {
		log.Error().Any("auth_id", authID).Any("action", "handlers_wallet.go_AddWalletRegister").
			Msg("invalid userid")
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid userid"})
		return
	}

	var wallet models.Wallet
	wallet.UserType = userType
	wallet.TotalMoney = req.Money
	wallet.UserID = userid

	err = db.AddWallet(wallet)
	if err != nil {
		log.Error().Err(err).Any("wallet", wallet).Any("action", "handlers_wallet.go_AddWalletRegister").
			Msg("error in adding wallet ")
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in adding wallet"})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"msg": "created"})

}
