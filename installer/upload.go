package installer

import (
	"errors"
	"fmt"
	fileutils "github.com/NubeIO/lib-dirs/dirs"
	log "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type Upload struct {
	Name      string                `json:"name"`
	BuildName string                `json:"build_name"`
	Version   string                `json:"version"`
	File      *multipart.FileHeader `json:"file"`
}

type UploadResponse struct {
	FileName     string `json:"file_name"`
	TmpFile      string `json:"tmp_file"`
	UploadedFile string `json:"uploaded_file"`
}

func (inst *App) AddUploadEdgeApp(app *Upload) (*AppResponse, error) {
	var file = app.File
	resp, err := inst.Upload(file)
	if err != nil {
		return nil, err
	}
	installApp, err := inst.InstallEdgeApp(&Install{
		Name:      app.Name,
		BuildName: app.BuildName,
		Version:   app.Version,
		Source:    resp.UploadedFile,
	})
	if err != nil {
		return nil, err
	}
	return installApp, nil
}

// Upload upload a build
func (inst *App) Upload(zip *multipart.FileHeader) (*UploadResponse, error) {
	// make the dirs
	var err error
	if err := inst.MakeTmpDir(); err != nil {
		return nil, err
	}
	var tmpDir string
	if tmpDir, err = inst.MakeTmpDirUpload(); err != nil {
		return nil, err
	}
	log.Infof("upload build to tmp dir:%s", tmpDir)

	// save app in tmp dir
	zipSource, err := inst.SaveUploadedFile(zip, tmpDir)
	if err != nil {
		return nil, err
	}
	return &UploadResponse{
		FileName:     zip.Filename,
		TmpFile:      tmpDir,
		UploadedFile: zipSource,
	}, err
}

func (inst *App) UploadServiceFile(app *Upload) (*UploadResponse, error) {
	var appName = app.Name
	var appBuildName = app.BuildName
	var version = app.Version
	var file = app.File
	return inst.uploadServiceFile(appName, appBuildName, version, file)
}

// uploadApp
func (inst *App) uploadServiceFile(appName, appBuildName, version string, file *multipart.FileHeader) (*UploadResponse, error) {
	// make the dirs
	var err error
	if filepath.Ext(file.Filename) != ".service" {
		return nil, errors.New(fmt.Sprintf("service file provided:%s, did not have correct file extension must be (.service)", file.Filename))
	}
	if err := inst.MakeTmpDir(); err != nil {
		return nil, err
	}
	var tmpDir string
	if tmpDir, err = inst.MakeTmpDirUpload(); err != nil {
		return nil, err
	}
	log.Infof("upload service to tmp dir:%s", tmpDir)
	log.Infof("app:%s buildName:%s version:%s", appName, appBuildName, version)
	// save app in tmp dir
	zipSource, err := inst.SaveUploadedFile(file, tmpDir)
	if err != nil {
		return nil, err
	}
	return &UploadResponse{
		TmpFile:      tmpDir,
		UploadedFile: zipSource,
	}, err
}

func (inst *App) unZip(source, destination string) ([]string, error) {
	source = filePath(source)
	destination = filePath(destination)
	return fileutils.New().UnZip(source, destination, os.FileMode(inst.FilePerm))
}

// SaveUploadedFile uploads the form file to specific dst.
// combination's of file name and the destination and will save file as: /data/my-file
// returns the filename and path as a string and any error
func (inst *App) SaveUploadedFile(file *multipart.FileHeader, dest string) (destination string, err error) {
	destination = fmt.Sprintf("%s/%s", dest, file.Filename)
	destination = filePath(destination)
	src, err := file.Open()
	if err != nil {
		return destination, err
	}
	defer src.Close()
	out, err := os.Create(destination)
	if err != nil {
		return destination, err
	}
	defer out.Close()
	_, err = io.Copy(out, src)
	return destination, err
}
