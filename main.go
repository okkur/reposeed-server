package main

import (
	"errors"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/okkur/reposeed-server/generator"
	"github.com/okkur/reposeed/cmd/reposeed/config"
)

// SupportedConfigVersion Supported reposeed config version
const SupportedConfigVersion = "v1"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Couldn't find .env file. Reading environment variables from system")
	}
	app := gin.Default()
	config := &config.Config{}
	app.POST("/generate", func(ctx *gin.Context) {
		err = ctx.BindJSON(config)
		if err != nil {
			ctx.AbortWithError(422, errors.New("couldn't parse the given config"))
		}
		if config.Project.Version == SupportedConfigVersion {
			filename, err := generator.CreateFiles(*config)
			if err != nil {
				ctx.AbortWithError(400, err)
			} else {
				ctx.File(filename)
			}
		} else {
			ctx.JSON(422, "ConfigVersion: Invalid config version")
		}
	})
	app.Run(os.Getenv("PORT"))
}
