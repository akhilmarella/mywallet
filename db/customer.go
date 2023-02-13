package db

import (
	"fmt"
	"mywallet/models"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
)

func AddCustomer(customer models.Customer) (int64, error) {
	customer.CreatedAt = time.Now()
	var check models.Check

	em := DB.Raw("select email from customers where email = ?", customer.Email).Scan(&check)
	if em.Error != nil && em.Error != gorm.ErrRecordNotFound {
		log.Error().Err(em.Error).Any("email", check).
			Any("action:", "db_customer.go_AddCustomer").
			Msg("error in email")
		return 0, em.Error
	}

	if check.Email != "" {
		log.Error().Any("email", check).
			Any("action:", "db_customer.go_AddCustomer").
			Msg("email already registered")
		return 0, fmt.Errorf("email is already exist :%v", check.Email)
	}

	tx := DB.Create(&customer).Model(models.Customer{}).Scan(&customer)
	if tx.Error != nil {
		log.Error().Err(tx.Error).Any("customer", customer).
			Any("action:", "db_customer.go_AddCustomer").
			Msg("error in creating user")
		return 0, tx.Error
	}

	return customer.ID, nil
}
