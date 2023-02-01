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
	auth.CreatedAt = time.Now()
	tx := DB.Create(&auth).Model(models.AuthDetails{})
	if tx.Error != nil {
		log.Error().Err(tx.Error).Any("auth", auth).
			Msg("error in creating auth")
		return tx.Error
	}

	return nil
}

func ReadUser(email, password, role string) (uint64, error) {
	var res models.AuthDetails
	tx := DB.Raw("select * from auth_details where email = ? ", email).Scan(&res)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		log.Error().Err(tx.Error).Any("email", email).
			Msg("error in fetching details")
		return 0, tx.Error
	}

	if email != res.Email {
		log.Error().Any("old_email", res.Email).Any("new_email", email).
			Msg("email not found")
		return 0, fmt.Errorf("email is incorrect %v ", email)
	}

	if role != res.UserType {
		log.Error().Any("usertype", res.UserType).Any("email", res.Email).
			Any("id", res.ID).Msg("usertype is incorrect")
		return 0, fmt.Errorf("usertype is wrong")
	}

	if password != res.Password {
		log.Error().Any("password", res.Password).Any("req_password:", password).
			Msg("password is incorrect")
		return 0, fmt.Errorf("password is wrong")
	}

	return res.ID, nil
}

func ChangePassword(email, password, role string) error {
	var res models.AuthDetails
	tx := DB.Raw("select * from auth_details where email = ? and user_type = ?", email, role).Scan(&res)
	if tx.Error != nil {
		log.Error().Err(tx.Error).Any("email", email).Any("role", role).
			Msg("couldn't get email (or) user-type")
		return tx.Error
	}

	if res.Password == password {
		log.Error().Any("old_password", res.Password).Any("new_pasword", password).
			Msg("entered the same password")
		return fmt.Errorf("entered the previous password")
	}

	res.Password = password
	res.LastUpdated = time.Now()

	pwd := DB.Model(models.AuthDetails{}).Save(&res)
	if pwd.Error != nil {
		log.Error().Err(pwd.Error).Any("password", res.Password).Any("email", res.Email).
			Any("user_type", res.UserType).Msg("entered the wrong password")
		return pwd.Error
	}
	return nil
}
