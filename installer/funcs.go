package installer

import (
	"archive/zip"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (inst *App) listFiles(file string) ([]string, error) {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return nil, err
	}
	var dirContent []string
	if fileInfo.IsDir() {
		files, err := ioutil.ReadDir(file)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			dirContent = append(dirContent, file.Name())
		}
	}
	return dirContent, nil
}

func timestamp() string {
	t := time.Now().Format("2006-01-02 15:04:05")
	return t
}

func readZip(src string) ([]string, error) {
	read, err := zip.OpenReader(src)
	if err != nil {
		msg := "Failed to open: %s"
		log.Fatalf(msg, err)
	}
	defer read.Close()
	var resp []string
	for _, file := range read.File {
		fileName, err := checkZipContents(file)
		if err != nil {
			return nil, err
		}
		resp = append(resp, fileName)
	}
	return resp, nil
}

func checkZipContents(file *zip.File) (string, error) {
	fileRead, err := file.Open()
	if err != nil {
		msg := "failed to open zip %s for reading: %s"
		return "", fmt.Errorf(msg, file.Name, err)
	}
	defer fileRead.Close()
	return file.Name, nil
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()
	// TODO: loop defer
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
		} else {
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}
			err = os.MkdirAll(fdir, f.Mode())
			if err != nil {
				log.Fatal(err)
				return err
			}
			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
