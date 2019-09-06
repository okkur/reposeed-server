package generator

import (
	"os"
	"testing"

	"github.com/rs/xid"
	"go.okkur.org/reposeed/cmd/config"
	templates "go.okkur.org/reposeed/cmd/templates"
)

func Test_generateFile(t *testing.T) {
	config := config.Config{}
	guid := xid.New()
	temps := parseTemplates(templates.GetTemplates())
	os.Setenv("STORAGE", "/tmp/storage/")
	projectPath := os.Getenv("STORAGE") + guid.String() + "/"
	err := os.MkdirAll(projectPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	zip, writer, err := initializeZipWriter(projectPath + config.Project.Name + ".zip")
	defer zip.Close()
	err = generateFile(config, temps, "README.md", projectPath, writer)
	if err != nil {
		t.Fatalf("Couldn't generate the file: %s", err.Error())
	}
}
