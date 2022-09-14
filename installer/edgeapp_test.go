package installer

import (
	"fmt"
	"github.com/NubeIO/lib-rubix-installer/pprint"
	"testing"
)

func Test_GetAppVersion(t *testing.T) {
	var appName = "flow-framework"
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FileMode: fileMode})
	version := app.GetAppVersion(appName)
	fmt.Println(version)
}

func Test_ListApps(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FileMode: fileMode})
	installed, err := app.ListApps()
	fmt.Println(err)
	if err != nil {
		return
	}
	pprint.PrintJSON(installed)
}

func Test_ListAppsStatus(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FileMode: fileMode})
	appServiceMapping := map[string]string{}
	appServiceMapping["flow-framework"] = "nubeio-flow-framework.service"
	installed, err := app.ListAppsStatus(appServiceMapping)
	fmt.Println(installed, err)
	if err != nil {
		return
	}
	pprint.PrintJSON(installed)
}
