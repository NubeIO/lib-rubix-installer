package installer

import (
	"errors"
	"fmt"
	fileutils "github.com/NubeIO/lib-dirs/dirs"
	"github.com/NubeIO/lib-uuid/uuid"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strings"
)

type MessageResponse struct {
	Message string `json:"message"`
}

type RestoreResponse struct {
	Message        string `json:"message,omitempty"`
	TakeBackupPath string `json:"take_backup_path,omitempty"`
}

/*
RESTORE A BACK-UPS
*/

type RestoreBackup struct {
	AppName      string                `json:"app_name"`
	DeviceName   string                `json:"device_name"`
	TakeBackup   bool                  `json:"take_backup"`
	RebootDevice bool                  `json:"reboot_device"`
	File         *multipart.FileHeader `json:"file"`
}

// RestoreBackup restore a backup data dir /data
func (inst *App) RestoreBackup(back *RestoreBackup) (*RestoreResponse, error) {
	var file = back.File
	var takeBackup = back.TakeBackup
	var deiceName = back.DeviceName
	var destination = "/"
	var deleteDirName = inst.DataDir
	var checkDirName = "data"
	// delete the existing data dir
	resp := &RestoreResponse{}
	if takeBackup {
		backup, err := inst.FullBackUp(deiceName)
		if err != nil {
			return nil, err
		}
		resp.TakeBackupPath = backup
	}
	restore, err := inst.restoreBackup(file, destination, deleteDirName, checkDirName, "")
	if err != nil {
		return nil, err
	}
	resp.Message = fmt.Sprintf("retored backup ok from: %s", restore.UploadedFile)
	return resp, nil
}

// RestoreAppBackup restore a backup of an app /data/flow-framework
func (inst *App) RestoreAppBackup(back *RestoreBackup) (*RestoreResponse, error) {
	var file = back.File
	var takeBackup = back.TakeBackup
	var appName = back.AppName
	var deiceName = back.DeviceName
	var checkDirName = appName
	var deleteDirName = fmt.Sprintf("%s/%s", inst.DataDir, appName)
	resp := &RestoreResponse{}
	if takeBackup {
		backup, err := inst.BackupApp(appName, deiceName)
		if err != nil {
			return nil, err
		}
		resp.TakeBackupPath = backup
	}
	version := inst.GetAppVersion(appName)
	if version == "" {
		return nil, errors.New("app version was not found")
	}
	restore, err := inst.restoreBackup(file, inst.DataDir, deleteDirName, checkDirName, version)
	if err != nil {
		return nil, err
	}
	resp.Message = fmt.Sprintf("retored backup ok from: %s", restore.UploadedFile)
	return resp, nil
}

// Upload upload a build
func (inst *App) restoreBackup(file *multipart.FileHeader, destination, deleteDirName, checkDirName, appVersion string) (*UploadResponse, error) {
	// make the dirs
	var err error
	if destination == "" {
		return nil, errors.New("destination can not be empty")
	}
	if deleteDirName == "" {
		return nil, errors.New("destination can not be empty")
	}
	log.Infof("restore backup delete exiting dir to destination dir:%s", deleteDirName)
	var tmpDir string
	if tmpDir, err = inst.MakeBackupTmpDirUpload(); err != nil {
		return nil, err
	}
	// save app in tmp dir
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
			return nil, errors.New(fmt.Sprintf("no mathcing path name in the uploaded zip folder equal to:%s", checkDirName))
		}
	}
	if appVersion != "" {
		parts := strings.Split(zipSource, "/")
		for _, part := range parts {
			fmt.Println(strings.Contains(part, ".zip"), part)
			if strings.Contains(part, ".zip") {

			}
		}

		fmt.Println(parts)
	}
	err = inst.RmRF(deleteDirName)
	if err != nil {
		return nil, err
	}
	log.Infof("restore backup unzip backup from source:%s", zipSource)
	log.Infof("restore backup unzip backup to destination dir:%s", destination)
	err = unzip(zipSource, destination)
	if err != nil {
		log.Errorf("restore backup unzip backup to dir:%s err:%s", destination, err.Error())
		return nil, err
	}
	return &UploadResponse{
		FileName:     file.Filename,
		UploadedFile: zipSource,
	}, err
}

/*
LIST BACK-UPS
*/

// ListFullBackups list all the backups taken for the data dir /data
func (inst *App) ListFullBackups() ([]ListBackups, error) {
	path, err := inst.generateHomeFullBackupFolderName()
	if err != nil {
		return nil, err
	}
	return inst.listFilesAndPath(path)

}

// ListAppBackupsDirs list all the folder for each app
func (inst *App) ListAppBackupsDirs() ([]string, error) {
	path, err := inst.generateAppHomeBackupsFolderName()
	if err != nil {
		return nil, err
	}
	return inst.listFiles(path)
}

// ListBackupsByApp list all the backups taken for each app
func (inst *App) ListBackupsByApp(appName string) ([]ListBackups, error) {
	if appName == "" {
		return nil, errors.New("app name can not be empty")
	}
	path, err := inst.generateAppHomeBackupFolderName(appName)
	if err != nil {
		return nil, err
	}
	apps, err := inst.listFilesAndPath(path)
	return apps, err
}

type ListBackups struct {
	Path        string `json:"path"`
	BackupName  string `json:"name,omitempty"`
	PathAndName string `json:"path_and_name,omitempty"`
	Message     string `json:"message,omitempty"`
}

func (inst *App) listFilesAndPath(path string) ([]ListBackups, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	var dirContents []ListBackups
	var dirContent ListBackups
	dirContent.Path = path
	if fileInfo.IsDir() {
		files, err := ioutil.ReadDir(path)
		if len(files) == 0 {
			dirContent.Message = "no apps found"
			dirContents = append(dirContents, dirContent)
		}
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			dirContent.BackupName = file.Name()
			dirContent.PathAndName = fmt.Sprintf("%s/%s", path, file.Name())
			dirContents = append(dirContents, dirContent)
		}
	}

	return dirContents, nil
}

/*
DELETE BACK-UPS
*/

//DeleteAllFullBackups will delete a full backup of the data dir /data
func (inst *App) DeleteAllFullBackups() (*MessageResponse, error) {
	resp := &MessageResponse{}
	path, err := inst.generateHomeFullBackupFolderName()
	if err != nil {
		resp.Message = "failed to find backup path"
		return resp, err
	}
	//err = inst.RmRF(path)
	if err != nil {
		resp.Message = fmt.Sprintf("failed to delete:%s", path)
		return resp, err
	}
	resp.Message = fmt.Sprintf("deleted ok:%s", path)
	return resp, nil
}

// DeleteAllAppBackups delete all the app backups /user/home/backup/apps
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

// WipeBackups delete all the backups
func (inst *App) WipeBackups() (*MessageResponse, error) {
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

/*
RUN A BACK-UP
*/

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
	if appName == "" {
		return "", errors.New("app name can not be empty")
	}
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

//MakeBackupTmpDirUpload  => /user/home/backup/tmp/tmp_dir_name
func (inst *App) MakeBackupTmpDirUpload() (string, error) {
	home, err := inst.backUpHome()
	if err != nil {
		return "", err
	}
	tmpDir := fmt.Sprintf("%s/tmp/%s", home, uuid.ShortUUID("tmp"))
	err = makeDirectoryIfNotExists(tmpDir, os.FileMode(inst.FilePerm))
	return tmpDir, err
}

// backUpHome backup home dir /user/home/backup
func (inst *App) backUpHome() (string, error) {
	return inst.BackupsDir, nil
}

// backUpHome backup home dir /user/home/backup/apps
func (inst *App) generateAppHomeBackupsFolderName() (string, error) {
	home, err := inst.backUpHome()
	path := fmt.Sprintf("%s/apps", home)
	return path, err
}

// generateHomeFullBackupFolderName backup an app  /user/home/backup/full/
func (inst *App) generateHomeFullBackupFolderName() (string, error) {
	home, err := inst.backUpHome()
	path := fmt.Sprintf("%s/full", home)
	return path, err
}

// generateAppHomeBackupFolderName backup an app  /user/home/backup/flow-framework/v0.0.1
func (inst *App) generateAppHomeBackupFolderName(appName string) (string, error) {
	home, err := inst.generateAppHomeBackupsFolderName()
	path := fmt.Sprintf("%s/%s", home, appName)
	return path, err
}
