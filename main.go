package main

import (
	"context"
	"log"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"

	"micrach/build"
	"micrach/config"
	"micrach/controllers"
	"micrach/db"
	"micrach/gateway"
	"micrach/repositories"
	"micrach/templates"
	"micrach/utils"
)

func skipRateLimit(isLocal bool, path string) bool {
	isRequestForStatic := strings.Contains(path, "/static") ||
		strings.Contains(path, "/uploads") ||
		strings.Contains(path, "/captcha")
	return isRequestForStatic || isLocal
}

func registerStaticRoutes(app *fiber.App) {
	app.Static("/uploads", "./uploads")
	app.Static("/static", "./static")
}

func main() {
	ctx := context.Background()
	config.Init()
	db.Init(ctx)
	db.Migrate(ctx)
	defer db.Pool.Close()

	if config.App.IsDbSeeded {
		repositories.Seed(ctx)
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
	if config.App.IsRateLimiterEnabled {
		app.Use(limiter.New(limiter.Config{
			Next: func(c *fiber.Ctx) bool {
				return skipRateLimit(c.IsFromLocal(), c.Path())
			},
			Max: 50,
		}))
	}
	app.Use(compress.New())

	registerStaticRoutes(app)

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
