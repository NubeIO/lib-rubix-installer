package installer

import (
	"github.com/NubeIO/lib-uuid/uuid"
	"path"
)

// GetAppDataPath get the full app install path => /data/rubix-wires
func (inst *App) GetAppDataPath(appName string) string {
	dataDirName := inst.GetDataDirNameFromAppName(appName)
	return path.Join(inst.DataDir, dataDirName)
}

// GetAppDataDataPath get the full app install path => /data/rubix-wires/data
func (inst *App) GetAppDataDataPath(appName string) string {
	dataDirName := inst.GetDataDirNameFromAppName(appName)
	return path.Join(inst.DataDir, dataDirName, "data")
}

// GetAppDataConfigPath get the full app path =>  /data/rubix-wires/config
func (inst *App) GetAppDataConfigPath(appName string) string {
	dataDirName := inst.GetDataDirNameFromAppName(appName)
	return path.Join(inst.DataDir, dataDirName, "config")
}

// GetAppInstallPath get the full app install path and version => /data/rubix-service/apps/install/wires-builds
func (inst *App) GetAppInstallPath(appName string) string {
	repoName := inst.GetRepoNameFromAppName(appName)
	return path.Join(inst.AppsInstallDir, repoName)
}

// GetAppInstallPathWithVersionPath get the full app install path and version => /data/rubix-service/apps/install/wires-builds/v0.0.1
func (inst *App) GetAppInstallPathWithVersionPath(appName, version string) string {
	repoName := inst.GetRepoNameFromAppName(appName)
	return path.Join(inst.AppsInstallDir, repoName, version)
}

// GetStoreAppsDir get store dir
func (inst *App) GetStoreAppsDir() string {
	return path.Join(inst.StoreDir, "apps")
}

// GetStoreAppPathAndVersion get the full app install path and version => /data/store/apps/rubix-wires/v0.0.1
func (inst *App) GetStoreAppPathAndVersion(appName, version string) string {
	return path.Join(inst.StoreDir, appName, version)
}

func (inst *App) CreateTmpPath() string {
	return path.Join(inst.TmpDir, uuid.ShortUUID("tmp"))
}

func (inst *App) GetRubixServiceDataDataPath() string {
	return path.Join(inst.RubixServiceDir, "data")
}
