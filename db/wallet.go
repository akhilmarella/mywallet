package db

import (
	"mywallet/models"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
)

func AddWallet(wallet models.Wallet) error {

	var currentWallet models.Wallet

	em := DB.Raw("select * from wallets where user_id = ? and user_type = ? ", wallet.UserID,
		wallet.UserType).Scan(&currentWallet)
	if em.Error != nil && em.Error != gorm.ErrRecordNotFound {
		log.Error().Err(em.Error).Any("user_id", wallet.UserID).Any("user_type", wallet.UserType).
			Any("action", "db_wallet.go_AddWallet").Msg("error in fetching details from wallet")
		return em.Error
	}

	if currentWallet.UserID == wallet.UserID {
		money := wallet.TotalMoney + currentWallet.TotalMoney
		// money = new_money + previous_money

		// tx := DB.Model(&currentWallet).Where("id = ? ", wallet.Id).
		// 	Updates(models.Wallet{TotalMoney: money, LastUpdated: time.Now()})
		// tx := DB.Exec("update wallets set total_money = ? , last_updated = ? where id = ? ",
		// 	money, time.Now(), wallet.Id)
		currentWallet.TotalMoney = money
		currentWallet.LastUpdated = time.Now()

		tx := DB.Save(&currentWallet)
		if tx.Error != nil {
			log.Error().Err(tx.Error).Any("previous_money", currentWallet.UserID).
				Any("action", "db_wallet.go_AddWallet").
				Msg("error in adding money  ")
			return tx.Error
		}
		
		return nil
	}

	wallet.CreatedAt = time.Now()
	tx := DB.Create(&wallet).Model(models.Wallet{}).Scan(&wallet)
	if tx.Error != nil {
		log.Error().Err(tx.Error).Any("wallet", wallet).Any("action", "db_wallet.go_AddWallet").
			Msg("error in inserting wallet")
		return tx.Error
	}

	return nil
}
