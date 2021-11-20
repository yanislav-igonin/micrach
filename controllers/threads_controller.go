package controllers

import (
	"context"
	"log"
	"math"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"

	Config "micrach/config"
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

	captchaID := captcha.New()
	htmlData := Repositories.GetThreadsHtmlData{
		Threads:    threads,
		PagesCount: pagesCount,
		Page:       page,
		FormData: Repositories.HtmlFormData{
			CaptchaID: captchaID,
		},
	}
	c.HTML(http.StatusOK, "index.html", htmlData)
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

	firstPost := thread[0]
	captchaID := captcha.New()
	htmlData := Repositories.GetThreadHtmlData{
		Thread: thread,
		FormData: Repositories.HtmlFormData{
			FirstPostID: firstPost.ID,
			CaptchaID:   captchaID,
		},
	}
	c.HTML(http.StatusOK, "thread.html", htmlData)
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
	validationErrorMessage := Utils.ValidatePost(title, text, filesInRequest)
	if validationErrorMessage != "" {
		errorHtmlData := Repositories.BadRequestHtmlData{
			Message: validationErrorMessage,
		}
		c.HTML(http.StatusInternalServerError, "400.html", errorHtmlData)
		return
	}

	captchaID := form.Value["captchaId"][0]
	captchaString := form.Value["captcha"][0]
	isCaptchaValid := captcha.VerifyString(captchaID, captchaString)
	if !isCaptchaValid {
		errorHtmlData := Repositories.BadRequestHtmlData{
			Message: Repositories.InvalidCaptchaErrorMessage,
		}
		c.HTML(http.StatusInternalServerError, "400.html", errorHtmlData)
		return
	}

	conn, err := Db.Pool.Acquire(context.TODO())
	if err != nil {
		log.Println("error:", err)
		c.HTML(http.StatusInternalServerError, "500.html", nil)
		return
	}
	defer conn.Release()

	threadsCount, err := Repositories.Posts.GetCount()
	if err != nil {
		log.Println("error:", err)
		c.HTML(http.StatusInternalServerError, "500.html", nil)
		return
	}

	if threadsCount >= Config.App.ThreadsMaxCount {
		oldestThreadUpdatedAt, err := Repositories.Posts.GetOldestThreadUpdatedAt()
		if err != nil {
			log.Println("error:", err)
			c.HTML(http.StatusInternalServerError, "500.html", nil)
			return
		}
		err = Repositories.Posts.ArchiveThreadsFrom(oldestThreadUpdatedAt)
		if err != nil {
			log.Println("error:", err)
			c.HTML(http.StatusInternalServerError, "500.html", nil)
			return
		}
	}

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
	threadID, err := Repositories.Posts.CreateInTx(tx, post)
	if err != nil {
		log.Println("error:", err)
		c.HTML(http.StatusInternalServerError, "500.html", nil)
		return
	}

	err = Utils.CreateThreadFolder(threadID)
	if err != nil {
		log.Println("error:", err)
		c.HTML(http.StatusInternalServerError, "500.html", nil)
		return
	}

	for _, fileInRequest := range filesInRequest {
		file := Repositories.File{
			PostID: threadID,
			Name:   fileInRequest.Filename,
			// image/jpeg -> jpeg
			Ext:  strings.Split(fileInRequest.Header["Content-Type"][0], "/")[1],
			Size: int(fileInRequest.Size),
		}

		fileID, err := Repositories.Files.CreateInTx(tx, file)
		if err != nil {
			log.Println("error:", err)
			c.HTML(http.StatusInternalServerError, "500.html", nil)
			return
		}

		path := filepath.Join(
			Utils.UPLOADS_DIR_PATH,
			strconv.Itoa(threadID),
			"o",
			strconv.Itoa(fileID)+"."+file.Ext,
		)
		err = c.SaveUploadedFile(fileInRequest, path)
		if err != nil {
			log.Println("error:", err)
			c.HTML(http.StatusInternalServerError, "500.html", nil)
			return
		}
		// creating thumbnail
		thumbImg, err := Utils.MakeImageThumbnail(path, file.Ext, threadID, fileID)
		if err != nil {
			log.Println("error:", err)
			c.HTML(http.StatusInternalServerError, "500.html", nil)
			return
		}
		// saving thumbnail
		err = Utils.SaveImageThumbnail(thumbImg, threadID, fileID, file.Ext)
		if err != nil {
			log.Println("error:", err)
			c.HTML(http.StatusInternalServerError, "500.html", nil)
			return
		}
	}

	tx.Commit(context.TODO())

	c.Redirect(http.StatusFound, "/"+strconv.Itoa(threadID))
}

func UpdateThread(c *gin.Context) {
	threadIDString := c.Param("threadID")
	threadID, err := strconv.Atoi(threadIDString)
	if err != nil {
		c.HTML(http.StatusNotFound, "500.html", nil)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		log.Println("error:", err)
		c.HTML(http.StatusInternalServerError, "500.html", nil)
		return
	}

	// TODO: dat shit crashes if no fields in request
	text := form.Value["text"][0]
	filesInRequest := form.File["files"]
	validationErrorMessage := Utils.ValidatePost("", text, filesInRequest)
	if validationErrorMessage != "" {
		errorHtmlData := Repositories.BadRequestHtmlData{
			Message: validationErrorMessage,
		}
		c.HTML(http.StatusInternalServerError, "400.html", errorHtmlData)
		return
	}

	captchaID := form.Value["captchaId"][0]
	captchaString := form.Value["captcha"][0]
	isCaptchaValid := captcha.VerifyString(captchaID, captchaString)
	if !isCaptchaValid {
		errorHtmlData := Repositories.BadRequestHtmlData{
			Message: Repositories.InvalidCaptchaErrorMessage,
		}
		c.HTML(http.StatusInternalServerError, "400.html", errorHtmlData)
		return
	}

	isSageField := form.Value["sage"]
	var isSageString string
	if len(isSageField) != 0 {
		isSageString = isSageField[0]
	}
	isSage := isSageString == "on"

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
		IsParent: false,
		ParentID: threadID,
		Title:    "",
		Text:     text,
		IsSage:   isSage,
	}
	postID, err := Repositories.Posts.CreateInTx(tx, post)
	if err != nil {
		log.Println("error:", err)
		c.HTML(http.StatusInternalServerError, "500.html", nil)
		return
	}

	for _, fileInRequest := range filesInRequest {
		file := Repositories.File{
			PostID: postID,
			Name:   fileInRequest.Filename,
			// image/jpeg -> jpeg
			Ext:  strings.Split(fileInRequest.Header["Content-Type"][0], "/")[1],
			Size: int(fileInRequest.Size),
		}

		fileID, err := Repositories.Files.CreateInTx(tx, file)
		if err != nil {
			log.Println("error:", err)
			c.HTML(http.StatusInternalServerError, "500.html", nil)
			return
		}

		path := filepath.Join(
			Utils.UPLOADS_DIR_PATH,
			strconv.Itoa(threadID),
			"o",
			strconv.Itoa(fileID)+"."+file.Ext,
		)
		err = c.SaveUploadedFile(fileInRequest, path)
		if err != nil {
			log.Println("error:", err)
			c.HTML(http.StatusInternalServerError, "500.html", nil)
			return
		}
		// creating thumbnail
		thumbImg, err := Utils.MakeImageThumbnail(path, file.Ext, threadID, fileID)
		if err != nil {
			log.Println("error:", err)
			c.HTML(http.StatusInternalServerError, "500.html", nil)
			return
		}
		// saving thumbnail
		err = Utils.SaveImageThumbnail(thumbImg, threadID, fileID, file.Ext)
		if err != nil {
			log.Println("error:", err)
			c.HTML(http.StatusInternalServerError, "500.html", nil)
			return
		}
	}

	tx.Commit(context.TODO())

	c.Header("Refresh", "0")
}
