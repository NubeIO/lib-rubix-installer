package installer

import (
	"errors"
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/lib-systemctl-go/systemd"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
)

type Install struct {
	Name                            string `json:"name"`
	Version                         string `json:"version"`
	Source                          string `json:"source"`
	MoveExtractedFileToNameApp      bool   `json:"move_extracted_file_to_name_app"`
	MoveOneLevelInsideFileToOutside bool   `json:"move_one_level_inside_file_to_outside"`
	DeleteAppDataDir                bool   `json:"delete_app_data_dir"` // this will delete for example the db, plugins and config
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
		return nil, errors.New(ErrEmptyAppName)
	}
	if app.Version == "" {
		return nil, errors.New(ErrEmptyAppVersion)
	}
	if app.Source == "" {
		return nil, errors.New("app build source can not be empty, try: /data/tmp/tmp_1234/flow-framework.zip")
	}

	log.Infof("remove existing app from the install dir before the install is started...")
	serviceName := inst.GetServiceNameFromAppName(app.Name)
	systemdService := systemd.New(serviceName, false, inst.DefaultTimeout)
	uninstallResponse := systemdService.Uninstall()

	err := inst.CreateInstallAppDirs(app.Name, app.Version)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("install edge app make dirs: %s", err.Error()))
	}
	log.Infof("made all dirs succefully for app: %s, version: %s", app.Name, app.Version)
	destination := inst.GetAppInstallPathWithVersionPath(app.Name, app.Version)
	log.Infof("app zip source: %s", app.Source)
	log.Infof("app zip destination: %s", destination)
	// unzip the build to the app dir  /data/rubix-service/install/apps/<name>/<version>
	_, err = inst.unzip(app.Source, destination) // unzip the build
	if err != nil {
		log.Errorf("install edge app unzip source: %s, dest: %s, err: %s", app.Source, destination, err.Error())
		return nil, errors.New(fmt.Sprintf("install edge app unzip err: %s", err.Error()))
	}
	// rename the extracted file into app, it's only for those apps which is not frontend and executable
	if app.MoveExtractedFileToNameApp {
		files, err := fileutils.ListFiles(destination)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("install edge app list files err: %s", err.Error()))
		}
		if len(files) > 0 {
			existingFile := path.Join(destination, files[0])
			newFile := path.Join(destination, "app")
			log.Infof("Existing file: %s renaming into: %s", existingFile, newFile)
			if knownBuildNames(files[0]) {
				err = fileutils.MoveFile(existingFile, newFile) // rename the build
				if err != nil {
					return nil, errors.New(fmt.Sprintf("install edge app rename file err: %s", err.Error()))
				}
				err = os.Chmod(newFile, os.FileMode(inst.FileMode))
				if err != nil {
					return nil, errors.New(fmt.Sprintf("install edge app giving permission err: %s", err.Error()))
				}
			}
		}
	}

	if app.MoveOneLevelInsideFileToOutside {
		err = MoveOneLevelInsideFileToOutside(destination)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("install edge app move one level inside file to outside: %s", err.Error()))
		}
	}

	installed, err := inst.ConfirmAppInstalled(app.Name)
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
