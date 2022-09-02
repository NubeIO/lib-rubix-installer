package installer

import (
	"errors"
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/lib-systemctl-go/systemctl"
	"github.com/NubeIO/lib-systemctl-go/systemd"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

type AppResponse struct {
	Name              string                     `json:"app"`
	Version           string                     `json:"version,omitempty"`
	AppStatus         *systemctl.SystemState     `json:"app_status,omitempty"`
	Error             string                     `json:"error,omitempty"`
	UninstallResponse *systemd.UninstallResponse `json:"remove_response"`
}

type Apps struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
	Path    string `json:"path,omitempty"`
}

// ListApps apps by listed in the installation (/data/rubix-service/apps/install)
func (inst *App) ListApps() ([]Apps, error) {
	rootDir := inst.AppsInstallDir
	var apps []Apps
	var app Apps
	files, err := ioutil.ReadDir(rootDir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		app.Name = file.Name()
		app.Version = inst.GetAppVersion(app.Name)
		app.Path = fmt.Sprintf("%s/apps/%s", rootDir, file.Name())
		apps = append(apps, app)
	}
	return apps, err
}

type AppsStatus struct {
	AppName     string                 `json:"app_name,omitempty"`
	Version     string                 `json:"version,omitempty"`
	ServiceName string                 `json:"service_name,omitempty"`
	AppState    *systemctl.SystemState `json:"app_state,omitempty"`
}

// ListAppsStatus get all the apps by listed in the installation (/data/rubix-service/apps/install) dir and then check the service
func (inst *App) ListAppsStatus(appServiceMapping map[string]string) ([]AppsStatus, error) {
	apps, err := inst.ListApps()
	if err != nil {
		return nil, err
	}
	var installedServices []AppsStatus
	for _, app := range apps {
		var installedService AppsStatus
		installedService.AppName = app.Name
		installedService.Version = app.Version
		serviceName, exist := appServiceMapping[app.Name]
		if exist {
			installedService.ServiceName = serviceName
			installed, err := inst.SystemCtl.State(serviceName)
			if err != nil {
				log.Errorf("service is not intalled: %s", serviceName)
			}
			installedService.AppState = &installed
		}
		installedServices = append(installedServices, installedService)
	}
	return installedServices, nil
}

func (inst *App) ConfirmAppInstalled(appName string, serviceFileName string) (*AppResponse, error) {
	if appName == "" {
		return nil, errors.New("app name can not be empty")
	}
	version := inst.GetAppVersion(appName)
	if version == "" {
		return nil, errors.New("failed to find app version")
	}
	ctl := systemctl.New(false, defaultTimeout)
	state, err := ctl.State(serviceFileName)
	if err != nil {
		return nil, err
	}
	return &AppResponse{
		Name:      appName,
		Version:   version,
		AppStatus: &state,
	}, err
}

func (inst *App) ConfirmAppDir(appName string) bool {
	return fileutils.DirExists(fmt.Sprintf("%s/%s", inst.DataDir, appName))
}

func (inst *App) ConfirmAppInstallDir(appName string) bool {
	return fileutils.DirExists(fmt.Sprintf("%s/%s", inst.AppsInstallDir, appName))
}

func (inst *App) DirExists(dir string) bool {
	return fileutils.DirExists(dir)
}

func (inst *App) FileExists(dir string) bool {
	return fileutils.FileExists(dir)
}

func (inst *App) ConfirmStoreDir() bool {
	return fileutils.DirExists(inst.GetStoreDir())
}

func (inst *App) ConfirmStoreAppDir(appName string) bool {
	return fileutils.DirExists(fmt.Sprintf("%s/apps/%s", inst.GetStoreDir(), appName))
}

func (inst *App) ConfirmStoreAppVersionDir(appName, version string) bool {
	return fileutils.DirExists(fmt.Sprintf("%s/apps/%s/%s", inst.GetStoreDir(), appName, version))
}

func (inst *App) GetAppVersion(appName string) string {
	file := fmt.Sprintf("%s/%s", inst.AppsInstallDir, appName)
	fileInfo, err := os.Stat(file)
	if err != nil {
		return ""
	}
	if fileInfo.IsDir() {
		files, err := ioutil.ReadDir(file)
		if err != nil {
			return ""
		}
		for _, file := range files {
			if checkVersionBool(file.Name()) {
				return file.Name()
			}
		}
	}
	return ""
}
