package installer

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Install struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	ServiceName  string `json:"service_name"`
	Source       string `json:"source"`
	DeleteAppDir bool   `json:"delete_app_dir"` // this will delete for example the db, plugins and config
}

type Response struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func (inst *App) InstallEdgeApp(app *Install) (*AppResponse, error) {
	if app == nil {
		return nil, errors.New("app install body can not be empty")
	}
	var appName = app.Name
	var version = app.Version
	var source = app.Source
	if appName == "" {
		return nil, errors.New("app name can not be empty")
	}
	if version == "" {
		return nil, errors.New("app version can not be empty")
	}
	if source == "" {
		return nil, errors.New("app build source can not be empty, try: /data/tmp/tmp_1223/flow-framework.zip")
	}
	return inst.installEdgeApp(appName, version, source, app.DeleteAppDir)
}

// InstallApp make all the required dirs and unzip build
//	zip, pass in the zip folder, or you can pass in a local path to param localZip
func (inst *App) installEdgeApp(appName, version, source string, deleteApp bool) (*AppResponse, error) {
	log.Infof("remove existing app from the install dir before the install is started")
	uninstallApp, err := inst.UninstallApp(appName, deleteApp)
	if err != nil {
		log.Errorf("remove app install dir:%s", err.Error())
	}
	// make the dirs
	err = inst.DirsInstallApp(appName, version)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("install edge app make dirs:%s", err.Error()))
	}
	log.Infof("made all dirs for app:%s, version:%s", appName, version)
	dest := inst.getAppInstallPathAndVersion(appName, version)
	log.Infof("app zip source:%s", source)
	log.Infof("app zip dest:%s", dest)
	// unzip the build to the app dir  /data/rubix-service/install/wires-build
	_, err = inst.unZip(source, dest) // unzip the build
	if err != nil {
		log.Errorf("install edge app unzip source:%s dest:%s err:%s", source, dest, err.Error())
		return nil, errors.New(fmt.Sprintf("install edge app unzip err:%s", err.Error()))
	}
	if appName != "rubix-wires" {
		files, err := inst.listFiles(dest)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("install edge app list files err:%s", err.Error()))
		}
		if len(files) > 0 {
			for _, file := range files {
				existingFile := fmt.Sprintf("%s/%s", dest, file)
				newFile := fmt.Sprintf("%s/app", dest)
				log.Infof("RENAME BUILD-EXISTSING %s", existingFile)
				log.Infof("RENAME BUILD-NEW %s", newFile)
				if knownBuildNames(file) {
					err = inst.MoveFile(existingFile, newFile, true) // rename the build
					os.Chmod(newFile, os.FileMode(inst.FilePerm))
					if err != nil {
						return nil, errors.New(fmt.Sprintf("install edge app rename file err:%s", err.Error()))
					}
				}
			}
		}
	}

	installed, err := inst.ConfirmAppInstalled(appName)
	if err != nil {
		return nil, err
	}
	if installed != nil {
		installed.RemoveRes = uninstallApp
	}
	return installed, err
}

func knownBuildNames(file string) bool {
	const (
		nubeio   = "nubeio"
		py       = "py"
		appAmd64 = "app-amd64"
		appArmv7 = "app-armv7"
	)
	if file == appArmv7 { // eg flow-framework
		return true
	}
	if file == appAmd64 { // eg flow-framework
		return true
	}
	if strings.Contains(file, nubeio) { // nubeio-rubix-app-lora-serial-py
		if strings.Contains(file, py) {
			return true
		}
	}
	return false
}
