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

func GetVendor(id int64) (*models.Vendor, error) {
	var vendor models.Vendor
	tx := DB.Raw("select * from vendors where id = ? ", id).Scan(&vendor)
	if tx.Error != nil {
		log.Error().Err(tx.Error).Any("vendor_id", id).Any("action", "db_vendor.go_GetVendor").
			Msg("error in reading vendor details")
		return nil, tx.Error
	}

	return &vendor, nil
}

func UpdateVendor(vendor models.Vendor) (*models.Vendor, error) {
	var currentVendor models.Vendor

	tx := DB.Raw("select * from vendors where id = ? ", vendor.ID).Scan(&currentVendor)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		log.Error().Err(tx.Error).Any("id", vendor.ID).Any("action", "db_vendor.go_UpdateVendor").
			Msg("user not found")
		return nil, tx.Error
	}

	if vendor.CompanyName == "" {
		vendor.CompanyName = currentVendor.CompanyName
	}

	if vendor.Name == "" {
		vendor.Name = currentVendor.Name
	}

	if vendor.Email == "" {
		vendor.Email = currentVendor.Email
	}

	if vendor.PhoneNumber == "" {
		vendor.PhoneNumber = currentVendor.PhoneNumber
	}

	vendor.AddressID = currentVendor.AddressID
	vendor.LastUpdated = time.Now()

	em := DB.Save(&vendor)
	if em.Error != nil {
		log.Error().Err(em.Error).Any("vendor_details", vendor).
			Any("action", "db_vendor.go_UpdateVendor").Msg("error in updating vendor details")
		return nil, em.Error
	}

	return &models.Vendor{
		ID:          vendor.ID,
		CompanyName: vendor.CompanyName,
		Name:        vendor.Name,
		Email:       vendor.Email,
		PhoneNumber: vendor.PhoneNumber,
		AddressID:   vendor.AddressID,
		CreatedAt:   vendor.CreatedAt,
		LastUpdated: vendor.LastUpdated,
	}, nil
}
