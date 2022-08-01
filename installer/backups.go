package installer

import (
	"errors"
	"fmt"
	fileutils "github.com/NubeIO/lib-dirs/dirs"
	"os"
)

type MessageResponse struct {
	Message string `json:"message"`
}

func (inst *App) ListFullBackups() ([]string, error) {
	path, err := inst.generateHomeFullBackupFolderName()
	if err != nil {
		return nil, err
	}
	return inst.listFiles(path)

}

func (inst *App) ListAppBackupsDirs() ([]string, error) {
	path, err := inst.generateAppHomeBackupsFolderName()
	if err != nil {
		return nil, err
	}
	return inst.listFiles(path)
}

func (inst *App) ListBackupsByApp(appName string) ([]string, error) {
	if appName == "" {
		return nil, errors.New("app name can not be empty")
	}
	path, err := inst.generateAppHomeBackupFolderName(appName)
	if err != nil {
		return nil, err
	}
	return inst.listFiles(path)
}

func (inst *App) DeleteAllFullBackups() (*MessageResponse, error) {
	resp := &MessageResponse{}
	path, err := inst.generateHomeFullBackupFolderName()
	if err != nil {
		resp.Message = "failed to find backup path"
		return resp, err
	}
	err = inst.RmRF(path)
	if err != nil {
		resp.Message = fmt.Sprintf("failed to delete:%s", path)
		return resp, err
	}
	resp.Message = fmt.Sprintf("deleted ok:%s", path)
	return resp, nil
}

func (inst *App) DeleteAllAppBackups() (*MessageResponse, error) {
	resp := &MessageResponse{}
	path, err := inst.generateAppHomeBackupsFolderName()
	if err != nil {
		resp.Message = "failed to find backup path"
		return resp, err
	}
	err = inst.RmRF(path)
	if err != nil {
		resp.Message = fmt.Sprintf("failed to delete:%s", path)
		return resp, err
	}
	resp.Message = fmt.Sprintf("deleted ok:%s", path)
	return resp, nil
}

// DeleteAppAllBackUpByName delete all apps backup, eg /data/backups/apps/flow-framework
func (inst *App) DeleteAppAllBackUpByName(appName string) (*MessageResponse, error) {
	resp := &MessageResponse{}
	path, err := inst.generateAppHomeBackupsFolderName()
	if err != nil {
		resp.Message = "failed to find backup path"
		return resp, err
	}
	path = fmt.Sprintf("%s/%s", path, appName)
	err = inst.RmRF(path)
	if err != nil {
		resp.Message = fmt.Sprintf("failed to delete:%s", path)
		return resp, err
	}
	resp.Message = fmt.Sprintf("deleted ok:%s", path)
	return resp, nil
}

// DeleteAppOneBackUpByName delete an app backup, eg /data/backups/apps/flow-framework/appName-version-2022-07-31 12:02:01
func (inst *App) DeleteAppOneBackUpByName(appName, backupFolder string) (*MessageResponse, error) {
	resp := &MessageResponse{}
	path, err := inst.generateAppHomeBackupsFolderName()
	if err != nil {
		resp.Message = "failed to find backup path"
		return resp, err
	}
	path = fmt.Sprintf("%s/%s/%s", path, appName, backupFolder)
	err = inst.Rm(path)
	if err != nil {
		resp.Message = fmt.Sprintf("failed to delete:%s", path)
		return resp, err
	}
	resp.Message = fmt.Sprintf("deleted ok:%s", path)
	return resp, nil
}

func (inst *App) DeleteAllBackups() (*MessageResponse, error) {
	resp := &MessageResponse{}
	path, err := inst.backUpHome()
	if err != nil {
		resp.Message = "failed to find backup path"
		return resp, err
	}
	err = inst.RmRF(path)
	if err != nil {
		resp.Message = fmt.Sprintf("failed to delete:%s", path)
		return resp, err
	}
	resp.Message = fmt.Sprintf("deleted ok:%s", path)
	return resp, nil
}

// FullBackUp make a backup of the whole edge device from the /data
func (inst *App) FullBackUp(deiceName ...string) (string, error) {
	found := inst.DirExists(inst.DataDir)
	if !found {
		return "", errors.New("failed to find /data")
	}
	source := inst.DataDir
	path, err := inst.generateHomeFullBackupFolderName()
	if err != nil {
		return "", err
	}
	err = inst.MakeDirectoryIfNotExists(path, os.FileMode(inst.FilePerm))
	if err != nil {
		return "", err
	}
	zipName := fmt.Sprintf("%s/full-backup-%s.zip", path, timestamp())
	if len(deiceName) > 0 {
		if deiceName[0] != "" {
			zipName = fmt.Sprintf("%s/%s-full-backup-%s.zip", path, deiceName[0], timestamp())
		}
	}
	return zipName, fileutils.New().RecursiveZip(source, zipName)
}

// BackupApp backup an app  /data/backups/apps/appName/appName-version-2022-07-31 12:02:01
func (inst *App) BackupApp(appName string, deiceName ...string) (string, error) {
	found := inst.ConfirmAppDir(appName)
	if !found {
		return "", errors.New(fmt.Sprintf("failed to find app:%s", appName))
	}
	version := inst.GetAppVersion(appName)
	if version == "" {
		return "", errors.New("failed to find app version")
	}
	source := fmt.Sprintf("%s/%s", inst.DataDir, appName)
	path, err := inst.generateAppHomeBackupFolderName(appName)
	if err != nil {
		return "", err
	}
	err = inst.MakeDirectoryIfNotExists(path, os.FileMode(inst.FilePerm))
	if err != nil {
		return "", err
	}
	zipName := fmt.Sprintf("%s/backup-%s-%s-%s.zip", path, appName, version, timestamp())
	if len(deiceName) > 0 {
		if deiceName[0] != "" {
			zipName = fmt.Sprintf("%s/%s-backup-%s-%s-%s.zip", path, deiceName[0], appName, version, timestamp())
		}
	}
	return zipName, fileutils.New().RecursiveZip(source, zipName)
}

// backUpHome backup home dir /user/home/backup
func (inst *App) backUpHome() (string, error) {
	home, err := fileutils.Dir()
	path := fmt.Sprintf("%s/backup", home)
	return path, err
}

// backUpHome backup home dir /user/home/backup/apps
func (inst *App) generateAppHomeBackupsFolderName() (string, error) {
	home, err := inst.backUpHome()
	path := fmt.Sprintf("%s/apps", home)
	return path, err
}

// generateHomeBackupFolderName backup an app  /user/home/backup/full/
func (inst *App) generateHomeFullBackupFolderName() (string, error) {
	home, err := inst.backUpHome()
	path := fmt.Sprintf("%s/full", home)
	return path, err
}

// generateHomeBackupFolderName backup an app  /user/home/backup/flow-framework/v0.0.1
func (inst *App) generateAppHomeBackupFolderName(appName string) (string, error) {
	home, err := inst.generateAppHomeBackupsFolderName()
	path := fmt.Sprintf("%s/%s", home, appName)
	return path, err
}
