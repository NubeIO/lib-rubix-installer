package installer

import (
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/lib-systemctl-go/systemctl"
	"path"
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
		app.DataDir = "/data"
	}
	if app.FileMode == 0 {
		app.FileMode = fileMode
	}
	if app.DefaultTimeout == 0 {
		app.DefaultTimeout = defaultTimeout
	}
	if app.StoreDir == "" {
		app.StoreDir = path.Join(app.DataDir, "store")
	}
	if app.UserRubixHome == "" {
		app.UserRubixHome = path.Join(homeDir, "rubix")
	}
	if app.RubixServiceDir == "" {
		app.RubixServiceDir = path.Join(app.DataDir, "rubix-service")
	}
	if app.AppsInstallDir == "" {
		app.AppsInstallDir = path.Join(app.DataDir, "rubix-service/apps/install")
	}
	if app.TmpDir == "" {
		app.TmpDir = path.Join(app.DataDir, "tmp")
	}
	if app.BackupsDir == "" {
		app.BackupsDir = path.Join(homeDir, "backup")
	}
	app.SystemCtl = systemctl.New(false, app.DefaultTimeout)
	return app
}
