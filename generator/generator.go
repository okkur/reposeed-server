package generator

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/gobuffalo/packr"
	"github.com/okkur/reposeed-server/config"
)

func generateFile(config config.Config, fileContent []byte, newPath string, overwrite bool, title string, fileNames *[]string) error {
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

func ZipFiles(file string, fileNames *[]string) (string, error) {
	outFile, err := os.Create("./storage/zips/" + file)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	zipWriter := zip.NewWriter(outFile)
	for _, file := range *fileNames {
		// fmt.Println(file)
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

func CreateFiles(config config.Config, path string, title string) string {
	box := packr.NewBox(path)
	templatesName := box.List()
	filesNames := []string{}
	filen := &filesNames
	for _, templateName := range templatesName {
		file, _ := box.Open(templateName)
		fileStat, _ := file.Stat()
		fileContent := box.Bytes(templateName)

		if !fileStat.IsDir() {
			err := generateFile(config, fileContent, templateName, true, title, filen)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	zipName, err := ZipFiles(title+".zip", filen)
	if err != nil {
		log.Fatal(err)
	}
	for _, filesName := range filesNames {
		os.Remove(filesName)
	}
	return zipName
}
