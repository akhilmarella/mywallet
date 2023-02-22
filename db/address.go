package db

import (
	"fmt"
	"mywallet/models"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
)

func AddAddress(address models.Address) error {
	var currentAddress models.Address
	add := DB.Raw("select *  from addresses where auth_id = ? ", address.AuthID).Scan(&currentAddress)
	if add.Error != nil && add.Error != gorm.ErrRecordNotFound {
		log.Error().Err(add.Error).Any("auth_id", address.AuthID).Any("action", "db_address.go_AddAddress").
			Msg("error in fetching details")
		return add.Error
	}

	if currentAddress.AuthID == address.AuthID {
		log.Error().Any("address", currentAddress).Any("action ", "db_address.go_AddAddress").
			Msg("email adready registered")
		return fmt.Errorf("email is already registered ")
	}
	address.CreatedAt = time.Now()

	tx := DB.Create(&address).Model(models.Address{}).Scan(&address)
	if tx.Error != nil {
		log.Error().Err(tx.Error).Any("address", address).
			Any("acion:", "db_address.go_AddAddress").
			Msg("error in inserting address")
		return tx.Error
	}

	type details struct {
		AccountID int64
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
		em := DB.Model(&customer).Where("id = ?", address.ID).
			UpdateColumn(models.Customer{AddressID: det.AccountID})
		//	em := DB.Raw("update customers set address_id = ? where id = ? ", address.ID, det.accountID)
		if em.Error != nil {
			log.Error().Err(em.Error).Any("id", address.ID).Any("auth_id", det.AccountID).
				Any("acion:", "db_address.go_AddAddress").
				Msg("error in updating customers for adding address")
			return em.Error
		}
		return nil
	}

	if address.UserType == "vendor" {
		var vendor *models.Vendor
		em := DB.Model(&vendor).Where("id = ?", address.ID).
			UpdateColumn(models.Vendor{AddressID: det.AccountID})
		//	em := DB.Exec("update vendors set address_id = ? where id = ? ", address.ID, det.accountID)
		if em.Error != nil {
			log.Error().Err(em.Error).Any("id", address.ID).Any("auth_id", det.AccountID).
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

func UpdateAddress(address models.Address) (*models.Address, error) {
	var currentAddress models.Address
	tx := DB.Raw("select  * from addresses where id = ? ", address.ID).Scan(&currentAddress)
	if tx.Error != nil {
		log.Error().Err(tx.Error).Any("id", address.ID).Any("action", "db_address.go_UpdateAddress").
			Msg("user not found")
		return nil, tx.Error
	}

	if address.StreetNo == "" {
		address.StreetNo = currentAddress.StreetNo
	}

	if address.Area == "" {
		address.Area = currentAddress.Area
	}

	if address.Place == "" {
		address.Place = currentAddress.Place
	}

	if address.District == "" {
		address.District = currentAddress.District
	}

	if address.State == "" {
		address.State = currentAddress.State
	}

	if address.PinCode == 0 {
		address.PinCode = currentAddress.PinCode
	}

	address.UserType = currentAddress.UserType
	address.AuthID = currentAddress.AuthID
	address.LastUpdated = time.Now()

	add := DB.Save(&address)
	if add.Error != nil {
		log.Error().Err(add.Error).Any("address", address).Any("action", "db_address.go_UpdateAddress").
			Msg("error in updating address")
		return nil, add.Error
	}

	return &models.Address{
		ID:          address.ID,
		UserType:    address.UserType,
		StreetNo:    address.StreetNo,
		AuthID:      address.AuthID,
		Area:        address.Area,
		Place:       address.Place,
		District:    address.District,
		State:       address.State,
		PinCode:     address.PinCode,
		CreatedAt:   address.CreatedAt,
		LastUpdated: address.LastUpdated,
	}, nil
}

