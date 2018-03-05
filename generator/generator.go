package generator

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/alecthomas/template"
	"github.com/gobuffalo/packr"
	"github.com/okkur/reposeed-server/config"
)

func generateFile(config config.Config, fileContent []byte, newPath string, overwrite bool, zip zip.Writer, title string) error {
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

	ZipFiles(tmpfile.Name(), zip)
	return nil
}

func ZipFiles(file string, zipWriter zip.Writer) error {
	fmt.Println(file + "----------------")
	zipfile, err := os.Open(file)
	if err != nil {
		log.Panicln(err)
		return err
	}
	defer zipfile.Close()

	info, err := zipfile.Stat()
	if err != nil {
		log.Panicln(err)
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		log.Panicln(err)
		return err
	}

	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		log.Panicln(err)
		return err
	}
	// reader, err := zip.OpenReader(zipfile.Name())

	_, err = io.Copy(writer, zipfile)
	if err != nil {
		log.Panicln(err)
		return err
	}
	return nil
}

func CreateFiles(config config.Config, path string, title string) *zip.Writer {
	newfile, err := os.Create(title + ".zip")
	if err != nil {
		log.Fatal(err)
	}
	defer newfile.Close()
	zw := zip.NewWriter(newfile)
	box := packr.NewBox(path)
	templatesName := box.List()
	// log.Println(templatesName)
	for _, templateName := range templatesName {
		file, _ := box.Open(templateName)
		fileStat, _ := file.Stat()
		fileContent := box.Bytes(templateName)

		if !fileStat.IsDir() {
			err := generateFile(config, fileContent, templateName, true, *zw, title)
			if err != nil {
				log.Println(err)
			}
		}
	}
	return zw
}
