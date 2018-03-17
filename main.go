package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/okkur/reposeed-server/generator"
	"github.com/okkur/reposeed/cmd/reposeed/config"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Couldn't find .env file. Reading environment variables from system")
	}
	app := gin.Default()
	config := &config.Config{}
	app.POST("/generate", func(ctx *gin.Context) {
		ctx.BindJSON(config)
		filename, err := generator.CreateFiles(*config, "../templates", config.Project.Name, os.Getenv("STORAGE"))
		if err.Code != 200 {
			ctx.JSON(400, err)
			ctx.Abort()
		} else {
			ctx.File(filename)
		}
	})
	app.Run(os.Getenv("PORT"))
}
