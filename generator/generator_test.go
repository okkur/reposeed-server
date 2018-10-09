package generator

import (
	"os"
	"testing"

	"github.com/okkur/reposeed/cmd/reposeed/config"
	templates "github.com/okkur/reposeed/cmd/reposeed/templates"
	"github.com/rs/xid"
)

func Test_generateFile(t *testing.T) {
	fileNames := []string{}
	config := config.Config{}
	guid := xid.New()
	temps := parseTemplates(templates.GetTemplates())
	os.Setenv("STORAGE", "/tmp/storage/")
	err := generateFile(config, temps, "README.md", &guid, &fileNames)
	if err != nil {
		panic(err)
	}
	defer os.Remove(fileNames[0])
	_, err = os.Open(fileNames[0])
	if err != nil {
		t.Errorf("Couldn't open file, %s", "testfile")
	}
}
