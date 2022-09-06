package installer

import (
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/lib-systemctl-go/systemctl"
)

const fileMode = 0755
const defaultTimeout = 30

type App struct {
	DataDir         string `json:"data_dir"`          // /data
	StoreDir        string `json:"store_dir"`         // <data_dir>/store
	TmpDir          string `json:"tmp_dir"`           // <data_dir>/tmp
	FileMode        int    `json:"file_mode"`         // 0755
	DefaultTimeout  int    `json:"default_timeout"`   // 30
	UserRubixHome   string `json:"user_rubix_home"`   // ~/rubix
	RubixServiceDir string `json:"rubix_service_dir"` // <data_dir>/rubix-service
	AppsInstallDir  string `json:"apps_install_dir"`  // <data_dir>/rubix-service/apps/install
	BackupsDir      string `json:"backups_dir"`       // ~/backup
	SystemCtl       *systemctl.SystemCtl
}

func New(app *App) *App {
	homeDir, _ := fileutils.HomeDir()
	if app == nil {
		app = &App{}
	}
	if app.DataDir == "" {
		app.DataDir = filePath("/data")
	}
	if app.FileMode == 0 {
		app.FileMode = fileMode
	}
	if app.DefaultTimeout == 0 {
		app.DefaultTimeout = defaultTimeout
	}
	if app.StoreDir == "" {
		app.StoreDir = filePath(fmt.Sprintf("%s/store", app.DataDir))
	}
	if app.UserRubixHome == "" {
		app.UserRubixHome = fmt.Sprintf("%s/rubix", homeDir)
	}
	if app.RubixServiceDir == "" {
		app.RubixServiceDir = filePath(fmt.Sprintf("%s/rubix-service", app.DataDir))
	}
	if app.AppsInstallDir == "" {
		app.AppsInstallDir = filePath(fmt.Sprintf("%s/rubix-service/apps/install", app.DataDir))
	}
	if app.TmpDir == "" {
		app.TmpDir = filePath(fmt.Sprintf("%s/tmp", app.DataDir))
	}
	if app.BackupsDir == "" {
		app.BackupsDir = fmt.Sprintf("%s/backup", homeDir)
	}
	app.SystemCtl = systemctl.New(false, app.DefaultTimeout)
	return app
}
