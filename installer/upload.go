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
	Name    string                `json:"name"`
	Version string                `json:"version"`
	Product string                `json:"product"`
	Arch    string                `json:"arch"`
	File    *multipart.FileHeader `json:"file"`
}

type UploadResponse struct {
	FileName     string `json:"file_name,omitempty"`
	TmpFile      string `json:"tmp_file,omitempty"`
	UploadedFile string `json:"uploaded_file,omitempty"`
}

func (inst *App) CompareBuildToArch(buildZipName, productType string) error {
	check := inst.GetZipBuildDetails(buildZipName)
	productInfo, err := inst.GetProduct() // same api as 0.0.0.0:1661/api/system/product check the arch type
	if err != nil {
		log.Errorf("upload build get product type err: %s", err.Error())
		return err
	}
	if check.MatchedArch != productInfo.Arch {
		errMsg := fmt.Sprintf("upload build incorrect arch type was uploaded build arch: %s host arch: %s", check.MatchedArch, productInfo.Arch)
		log.Errorf(errMsg)
		return errors.New(errMsg)
	}
	if productType != productInfo.Product {
		errMsg := fmt.Sprintf("upload build incorrect product type was uploaded build arch: %s host product: %s", productType, productInfo.Product)
		log.Errorf(errMsg)
		return errors.New(errMsg)
	}
	return nil
}

func (inst *App) AddUploadEdgeApp(app *Upload) (*AppResponse, error) {
	var file = app.File
	var appName = app.Name
	var version = app.Version
	var archType = app.Arch
	var productType = app.Product
	if appName == "" {
		return nil, errors.New("app name can not be empty")
	}
	if version == "" {
		return nil, errors.New("app version can not be empty")
	}
	if archType == "" {
		return nil, errors.New("arch type can not be empty, try armv7 amd64")
	}
	if productType == "" {
		return nil, errors.New("product type can not be empty, try RubixCompute, RubixComputeIO, RubixCompute5, Server, Edge28, Nuc")
	}
	if app.Name == "rubix-wires" { // wires don't care about the arch
	} else {
		err := inst.CompareBuildToArch(file.Filename, productType)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("upload edge app check arch err: %s", err.Error()))
		}
	}
	resp, err := inst.Upload(file) // save app in tmp dir
	if err != nil {
		return nil, errors.New(fmt.Sprintf("upload edge app unzip err: %s", err.Error()))
	}
	serviceFile, err := inst.GetNubeServiceFileName(appName)
	if err != nil {
		serviceFile = fmt.Sprintf("nubeio-%s.service", appName)
	}
	log.Infof("try and stop service: %s", serviceFile)
	action, err := inst.Ctl.CtlAction("stop", serviceFile, inst.DefaultTimeout) // try and stop the app as when updating and trying to delete the existing instance linux can throw and error saying `file is busy`
	if action != nil {
		if action.Ok {
			log.Infof("failed to stop service: %s", serviceFile)
		} else {
			log.Infof("was able to stop service: %s", serviceFile)
		}
	} else {
		log.Infof("was able to stop service: %s", serviceFile)
	}
	installApp, err := inst.InstallEdgeApp(&Install{
		Name:    app.Name,
		Version: app.Version,
		Source:  resp.UploadedFile,
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
	log.Infof("upload build to tmp dir: %s", tmpDir)
	zipSource, err := inst.SaveUploadedFile(zip, tmpDir) // save app in tmp dir
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
	var version = app.Version
	var file = app.File
	return inst.uploadServiceFile(appName, version, file)
}

// uploadApp
func (inst *App) uploadServiceFile(appName, version string, file *multipart.FileHeader) (*UploadResponse, error) {
	// make the dirs
	var err error
	if filepath.Ext(file.Filename) != ".service" {
		return nil, errors.New(fmt.Sprintf("service file provided: %s, did not have correct file extension must be (.service)", file.Filename))
	}
	if err := inst.MakeTmpDir(); err != nil {
		return nil, err
	}
	var tmpDir string
	if tmpDir, err = inst.MakeTmpDirUpload(); err != nil {
		return nil, err
	}
	log.Infof("upload service to tmp dir: %s", tmpDir)
	log.Infof("app: %s version: %s", appName, version)
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
	fmt.Println("SaveUploadedFile dest", destination)
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
