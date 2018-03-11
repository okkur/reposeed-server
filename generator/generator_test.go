package generator

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/okkur/reposeed-server/config"
)

func TestZipFiles(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "testfile")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	filenames := []string{tmpfile.Name()}
	filen := &filenames
	zipPath, err := ZipFiles("zipfile.zip", filen)
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
	filen := &fileNames
	config := config.Config{}
	generateFile(config, []byte("test content"), "testfile", true, filen)
	defer os.Remove("testfile")
	_, err := os.Open("testfile")
	if err != nil {
		t.Errorf("Couldn't open file, %s", "testfile")
	}
}
