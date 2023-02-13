package db

import (
	"fmt"
	"mywallet/models"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
)

func AddVendor(vendor models.Vendor) (int64, error) {
	vendor.CreatedAt = time.Now()
	var check models.Check

	em := DB.Raw("select email from vendors where email = ?", vendor.Email).Scan(&check)
	if em.Error != nil && em.Error != gorm.ErrRecordNotFound {
		log.Error().Err(em.Error).Any("email", check).
			Any("action:", "db_vendor.go_AddVendor").
			Msg("error in email")
		return 0, em.Error
	}

	if check.Email != "" {
		log.Error().Any("email", check).
		Any("action:", "db_vendor.go_AddVendor").
			Msg("email already registered")
		return 0, fmt.Errorf("email is already exist :%v", check.Email)
	}

	tx := DB.Create(&vendor).Model(models.Vendor{}).Scan(&vendor)
	if tx.Error != nil {
		log.Error().Err(tx.Error).Any("vendor", vendor).
		Any("action:", "db_vendor.go_AddVendor").
			Msg("error in creating user")
		return 0, tx.Error
	}

	return vendor.ID, nil
}
