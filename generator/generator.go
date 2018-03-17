package generator

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/okkur/reposeed/cmd/reposeed/config"
	templates "github.com/okkur/reposeed/cmd/reposeed/templates"
)

type JSONerror struct {
	Code    int
	Message string
}

func generateFile(config config.Config, fileContent []byte, newPath string, overwrite bool, fileNames *[]string) error {
	tmpfile, err := ioutil.TempFile("", "template")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.Write(fileContent); err != nil {
		log.Fatal(err)
	}
	if _, e := os.Stat(newPath); os.IsNotExist(e) {
		os.MkdirAll(filepath.Dir(newPath), os.ModePerm)
	}

	if !overwrite {
		if _, e := os.Stat(newPath); !os.IsNotExist(e) {
			return fmt.Errorf("file %s not overwritten", newPath)
		}
	}

	file, err := os.Create(newPath)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("unable to create file: %s", err)
	}

	temp, err := template.ParseFiles(tmpfile.Name())
	if err != nil {
		return fmt.Errorf("unable to parse file: %s", err)
	}

	err = temp.Execute(file, config)
	if err != nil {
		return fmt.Errorf("unable to parse template: %s", err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
	*fileNames = append(*fileNames, file.Name())
	return nil
}

func ZipFiles(file string, fileNames *[]string, storagePath string) (string, error) {
	outFile, err := os.Create(storagePath + file)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	zipWriter := zip.NewWriter(outFile)
	for _, file := range *fileNames {
		fileContent, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		fileWriter, err := zipWriter.Create(file)
		if err != nil {
			log.Fatal(err)
		}
		_, err = fileWriter.Write(fileContent)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = zipWriter.Close()
	if err != nil {
		log.Fatal(err)
	}
	return outFile.Name(), nil
}

func CreateFiles(config config.Config, path string, title string, storagePath string) (string, JSONerror) {
	box := templates.GetTemplates()
	templatesName := box.List()
	filesNames := []string{}
	for _, templateName := range templatesName {
		file, _ := box.Open(templateName)
		fileStat, _ := file.Stat()
		fileContent := box.Bytes(templateName)

		if !fileStat.IsDir() {
			err := generateFile(config, fileContent, templateName, true, &filesNames)
			if err != nil {
				return "", JSONerror{400, err.Error()}
			}
		}
	}
	zipName, err := ZipFiles(title+".zip", &filesNames, storagePath)
	if err != nil {
		return "", JSONerror{400, err.Error()}
	}
	for _, filesName := range filesNames {
		os.Remove(filesName)
	}
	return zipName, JSONerror{200, ""}
}
