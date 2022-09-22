package installer

import (
	"errors"
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	log "github.com/sirupsen/logrus"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

type MessageResponse struct {
	Message string `json:"message"`
}

type RestoreResponse struct {
	Message        string `json:"message,omitempty"`
	TakeBackupPath string `json:"take_backup_path,omitempty"`
}

type RestoreBackup struct {
	AppName      string                `json:"app_name"`
	DeviceName   *string               `json:"device_name"`
	TakeBackup   bool                  `json:"take_backup"`
	RebootDevice bool                  `json:"reboot_device"`
	File         *multipart.FileHeader `json:"file"`
}

// ------------------
// RESTORE A BACK-UPS
// ------------------

// RestoreBackup restore a backup data dir /data
func (inst *App) RestoreBackup(back *RestoreBackup) (*RestoreResponse, error) {
	var destination = "/"
	var checkDirName = "data"
	resp := &RestoreResponse{}
	if back.TakeBackup {
		backup, err := inst.FullBackUp(back.DeviceName)
		if err != nil {
			return nil, err
		}
		resp.TakeBackupPath = backup
	}
	restore, err := inst.restoreBackup(back.File, destination, inst.DataDir, checkDirName)
	if err != nil {
		return nil, err
	}
	resp.Message = fmt.Sprintf("restored backup sucessfully from: %s", restore.UploadedFile)
	return resp, nil
}

// RestoreAppBackup restore a backup of an app /data/flow-framework
func (inst *App) RestoreAppBackup(back *RestoreBackup) (*RestoreResponse, error) {
	var deleteDirName = inst.GetAppDataPath(back.AppName)
	resp := &RestoreResponse{}
	if back.TakeBackup {
		backup, err := inst.BackupApp(back.AppName, back.DeviceName)
		if err != nil {
			return nil, err
		}
		resp.TakeBackupPath = backup
	}
	restore, err := inst.restoreBackup(back.File, inst.DataDir, deleteDirName, back.AppName)
	if err != nil {
		return nil, err
	}
	resp.Message = fmt.Sprintf("retored backup sucessfully from: %s", restore.UploadedFile)
	return resp, nil
}

func (inst *App) restoreBackup(file *multipart.FileHeader, destination, deleteDirName, checkDirName string) (*UploadResponse, error) {
	if destination == "" {
		return nil, errors.New("destination can not be empty")
	}
	if deleteDirName == "" {
		return nil, errors.New("delete_dir_name can not be empty")
	}
	log.Infof("restore backup delete exiting dir to destination dir: %s", deleteDirName)
	tmpDir, err := inst.MakeBackupTmpDirUpload()
	if err != nil {
		return nil, err
	}
	zipSource, err := inst.SaveUploadedFile(file, tmpDir)
	if err != nil {
		return nil, err
	}
	if checkDirName != "" { // this is a basic check to make sure the upload has a name for example of flow-framework
		zip, err := readZip(zipSource)
		if err != nil {
			return nil, err
		}
		var hasCorrectPath bool
		for _, name := range zip {
			parts := strings.Split(name, "/")
			for _, part := range parts {
				hasCorrectPath = strings.Contains(part, checkDirName)
				if hasCorrectPath {
					break
				}
			}
		}
		if !hasCorrectPath {
			return nil, errors.New(fmt.Sprintf("no mathcing path name in the uploaded zip folder equal to: %s", checkDirName))
		}
	}
	err = fileutils.RmRF(deleteDirName)
	if err != nil {
		return nil, err
	}
	log.Infof("restore backup unzip backup from source: %s to destination dir: %s", zipSource, destination)
	err = unzip(zipSource, destination)
	if err != nil {
		log.Errorf("restore backup unzip backup to dir: %s err: %s", destination, err.Error())
		return nil, err
	}
	return &UploadResponse{
		FileName:     file.Filename,
		UploadedFile: zipSource,
	}, nil
}

// -------------
// LIST BACK-UPS
// -------------

// ListFullBackups list all the backups taken for the data dir /data
func (inst *App) ListFullBackups() ([]string, error) {
	fullBackupDir := inst.getFullBackupDir()
	return fileutils.ListFiles(fullBackupDir)
}

// ListAppsBackups list all the folder for each app
func (inst *App) ListAppsBackups() ([]string, error) {
	appsBackupDir := inst.getAppsBackupDir()
	return fileutils.ListFiles(appsBackupDir)
}

// ListAppBackups list all the backups taken for each app
func (inst *App) ListAppBackups(appName string) ([]string, error) {
	if appName == "" {
		return nil, errors.New(ErrEmptyAppName)
	}
	appBackupDir := inst.getAppBackupDir(appName)
	return fileutils.ListFiles(appBackupDir)
}

// ---------------
// DELETE BACK-UPS
// ---------------

// DeleteAllFullBackups will delete a full backup of the data dir ~/backup/full
func (inst *App) DeleteAllFullBackups() (*MessageResponse, error) {
	resp := &MessageResponse{}
	fullBackupDir := inst.getFullBackupDir()
	err := fileutils.RmRF(fullBackupDir)
	if err != nil {
		resp.Message = fmt.Sprintf("failed to delete: %s", fullBackupDir)
		return resp, err
	}
	resp.Message = fmt.Sprintf("deleted successfully: %s", fullBackupDir)
	return resp, nil
}

// DeleteAllAppsBackups delete all the app backups ~/backups/apps
func (inst *App) DeleteAllAppsBackups() (*MessageResponse, error) {
	resp := &MessageResponse{}
	appsBackupDir := inst.getAppsBackupDir()
	err := fileutils.RmRF(appsBackupDir)
	if err != nil {
		resp.Message = fmt.Sprintf("failed to delete: %s", appsBackupDir)
		return resp, err
	}
	resp.Message = fmt.Sprintf("deleted succefully: %s", appsBackupDir)
	return resp, nil
}

// DeleteAllAppBackups delete all apps backup, eg ~/backups/apps/flow-framework
func (inst *App) DeleteAllAppBackups(appName string) (*MessageResponse, error) {
	resp := &MessageResponse{}
	appBackupDir := inst.getAppBackupDir(appName)
	err := fileutils.RmRF(appBackupDir)
	if err != nil {
		resp.Message = fmt.Sprintf("failed to delete: %s", appBackupDir)
		return resp, err
	}
	resp.Message = fmt.Sprintf("deleted successfully: %s", appBackupDir)
	return resp, nil
}

// DeleteOneAppBackup delete an app backup, eg ~/backups/apps/flow-framework/flow-framework-v0.0.1-2022-07-31T12:02:01
func (inst *App) DeleteOneAppBackup(appName, zipFile string) (*MessageResponse, error) {
	resp := &MessageResponse{}
	appBackupWithBackupFolderDir := inst.getAppBackupWithZipFile(appName, zipFile)
	err := fileutils.Rm(appBackupWithBackupFolderDir)
	if err != nil {
		resp.Message = fmt.Sprintf("failed to delete: %s", appBackupWithBackupFolderDir)
		return resp, err
	}
	resp.Message = fmt.Sprintf("deleted successfully: %s", appBackupWithBackupFolderDir)
	return resp, nil
}

// WipeBackups delete all the backups
func (inst *App) WipeBackups() (*MessageResponse, error) {
	resp := &MessageResponse{}
	backupDir := inst.getBackupDir()
	err := fileutils.RmRF(backupDir)
	if err != nil {
		resp.Message = fmt.Sprintf("failed to delete: %s", backupDir)
		return resp, err
	}
	resp.Message = fmt.Sprintf("deleted successfully: %s", backupDir)
	return resp, nil
}

// -------------
// RUN A BACK-UP
// -------------

// FullBackUp make a backup of the whole edge device from the DataDir, i.e. /data by default
func (inst *App) FullBackUp(deviceName *string) (string, error) {
	found := fileutils.DirExists(inst.DataDir)
	if !found {
		return "", errors.New(fmt.Sprintf("failed to find %s", inst.DataDir))
	}
	fullBackupDir := inst.getFullBackupDir()
	err := os.MkdirAll(fullBackupDir, os.FileMode(inst.FileMode))
	if err != nil {
		return "", err
	}
	zipFile := inst.generateFullBackupZipFile(deviceName)
	return zipFile, fileutils.RecursiveZip(inst.DataDir, zipFile)
}

// BackupApp backup an app ~/apps/appName/appName-version-2022-07-31T12:02:01
func (inst *App) BackupApp(appName string, deviceName *string) (string, error) {
	if appName == "" {
		return "", errors.New(ErrEmptyAppName)
	}
	found := fileutils.DirExists(inst.GetAppDataPath(appName))
	if !found {
		return "", errors.New(fmt.Sprintf("failed to find app: %s", appName))
	}
	version := inst.GetAppVersion(appName)
	if version == "" {
		return "", errors.New("failed to find app version")
	}
	source := inst.GetAppDataPath(appName)
	appBackupDir := inst.getAppBackupDir(appName)
	err := os.MkdirAll(appBackupDir, os.FileMode(inst.FileMode))
	if err != nil {
		return "", err
	}
	zipFile := inst.generateAppBackupZipFile(appName, version, deviceName)
	return zipFile, fileutils.RecursiveZip(source, zipFile)
}

// getBackupDir backup home dir ~/backup
func (inst *App) getBackupDir() string {
	return inst.BackupsDir
}

// getFullBackupDir backup an app  ~/backup/full
func (inst *App) getFullBackupDir() string {
	home := inst.getBackupDir()
	return path.Join(home, "full")
}

// getAppsBackupDir backup home dir ~/backup/apps
func (inst *App) getAppsBackupDir() string {
	home := inst.getBackupDir()
	return path.Join(home, "apps")
}

// getAppBackupDir backup an app  ~/backup/apps/flow-framework
func (inst *App) getAppBackupDir(appName string) string {
	appsBackupDir := inst.getAppsBackupDir()
	return path.Join(appsBackupDir, appName)
}

// getAppBackupWithZipFile backup an app  ~/backup/apps/flow-framework/flow-framework-v0.0.1-2022-07-31T12:02:01.zip
func (inst *App) getAppBackupWithZipFile(appName, zipFile string) string {
	appBackupDir := inst.getAppBackupDir(appName)
	return path.Join(appBackupDir, zipFile)
}

func (inst *App) generateFullBackupZipFile(deviceName *string) string {
	zipFileName := ""
	if deviceName == nil {
		zipFileName = fmt.Sprintf("full-backup-%s.zip", timestamp())
	} else {
		zipFileName = fmt.Sprintf("%s-full-backup-%s.zip", *deviceName, timestamp())
	}
	return path.Join(inst.getFullBackupDir(), zipFileName)
}

func (inst *App) generateAppBackupZipFile(appName, version string, deviceName *string) string {
	zipFileName := ""
	if deviceName == nil {
		zipFileName = fmt.Sprintf("backup-%s-%s-%s.zip", appName, version, timestamp())
	} else {
		zipFileName = fmt.Sprintf("%s-backup-%s-%s-%s.zip", *deviceName, appName, version, timestamp())
	}
	return path.Join(inst.getAppBackupDir(appName), zipFileName)
}
