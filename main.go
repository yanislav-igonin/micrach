package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"

	Config "micrach/config"
	Controllers "micrach/controllers"
	Db "micrach/db"
	// Utils "micrach/utils"
)

func main() {
	Config.Init()
	Db.Init()
	defer Db.Pool.Close()
	gin.SetMode(Config.App.Env)

	router := gin.Default()
	router.GET("/", Controllers.GetThreads)
	router.POST("/", Controllers.CreateThread)
	router.GET("/:threadId", Controllers.GetThread)
	router.POST("/:threadId", Controllers.UpdateThread)

	log.Println("all systems nominal")

	router.Run("localhost:" + strconv.Itoa(Config.App.Port))
}
