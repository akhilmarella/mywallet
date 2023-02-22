package db

import (
	"fmt"
	"mywallet/config"
	"mywallet/models"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

var DB *gorm.DB

func InitDB(conf config.Configuration) {
	db, err := gorm.Open("postgres", fmt.Sprintf(
		"host=%s port=%v dbname=%s user=%s password=%s sslmode=disable",
		conf.Db.Host,
		conf.Db.Port,
		conf.Db.Name,
		conf.Db.User,
		conf.Db.Password,
	))
	if err != nil {
		log.Error().Err(err).Any("action", "db_db.go_InitDB").Msg("DB init fail")
		os.Exit(1)
	}
	err = db.Debug().AutoMigrate(models.Vendor{}, models.AuthDetails{}, models.Customer{},
		models.Address{}, models.Wallet{}).Error
	if err != nil {
		log.Error().Err(err).
			Any("action:", "db_db.go_InitDB").
			Msg("error in AutoMigrate Table")
		return
	}
	DB = db
}
