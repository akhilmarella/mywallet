package db

import (
	"fmt"
	"mywallet/models"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
)

type Check struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

func AddCustomer(customer models.Customer) (uint64, error) {
	customer.CreatedAt = time.Now()
	var check Check

	em := DB.Raw("select email from customers where email = ?", customer.Email).Scan(&check)
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

	tx := DB.Create(&customer).Model(models.Customer{}).Scan(&customer)
	if tx.Error != nil {
		log.Error().Err(tx.Error).Any("customer", customer).
			Msg("error in creating user")
		return 0, tx.Error
	}

	return customer.ID, nil
}
