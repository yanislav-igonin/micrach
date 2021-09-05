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
		log.Println("error:", err)
		c.HTML(http.StatusOK, "500.html", nil)
		return
	}
	c.HTML(http.StatusOK, "index.html", threads)
}

func GetThread(c *gin.Context) {
	threadIDString := c.Param("threadID")
	threadID, err := strconv.Atoi(threadIDString)
	if err != nil {
		c.HTML(http.StatusOK, "404.html", nil)
		return
	}
	thread, err := Repositories.Posts.GetThreadByPostID(threadID)
	if err != nil {
		log.Println("error:", err)
		c.HTML(http.StatusOK, "500.html", nil)
		return
	}
	if thread == nil {
		c.HTML(http.StatusOK, "404.html", nil)
		return
	}
	c.HTML(http.StatusOK, "thread.html", thread)
}

func CreateThread(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Problem uploading file!",
		})
		return
	}

	// TODO: dat shit crashes if no fields in request
	// text := form.Value["text"][0]
	// title := form.Value["title"][0]
	// isSageString := form.Value["isSage"][0]
	// isSage, err := strconv.ParseBool(isSageString)
	// if err != nil {
	// 	// TODO: validation error
	// 	response := Dto.GetInternalServerErrorResponse()
	// 	c.JSON(http.StatusInternalServerError, response)
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{"route": form})
}

func UpdateThread(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"route": "update thread"})
}
