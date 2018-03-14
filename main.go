package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kataras/iris"
	"github.com/okkur/reposeed-server/generator"
	"github.com/okkur/reposeed/cmd/reposeed/config"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Couldn't find .env file. Reading environment variables from system")
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
