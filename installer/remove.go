package installer

import (
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/lib-systemctl-go/systemd"
)

type RemoveRes struct {
	ServiceWasInstalled bool   `json:"service_was_installed"`
	Stop                bool   `json:"stop"`
	DaemonReload        bool   `json:"daemon_reload"`
	UnlinkServiceFile   bool   `json:"unlink_service_file"`
	DeleteServiceFile   bool   `json:"delete_service_file"`
	DeleteAppInstallDir string `json:"delete_app_install_dir"`
	DeleteAppDataDir    string `json:"delete_app_data_dir"`
	Error               string `json:"error,omitempty"`
}

// UninstallApp full removal of an app, including removing the linux service
func (inst *App) UninstallApp(appName string, deleteAppDataDir bool) (*RemoveRes, error) {
	serviceName := inst.setServiceFileName(appName) // TODO: we don't use this convention, we directly pass service_name
	systemdService := systemd.New(serviceName, false, inst.DefaultTimeout)
	remove := systemdService.Remove()
	resp := &RemoveRes{
		ServiceWasInstalled: remove.ServiceWasInstalled,
		Stop:                remove.Stop,
		DaemonReload:        remove.DaemonReload,
		UnlinkServiceFile:   remove.UnlinkServiceFile,
		DeleteServiceFile:   remove.DeleteServiceFile,
	}

	var deleteInstallDir = "deleted app from install dir"
	err := inst.DeleteAppInstall(appName)
	if err != nil {
		resp.Error = err.Error()
		deleteInstallDir = fmt.Sprintf("failed to delete app from install dir")
	}
	resp.DeleteAppInstallDir = deleteInstallDir

	if deleteAppDataDir { // delete app from app install dir
		var removeApp = "deleted app data dir"
		err := inst.DeleteAppData(appName)
		if err != nil {
			resp.Error = err.Error()
			removeApp = fmt.Sprintf("failed to delete app data dir")
		}
		resp.DeleteAppDataDir = removeApp
	} else {
		resp.DeleteAppDataDir = "app data is not deleted"
	}
	return resp, nil
}

func (inst *App) DeleteAppData(appName string) error {
	return fileutils.RmRF(inst.getAppDataPath(appName))
}

func (inst *App) DeleteAppInstall(appName string) error {
	return fileutils.RmRF(inst.getAppInstallPath(appName))
}
