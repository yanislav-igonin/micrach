package controlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	Repositories "micrach/repositories"
)

func GetThreads(c *gin.Context) {
	threads, err := Repositories.Threads.Get(10, 10)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": true})
		return
	}
	c.HTML(http.StatusOK, "index.html", threads)
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
