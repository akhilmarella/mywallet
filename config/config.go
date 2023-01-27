package config

import (
	"github.com/joho/godotenv"
)

type Configuration struct {
	DB
}

type DB struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func LoadConfig() (Configuration, error) {
	err := godotenv.Load("/home/akhil/Github/mywallet/.env")
	if err != nil {
		return Configuration{}, err
	}

	return Configuration{
		DB: DB{
			Host:     "localhost",
			Port:     "5432",
			Name:     "golang",
			User:     "akhil",
			Password: "marella12",
		},
	}, err
}
