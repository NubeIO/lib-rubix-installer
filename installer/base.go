package installer

import (
	"fmt"
	fileutils "github.com/NubeIO/lib-dirs/dirs"
)

const nonRoot = 0700
const root = 0777

var FilePerm = root
var DataDir = "/data"
var TmpDir = ""
var AppsInstallDir = ""
var AppBuildName = ""
var HostDownloadPath = ""
var AppsDownloadDir = ""
var LibSystemPath = "/lib/systemd"
var EtcSystemPath = "/etc/systemd"
var DefaultTimeout = 30

type App struct {
	Name             string `json:"app"`            // rubix-wires
	AppBuildName     string `json:"app_build_name"` // wires-builds
	Version          string `json:"version"`        // v1.1.1
	DataDir          string `json:"data_dir"`       // /data
	Perm             int    // file permissions
	HostDownloadPath string `json:"host_download_path"` // home/user/downloads
	ServiceName      string `json:"service_name"`       // nubeio-rubix-wires
	LibSystemPath    string `json:"lib_system_path"`    // /lib/systemd/
	EtcSystemPath    string `json:"etc_system_path"`    // /etc/systemd/
	DefaultTimeout   int    `json:"default_timeout"`
}

func New(app *App) *App {
	if app == nil {
		app = &App{}
	}
	if app.DataDir == "" {
		app.DataDir = DataDir
	}
	if app.Perm == 0 {
		app.Perm = FilePerm
	}
	if app.DefaultTimeout == 0 {
		app.DefaultTimeout = DefaultTimeout
	}
	if app.LibSystemPath == "" {
		app.LibSystemPath = LibSystemPath
	}
	if app.EtcSystemPath == "" {
		app.EtcSystemPath = EtcSystemPath
	}
	if app.HostDownloadPath == "" {
		homeDir, _ := fileutils.Dir()
		app.HostDownloadPath = fmt.Sprintf("%s/Downloads", homeDir)
	}
	DataDir = app.DataDir
	AppBuildName = app.AppBuildName
	HostDownloadPath = app.HostDownloadPath
	TmpDir = fmt.Sprintf("%s/tmp", DataDir)
	AppsInstallDir = fmt.Sprintf("%s/rubix-service/apps/install", DataDir)
	AppsDownloadDir = fmt.Sprintf("%s/rubix-service/apps/download", DataDir)
	return app
}
