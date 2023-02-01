package main

import (
	"mywallet/config"
	"mywallet/db"
	"mywallet/handlers"
	"mywallet/middleware"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {

	conf, err := config.LoadConfig()
	if err != nil {
		log.Error().Err(err).Msg("error in loadingConfig")
		return
	}
	db.InitDB(conf)

	router := gin.Default()
	router.GET("/", handlers.HealthCheck)
	router.POST("/vendor-register", handlers.VendorRegister)
	router.POST("/customer-register", handlers.CustomerRegister)
	router.POST("/login", handlers.Login)
	router.PUT("/reset", middleware.IsAuthorized(), handlers.ResetPassword)

	router.Run("localhost:8080")
}
