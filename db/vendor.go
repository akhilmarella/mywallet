package db

import (
	"fmt"
	"mywallet/models"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
)

type RegCheck struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

func AddVendor(vendor models.Vendor) (uint64, error) {
	vendor.CreatedAt = time.Now()
	var check RegCheck

	em := DB.Raw("select email from vendors where email = ?", vendor.Email).Scan(&check)
	if em.Error != nil && em.Error != gorm.ErrRecordNotFound {
		log.Error().Err(em.Error).Any("email", check).
			Msg("error in email")
		return 0, em.Error
	}

	if check.Email != "" {
		log.Error().Any("email", check).
			Msg("email already registered")
		return 0, fmt.Errorf("email is already exist :%v", check.Email)
	}

	tx := DB.Create(&vendor).Model(models.Vendor{}).Scan(&vendor)
	if tx.Error != nil {
		log.Error().Err(tx.Error).Any("vendor", vendor).
			Msg("error in creating user")
		return 0, tx.Error
	}

	return vendor.ID, nil
}

func AddAuth(auth models.AuthDetails) error {
	var check RegCheck
	em := DB.Raw("select email from auth_details where email = ?", auth.Email).Scan(&check)
	if em.Error != nil && em.Error != gorm.ErrRecordNotFound {
		log.Error().Err(em.Error).Any("email", check).
			Msg("error in email")
		return em.Error
	}

	if check.Email != "" {
		log.Error().Any("email", check).
			Msg("already registered")
		return fmt.Errorf("email is already exist :%v", check.Email)
	}

	auth.CreatedAt = time.Now()
	tx := DB.Create(&auth).Model(models.AuthDetails{})
	if tx.Error != nil {
		log.Error().Err(tx.Error).Any("auth", auth).
			Msg("error in creating auth")
		return tx.Error
	}

	return nil
}

