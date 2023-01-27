package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	t := time.Now()
	// fmt.Println("current time:",t)
	c.JSON(200,t)
}
