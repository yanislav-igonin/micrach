package controlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetThreads(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"route": "get threads"})
}

func GetThread(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"route": "get thread"})
}

func CreateThread(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"route": "create thread"})
}

func UpdateThread(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"route": "update thread"})
}
