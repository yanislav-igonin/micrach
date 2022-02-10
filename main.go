package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	limiter "github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"

	Build "micrach/build"
	Config "micrach/config"
	Controllers "micrach/controllers"
	Db "micrach/db"
	Repositories "micrach/repositories"
	Utils "micrach/utils"
)

func main() {
	Config.Init()
	Db.Init()
	Db.Migrate()
	defer Db.Pool.Close()
	gin.SetMode(Config.App.Env)
	if Config.App.IsDbSeeded {
		Repositories.Seed()
	}

	if Config.App.Env == "release" {
		Build.RenameCss()
	}

	err := Utils.CreateUploadsFolder()
	if err != nil {
		log.Panicln(err)
	}

	rate := limiter.Rate{
		Period: 1 * time.Hour,
		Limit:  1000,
	}
	rateLimiterStore := memory.NewStore()
	instance := limiter.New(rateLimiterStore, rate)
	middleware := mgin.NewMiddleware(instance)

	router := gin.New()
	router.Use(gin.Recovery())

	router.SetFuncMap(template.FuncMap{
		"Iterate": func(count int) []int {
			var i int
			var Items []int
			for i = 1; i < count+1; i++ {
				Items = append(Items, i)
			}
			return Items
		},
	})
	router.LoadHTMLGlob("templates/**/*")
	router.ForwardedByClientIP = true
	if Config.App.IsRateLimiterEnabled {
		router.Use(middleware)
	}
	router.Static("/uploads", "./uploads")
	router.Static("/static", "./static")
	if Config.App.Gateway.Url != "" {
		router.GET("/ping", Controllers.Ping)
		// make http request to gateway to tell the board id and description
		requestBody, _ := json.Marshal(map[string]string{
			"id":          Config.App.Gateway.BoardId,
			"description": Config.App.Gateway.BoardDescription,
		})
		url := Config.App.Gateway.Url + "/boards"
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
		if err != nil {
			log.Panicln(err)
		}
		req.Header.Set("Authorization", Config.App.Gateway.ApiKey)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Panicln(err)
		}
		//We Read the response body on the line below.
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Panicln(err)
		}
	}
	router.GET("/", Controllers.GetThreads)
	router.POST("/", Controllers.CreateThread)
	router.GET("/:threadID", Controllers.GetThread)
	router.POST("/:threadID", Controllers.UpdateThread)
	router.GET("/captcha/:captchaID", Controllers.GetCaptcha)

	log.Println("port", Config.App.Port, "- online")
	log.Println("all systems nominal")

	router.Run(":" + strconv.Itoa(Config.App.Port))
}
