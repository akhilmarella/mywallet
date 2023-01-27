package db

import (
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func InitDB() {
	// db, err := gorm.Open("postgres", fmt.Sprintf(
	// 	"host=%s port=%v dbname=%s user=%s password=%s sslmode=disable",
	// 	"localhost",
	// 	"5432",
	// 	"golang",
	// 	"akhil",
	// 	"marella12",
	// ))
	// if err!=nil{
	// 	log.Println("DB init fail",err)
	// 	os.Exit(1)
	// }

}
