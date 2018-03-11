package main

import (
	"github.com/kataras/iris"
	"github.com/okkur/reposeed-server/config"
	"github.com/okkur/reposeed-server/generator"
)

func main() {
	app := iris.New()
	config := &config.Config{}
	app.Post("/generate", func(ctx iris.Context) {
		ctx.ReadJSON(config)
		filename := generator.CreateFiles(*config, "../templates", config.Project.Name)
		ctx.SendFile(filename, filename)
	})
	app.Run(iris.Addr(":8080"))
}
