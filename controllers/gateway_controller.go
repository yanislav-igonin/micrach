package controllers

import (
	Config "micrach/config"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	headerKey := c.Request.Header.Get("Authorization")
	if Config.App.Gateway.ApiKey != headerKey {
		c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
