package installer

import (
	"errors"
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	log "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
)

type Upload struct {
	Name              string                `json:"name"`
	Version           string                `json:"version"`
	Product           string                `json:"product"`
	Arch              string                `json:"arch"`
	DoNotValidateArch bool                  `json:"do_not_validate_arch"`
	File              *multipart.FileHeader `json:"file"`
}

type UploadResponse struct {
	FileName     string `json:"file_name,omitempty"`
	TmpFile      string `json:"tmp_file,omitempty"`
	UploadedFile string `json:"uploaded_file,omitempty"`
}

func (inst *App) CompareBuildToArch(buildZipName, productType string) error {
	zipBuildInfo := inst.GetZipBuildDetails(buildZipName)
	productInfo, err := inst.GetProduct() // same api as 0.0.0.0:1661/api/system/product zipBuildInfo the arch type
	if err != nil {
		log.Errorf("upload build get product type err: %s", err.Error())
		return err
	}
	if zipBuildInfo.MatchedArch != productInfo.Arch {
		errMsg := fmt.Sprintf("upload build incorrect arch type, was uploaded build arch: %s & host arch: %s", zipBuildInfo.MatchedArch, productInfo.Arch)
		log.Errorf(errMsg)
		return errors.New(errMsg)
	}
	if productType != productInfo.Product {
		errMsg := fmt.Sprintf("upload build incorrect product type, was uploaded build arch: %s & host product: %s", productType, productInfo.Product)
		log.Errorf(errMsg)
		return errors.New(errMsg)
	}
	return nil
}

func (inst *App) UploadEdgeApp(app *Upload) (*AppResponse, error) {
	if app.Name == "" {
		return nil, errors.New("app name can not be empty")
	}
	if app.Version == "" {
		return nil, errors.New("app version can not be empty")
	}
	if app.Arch == "" {
		return nil, errors.New("arch type can not be empty, try armv7 amd64")
	}
	if app.Product == "" {
		return nil, errors.New("product type can not be empty, try RubixCompute, RubixComputeIO, RubixCompute5, Server, Edge28, Nuc")
	}
	if !app.DoNotValidateArch { // wires don't care about the arch
		err := inst.CompareBuildToArch(app.File.Filename, app.Product)
		if err != nil {
			return nil, err
		}
	}
	resp, err := inst.Upload(app.File) // save app in tmp dir
	if err != nil {
		return nil, errors.New(fmt.Sprintf("upload edge app unzip err: %s", err.Error()))
	}
	serviceName := inst.GetServiceNameFromAppName(app.Name)
	err = inst.SystemCtl.Stop(serviceName) // try and stop the app as when updating and trying to delete the existing instance linux can throw and error saying `file is busy`
	if err != nil {
		log.Infof("was able to stop service: %s", serviceName)
	}
	response, err := inst.InstallEdgeApp(&Install{
		Name:    app.Name,
		Version: app.Version,
		Source:  resp.UploadedFile,
	})
	return response, err
}

// Upload upload a build
func (inst *App) Upload(zip *multipart.FileHeader) (*UploadResponse, error) {
	if err := inst.MakeTmpDir(); err != nil {
		return nil, err
	}
	tmpDir, err := inst.MakeTmpDirUpload()
	if err != nil {
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
	}, nil
}

func (inst *App) UploadServiceFile(app *Upload) (*UploadResponse, error) {
	var appName = app.Name
	var version = app.Version
	var file = app.File
	return inst.uploadServiceFile(appName, version, file)
}

func (inst *App) uploadServiceFile(appName, version string, file *multipart.FileHeader) (*UploadResponse, error) {
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

func (inst *App) unzip(source, destination string) ([]string, error) {
	return fileutils.UnZip(source, destination, os.FileMode(inst.FileMode))
}

// SaveUploadedFile uploads the form file to specific dst.
// combination's of file name and the destination and will save file as: /data/my-file
// returns the filename and path as a string and any error
func (inst *App) SaveUploadedFile(file *multipart.FileHeader, destination string) (uploadedFile string, err error) {
	uploadedFile = path.Join(destination, file.Filename)
	fmt.Println("SaveUploadedFile destination", uploadedFile)
	src, err := file.Open()
	if err != nil {
		return uploadedFile, err
	}
	defer src.Close()
	out, err := os.Create(uploadedFile)
	if err != nil {
		return uploadedFile, err
	}
	defer out.Close()
	_, err = io.Copy(out, src)
	return uploadedFile, err
}
