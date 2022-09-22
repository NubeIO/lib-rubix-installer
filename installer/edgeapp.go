package installer

import (
	"errors"
	"github.com/NubeIO/lib-systemctl-go/systemctl"
	"github.com/NubeIO/lib-systemctl-go/systemd"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
)

type AppResponse struct {
	Name              string                     `json:"name"`
	Version           string                     `json:"version,omitempty"`
	State             *systemctl.SystemState     `json:"state,omitempty"`
	Error             string                     `json:"error,omitempty"`
	UninstallResponse *systemd.UninstallResponse `json:"remove_response"`
}

type Apps struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
	Path    string `json:"path,omitempty"`
}

type AppsStatus struct {
	Name        string                 `json:"name,omitempty"`
	Version     string                 `json:"version,omitempty"`
	ServiceName string                 `json:"service_name,omitempty"`
	State       *systemctl.SystemState `json:"state,omitempty"`
}

// ListApps apps by listed in the installation (/data/rubix-service/apps/install)
func (inst *App) ListApps() ([]Apps, error) {
	var apps []Apps
	var app Apps
	files, err := ioutil.ReadDir(inst.AppsInstallDir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		app.Name = inst.GetAppNameFromRepoName(file.Name())
		app.Version = inst.GetAppVersion(app.Name)
		app.Path = path.Join(inst.AppsInstallDir, file.Name())
		apps = append(apps, app)
	}
	return apps, err
}

// ListAppsStatus get all the apps by listed in the installation (/data/rubix-service/apps/install) dir and then check the service
func (inst *App) ListAppsStatus() ([]AppsStatus, error) {
	apps, err := inst.ListApps()
	if err != nil {
		return nil, err
	}
	var installedServices []AppsStatus
	for _, app := range apps {
		var installedService AppsStatus
		installedService.Name = app.Name
		installedService.Version = app.Version
		serviceName := inst.GetServiceNameFromAppName(app.Name)
		installedService.ServiceName = serviceName
		installed, err := inst.SystemCtl.State(serviceName)
		if err != nil {
			log.Errorf("service is not intalled: %s", serviceName)
		}
		installedService.State = &installed
		installedServices = append(installedServices, installedService)
	}
	return installedServices, nil
}

func (inst *App) ConfirmAppInstalled(appName string) (*AppResponse, error) {
	if appName == "" {
		return nil, errors.New(ErrEmptyAppName)
	}
	version := inst.GetAppVersion(appName)
	if version == "" {
		return nil, errors.New("failed to find app version")
	}
	ctl := systemctl.New(false, defaultTimeout)
	serviceName := inst.GetServiceNameFromAppName(appName)
	state, err := ctl.State(serviceName)
	if err != nil {
		return nil, err
	}
	return &AppResponse{
		Name:    appName,
		Version: version,
		State:   &state,
	}, err
}

func (inst *App) GetAppVersion(appName string) string {
	file := inst.GetAppInstallPath(appName)
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
			if CheckVersionBool(file.Name()) {
				return file.Name()
			}
		}
	}
	return ""
}
