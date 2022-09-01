package installer

import (
	"errors"
	"fmt"
	"github.com/NubeIO/lib-systemctl-go/systemd"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Install struct {
	Name             string `json:"name"`
	ServiceName      string `json:"service_name"`
	Version          string `json:"version"`
	Source           string `json:"source"`
	DeleteAppDataDir bool   `json:"delete_app_data_dir"` // this will delete for example the db, plugins and config
}

type Response struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func (inst *App) InstallEdgeApp(app *Install) (*AppResponse, error) {
	if app == nil {
		return nil, errors.New("app install body can not be empty")
	}
	if app.Name == "" {
		return nil, errors.New("app name can not be empty")
	}
	if app.Version == "" {
		return nil, errors.New("app version can not be empty")
	}
	if app.ServiceName == "" {
		return nil, errors.New("app service_name can not be empty")
	}
	if app.Source == "" {
		return nil, errors.New("app build source can not be empty, try: /data/tmp/tmp_1223/flow-framework.zip")
	}

	log.Infof("remove existing app from the install dir before the install is started")
	systemdService := systemd.New(app.ServiceName, false, inst.DefaultTimeout)
	uninstallResponse := systemdService.Uninstall()

	err := inst.DirsInstallApp(app.Name, app.Version)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("install edge app make dirs: %s", err.Error()))
	}
	log.Infof("made all dirs succefully for app: %s, version: %s", app.Name, app.Version)
	destination := inst.getAppInstallPathAndVersion(app.Name, app.Version)
	log.Infof("app zip source: %s", app.Source)
	log.Infof("app zip destination: %s", destination)
	// unzip the build to the app dir  /data/rubix-service/install/wires-build
	_, err = inst.unzip(app.Source, destination) // unzip the build
	if err != nil {
		log.Errorf("install edge app unzip source: %s dest: %s err: %s", app.Source, destination, err.Error())
		return nil, errors.New(fmt.Sprintf("install edge app unzip err: %s", err.Error()))
	}
	if app.Name != "rubix-wires" {
		files, err := inst.listFiles(destination)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("install edge app list files err: %s", err.Error()))
		}
		if len(files) > 0 {
			for _, file := range files {
				existingFile := fmt.Sprintf("%s/%s", destination, file)
				newFile := fmt.Sprintf("%s/app", destination)
				log.Infof("RENAME BUILD-EXISTSING %s", existingFile)
				log.Infof("RENAME BUILD-NEW %s", newFile)
				if knownBuildNames(file) {
					err = inst.MoveFile(existingFile, newFile, true) // rename the build
					os.Chmod(newFile, os.FileMode(inst.FilePerm))
					if err != nil {
						return nil, errors.New(fmt.Sprintf("install edge app rename file err: %s", err.Error()))
					}
				}
			}
		}
	}

	installed, err := inst.ConfirmAppInstalled(app.Name, app.ServiceName)
	if err != nil {
		return nil, err
	}
	if installed != nil {
		installed.UninstallResponse = uninstallResponse
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
