package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kataras/iris"
	"github.com/okkur/reposeed-server/config"
	"github.com/okkur/reposeed-server/generator"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app := iris.New()
	config := &config.Config{}
	app.Post("/generate", func(ctx iris.Context) {
		ctx.ReadJSON(config)
		filename := generator.CreateFiles(*config, "../templates", config.Project.Name, os.Getenv("STORAGE"))
		ctx.SendFile(filename, filename)
	})
	app.Run(iris.Addr(os.Getenv("PORT")))
}
