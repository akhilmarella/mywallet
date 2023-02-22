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

func GetCustomer(id int64) (*models.Customer, error) {
	var customer models.Customer
	tx := DB.Raw("select * from customers where id = ? ", id).Scan(&customer)
	if tx.Error != nil {
		log.Error().Err(tx.Error).Any("customer_id", id).Any("action", "db_customer.go_GetCustomer").
			Msg("error in reading customer details")
		return nil, tx.Error
	}
	return &customer, nil
}

func UpdateCustomer(customer models.Customer) (*models.Customer, error) {
	var currentCustomer models.Customer

	tx := DB.Raw("select * from customers where id = ? ", customer.ID).Scan(&currentCustomer)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		log.Error().Err(tx.Error).Any("id", customer.ID).Any("action", "db_customer.go_UpdateCustomer").
			Msg("user not found")
		return nil, tx.Error
	}

	if customer.FirstName == "" {
		customer.FirstName = currentCustomer.FirstName
	}

	if customer.LastName == "" {
		customer.LastName = currentCustomer.LastName
	}

	if customer.UserName == "" {
		customer.UserName = currentCustomer.UserName
	}

	if customer.Email == "" {
		customer.Email = currentCustomer.Email
	}

	if customer.PhoneNumber == "" {
		customer.PhoneNumber = currentCustomer.PhoneNumber
	}

	if customer.DOB == "" {
		customer.DOB = currentCustomer.DOB
	}
	customer.AddressID = currentCustomer.AddressID
	customer.LastUpdated = time.Now()
	//em := DB.Model(models.Customer{}).Save(&currentCustomer)
	em := DB.Save(&customer)
	if em.Error != nil {
		log.Error().Err(em.Error).Any("coustmer_details", customer).
			Any("action", "db_customer.go_updateCustomer").
			Msg("error in updating customer details")
		return nil, em.Error
	}

	return &models.Customer{
		ID:          customer.ID,
		FirstName:   customer.FirstName,
		LastName:    customer.LastName,
		UserName:    customer.UserName,
		Email:       customer.Email,
		PhoneNumber: customer.PhoneNumber,
		DOB:         customer.DOB,
		AddressID:   customer.AddressID,
		CreatedAt:   customer.CreatedAt,
		LastUpdated: customer.LastUpdated,
	}, nil
}
