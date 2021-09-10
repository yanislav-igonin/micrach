package main

import (
	"html/template"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	limiter "github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"

	Config "micrach/config"
	Controllers "micrach/controllers"
	Db "micrach/db"
	Repositories "micrach/repositories"
	Utils "micrach/utils"
)

func main() {
	Config.Init()
	Db.Init()
	defer Db.Pool.Close()
	gin.SetMode(Config.App.Env)
	if Config.App.SeedDb {
		Repositories.Seed()
	}

	err := Utils.CreateUploadsFolder()
	if err != nil {
		log.Panicln(err)
	}

	rate := limiter.Rate{
		Period: 1 * time.Hour,
		Limit:  1000,
	}
	store := memory.NewStore()
	instance := limiter.New(store, rate)
	middleware := mgin.NewMiddleware(instance)

	router := gin.Default()
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
	router.LoadHTMLGlob("templates/*.html")
	router.ForwardedByClientIP = true
	if Config.App.IsRateLimiterEnabled {
		router.Use(middleware)
	}
	router.Static("/uploads", "./uploads")
	router.Static("/static", "./static")
	router.GET("/", Controllers.GetThreads)
	router.POST("/", Controllers.CreateThread)
	router.GET("/:threadID", Controllers.GetThread)
	router.POST("/:threadID", Controllers.UpdateThread)

	log.Println("port", Config.App.Port, "- online")
	log.Println("all systems nominal")

	router.Run(":" + strconv.Itoa(Config.App.Port))
}
