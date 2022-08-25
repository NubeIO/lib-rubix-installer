package installer

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"path/filepath"
)

func setWiresName(appName string) string {
	if appName == "rubix-wires" {
		appName = "wires-builds"
	}
	return appName
}

func setWiresBuildName(appName string) string {
	if appName == "wires-builds" {
		appName = "rubix-wires"
	}
	return appName
}

// filePath make the file path work for unix or windows
func filePath(path string, debug ...bool) string {
	updated := filepath.FromSlash(path)
	if len(debug) > 0 {
		if debug[0] {
			log.Infof("existing-path: %s", path)
			log.Infof("updated-path: %s", updated)
		}
	}
	return filepath.FromSlash(updated)
}

// GetAppPath get the full app install path => /data/rubix-wires
func (inst *App) GetAppPath(appName string) string {
	return inst.getAppPath(appName)
}

// GetAppConfigPath get the full app path =>  /data/rubix-wires
func (inst *App) GetAppConfigPath(appName string) string {
	return filePath(fmt.Sprintf("%s/%s/config", inst.DataDir, appName))
}

// GetAppInstallPath get the full app install path and version => /data/rubix-wires
func (inst *App) getAppPath(appName string) string {
	return filePath(fmt.Sprintf("%s/%s", inst.DataDir, appName))
}

// GetAppInstallPath get the full app install path and version => /data/rubix-service/apps/install/wires-builds
func (inst *App) GetAppInstallPath(appName string) string {
	appName = setWiresName(appName)
	return inst.getAppInstallPath(appName)
}

// GetAppInstallPath get the full app install path and version => /data/rubix-service/apps/install/wires-builds
func (inst *App) getAppInstallPath(appName string) string {
	appName = setWiresName(appName)
	return filePath(fmt.Sprintf("%s/%s", inst.AppsInstallDir, appName))
}

// GetStoreDir get store dir
func (inst *App) GetStoreDir() string {
	return filePath(inst.StoreDir)
}

// GetStoreAppPathAndVersion get the full app install path and version => /data/store/apps/rubix-wires/v0.0.1
func (inst *App) GetStoreAppPathAndVersion(appName, version string) string {
	return filePath(fmt.Sprintf("%s/%s/%s", inst.StoreDir, appName, version))
}

// GetAppInstallPathAndVersion get the full app install path and version => /data/rubix-service/apps/install/wires-builds/v0.0.1
func (inst *App) GetAppInstallPathAndVersion(appName, version string) string {
	appName = setWiresName(appName)
	return inst.getAppInstallPathAndVersion(appName, version)
}

// GetAppInstallPathAndVersion get the full app install path and version => /data/rubix-service/apps/install/wires-builds/v0.0.1
func (inst *App) getAppInstallPathAndVersion(appName, version string) string {
	appName = setWiresName(appName)
	return filePath(fmt.Sprintf("%s/%s/%s", inst.AppsInstallDir, appName, version))
}
