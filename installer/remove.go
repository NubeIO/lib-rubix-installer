package installer

import (
	"fmt"
	fileutils "github.com/NubeIO/lib-dirs/dirs"
	"github.com/NubeIO/lib-systemctl-go/ctl"
	"github.com/NubeIO/lib-systemctl-go/systemctl"
)

type RemoveRes struct {
	DeleteAppDir         string `json:"delete_app_dir"`
	DeleteAppInstallDir  string `json:"delete_app_install_dir"`
	ServiceWasInstalled  string `json:"service_was_installed"`
	RemoveServiceErr     string `json:"remove_service_err,omitempty"`
	Stop                 bool   `json:"stop"`
	Disable              bool   `json:"disable"`
	DaemonReload         bool   `json:"daemon_reload"`
	RestartFailed        bool   `json:"restart_failed"`
	DeleteServiceFile    bool   `json:"delete_service_file"`
	DeleteServiceFileUsr bool   `json:"delete_service_file_usr"`
	Error                string `json:"error,omitempty"`
}

/*
- stop, disable service
- remove service file
*/

// UninstallApp full removal of an app, including removing the linux service
func (inst *App) UninstallApp(appName, serviceName string, deleteApp bool) (*RemoveRes, error) {
	service := ctl.New(serviceName, "")
	service.InstallOpts = ctl.InstallOpts{
		Options: systemctl.Options{Timeout: inst.DefaultTimeout},
	}
	remove, err := service.Remove()
	resp := &RemoveRes{
		ServiceWasInstalled:  remove.ServiceWasInstalled,
		Stop:                 remove.Stop,
		Disable:              remove.Disable,
		DaemonReload:         remove.DaemonReload,
		RestartFailed:        remove.RestartFailed,
		DeleteServiceFile:    remove.DeleteServiceFile,
		DeleteServiceFileUsr: remove.DeleteServiceFileUsr,
	}
	if err != nil {
		resp.RemoveServiceErr = err.Error()
		err = nil
	}

	if deleteApp {
		err := inst.RemoveApp(appName)
		var removeApp = "removed app from data dir ok"
		var removeAppInstall = "removed app from install dir ok"
		if err != nil {
			resp.Error = err.Error()
			removeApp = fmt.Sprintf("failed to delete app from data dir")
		}
		err = inst.RemoveAppInstall(appName)
		if err != nil {
			resp.Error = err.Error()
			removeAppInstall = fmt.Sprintf("failed to delete app from install dir")
		}
		resp.DeleteAppDir = removeApp
		resp.DeleteAppInstallDir = removeAppInstall
	} else {
		resp.DeleteAppDir = "app was not deleted"
		resp.DeleteAppInstallDir = "app install dir was not deleted"
	}
	return resp, nil
}

// RemoveApp delete app
func (inst *App) RemoveApp(appName string) error {
	return inst.RmRF(inst.getAppPath(appName))
}

// RemoveAppInstall delete app install path
func (inst *App) RemoveAppInstall(appName string) error {
	return inst.RmRF(inst.getAppInstallPath(appName))
}

// RmRF remove file and all files inside
func (inst *App) RmRF(path string) error {
	return fileutils.New().RmRF(path)
}

// Rm remove file
func (inst *App) Rm(path string) error {
	return fileutils.New().Rm(path)
}
