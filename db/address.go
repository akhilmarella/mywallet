package db

import (
	"mywallet/models"
	"time"

	"github.com/rs/zerolog/log"
)

func AddAddress(address models.Address) error {

	address.CreatedAt = time.Now()

	tx := DB.Create(&address).Model(models.Address{}).Scan(&address)
	if tx.Error != nil {
		log.Error().Err(tx.Error).Any("address", address).
			Any("acion:", "db_address.go_AddAddress").
			Msg("error in inserting address")
		return tx.Error
	}

	type details struct {
		accountID int64
	}

	var det details

	em := DB.Raw("select account_id from auth_details where id = ? ", address.AuthID).Scan(&det)
	if em.Error != nil {
		log.Error().Err(em.Error).Any("auth_id", address.ID).
			Any("acion:", "db_address.go_AddAddress").
			Msg("error in storing id in address")
		return em.Error
	}

	if address.UserType == "customer" {
		var customer *models.Customer
		em := DB.Model(&customer).Where("id = ?", address.ID).UpdateColumn(models.Customer{AddressID: address.ID})
		//	em := DB.Raw("update customers set address_id = ? where id = ? ", address.ID, det.accountID)
		if em.Error != nil {
			log.Error().Err(em.Error).Any("id", address.ID).Any("auth_id", det.accountID).
				Any("acion:", "db_address.go_AddAddress").
				Msg("error in updating customers for adding address")
			return em.Error
		}
		return nil
	}

	if address.UserType == "vendor" {
		var vendor *models.Vendor
		em := DB.Model(&vendor).Where("id = ?", address.ID).UpdateColumn(models.Vendor{AddressID: address.ID})
		//	em := DB.Exec("update vendors set address_id = ? where id = ? ", address.ID, det.accountID)
		if em.Error != nil {
			log.Error().Err(em.Error).Any("id", address.ID).Any("auth_id", det.accountID).
				Any("acion:", "db_address.go_AddAddress").
				Msg("error in updating vendors for adding address")
			return em.Error
		}
		return nil
	}

	return nil
}

func GetAddress(id int64) (*models.Address, error) {
	var address models.Address
	tx := DB.Raw("select * from addresses where id = ? ", id).Scan(&address)
	if tx.Error != nil {
		log.Error().Err(tx.Error).Any("id", id).
			Any("action", "db_address.go_GetAddress").Msg("error in reading address details")
		return nil, tx.Error
	}
	return &address, nil
}
