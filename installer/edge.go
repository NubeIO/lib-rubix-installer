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
	"strings"
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

// ListAppsAndService get all the apps by listed in the installation (/data/rubix-service/apps/install) dir and then check the service
func (inst *App) ListAppsAndService() ([]InstalledServices, error) {
	apps, err := inst.ListApps()
	if err != nil {
		return nil, err
	}
	var installedServices []InstalledServices
	var installedService InstalledServices
	for _, app := range apps {
		name, err := inst.GetNubeServiceFileName(app.Name)
		if err != nil {
			return nil, err
		}
		systemCtl := systemctl.New(false, inst.DefaultTimeout)
		installedService.AppName = app.Name
		installedService.ServiceName = name
		installed, err := systemCtl.State(name)
		if err != nil {
			log.Errorf("service is not isntalled: %s", name)
		}
		installedService.AppStatus = installed
		installedServices = append(installedServices, installedService)
	}
	return installedServices, nil
}

type InstalledServices struct {
	AppName     string                `json:"app_name,omitempty"`
	ServiceName string                `json:"service_name,omitempty"`
	AppStatus   systemctl.SystemState `json:"app_status,omitempty"`
}

// ListNubeServices list all the services by filtering all the service files with name nubeio
func (inst *App) ListNubeServices() ([]InstalledServices, error) {
	files, err := inst.ListNubeServiceFiles()
	var installedServices []InstalledServices
	var installedService InstalledServices
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		systemCtl := systemctl.New(false, inst.DefaultTimeout)
		installedService.ServiceName = file
		installed, err := systemCtl.State(file)
		if err != nil {
			log.Errorf("service is not isntalled: %s", file)
		}
		installedService.AppStatus = installed
		installedServices = append(installedServices, installedService)
	}
	return installedServices, err
}

func (inst *App) ListNubeServiceFiles() ([]string, error) {
	var resp []string
	files, err := inst.listFiles("/lib/systemd/system")
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if strings.Contains(file, "nubeio") {
			resp = append(resp, file)
		}
	}
	return resp, err
}

func (inst *App) GetNubeServiceFileName(appName string) (string, error) {
	var resp string
	files, err := inst.listFiles("/lib/systemd/system")
	if err != nil {
		return "", err
	}
	for _, file := range files {
		if strings.Contains(file, "nubeio") {
			if strings.Contains(file, appName) {
				resp = file
			}
		}
	}
	return resp, err
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
