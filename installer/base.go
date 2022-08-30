package installer

import (
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/lib-systemctl-go/systemctl"
)

const filePerm = 0755
const defaultTimeout = 30

var libSystemPath = "/lib/systemd"
var etcSystemPath = "/etc/systemd"

type App struct {
	Name             string `json:"name"`               // rubix-wires
	Version          string `json:"version"`            // v1.1.1
	DataDir          string `json:"data_dir"`           // /data
	HostDownloadPath string `json:"host_download_path"` // home/user/downloads
	StoreDir         string `json:"store_dir"`
	TmpDir           string `json:"tmp_dir"`
	UserRubixHome    string `json:"user_rubix_home"`
	FilePerm         int    `json:"file_perm"`       // file permissions
	ServiceName      string `json:"service_name"`    // nubeio-rubix-wires
	LibSystemPath    string `json:"lib_system_path"` // /lib/systemd/
	EtcSystemPath    string `json:"etc_system_path"` // /etc/systemd/
	DefaultTimeout   int    `json:"default_timeout"`
	RubixServiceDir  string `json:"rubix_service_dir"`
	AppsInstallDir   string `json:"apps_install_dir"` // /rubix-service/apps/install
	AppsDownloadDir  string `json:"apps_download_dir"`
	BackupsDir       string `json:"backups_dir"`
	Ctl              *systemctl.Ctl
}

func New(app *App) *App {
	homeDir, _ := fileutils.HomeDir()
	if app == nil {
		app = &App{}
	}
	if app.DataDir == "" {
		app.DataDir = filePath("/data")
	}
	if app.FilePerm == 0 {
		app.FilePerm = filePerm
	}
	if app.DefaultTimeout == 0 {
		app.DefaultTimeout = defaultTimeout
	}
	if app.LibSystemPath == "" {
		app.LibSystemPath = libSystemPath
	}
	if app.EtcSystemPath == "" {
		app.EtcSystemPath = etcSystemPath
	}
	if app.StoreDir == "" {
		app.StoreDir = filePath(fmt.Sprintf("%s/store", app.DataDir))
	}
	if app.UserRubixHome == "" {
		app.HostDownloadPath = fmt.Sprintf("%s/rubix", homeDir)
	}
	if app.HostDownloadPath == "" {
		app.HostDownloadPath = filePath(fmt.Sprintf("%s/Downloads", homeDir))
	}
	if app.RubixServiceDir == "" {
		app.RubixServiceDir = filePath(fmt.Sprintf("%s/rubix-service", app.DataDir))
	}
	if app.AppsDownloadDir == "" {
		app.AppsDownloadDir = filePath(fmt.Sprintf("%s/rubix-service/apps/download", app.DataDir))
	}
	if app.AppsInstallDir == "" {
		app.AppsInstallDir = filePath(fmt.Sprintf("%s/rubix-service/apps/install", app.DataDir))
	}
	if app.AppsDownloadDir == "" {
		app.AppsDownloadDir = filePath(fmt.Sprintf("%s/rubix-service/apps/download", app.DataDir))
	}
	if app.TmpDir == "" {
		app.TmpDir = "/data/tmp"
	}
	if app.BackupsDir == "" {
		app.BackupsDir = fmt.Sprintf("%s/backup", homeDir)
	}
	app.Ctl = systemctl.New(false, app.DefaultTimeout)
	return app
}
