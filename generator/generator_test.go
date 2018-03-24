package generator

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/okkur/reposeed/cmd/reposeed/config"
	"github.com/rs/xid"
)

func TestZipFiles(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "testfile")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	filenames := []string{tmpfile.Name()}
	guid := xid.New()
	os.Setenv("STORAGE", "../storage/")
	projectsPath := os.Getenv("STORAGE") + guid.String() + "/"
	err = createDir(projectsPath, tmpfile.Name())
	zipPath, err := ZipFiles("zipfile.zip", &filenames, os.Getenv("STORAGE"), guid.String())
	defer os.Remove(zipPath)
	zipName := strings.Split(zipPath, "/")
	if err != nil {
		t.Error(err)
	}
	if zipName[len(zipName)-1] != "zipfile.zip" {
		t.Errorf("Expected %s, got %s", "zipfile.zip", zipName[len(zipName)-1])
	}
}

func Test_generateFile(t *testing.T) {
	fileNames := []string{}
	config := config.Config{}
	guid := xid.New()
	os.Setenv("STORAGE", "../storage/")
	generateFile(config, []byte("test content"), "test/testfile", true, &fileNames, &guid)
	defer os.Remove(fileNames[0])
	_, err := os.Open(fileNames[0])
	if err != nil {
		t.Errorf("Couldn't open file, %s", "testfile")
	}
}
