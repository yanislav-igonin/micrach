package controllers

import (
	"context"
	"log"
	"math"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dchest/captcha"
	"github.com/gofiber/fiber/v2"

	"micrach/config"
	"micrach/db"
	"micrach/repositories"
	"micrach/utils"
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

func CreateThread(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		log.Println("error:", err)
		return c.Status(fiber.StatusInternalServerError).Render("pages/500", nil)
	}

	// TODO: dat shit crashes if no fields in request
	text := form.Value["text"][0]
	title := form.Value["title"][0]
	filesInRequest := form.File["files"]
	validationErrorMessage := utils.ValidatePost(title, text, filesInRequest)
	if validationErrorMessage != "" {
		errorHtmlData := repositories.BadRequestHtmlData{
			Message: validationErrorMessage,
		}
		return c.Status(fiber.StatusBadRequest).Render("pages/400", errorHtmlData)
	}

	if config.App.IsCaptchaActive {
		captchaID := form.Value["captchaId"][0]
		captchaString := form.Value["captcha"][0]
		isCaptchaValid := captcha.VerifyString(captchaID, captchaString)
		if !isCaptchaValid {
			errorHtmlData := repositories.BadRequestHtmlData{
				Message: repositories.InvalidCaptchaErrorMessage,
			}
			return c.Status(fiber.StatusBadRequest).Render("pages/400", errorHtmlData)
		}
	}

	conn, err := db.Pool.Acquire(context.TODO())
	if err != nil {
		log.Println("error:", err)
		return c.Status(fiber.StatusInternalServerError).Render("pages/500", nil)
	}
	defer conn.Release()

	threadsCount, err := repositories.Posts.GetThreadsCount()
	if err != nil {
		log.Println("error:", err)
		return c.Status(fiber.StatusInternalServerError).Render("pages/500", nil)
	}

	if threadsCount >= config.App.ThreadsMaxCount {
		oldestThreadUpdatedAt, err := repositories.Posts.GetOldestThreadUpdatedAt()
		if err != nil {
			log.Println("error:", err)
			return c.Status(fiber.StatusInternalServerError).Render("pages/500", nil)
		}
		err = repositories.Posts.ArchiveThreadsFrom(oldestThreadUpdatedAt)
		if err != nil {
			log.Println("error:", err)
			return c.Status(fiber.StatusInternalServerError).Render("pages/500", nil)
		}
	}

	tx, err := conn.Begin(context.TODO())
	if err != nil {
		log.Println("error:", err)
		return c.Status(fiber.StatusInternalServerError).Render("pages/500", nil)
	}
	defer tx.Rollback(context.TODO())

	post := repositories.Post{
		IsParent: true,
		Title:    title,
		Text:     text,
		IsSage:   false,
	}
	threadID, err := repositories.Posts.CreateInTx(tx, post)
	if err != nil {
		log.Println("error:", err)
		return c.Status(fiber.StatusInternalServerError).Render("pages/500", nil)
	}

	err = utils.CreateThreadFolder(threadID)
	if err != nil {
		log.Println("error:", err)
		return c.Status(fiber.StatusInternalServerError).Render("pages/500", nil)
	}

	for _, fileInRequest := range filesInRequest {
		file := repositories.File{
			PostID: threadID,
			Name:   fileInRequest.Filename,
			// image/jpeg -> jpeg
			Ext:  strings.Split(fileInRequest.Header["Content-Type"][0], "/")[1],
			Size: int(fileInRequest.Size),
		}

		fileID, err := repositories.Files.CreateInTx(tx, file)
		if err != nil {
			log.Println("error:", err)
			return c.Status(fiber.StatusInternalServerError).Render("pages/500", nil)
		}

		path := filepath.Join(
			utils.UPLOADS_DIR_PATH,
			strconv.Itoa(threadID),
			"o",
			strconv.Itoa(fileID)+"."+file.Ext,
		)
		err = c.SaveFile(fileInRequest, path)
		if err != nil {
			log.Println("error:", err)
			return c.Status(fiber.StatusInternalServerError).Render("pages/500", nil)
		}
		// creating thumbnail
		thumbImg, err := utils.MakeImageThumbnail(path, file.Ext, threadID, fileID)
		if err != nil {
			log.Println("error:", err)
			return c.Status(fiber.StatusInternalServerError).Render("pages/500", nil)
		}
		// saving thumbnail
		err = utils.SaveImageThumbnail(thumbImg, threadID, fileID, file.Ext)
		if err != nil {
			log.Println("error:", err)
			return c.Status(fiber.StatusInternalServerError).Render("pages/500", nil)
		}
	}

	tx.Commit(context.TODO())

	path := "/" + strconv.Itoa(threadID)
	return c.Redirect(path, fiber.StatusFound)
}

// Add new post in thread
func UpdateThread(c *fiber.Ctx) error {
	threadID, err := c.ParamsInt("threadID")
	if err != nil {
		return c.Status(fiber.StatusNotFound).Render("pages/404", nil)
	}

	isArchived, err := repositories.Posts.GetIfThreadIsArchived(threadID)
	if isArchived {
		errorHtmlData := repositories.BadRequestHtmlData{
			Message: repositories.ThreadIsArchivedErrorMessage,
		}
		return c.Status(fiber.StatusBadRequest).Render("pages/400", errorHtmlData)
	}
	if err != nil {
		log.Println("error:", err)
		return c.Status(fiber.StatusInternalServerError).Render("pages/500", nil)
	}

	form, err := c.MultipartForm()
	if err != nil {
		log.Println("error:", err)
		return c.Status(fiber.StatusInternalServerError).Render("pages/500", nil)
	}

	// TODO: dat shit crashes if no fields in request
	text := form.Value["text"][0]
	filesInRequest := form.File["files"]

	validationErrors := utils.ValidatePost2("", text, filesInRequest)
	if validationErrors != nil {
		thread, err := repositories.Posts.GetThreadByPostID(threadID)
		if err != nil {
			log.Println("error:", err)
			return c.Status(fiber.StatusInternalServerError).Render("pages/500", nil)
		}

		firstPost := thread[0]
		captchaID := form.Value["captchaId"][0]
		htmlData := repositories.GetThreadHtmlData{
			Thread: thread,
			FormData: repositories.HtmlFormData{
				FirstPostID:     firstPost.ID,
				CaptchaID:       captchaID,
				IsCaptchaActive: config.App.IsCaptchaActive,
				Errors:          *validationErrors,
			},
		}

		return c.Render("pages/thread", htmlData)
	}

	if config.App.IsCaptchaActive {
		captchaID := form.Value["captchaId"][0]
		captchaString := form.Value["captcha"][0]
		isCaptchaValid := captcha.VerifyString(captchaID, captchaString)
		if !isCaptchaValid {
			errorHtmlData := repositories.BadRequestHtmlData{
				Message: repositories.InvalidCaptchaErrorMessage,
			}
			return c.Status(fiber.StatusBadRequest).Render("pages/400", errorHtmlData)
		}
	}

	isSageField := form.Value["sage"]
	var isSageString string
	if len(isSageField) != 0 {
		isSageString = isSageField[0]
	}
	isSage := isSageString == "on"

	conn, err := db.Pool.Acquire(context.TODO())
	if err != nil {
		log.Println("error:", err)
		return c.Status(fiber.StatusInternalServerError).Render("pages/500", nil)

	}
	defer conn.Release()

	tx, err := conn.Begin(context.TODO())
	if err != nil {
		log.Println("error:", err)
		return c.Status(fiber.StatusInternalServerError).Render("pages/500", nil)

	}
	defer tx.Rollback(context.TODO())

	post := repositories.Post{
		IsParent: false,
		ParentID: &threadID,
		Title:    "",
		Text:     text,
		IsSage:   isSage,
	}
	postID, err := repositories.Posts.CreateInTx(tx, post)
	if err != nil {
		log.Println("error:", err)
		return c.Status(fiber.StatusInternalServerError).Render("pages/500", nil)

	}

	postsCountInThread, err := repositories.Posts.GetThreadPostsCount(threadID)
	if err != nil {
		log.Println("error:", err)
		return c.Status(fiber.StatusInternalServerError).Render("pages/500", nil)

	}
	isBumpLimit := postsCountInThread >= config.App.ThreadBumpLimit
	isThreadBumped := !isBumpLimit && !isSage && !post.IsParent
	if isThreadBumped {
		err = repositories.Posts.BumpThreadInTx(tx, threadID)
		if err != nil {
			log.Println("error:", err)
			return c.Status(fiber.StatusInternalServerError).Render("pages/500", nil)

		}
	}

	for _, fileInRequest := range filesInRequest {
		file := repositories.File{
			PostID: postID,
			Name:   fileInRequest.Filename,
			// image/jpeg -> jpeg
			Ext:  strings.Split(fileInRequest.Header["Content-Type"][0], "/")[1],
			Size: int(fileInRequest.Size),
		}

		fileID, err := repositories.Files.CreateInTx(tx, file)
		if err != nil {
			log.Println("error:", err)
			return c.Status(fiber.StatusInternalServerError).Render("pages/500", nil)

		}

		path := filepath.Join(
			utils.UPLOADS_DIR_PATH,
			strconv.Itoa(threadID),
			"o",
			strconv.Itoa(fileID)+"."+file.Ext,
		)
		err = c.SaveFile(fileInRequest, path)
		if err != nil {
			log.Println("error:", err)
			return c.Status(fiber.StatusInternalServerError).Render("pages/500", nil)

		}
		// creating thumbnail
		thumbImg, err := utils.MakeImageThumbnail(path, file.Ext, threadID, fileID)
		if err != nil {
			log.Println("error:", err)
			return c.Status(fiber.StatusInternalServerError).Render("pages/500", nil)

		}
		// saving thumbnail
		err = utils.SaveImageThumbnail(thumbImg, threadID, fileID, file.Ext)
		if err != nil {
			log.Println("error:", err)
			return c.Status(fiber.StatusInternalServerError).Render("pages/500", nil)

		}
	}

	tx.Commit(context.TODO())

	path := "/" + strconv.Itoa(threadID) + "#" + strconv.Itoa(postID)
	return c.Redirect(path)
}
