package installer

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"path/filepath"
)

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

// GetAppPath get the full app install path and version => /data/rubix-service/apps/install/wires-builds
func (inst *App) GetAppPath(appName string) string {
	return inst.getAppPath(appName)
}

// GetAppInstallPath get the full app install path and version => /data/rubix-service/apps/install/wires-builds
func (inst *App) getAppPath(appName string) string {
	return filePath(fmt.Sprintf("%s/%s", inst.DataDir, appName))
}

// GetAppInstallPath get the full app install path and version => /data/rubix-service/apps/install/wires-builds
func (inst *App) GetAppInstallPath(appBuildName string) string {
	return inst.getAppInstallPath(appBuildName)
}

// GetAppInstallPath get the full app install path and version => /data/rubix-service/apps/install/wires-builds
func (inst *App) getAppInstallPath(appBuildName string) string {
	return filePath(fmt.Sprintf("%s/%s", inst.AppsInstallDir, appBuildName))
}

// GetStoreDir get store dir
func (inst *App) GetStoreDir() string {
	return filePath(inst.StoreDir)
}

// GetAppInstallPathAndVersion get the full app install path and version => /data/rubix-service/apps/install/wires-builds/v0.0.1
func (inst *App) GetAppInstallPathAndVersion(appBuildName, version string) string {
	return inst.getAppInstallPathAndVersion(appBuildName, version)
}

// GetAppInstallPathAndVersion get the full app install path and version => /data/rubix-service/apps/install/wires-builds/v0.0.1
func (inst *App) getAppInstallPathAndVersion(appBuildName, version string) string {
	return filePath(fmt.Sprintf("%s/%s/%s", inst.AppsInstallDir, appBuildName, version))
}
