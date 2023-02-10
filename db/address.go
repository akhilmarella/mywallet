package db

import (
	"fmt"
	"mywallet/models"
	"time"

	"github.com/rs/zerolog/log"
)

func AddAddress(address models.Address) error {

	address.CreatedAt = time.Now()
	tx := DB.Exec(`insert into addresses (user_type, street_no, user_id, area, place, district, state,
		pin_code, created_at, last_updated)
	values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, address.UserType, address.StreetNo, address.UserID,
		address.Area, address.Place, address.District, address.State,
		address.PinCode, time.Now(), time.Now())
	if tx.Error != nil {
		log.Error().Err(tx.Error).Any("address", address).
			Msg("error in inserting address")
		return tx.Error
	}

	fmt.Println("address_id:", address.ID)
	type details struct {
		accountID uint64
	}

	var det details

	em := DB.Raw("select account_id from auth_details where id = ? ", address.ID).Scan(&det)
	if em.Error != nil {
		log.Error().Err(em.Error).Any("auth_id", address.ID).
			Any("action", "db_address.go_AddAddress").Msg("error in storing id in address")
		return em.Error
	}

	if address.UserType == "customer" {
		em := DB.Raw("update customers set address_id = ? where id = ? ", address.ID, address.UserID)
		if em.Error != nil {
			log.Error().Err(em.Error).Any("id", address.ID).Any("user_id", address.UserID).
				Msg("error in updating customers for adding address")
			return em.Error
		}
		return nil
	}

	if address.UserType == "vendor" {
		em := DB.Raw("update vendors set address_id = ? where id = ? ", address.ID, address.UserID)
		if em.Error != nil {
			log.Error().Err(em.Error).Any("id", address.ID).Any("user_id", address.UserID).
				Msg("error in updating vendors for adding address")
			return em.Error
		}
		return nil
	}

	return nil
}
