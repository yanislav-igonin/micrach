package controlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	Repositories "micrach/repositories"
)

func GetThreads(c *gin.Context) {
	threads, err := Repositories.Posts.Get(10, 10)
	if err != nil {
		// TODO: наверное, тут 500 будет
		log.Println("error:", err)
		c.JSON(http.StatusOK, gin.H{"error": true})
		return
	}
	c.HTML(http.StatusOK, "index.html", threads)
}

func GetThread(c *gin.Context) {
	threadIDString := c.Param("threadID")
	threadID, err := strconv.Atoi(threadIDString)
	if err != nil {
		// TODO: тут 400 будет
		log.Println("error:", err)
		c.JSON(http.StatusOK, gin.H{"error": true})
		return
	}
	thread, err := Repositories.Posts.GetThreadByPostID(threadID)
	if err != nil {
		// TODO: наверное, тут 500 будет
		log.Println("error:", err)
		c.JSON(http.StatusOK, gin.H{"error": true})
		return
	}
	c.HTML(http.StatusOK, "thread.html", thread)
}

func CreateThread(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"route": "create thread"})
}

func UpdateThread(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"route": "update thread"})
}
