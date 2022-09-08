package installer

import (
	"github.com/NubeIO/lib-files/fileutils"
	"path"
)

func (inst *App) ConfirmAppDataDir(appName string) bool {
	return fileutils.DirExists(path.Join(inst.DataDir, appName))
}

func (inst *App) ConfirmAppInstallDir(appName string) bool {
	return fileutils.DirExists(path.Join(inst.AppsInstallDir, appName))
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
	return fileutils.DirExists(path.Join(inst.GetStoreDir(), "apps", appName))
}

func (inst *App) ConfirmStoreAppVersionDir(appName, version string) bool {
	return fileutils.DirExists(path.Join(inst.GetStoreDir(), "apps", appName, version))
}
