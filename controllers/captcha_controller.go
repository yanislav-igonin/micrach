package controllers

import (
	"bytes"
	"log"

	"github.com/dchest/captcha"
	"github.com/gofiber/fiber/v3"
)

func GetCaptcha(c fiber.Ctx) error {
	ID := c.Params("captchaID")
	var content bytes.Buffer
	err := captcha.WriteImage(&content, ID, captcha.StdWidth, captcha.StdHeight)
	if err != nil {
		log.Println("error:", err)
		return c.Status(fiber.StatusInternalServerError).Render("pages/500", nil)
	}

	c.RequestCtx().SetContentType("image/png")
	return c.Send(content.Bytes())
}
