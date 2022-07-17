package installer

import (
	"errors"
	"fmt"
	"github.com/NubeIO/lib-uuid/uuid"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

// DirsInstallApp make all the installation dirs
//	appDirName => rubix-wires
//	appInstallName => wires-builds
func (inst *App) DirsInstallApp(appName, appBuildName, version string) error {
	err := inst.MakeAllDirs()
	if err != nil {
		return err
	}
	err = inst.MakeAppDir(appName)
	if err != nil {
		return err
	}
	err = inst.MakeAppInstallDir(appBuildName)
	if err != nil {
		return err
	}
	err = inst.MakeAppVersionDir(appBuildName, version)
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
	err = inst.MakeDownloadDir()
	if err != nil {
		return err
	}
	return nil
}

//MakeTmpDir  => /data/tmp
func (inst *App) MakeTmpDir() error {
	if err := checkDir(inst.DataDir); err != nil {
		return errors.New(fmt.Sprintf("dir not exists %s", inst.DataDir))
	}
	return makeDirectoryIfNotExists(TmpDir, os.FileMode(FilePerm))
}

//MakeTmpDirUpload  => /data/tmp
func (inst *App) MakeTmpDirUpload() (string, error) {
	if err := checkDir(inst.DataDir); err != nil {
		return "", errors.New(fmt.Sprintf("dir not exists %s", inst.DataDir))
	}
	tmpDir := fmt.Sprintf("%s/%s", TmpDir, uuid.ShortUUID("tmp"))
	err := makeDirectoryIfNotExists(tmpDir, os.FileMode(FilePerm))
	return tmpDir, err
}

//MakeAppInstallDir  => /data/rubix-service/apps/install/wires-builds
func (inst *App) MakeAppInstallDir(appBuildName string) error {
	if err := emptyPath(appBuildName); err != nil {
		return err
	}
	appInstallDir := fmt.Sprintf("%s/%s", AppsInstallDir, appBuildName)
	err := inst.RmRF(appInstallDir)
	if err != nil {
		log.Errorf("delete existing install dir: %s", err.Error())
	}
	return makeDirectoryIfNotExists(fmt.Sprintf("%s/%s", AppsInstallDir, appBuildName), os.FileMode(FilePerm))
}

//MakeAppVersionDir  => /data/rubix-service/apps/install/wires-builds/v0.0.1
func (inst *App) MakeAppVersionDir(appBuildName, version string) error {
	if err := emptyPath(appBuildName); err != nil {
		return err
	}
	if err := checkVersion(version); err != nil {
		return err
	}
	appDir := fmt.Sprintf("%s/%s/%s", AppsInstallDir, appBuildName, version)
	return makeDirectoryIfNotExists(appDir, os.FileMode(FilePerm))
}

//MakeAppDir  => /data/flow-framework
func (inst *App) MakeAppDir(appName string) error {
	if err := emptyPath(appName); err != nil {
		return err
	}
	if err := checkDir(inst.DataDir); err != nil {
		return errors.New(fmt.Sprintf("dir not exists %s", inst.DataDir))
	}
	return makeDirectoryIfNotExists(fmt.Sprintf("%s/%s", inst.DataDir, appName), os.FileMode(FilePerm))
}

//MakeInstallDir  => /data/rubix-service/install
func (inst *App) MakeInstallDir() error {
	if AppsInstallDir == "" {
		return errors.New("path can not be empty")
	}
	return mkdirAll(AppsInstallDir, os.FileMode(FilePerm))
}

//MakeDownloadDir  => /data/rubix-service/download
func (inst *App) MakeDownloadDir() error {
	if AppsInstallDir == "" {
		return errors.New("path can not be empty")
	}
	return mkdirAll(AppsDownloadDir, os.FileMode(FilePerm))
}

//MakeDataDir  => /data
func (inst *App) MakeDataDir() error {
	if inst.DataDir == "" {
		return errors.New("path can not be empty")
	}
	return makeDirectoryIfNotExists(inst.DataDir, os.FileMode(FilePerm))
}

// mkdirAll all dirs
func mkdirAll(path string, perm os.FileMode) error {
	return mkdirAll(path, perm)
}

// makeDirectoryIfNotExists if not exist make dir
func makeDirectoryIfNotExists(path string, perm os.FileMode) error {
	path = filePath(path)
	if perm == 0 {
		perm = 0755
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return mkdirAll(path, os.ModeDir|perm)
	}
	return nil
}

func empty(name string) error {
	if name == "" {
		return errors.New("can not be empty")
	}
	return nil
}

func emptyPath(path string) error {
	if path == "" {
		return errors.New("path can not be empty")
	}
	return nil
}

func checkDir(path string) error {
	path = filePath(path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}
	return nil
}

func checkVersionBool(version string) bool {
	var hasV bool
	var correctLen bool
	if version[0:1] == "v" { // make sure have a v at the start v0.1.1
		hasV = true
	}
	p := strings.Split(version, ".")
	if len(p) >= 2 && len(p) < 4 {
		correctLen = true
	}
	if hasV && correctLen {
		return true
	}
	return false
}

func checkVersion(version string) error {
	if version[0:1] != "v" { // make sure have a v at the start v0.1.1
		return errors.New(fmt.Sprintf("incorrect provided:%s version number try: v1.2.3", version))
	}
	p := strings.Split(version, ".")
	if len(p) >= 2 && len(p) < 4 {
	} else {
		return errors.New(fmt.Sprintf("incorrect lenght provided:%s version number try: v1.2.3", version))
	}
	return nil
}
