package installer

import (
	"errors"
	"fmt"
	fileutils "github.com/NubeIO/lib-dirs/dirs"
	"github.com/NubeIO/lib-systemctl-go/systemctl"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type AppResponse struct {
	Name       string                 `json:"app"`
	Version    string                 `json:"version,omitempty"`
	IsAService bool                   `json:"is_service"`
	AppStatus  *systemctl.SystemState `json:"app_status,omitempty"`
	Error      string                 `json:"error,omitempty"`
}

var systemOpts = systemctl.Options{
	UserMode: false,
	Timeout:  defaultTimeout,
}

func (inst *App) ConfirmAppInstalled(appName, serviceName string) (*AppResponse, error) {
	if appName == "" {
		return nil, errors.New("app build/repo name can not be empty")
	}
	if serviceName == "" {
		return nil, errors.New("app service name can not be empty")
	}
	version := inst.GetAppVersion(appName)
	if version == "" {
		return nil, errors.New("failed to find app version")
	}
	installed, _ := inst.IsInstalled(serviceName, inst.DefaultTimeout)
	var isAService bool
	if installed != nil {
		isAService = installed.Is
	}
	return &AppResponse{
		Name:       appName,
		Version:    version,
		IsAService: isAService,
	}, nil

}

func (inst *App) ConfirmAppDir(appName string) bool {
	return fileutils.New().DirExists(fmt.Sprintf("%s/%s", inst.DataDir, appName))
}

func (inst *App) ConfirmAppInstallDir(appName string) bool {
	return fileutils.New().DirExists(fmt.Sprintf("%s/%s", inst.AppsInstallDir, appName))
}

func (inst *App) DirExists(dir string) bool {
	return fileutils.New().DirExists(dir)
}

func (inst *App) FileExists(dir string) bool {
	return fileutils.New().FileExists(dir)
}

func (inst *App) ConfirmStoreDir() bool {
	return fileutils.New().DirExists(inst.GetStoreDir())
}

func (inst *App) ConfirmStoreAppDir(appName string) bool {
	return fileutils.New().DirExists(fmt.Sprintf("%s/apps/%s", inst.GetStoreDir(), appName))
}

func (inst *App) ConfirmStoreAppVersionDir(appName, version string) bool {
	return fileutils.New().DirExists(fmt.Sprintf("%s/apps/%s/%s", inst.GetStoreDir(), appName, version))
}

func (inst *App) ConfirmServiceFile(serviceName string) bool {
	return fileutils.New().FileExists(fmt.Sprintf("%s/%s", inst.LibSystemPath, serviceName))
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

func (inst *App) listFiles(file string) ([]string, error) {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return nil, err
	}
	var dirContent []string
	if fileInfo.IsDir() {
		files, err := ioutil.ReadDir(file)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			dirContent = append(dirContent, file.Name())
		}
	}
	return dirContent, nil
}

type InstalledApps struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (inst *App) GetApps() ([]InstalledApps, error) {
	rootDir := inst.AppsInstallDir
	var files []InstalledApps
	app := InstalledApps{}
	err := filepath.WalkDir(rootDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && strings.Count(p, string(os.PathSeparator)) == 6 {
			parts := strings.Split(p, "/")
			if len(parts) >= 5 { // app name
				app.Name = parts[5]
			}
			if len(parts) >= 6 { // version
				app.Version = parts[6]
			}
			files = append(files, app)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
