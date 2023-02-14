package main

import (
	"mywallet/config"
	"mywallet/db"
	"mywallet/handlers"
	"mywallet/middleware"
	"mywallet/store"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Error().Err(err).
			Any("action:", "cmd/app_main.go_main").
			Msg("error in loadingConfig")
		return
	}

	db.InitDB(conf)
	store.InitRedis()
	router := gin.Default()
	router.GET("/", handlers.HealthCheck)
	router.POST("/vendor-register", handlers.VendorRegister)
	router.GET("/vendor/:id", middleware.IsAuthorized(), handlers.GetVendor)

	router.POST("/customer-register", handlers.CustomerRegister)
	router.GET("/customer/:id", middleware.IsAuthorized(), handlers.GetCustomer)

	router.POST("/address-register", middleware.IsAuthorized(), handlers.AddressRegister)

	router.POST("/login", handlers.Login)
	router.PUT("/reset", middleware.IsAuthorized(), handlers.ResetPassword)
	router.POST("/refresh", handlers.RefreshToken)
	router.DELETE("/remove", handlers.Logout)

	router.Run("localhost:8080")
}
