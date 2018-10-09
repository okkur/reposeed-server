package generator

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gobuffalo/packr"
	"github.com/okkur/reposeed/cmd/reposeed/config"
	templates "github.com/okkur/reposeed/cmd/reposeed/templates"
	"github.com/rs/xid"
)

type JSONerror struct {
	Code    int
	Message string
}

func createDir(storagePath string, filePath string) error {
	dir := strings.Split(filePath, "/")
	if len(dir) > 1 {
		dir = dir[:len(dir)-1]
		path := strings.Join(dir, "/")
		err := os.MkdirAll(storagePath+path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("unable to create path: %s", storagePath+path)
		}
		return nil
	}
	err := os.MkdirAll(storagePath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("unable to create path: %s", storagePath)
	}
	return nil
}

func ZipFiles(file string, fileNames *[]string, storagePath string, uuid string) (string, error) {
	zipFile, err := os.Create(storagePath + uuid + "/" + file)
	if err != nil {
		log.Fatal(err)
	}
	defer zipFile.Close()
	zipWriter := zip.NewWriter(zipFile)
	for _, generatedFile := range *fileNames {
		fileContent, err := ioutil.ReadFile(generatedFile)
		if err != nil {
			log.Fatal(err)
		}
		fileName := strings.Split(generatedFile, "/")
		fileWriter, err := zipWriter.Create(strings.Join(fileName[3:], "/"))
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
	return zipFile.Name(), nil
}

func parseTemplates(box packr.Box) *template.Template {
	templatesName := box.List()
	templates := &template.Template{}
	for _, templateName := range templatesName {
		templateFile, err := box.Open(templateName)
		if err != nil {
			log.Fatalf("could not open the template file: %s", templateName)
		}
		defer templateFile.Close()
		templateContent := box.String(templateName)
		templates.New(templateName).Parse(templateContent)
	}
	return templates
}

func generateFile(config config.Config, templates *template.Template, newPath string, guid *xid.ID, fileNames *[]string) error {
	if _, e := os.Stat(newPath); os.IsNotExist(e) {
		os.MkdirAll(filepath.Dir(newPath), os.ModePerm)
	}

	projectsPath := os.Getenv("STORAGE") + guid.String() + "/"
	err := createDir(projectsPath, newPath)
	if err != nil {
		return fmt.Errorf("unable to create path %s", err)
	}
	file, err := os.Create(projectsPath + newPath)
	if err != nil {
		return fmt.Errorf("unable to create file: %s", err)
	}
	defer file.Close()

	err = templates.Lookup(newPath).Execute(file, config)
	if err != nil {
		return fmt.Errorf("unable to parse template: %s", err)
	}
	*fileNames = append(*fileNames, file.Name())
	return nil
}

func CreateFiles(config config.Config, title string, storagePath string) (string, JSONerror) {
	box := templates.GetTemplates()
	temps := parseTemplates(box)
	guid := xid.New()
	filesNames := []string{}
	for _, templateName := range box.List() {
		file, _ := box.Open(templateName)
		fileStat, _ := file.Stat()

		if fileStat.IsDir() {
			continue
		}

		if strings.Contains(templateName, "partials/") {
			continue
		}

		err := generateFile(config, temps, templateName, &guid, &filesNames)
		if err != nil {
			return "", JSONerror{400, err.Error()}
		}
	}
	zipName, err := ZipFiles(title+".zip", &filesNames, storagePath, guid.String())
	if err != nil {
		return "", JSONerror{400, err.Error()}
	}
	return zipName, JSONerror{200, ""}
}
