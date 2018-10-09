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
	return nil
}

func initializeZipWriter(file string) (*os.File, *zip.Writer, error) {
	zipFile, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}
	zipWriter := zip.NewWriter(zipFile)
	return zipFile, zipWriter, nil
}

func addToZip(writer *zip.Writer, file string) error {
	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	fileName := strings.Split(file, "/")
	fileWriter, err := writer.Create(strings.Join(fileName[3:], "/"))
	if err != nil {
		return err
	}
	_, err = fileWriter.Write(fileContent)
	if err != nil {
		return err
	}
	return nil
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

func generateFile(config config.Config, templates *template.Template, newPath string, projectPath string, writer *zip.Writer) error {
	if _, e := os.Stat(newPath); os.IsNotExist(e) {
		os.MkdirAll(filepath.Dir(newPath), os.ModePerm)
	}

	err := createDir(projectPath, newPath)
	if err != nil {
		return fmt.Errorf("unable to create path %s", err)
	}
	file, err := os.Create(projectPath + newPath)
	if err != nil {
		return fmt.Errorf("unable to create file: %s", err)
	}
	defer file.Close()

	err = templates.Lookup(newPath).Execute(file, config)
	if err != nil {
		return fmt.Errorf("unable to parse template: %s", err)
	}
	err = addToZip(writer, file.Name())
	if err != nil {
		return fmt.Errorf("unable to add the generated file to zip: %s", err)
	}
	return nil
}

func CreateFiles(config config.Config) (string, JSONerror) {
	box := templates.GetTemplates()
	temps := parseTemplates(box)
	guid := xid.New()
	projectPath := os.Getenv("STORAGE") + guid.String() + "/"
	err := os.MkdirAll(projectPath, os.ModePerm)
	if err != nil {
		return "", JSONerror{400, err.Error()}
	}
	zip, writer, err := initializeZipWriter(projectPath + config.Project.Name + ".zip")
	defer zip.Close()
	if err != nil {
		return "", JSONerror{400, err.Error()}
	}
	for _, templateName := range box.List() {
		file, _ := box.Open(templateName)
		fileStat, _ := file.Stat()

		if fileStat.IsDir() {
			continue
		}

		if strings.Contains(templateName, "partials/") {
			continue
		}

		err := generateFile(config, temps, templateName, projectPath, writer)
		if err != nil {
			return "", JSONerror{400, err.Error()}
		}
	}
	err = writer.Close()
	if err != nil {
		log.Fatal(err)
	}
	return zip.Name(), JSONerror{200, ""}
}
