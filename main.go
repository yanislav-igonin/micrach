package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"

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
	Repositories.Seed()

	err := Utils.CreateUploadsFolder()
	if err != nil {
		log.Panicln(err)
	}

	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
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
