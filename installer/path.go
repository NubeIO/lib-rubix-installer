package installer

import (
	"path"
)

// GetAppDataPath get the full app install path => /data/rubix-wires
func (inst *App) GetAppDataPath(appName string) string {
	return path.Join(inst.DataDir, appName)
}

// GetAppConfigPath get the full app path =>  /data/rubix-wires
func (inst *App) GetAppConfigPath(appName string) string {
	return path.Join(inst.DataDir, appName, "config")
}

// GetAppInstallPath get the full app install path and version => /data/rubix-service/apps/install/wires-builds
func (inst *App) GetAppInstallPath(appName string) string {
	return path.Join(inst.AppsInstallDir, appName)
}

// GetStoreDir get store dir
func (inst *App) GetStoreDir() string {
	return inst.StoreDir
}

// GetStoreAppPathAndVersion get the full app install path and version => /data/store/apps/rubix-wires/v0.0.1
func (inst *App) GetStoreAppPathAndVersion(appName, version string) string {
	return path.Join(inst.StoreDir, appName, version)
}

// GetAppInstallPathAndVersion get the full app install path and version => /data/rubix-service/apps/install/wires-builds/v0.0.1
func (inst *App) GetAppInstallPathAndVersion(appName, version string) string {
	return path.Join(inst.AppsInstallDir, appName, version)
}
