package installer

import (
	log "github.com/sirupsen/logrus"
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
	_, err = inst.unZip(source, dest)
	if err != nil {
		return nil, err
	}
	return inst.ConfirmAppInstalled(appName, appBuildName), err

}
