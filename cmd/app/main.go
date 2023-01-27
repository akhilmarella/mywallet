package main

import (
	"mywallet/config"
	"mywallet/handlers"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	_, err := config.LoadConfig()
	if err != nil {
		log.Error().Err(err).Msg("error in loadingConfig")
		return
	}

	router := gin.Default()
	router.GET("/", handlers.HealthCheck)

	router.Run("localhost:8080")
}
