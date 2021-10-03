package controllers

import (
	"bytes"
	"log"
	"net/http"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

func GetCaptcha(c *gin.Context) {
	ID := c.Param("captchaID")
	var content bytes.Buffer
	err := captcha.WriteImage(&content, ID, captcha.StdWidth, captcha.StdHeight)
	if err != nil {
		log.Println("error:", err)
		c.HTML(http.StatusInternalServerError, "500.html", nil)
		return
	}

	c.Data(200, "image/png", content.Bytes())
}
