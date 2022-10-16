package installer

import (
	"archive/zip"
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func timestamp() string {
	t := time.Now().Format("2006-01-02T15:04:05")
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
		return "", fmt.Errorf("failed to open zip %s for reading: %s", file.Name, err)
	}
	defer fileRead.Close()
	return file.Name, nil
}

func unzip(src, destination string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		err = extractZipFile(f, destination)
		if err != nil {
			return err
		}
	}
	return nil
}

func extractZipFile(f *zip.File, destination string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	fpath := filepath.Join(destination, f.Name)
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
	return nil
}

// MoveOneLevelInsideFileToOutside needs to move all these files:
// - /data/rubix-service/apps/install/wires-builds/v2.7.4/NubeIO-wires-builds-5082d47/README.md
// - /data/rubix-service/apps/install/wires-builds/v2.7.4/NubeIO-wires-builds-5082d47/rubix-wires/*
// To: /data/rubix-service/apps/install/wires-builds/v2.7.4
func MoveOneLevelInsideFileToOutside(file string) error { // /data/rubix-service/apps/install/wires-builds/v2.7.4
	fileInfo, err := os.Stat(file)
	if err != nil {
		return err
	}
	if fileInfo.IsDir() { // /data/rubix-service/apps/install/wires-builds/v2.7.4
		files, err := ioutil.ReadDir(file)
		if err != nil {
			return err
		}
		for _, f := range files {
			if f.IsDir() { // /data/rubix-service/apps/install/wires-builds/v2.7.4/NubeIO-wires-builds-5082d47
				file2 := path.Join(file, f.Name()) // /data/rubix-service/apps/install/wires-builds/v2.7.4/NubeIO-wires-builds-5082d47
				fileInfo2, err := os.Stat(file2)
				if err != nil {
					return err
				}
				if fileInfo2.IsDir() { // /data/rubix-service/apps/install/wires-builds/v2.7.4/NubeIO-wires-builds-5082d47
					files2, err := ioutil.ReadDir(file2)
					if err != nil {
						return err
					}
					for _, f2 := range files2 { // /data/rubix-service/apps/install/wires-builds/v2.7.4/NubeIO-wires-builds-5082d47/*
						file3 := path.Join(file2, f2.Name())
						fileInfo3, err := os.Stat(file3)
						if err != nil {
							return err
						}
						if fileInfo3.IsDir() {
							source := file3     // /data/rubix-service/apps/install/wires-builds/v2.7.4/NubeIO-wires-builds-5082d47/rubix-wires
							destination := file // /data/rubix-service/apps/install/wires-builds/v2.7.4
							err = fileutils.Copy(source, destination)
							if err != nil {
								return err
							}
						} else {
							destination := path.Join(file, fileInfo3.Name())
							// /data/rubix-service/apps/install/wires-builds/v2.7.4/NubeIO-wires-builds-5082d47/README.md => /data/rubix-service/apps/install/wires-builds/v2.7.4/README.md
							err = fileutils.Copy(file3, destination)
							if err != nil {
								return err
							}
						}
					}
				}
				err = os.RemoveAll(file2)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
