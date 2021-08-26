package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	Config "micrach/config"
	// Controllers "micrach/controllers"
	Db "micrach/db"
	// Utils "micrach/utils"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Panicln(err)
	}

	Config.Init()
	Db.Init()
	defer Db.Pool.Close()
	// err = Utils.CreateUploadsFolder()
	if err != nil {
		log.Panicln(err)
	}
	gin.SetMode(Config.App.Env)

	router := gin.Default()
	// router.GET("/boards", Controllers.GetAllBoards)
	// router.GET("/threads/:boardId", Controllers.GetThreads)
	// router.POST("/threads/:boardId", Controllers.CreateThread)
	// router.POST("/posts/:boardId/:threadId", Controllers.CreatePost)

	log.Println("all systems nominal")

	router.Run("localhost:" + strconv.Itoa(Config.App.Port))
}
