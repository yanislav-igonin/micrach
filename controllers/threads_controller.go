package controllers

import (
	"log"
	"math"
	"strconv"

	"github.com/dchest/captcha"
	"github.com/gofiber/fiber/v2"

	"micrach/config"
	"micrach/repositories"
)

func GetThreads(c *fiber.Ctx) error {
	pageString := c.Query("page", "1")
	page, err := strconv.Atoi(pageString)
	if err != nil {
		return c.Status(fiber.StatusNotFound).Render("pages/404", nil)
	}

	if page <= 0 {
		return c.Status(fiber.StatusNotFound).Render("pages/404", nil)
	}

	limit := 10
	offset := limit * (page - 1)
	threads, err := repositories.Posts.Get(limit, offset)
	if err != nil {
		log.Println("error:", err)
		return c.Status(fiber.StatusInternalServerError).Render("pages/500", nil)
	}
	count, err := repositories.Posts.GetThreadsCount()
	if err != nil {
		log.Println("error:", err)
		return c.Status(fiber.StatusInternalServerError).Render("pages/500", nil)
	}

	pagesCount := int(math.Ceil(float64(count) / 10))
	if page > pagesCount && count != 0 {
		return c.Status(fiber.StatusNotFound).Render("pages/404", nil)
	}

	captchaID := captcha.New()
	htmlData := repositories.GetThreadsHtmlData{
		Threads: threads,
		Pagination: repositories.HtmlPaginationData{
			PagesCount: pagesCount,
			Page:       page,
		},
		FormData: repositories.HtmlFormData{
			CaptchaID:       captchaID,
			IsCaptchaActive: config.App.IsCaptchaActive,
		},
	}
	return c.Render("pages/index", htmlData)
}

func GetThread(c *fiber.Ctx) error {
	threadID, err := c.ParamsInt("threadID")
	if err != nil {
		return c.Status(fiber.StatusNotFound).Render("pages/404", nil)
	}
	thread, err := repositories.Posts.GetThreadByPostID(threadID)
	if err != nil {
		log.Println("error:", err)
		return c.Status(fiber.StatusInternalServerError).Render("pages/500", nil)
	}
	if thread == nil {
		return c.Status(fiber.StatusNotFound).Render("pages/404", nil)
	}

	firstPost := thread[0]
	captchaID := captcha.New()
	htmlData := repositories.GetThreadHtmlData{
		Thread: thread,
		FormData: repositories.HtmlFormData{
			FirstPostID:     firstPost.ID,
			CaptchaID:       captchaID,
			IsCaptchaActive: config.App.IsCaptchaActive,
		},
	}
	return c.Render("pages/thread", htmlData)
}

// func CreateThread(c *gin.Context) {
// 	form, err := c.MultipartForm()
// 	if err != nil {
// 		log.Println("error:", err)
// 		c.HTML(http.StatusInternalServerError, "500.html", nil)
// 		return
// 	}

// 	// TODO: dat shit crashes if no fields in request
// 	text := form.Value["text"][0]
// 	title := form.Value["title"][0]
// 	filesInRequest := form.File["files"]
// 	validationErrorMessage := Utils.ValidatePost(title, text, filesInRequest)
// 	if validationErrorMessage != "" {
// 		errorHtmlData := Repositories.BadRequestHtmlData{
// 			Message: validationErrorMessage,
// 		}
// 		c.HTML(http.StatusBadRequest, "400.html", errorHtmlData)
// 		return
// 	}

// 	if Config.App.IsCaptchaActive {
// 		captchaID := form.Value["captchaId"][0]
// 		captchaString := form.Value["captcha"][0]
// 		isCaptchaValid := captcha.VerifyString(captchaID, captchaString)
// 		if !isCaptchaValid {
// 			errorHtmlData := Repositories.BadRequestHtmlData{
// 				Message: Repositories.InvalidCaptchaErrorMessage,
// 			}
// 			c.HTML(http.StatusBadRequest, "400.html", errorHtmlData)
// 			return
// 		}
// 	}

// 	conn, err := Db.Pool.Acquire(context.TODO())
// 	if err != nil {
// 		log.Println("error:", err)
// 		c.HTML(http.StatusInternalServerError, "500.html", nil)
// 		return
// 	}
// 	defer conn.Release()

// 	threadsCount, err := Repositories.Posts.GetThreadsCount()
// 	if err != nil {
// 		log.Println("error:", err)
// 		c.HTML(http.StatusInternalServerError, "500.html", nil)
// 		return
// 	}

// 	if threadsCount >= Config.App.ThreadsMaxCount {
// 		oldestThreadUpdatedAt, err := Repositories.Posts.GetOldestThreadUpdatedAt()
// 		if err != nil {
// 			log.Println("error:", err)
// 			c.HTML(http.StatusInternalServerError, "500.html", nil)
// 			return
// 		}
// 		err = Repositories.Posts.ArchiveThreadsFrom(oldestThreadUpdatedAt)
// 		if err != nil {
// 			log.Println("error:", err)
// 			c.HTML(http.StatusInternalServerError, "500.html", nil)
// 			return
// 		}
// 	}

// 	tx, err := conn.Begin(context.TODO())
// 	if err != nil {
// 		log.Println("error:", err)
// 		c.HTML(http.StatusInternalServerError, "500.html", nil)
// 		return
// 	}
// 	defer tx.Rollback(context.TODO())

// 	post := Repositories.Post{
// 		IsParent: true,
// 		Title:    title,
// 		Text:     text,
// 		IsSage:   false,
// 	}
// 	threadID, err := Repositories.Posts.CreateInTx(tx, post)
// 	if err != nil {
// 		log.Println("error:", err)
// 		c.HTML(http.StatusInternalServerError, "500.html", nil)
// 		return
// 	}

// 	err = Utils.CreateThreadFolder(threadID)
// 	if err != nil {
// 		log.Println("error:", err)
// 		c.HTML(http.StatusInternalServerError, "500.html", nil)
// 		return
// 	}

// 	for _, fileInRequest := range filesInRequest {
// 		file := Repositories.File{
// 			PostID: threadID,
// 			Name:   fileInRequest.Filename,
// 			// image/jpeg -> jpeg
// 			Ext:  strings.Split(fileInRequest.Header["Content-Type"][0], "/")[1],
// 			Size: int(fileInRequest.Size),
// 		}

// 		fileID, err := Repositories.Files.CreateInTx(tx, file)
// 		if err != nil {
// 			log.Println("error:", err)
// 			c.HTML(http.StatusInternalServerError, "500.html", nil)
// 			return
// 		}

// 		path := filepath.Join(
// 			Utils.UPLOADS_DIR_PATH,
// 			strconv.Itoa(threadID),
// 			"o",
// 			strconv.Itoa(fileID)+"."+file.Ext,
// 		)
// 		err = c.SaveUploadedFile(fileInRequest, path)
// 		if err != nil {
// 			log.Println("error:", err)
// 			c.HTML(http.StatusInternalServerError, "500.html", nil)
// 			return
// 		}
// 		// creating thumbnail
// 		thumbImg, err := Utils.MakeImageThumbnail(path, file.Ext, threadID, fileID)
// 		if err != nil {
// 			log.Println("error:", err)
// 			c.HTML(http.StatusInternalServerError, "500.html", nil)
// 			return
// 		}
// 		// saving thumbnail
// 		err = Utils.SaveImageThumbnail(thumbImg, threadID, fileID, file.Ext)
// 		if err != nil {
// 			log.Println("error:", err)
// 			c.HTML(http.StatusInternalServerError, "500.html", nil)
// 			return
// 		}
// 	}

// 	tx.Commit(context.TODO())

// 	c.Redirect(http.StatusFound, "/"+strconv.Itoa(threadID))
// }

// // Add new post in thread
// func UpdateThread(c *gin.Context) {
// 	threadIDString := c.Param("threadID")
// 	threadID, err := strconv.Atoi(threadIDString)
// 	if err != nil {
// 		c.HTML(http.StatusNotFound, "500.html", nil)
// 		return
// 	}

// 	isArchived, err := Repositories.Posts.GetIfThreadIsArchived(threadID)
// 	if isArchived {
// 		errorHtmlData := Repositories.BadRequestHtmlData{
// 			Message: Repositories.ThreadIsArchivedErrorMessage,
// 		}
// 		c.HTML(http.StatusBadRequest, "400.html", errorHtmlData)
// 		return
// 	}
// 	if err != nil {
// 		log.Println("error:", err)
// 		c.HTML(http.StatusInternalServerError, "500.html", nil)
// 		return
// 	}

// 	form, err := c.MultipartForm()
// 	if err != nil {
// 		log.Println("error:", err)
// 		c.HTML(http.StatusInternalServerError, "500.html", nil)
// 		return
// 	}

// 	// TODO: dat shit crashes if no fields in request
// 	text := form.Value["text"][0]
// 	filesInRequest := form.File["files"]
// 	validationErrorMessage := Utils.ValidatePost("", text, filesInRequest)
// 	if validationErrorMessage != "" {
// 		errorHtmlData := Repositories.BadRequestHtmlData{
// 			Message: validationErrorMessage,
// 		}
// 		c.HTML(http.StatusBadRequest, "400.html", errorHtmlData)
// 		return
// 	}

// 	if Config.App.IsCaptchaActive {
// 		captchaID := form.Value["captchaId"][0]
// 		captchaString := form.Value["captcha"][0]
// 		isCaptchaValid := captcha.VerifyString(captchaID, captchaString)
// 		if !isCaptchaValid {
// 			errorHtmlData := Repositories.BadRequestHtmlData{
// 				Message: Repositories.InvalidCaptchaErrorMessage,
// 			}
// 			c.HTML(http.StatusBadRequest, "400.html", errorHtmlData)
// 			return
// 		}
// 	}

// 	isSageField := form.Value["sage"]
// 	var isSageString string
// 	if len(isSageField) != 0 {
// 		isSageString = isSageField[0]
// 	}
// 	isSage := isSageString == "on"

// 	conn, err := Db.Pool.Acquire(context.TODO())
// 	if err != nil {
// 		log.Println("error:", err)
// 		c.HTML(http.StatusInternalServerError, "500.html", nil)
// 		return
// 	}
// 	defer conn.Release()

// 	tx, err := conn.Begin(context.TODO())
// 	if err != nil {
// 		log.Println("error:", err)
// 		c.HTML(http.StatusInternalServerError, "500.html", nil)
// 		return
// 	}
// 	defer tx.Rollback(context.TODO())

// 	if err != nil {
// 		log.Println("error:", err)
// 		c.HTML(http.StatusInternalServerError, "500.html", nil)
// 		return
// 	}
// 	post := Repositories.Post{
// 		IsParent: false,
// 		ParentID: &threadID,
// 		Title:    "",
// 		Text:     text,
// 		IsSage:   isSage,
// 	}
// 	postID, err := Repositories.Posts.CreateInTx(tx, post)
// 	if err != nil {
// 		log.Println("error:", err)
// 		c.HTML(http.StatusInternalServerError, "500.html", nil)
// 		return
// 	}

// 	postsCountInThread, err := Repositories.Posts.GetThreadPostsCount(threadID)
// 	if err != nil {
// 		log.Println("error:", err)
// 		c.HTML(http.StatusInternalServerError, "500.html", nil)
// 		return
// 	}
// 	isBumpLimit := postsCountInThread >= Config.App.ThreadBumpLimit
// 	isThreadBumped := !isBumpLimit && !isSage && !post.IsParent
// 	if isThreadBumped {
// 		err = Repositories.Posts.BumpThreadInTx(tx, threadID)
// 		if err != nil {
// 			log.Println("error:", err)
// 			c.HTML(http.StatusInternalServerError, "500.html", nil)
// 			return
// 		}
// 	}

// 	for _, fileInRequest := range filesInRequest {
// 		file := Repositories.File{
// 			PostID: postID,
// 			Name:   fileInRequest.Filename,
// 			// image/jpeg -> jpeg
// 			Ext:  strings.Split(fileInRequest.Header["Content-Type"][0], "/")[1],
// 			Size: int(fileInRequest.Size),
// 		}

// 		fileID, err := Repositories.Files.CreateInTx(tx, file)
// 		if err != nil {
// 			log.Println("error:", err)
// 			c.HTML(http.StatusInternalServerError, "500.html", nil)
// 			return
// 		}

// 		path := filepath.Join(
// 			Utils.UPLOADS_DIR_PATH,
// 			strconv.Itoa(threadID),
// 			"o",
// 			strconv.Itoa(fileID)+"."+file.Ext,
// 		)
// 		err = c.SaveUploadedFile(fileInRequest, path)
// 		if err != nil {
// 			log.Println("error:", err)
// 			c.HTML(http.StatusInternalServerError, "500.html", nil)
// 			return
// 		}
// 		// creating thumbnail
// 		thumbImg, err := Utils.MakeImageThumbnail(path, file.Ext, threadID, fileID)
// 		if err != nil {
// 			log.Println("error:", err)
// 			c.HTML(http.StatusInternalServerError, "500.html", nil)
// 			return
// 		}
// 		// saving thumbnail
// 		err = Utils.SaveImageThumbnail(thumbImg, threadID, fileID, file.Ext)
// 		if err != nil {
// 			log.Println("error:", err)
// 			c.HTML(http.StatusInternalServerError, "500.html", nil)
// 			return
// 		}
// 	}

// 	tx.Commit(context.TODO())

// 	c.Header("Refresh", "0")
// }
