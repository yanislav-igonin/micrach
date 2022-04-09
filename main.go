package main

import (
	"log"
	"strconv"

	_ "github.com/joho/godotenv/autoload"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"

	"micrach/build"
	"micrach/config"
	"micrach/controllers"
	"micrach/db"
	"micrach/gateway"
	"micrach/repositories"
	"micrach/templates"
	"micrach/utils"
)

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

	engine := html.New("./templates", ".html")
	engine.AddFunc("Iterate", templates.Iterate)
	engine.AddFunc("NotNil", templates.NotNil)

	app := fiber.New(fiber.Config{Views: engine})

	app.Use(recover.New())
	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool { return c.IsFromLocal() },
	}))
	app.Use(compress.New())
	app.Use(etag.New())

	app.Static("/uploads", "./uploads")
	app.Static("/static", "./static")

	app.Get("/", controllers.GetThreads)
	app.Post("/", controllers.CreateThread)
	app.Get("/:threadID", controllers.GetThread)
	app.Post("/:threadID", controllers.UpdateThread)
	app.Get("/captcha/:captchaID", controllers.GetCaptcha)

	if config.App.Gateway.Url != "" {
		app.Get("/api/ping", gateway.Ping)
		gateway.Connect()
	}

	log.Println("app - online, port -", strconv.Itoa(config.App.Port))
	log.Println("all systems nominal")
	err = app.Listen(":" + strconv.Itoa(config.App.Port))
	if err != nil {
		log.Println("app - ofline")
		log.Panicln(err)
	}
}
