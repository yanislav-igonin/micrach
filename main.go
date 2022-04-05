package main

import (
	"log"
	"micrach/build"
	"micrach/config"
	"micrach/db"
	"micrach/repositories"
	"micrach/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
)

// import (
// 	"html/template"
// 	"log"
// 	"strconv"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	_ "github.com/joho/godotenv/autoload"
// 	limiter "github.com/ulule/limiter/v3"
// 	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
// 	memory "github.com/ulule/limiter/v3/drivers/store/memory"

// 	Build "micrach/build"
// 	Config "micrach/config"
// 	Controllers "micrach/controllers"
// 	Db "micrach/db"
// 	Gateway "micrach/gateway"
// 	Repositories "micrach/repositories"
// 	Templates "micrach/templates"
// 	Utils "micrach/utils"
// )

// func main() {
// 	Config.Init()
// 	Db.Init()
// 	Db.Migrate()
// 	defer Db.Pool.Close()
// 	gin.SetMode(Config.App.Env)
// 	if Config.App.IsDbSeeded {
// 		Repositories.Seed()
// 	}

// 	if Config.App.Env == "release" {
// 		Build.RenameCss()
// 	}

// 	err := Utils.CreateUploadsFolder()
// 	if err != nil {
// 		log.Panicln(err)
// 	}

// 	rate := limiter.Rate{
// 		Period: 1 * time.Hour,
// 		Limit:  1000,
// 	}
// 	rateLimiterStore := memory.NewStore()
// 	instance := limiter.New(rateLimiterStore, rate)
// 	middleware := mgin.NewMiddleware(instance)

// 	router := gin.New()
// 	router.Use(gin.Recovery())

// 	router.SetFuncMap(template.FuncMap{
// 		"Iterate": Templates.Iterate,
// 		"NotNil":  Templates.NotNil,
// 	})
// 	router.LoadHTMLGlob("templates/**/*")
// 	router.ForwardedByClientIP = true
// 	if Config.App.IsRateLimiterEnabled {
// 		router.Use(middleware)
// 	}
// 	router.Static("/uploads", "./uploads")
// 	router.Static("/static", "./static")
// 	if Config.App.Gateway.Url != "" {
// 		router.GET("/api/ping", Gateway.Ping)
// 		Gateway.Connect()
// 	}
// 	router.GET("/", Controllers.GetThreads)
// 	router.POST("/", Controllers.CreateThread)
// 	router.GET("/:threadID", Controllers.GetThread)
// 	router.POST("/:threadID", Controllers.UpdateThread)
// 	router.GET("/captcha/:captchaID", Controllers.GetCaptcha)

// 	log.Println("port", Config.App.Port, "- online")
// 	log.Println("all systems nominal")

// 	router.Run(":" + strconv.Itoa(Config.App.Port))
// }

func main() {
	config.Init()
	db.Init()
	db.Migrate()
	defer db.Pool.Close()

	if config.App.IsDbSeeded {
		repositories.Seed()
	}

	if config.App.Env == "production" {
		build.RenameCss()
	}

	err := utils.CreateUploadsFolder()
	if err != nil {
		log.Panicln(err)
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("get threads")
	})
	app.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("create thread")
	})
	app.Get("/:threadID", func(c *fiber.Ctx) error {
		return c.SendString("get thread by id")
	})
	app.Post("/threadID", func(c *fiber.Ctx) error {
		return c.SendString("create post in thread")
	})
	app.Get("/captcha/:captchaID", func(c *fiber.Ctx) error {
		return c.SendString("get captcha by id")
	})

	log.Println("app - online, port -", strconv.Itoa(config.App.Port))
	log.Println("all systems nominal")
	err = app.Listen(":" + strconv.Itoa(config.App.Port))
	if err != nil {
		log.Println("app - ofline")
		log.Panicln(err)
	}
}
