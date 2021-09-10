package controlers

import (
	"context"
	"log"
	"math"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"

	Db "micrach/db"
	Repositories "micrach/repositories"
	Utils "micrach/utils"
)

func GetThreads(c *gin.Context) {
	pageString := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageString)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}

	if page <= 0 {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}

	limit := 10
	offset := limit * (page - 1)
	threads, err := Repositories.Posts.Get(limit, offset)
	if err != nil {
		log.Println("error:", err)
		c.HTML(http.StatusInternalServerError, "500.html", nil)
		return
	}
	count, err := Repositories.Posts.GetCount()
	if err != nil {
		log.Println("error:", err)
		c.HTML(http.StatusInternalServerError, "500.html", nil)
		return
	}

	pagesCount := int(math.Ceil(float64(count) / 10))
	if page > pagesCount && count != 0 {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}
	data := Repositories.IndexPageData{
		Threads:    threads,
		PagesCount: pagesCount,
		Page:       page,
	}
	c.HTML(http.StatusOK, "index.html", data)
}

func GetThread(c *gin.Context) {
	threadIDString := c.Param("threadID")
	threadID, err := strconv.Atoi(threadIDString)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}
	thread, err := Repositories.Posts.GetThreadByPostID(threadID)
	if err != nil {
		log.Println("error:", err)
		c.HTML(http.StatusInternalServerError, "500.html", nil)
		return
	}
	if thread == nil {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}
	c.HTML(http.StatusOK, "thread.html", thread)
}

func CreateThread(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		log.Println("error:", err)
		c.HTML(http.StatusInternalServerError, "500.html", nil)
		return
	}

	// TODO: dat shit crashes if no fields in request
	text := form.Value["text"][0]
	title := form.Value["title"][0]
	filesInRequest := form.File["files"]

	conn, err := Db.Pool.Acquire(context.TODO())
	if err != nil {
		log.Println("error:", err)
		c.HTML(http.StatusInternalServerError, "500.html", nil)
		return
	}
	defer conn.Release()

	tx, err := conn.Begin(context.TODO())
	if err != nil {
		log.Println("error:", err)
		c.HTML(http.StatusInternalServerError, "500.html", nil)
		return
	}
	defer tx.Rollback(context.TODO())

	post := Repositories.Post{
		IsParent: true,
		Title:    title,
		Text:     text,
		IsSage:   false,
	}
	postID, err := Repositories.Posts.CreateInTx(tx, post)
	if err != nil {
		log.Println("error:", err)
		c.HTML(http.StatusInternalServerError, "500.html", nil)
		return
	}

	err = Utils.CreateThreadFolder(postID)
	if err != nil {
		log.Println("error:", err)
		c.HTML(http.StatusInternalServerError, "500.html", nil)
		return
	}

	for _, fileInRequest := range filesInRequest {
		path := filepath.Join(
			Utils.UPLOADS_DIR_PATH,
			strconv.Itoa(postID),
			fileInRequest.Filename,
		)
		log.Println(path)
		file := Repositories.File{
			PostID: postID,
			Name:   fileInRequest.Filename,
			Ext:    fileInRequest.Header["Content-Type"][0],
			Size:   int(fileInRequest.Size),
		}

		err := Repositories.Files.CreateInTx(tx, file)
		if err != nil {
			log.Println("error:", err)
			c.HTML(http.StatusInternalServerError, "500.html", nil)
			return
		}

		err = c.SaveUploadedFile(fileInRequest, path)
		if err != nil {
			log.Println("error:", err)
			c.HTML(http.StatusInternalServerError, "500.html", nil)
			return
		}
	}

	tx.Commit(context.TODO())

	c.Redirect(http.StatusFound, "/"+strconv.Itoa(postID))
}

func UpdateThread(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"route": "update thread"})
}
