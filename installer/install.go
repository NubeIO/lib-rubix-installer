package installer

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Install struct {
	Name        string `json:"name"`
	BuildName   string `json:"build_name"`
	Version     string `json:"version"`
	ServiceName string `json:"service_name"`
	Source      string `json:"source"`
}

type Response struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func (inst *App) InstallApp(app *Install) (*AppResponse, error) {
	var appName = app.Name
	var appBuildName = app.BuildName
	var version = app.Version
	var source = app.Source
	if appName == "" {
		return nil, errors.New("app name can not be empty")
	}
	if appBuildName == "" {
		return nil, errors.New("app build name can not be empty")
	}
	if version == "" {
		return nil, errors.New("app version can not be empty")
	}
	if source == "" {
		return nil, errors.New("app build source can not be empty, try: /data/tmp/tmp_1223/flow-framework.zip")
	}
	return inst.installApp(appName, appBuildName, version, source)
}

// InstallApp make all the required dirs and unzip build
//	zip, pass in the zip folder, or you can pass in a local path to param localZip
func (inst *App) installApp(appName, appBuildName, version string, source string) (*AppResponse, error) {
	// make the dirs
	err := inst.DirsInstallApp(appName, appBuildName, version)
	if err != nil {
		return nil, err
	}
	log.Infof("made all dirs for app:%s,  buildName:%s, version:%s", appName, appBuildName, version)
	dest := inst.getAppInstallPathAndVersion(appBuildName, version)
	log.Infof("app zip source:%s", source)
	log.Infof("app zip dest:%s", dest)
	// unzip the build to the app dir  /data/rubix-service/install/wires-build
	_, err = inst.unZip(source, dest) // unzip the build
	if err != nil {
		return nil, err
	}

	files, err := inst.listFiles(dest)
	if err != nil {
		return nil, err
	}
	if len(files) > 0 {
		for _, file := range files {
			existingFile := fmt.Sprintf("%s/%s", dest, file)
			newFile := fmt.Sprintf("%s/app", dest)
			log.Infof("RENAME BUILD-EXISTSING %s", existingFile)
			log.Infof("RENAME BUILD-NEW %s", newFile)
			if knownBuildNames(file) {
				err = inst.MoveFile(existingFile, newFile, true) // rename the build
				if err != nil {
					return nil, err
				}
			}

		}
	}

	return inst.ConfirmAppInstalled(appName, appBuildName), err

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
