package installer

import (
	"errors"
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/lib-uuid/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

// DirsInstallApp make all the installation dirs
//	appDirName => rubix-wires
//	appInstallName => wires-builds
func (inst *App) DirsInstallApp(appName, version string) error {
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

// MakeTmpDir  => /data/tmp
func (inst *App) MakeTmpDir() error {
	if err := checkDir(inst.DataDir); err != nil {
		return errors.New(fmt.Sprintf("dir not exists %s", inst.DataDir))
	}
	return makeDirectoryIfNotExists(inst.TmpDir, os.FileMode(inst.FileMode))
}

// MakeTmpDirUpload  => /data/tmp
func (inst *App) MakeTmpDirUpload() (string, error) {
	if err := checkDir(inst.DataDir); err != nil {
		return "", errors.New(fmt.Sprintf("dir not exists %s", inst.DataDir))
	}
	tmpDir := fmt.Sprintf("%s/%s", inst.TmpDir, uuid.ShortUUID("tmp"))
	err := makeDirectoryIfNotExists(tmpDir, os.FileMode(inst.FileMode))
	return tmpDir, err
}

// MakeAppInstallDir  => /data/rubix-service/apps/install/wires-builds
func (inst *App) MakeAppInstallDir(appName string, removeExisting ...bool) error {
	if err := emptyPath(appName); err != nil {
		return err
	}

	appInstallDir := fmt.Sprintf("%s/%s", inst.AppsInstallDir, appName)
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
	appDir := fmt.Sprintf("%s/%s/%s", inst.AppsInstallDir, appName, version)
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
	return makeDirectoryIfNotExists(fmt.Sprintf("%s/%s", inst.DataDir, appName), os.FileMode(inst.FileMode))
}

// MakeInstallDir  => /data/rubix-service/install
func (inst *App) MakeInstallDir() error {
	if inst.AppsInstallDir == "" {
		return errors.New("MakeDataDir path can not be empty")
	}
	rsDir := fmt.Sprintf("%s/data", inst.RubixServiceDir)
	err := mkdirAll(rsDir, os.FileMode(inst.FileMode))
	if err != nil {
		log.Errorf("error on making rubix-service dir %s", err.Error())
		return err
	}
	err = mkdirAll(inst.AppsInstallDir, os.FileMode(inst.FileMode))
	if err != nil {
		log.Errorf("error on making rubix-service app install dir %s", err.Error())
		return err
	}
	return err
}

// MakeDataDir  => /data
func (inst *App) MakeDataDir() error {
	if inst.DataDir == "" {
		return errors.New("/data path can not be empty")
	}
	return makeDirectoryIfNotExists(inst.DataDir, os.FileMode(inst.FileMode))
}

// MakeDirectoryIfNotExists make dir
func (inst *App) MakeDirectoryIfNotExists(path string, perm os.FileMode) error {
	return makeDirectoryIfNotExists(path, perm)
}

// MkdirAll make dir
func (inst *App) MkdirAll(path string, perm os.FileMode) error {
	return mkdirAll(path, perm)
}

// mkdirAll all dirs
func mkdirAll(path string, perm os.FileMode) error {
	err := os.MkdirAll(path, perm)
	if err != nil {
		return errors.New(fmt.Sprintf("path %s, err: %s", path, err.Error()))
	}
	return nil
}

// makeDirectoryIfNotExists if not exist make dir
func makeDirectoryIfNotExists(path string, perm os.FileMode) error {
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
		return errors.New(fmt.Sprintf("incorrect provided: %s version number try: v1.2.3", version))
	}
	p := strings.Split(version, ".")
	if len(p) >= 2 && len(p) < 4 {
	} else {
		return errors.New(fmt.Sprintf("incorrect lenght provided: %s version number try: v1.2.3", version))
	}
	return nil
}

func (inst *App) MoveFile(sourcePath, destPath string, deleteAfter bool) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	if deleteAfter {
		err = os.Remove(sourcePath)
		if err != nil {
			return fmt.Errorf("failed removing original file: %s", err)
		}
	}
	return nil
}
