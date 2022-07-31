package installer

import (
	"errors"
	"fmt"
	fileutils "github.com/NubeIO/lib-dirs/dirs"
	"os"
)

// BackupFullEdge make a backup of the whole edge device from the /data
func (inst *App) BackupFullEdge(appName string) (string, error) {
	found := inst.DirExists("/data")
	if !found {
		return "", errors.New("failed to find /data")
	}
	version := inst.GetAppVersion(appName)
	if version == "" {
		return "", errors.New("failed to find app version")
	}
	if err := inst.makeBackupAppDir(appName); err != nil {
		return "", err
	}
	path := fmt.Sprintf("/data/backup/%s-%s-%s", appName, version, timestamp())
	return path, fileutils.New().RecursiveZip("", path)
}

// BackupApp backup an app  /data/backups/appName/appName_2022-07-31 12:02:01
func (inst *App) BackupApp(appName string) (string, error) {
	found := inst.ConfirmAppDir(appName)
	if !found {
		return "", errors.New("failed to find app")
	}
	version := inst.GetAppVersion(appName)
	if version == "" {
		return "", errors.New("failed to find app version")
	}
	if err := inst.makeBackupAppDir(appName); err != nil {
		return "", err
	}
	path := fmt.Sprintf("/data/backup/%s-%s-%s", appName, version, timestamp())
	return path, fileutils.New().RecursiveZip("", path)
}

// makeHomeBackupDir backup an app  /user/home/backup
func (inst *App) makeHomeBackupDir(appName string) error {
	home, err := fileutils.Dir()
	if err != nil {
		return err
	}
	name := fmt.Sprintf("%s/backup/%s", home, appName)
	return makeDirectoryIfNotExists(name, os.FileMode(inst.FilePerm))
}

// getHomeBackupDir backup an app  /user/home/backup
func (inst *App) getHomeBackupDir() (string, error) {
	home, err := fileutils.Dir()
	return fmt.Sprintf("%s/backup", home), err
}

// makeBackupAppDir backup an app /data/backups/appName/appName_2022-07-31 12:02:01
func (inst *App) makeBackupAppDir(appName string) error {
	name := fmt.Sprintf("/data/backup/%s", appName)
	return makeDirectoryIfNotExists(name, os.FileMode(inst.FilePerm))
}
