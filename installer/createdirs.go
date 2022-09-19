package installer

import (
	"errors"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/lib-uuid/uuid"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

// CreateInstallAppDirs make all the installation dirs
func (inst *App) CreateInstallAppDirs(appName, appVersion string) error {
	if appName == "" {
		return errors.New(ErrEmptyAppName)
	}
	if appVersion == "" {
		return errors.New(ErrEmptyAppVersion)
	}
	err := inst.MakeAllDirs()
	log.Info("install app edge: MakeAllDirs")
	if err != nil {
		return err
	}
	err = inst.MakeAppDataDir(appName)
	log.Infof("install app edge: MakeAppDataDir app_name: %s", appName)
	if err != nil {
		return err
	}
	err = os.MkdirAll(inst.GetAppDataConfigPath(appName), os.FileMode(inst.FileMode)) // make the app config dir
	log.Infof("install app edge: MakeDirectoryIfNotExists app_name: %s", appName)
	if err != nil {
		return err
	}
	err = inst.MakeAppInstallDir(appName)
	log.Infof("install app edge: MakeAppInstallDir app_name: %s", appName)
	if err != nil {
		return err
	}
	err = inst.MakeAppVersionDir(appName, appVersion)
	log.Infof("install app edge: MakeAppInstallDir app_name: %s app_version: %s", appName, appVersion)
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
	return os.MkdirAll(inst.DataDir, os.FileMode(inst.FileMode))
}

// MakeTmpDir  => /data/tmp
func (inst *App) MakeTmpDir() error {
	return os.MkdirAll(inst.TmpDir, os.FileMode(inst.FileMode))
}

// MakeTmpDirUpload  => /data/tmp/tmp_45EA34EB
func (inst *App) MakeTmpDirUpload() (string, error) {
	tmpDir := inst.CreateTmpPath()
	err := os.MkdirAll(tmpDir, os.FileMode(inst.FileMode))
	return tmpDir, err
}

// MakeInstallDir  => /data/rubix-service/apps/install
func (inst *App) MakeInstallDir() error {
	if inst.AppsInstallDir == "" {
		return errors.New("MakeDataDir path can not be empty")
	}
	rsDataDataDir := inst.GetRubixServiceDataDataPath()
	err := os.MkdirAll(rsDataDataDir, os.FileMode(inst.FileMode))
	if err != nil {
		log.Errorf("error on making rubix-service data dir %s", err.Error())
		return err
	}
	err = os.MkdirAll(inst.AppsInstallDir, os.FileMode(inst.FileMode))
	if err != nil {
		log.Errorf("error on making rubix-service app install dir %s", err.Error())
	}
	return err
}

// MakeAppInstallDir  => /data/rubix-service/apps/install/wires-builds
func (inst *App) MakeAppInstallDir(appName string, removeExisting ...bool) error {
	if appName != "" {
		return errors.New(ErrEmptyAppName)
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
	return os.MkdirAll(appInstallDir, os.FileMode(inst.FileMode))
}

// MakeAppVersionDir  => /data/rubix-service/apps/install/wires-builds/v0.0.1
func (inst *App) MakeAppVersionDir(appName, version string) error {
	if appName != "" {
		return errors.New(ErrEmptyAppName)
	}
	if err := CheckVersion(version); err != nil {
		return err
	}
	appDir := inst.GetAppInstallPathWithVersionPath(appName, version)
	return os.MkdirAll(appDir, os.FileMode(inst.FileMode))
}

// MakeAppDataDir  => /data/flow-framework
func (inst *App) MakeAppDataDir(appName string) error {
	if appName != "" {
		return errors.New(ErrEmptyAppName)
	}
	dataPath := inst.GetAppDataPath(appName)
	return os.MkdirAll(dataPath, os.FileMode(inst.FileMode))
}

// MakeBackupTmpDirUpload  => ~/backup/tmp/tmp_AB34DA34
func (inst *App) MakeBackupTmpDirUpload() (string, error) {
	backupDir := inst.getBackupDir()
	tmpDir := path.Join(backupDir, "tmp", uuid.ShortUUID("tmp"))
	err := os.MkdirAll(tmpDir, os.FileMode(inst.FileMode))
	return tmpDir, err
}
