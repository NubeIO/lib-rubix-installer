package installer

import (
	"errors"
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/lib-uuid/uuid"
	log "github.com/sirupsen/logrus"
	"os"
)

// CreateInstallAppDirs make all the installation dirs
func (inst *App) CreateInstallAppDirs(appName, version string) error {
	if appName == "" {
		return errors.New("app name can not be empty")
	}
	if version == "" {
		return errors.New("app version can not be empty")
	}
	err := inst.MakeAllDirs()
	log.Info("install app edge: MakeAllDirs")
	if err != nil {
		return err
	}
	err = inst.MakeAppDataDir(appName)
	log.Infof("install app edge: MakeAppDataDir app: %s", appName)
	if err != nil {
		return err
	}
	err = inst.MakeDirectoryIfNotExists(fmt.Sprintf("%s/config", inst.GetAppDataPath(appName)), os.FileMode(inst.FileMode)) // make the app config dir
	log.Infof("install app edge: MakeDirectoryIfNotExists app: %s", appName)
	if err != nil {
		return err
	}
	err = inst.MakeAppInstallDir(appName)
	log.Infof("install app edge: MakeAppInstallDir app-build-name: %s", appName)
	if err != nil {
		return err
	}
	err = inst.MakeAppVersionDir(appName, version)
	log.Infof("install app edge: MakeAppInstallDir app-build-name: %s version: %s", appName, version)
	if err != nil {
		return err
	}
	return nil
}

// MakeAllDirs make all the required dirs
func (inst *App) MakeAllDirs() error {
	err := inst.MakeDataDir()
	if err != nil {
		return err
	}
	err = inst.MakeTmpDir()
	if err != nil {
		return err
	}
	err = inst.MakeInstallDir()
	if err != nil {
		return err
	}
	return nil
}

// MakeDataDir  => /data
func (inst *App) MakeDataDir() error {
	return makeDirectoryIfNotExists(inst.DataDir, os.FileMode(inst.FileMode))
}

// MakeTmpDir  => /data/tmp
func (inst *App) MakeTmpDir() error {
	if err := checkDir(inst.DataDir); err != nil {
		return errors.New(fmt.Sprintf("dir not exists %s", inst.DataDir))
	}
	return makeDirectoryIfNotExists(inst.TmpDir, os.FileMode(inst.FileMode))
}

// MakeTmpDirUpload  => /data/tmp/tmp_45EA34EB
func (inst *App) MakeTmpDirUpload() (string, error) {
	if err := checkDir(inst.DataDir); err != nil {
		return "", errors.New(fmt.Sprintf("dir not exists %s", inst.DataDir))
	}
	tmpDir := fmt.Sprintf("%s/%s", inst.TmpDir, uuid.ShortUUID("tmp"))
	err := makeDirectoryIfNotExists(tmpDir, os.FileMode(inst.FileMode))
	return tmpDir, err
}

// MakeInstallDir  => /data/rubix-service/install
func (inst *App) MakeInstallDir() error {
	if inst.AppsInstallDir == "" {
		return errors.New("MakeDataDir path can not be empty")
	}
	rsDir := fmt.Sprintf("%s/data", inst.RubixServiceDir)
	err := mkdirAll(rsDir, os.FileMode(inst.FileMode))
	if err != nil {
		log.Errorf("error on making rubix-service data dir %s", err.Error())
		return err
	}
	err = mkdirAll(inst.AppsInstallDir, os.FileMode(inst.FileMode))
	if err != nil {
		log.Errorf("error on making rubix-service app install dir %s", err.Error())
		return err
	}
	return err
}

// MakeAppInstallDir  => /data/rubix-service/apps/install/wires-builds
func (inst *App) MakeAppInstallDir(appName string, removeExisting ...bool) error {
	if err := emptyPath(appName); err != nil {
		return err
	}

	appInstallDir := inst.GetAppInstallPath(appName)
	if len(removeExisting) > 0 {
		if removeExisting[0] {
			err := fileutils.RmRF(appInstallDir)
			if err != nil {
				log.Errorf("delete existing install dir: %s", err.Error())
			}
		}
	}
	return makeDirectoryIfNotExists(appInstallDir, os.FileMode(inst.FileMode))
}

// MakeAppVersionDir  => /data/rubix-service/apps/install/wires-builds/v0.0.1
func (inst *App) MakeAppVersionDir(appName, version string) error {
	if err := emptyPath(appName); err != nil {
		return err
	}
	if err := checkVersion(version); err != nil {
		return err
	}
	appDir := inst.GetAppInstallPathWithVersionPath(appName, version)
	return makeDirectoryIfNotExists(appDir, os.FileMode(inst.FileMode))
}

// MakeAppDataDir  => /data/flow-framework
func (inst *App) MakeAppDataDir(appName string) error {
	if err := emptyPath(appName); err != nil {
		return err
	}
	if err := checkDir(inst.DataDir); err != nil {
		return errors.New(fmt.Sprintf("dir not exists %s", inst.DataDir))
	}
	dataPath := inst.GetAppDataPath(appName)
	return makeDirectoryIfNotExists(dataPath, os.FileMode(inst.FileMode))
}

// MakeDirectoryIfNotExists make dir
func (inst *App) MakeDirectoryIfNotExists(path string, perm os.FileMode) error {
	return makeDirectoryIfNotExists(path, perm)
}
