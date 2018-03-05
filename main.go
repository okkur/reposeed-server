package main

import (
	"fmt"
	"os/exec"

	"github.com/kataras/iris"
	"github.com/okkur/reposeed-server/config"
	"github.com/okkur/reposeed-server/generator"
)

func main() {
	app := iris.New()
	config := &config.Config{}
	app.Post("/generate", func(ctx iris.Context) {
		ctx.ReadJSON(config)
		// config, err := yaml.Marshal(&config)
		// if err != nil {
		// 	log.Fatalf("error: %v", err)
		// }

		// file, err = ioutil.TempFile("", "config")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer os.Remove(file.Name())
		// if _, err := file.Write(config); err != nil {
		// 	log.Fatal(err)
		// }
		zw := generator.CreateFiles(*config, "../templates", config.Project.Name)
		zw.Close()
		// generator.GenerateFile(config)
		fmt.Println(config.Repo.Type)
		// panic(config.Readme.UsageExample)
	})
	cmd := exec.Command("reposeed", "a-z", "A-Z")
	cmd.Run()
	app.Run(iris.Addr(":8080"))
}
