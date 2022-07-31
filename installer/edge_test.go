package installer

import (
	"fmt"
	pprint "github.com/NubeIO/lib-rubix-installer/print"
	"testing"
)

var appName = "flow-framework"
var serviceName = "nubeio-flow-framework"

func Test_ConfirmAppInstalled(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FilePerm: nonRoot})
	installed, err := app.ConfirmAppInstalled(appName, serviceName)
	fmt.Println(err)
	pprint.PrintJOSN(installed)

}

func Test_GetAppVersion(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FilePerm: nonRoot})
	version := app.GetAppVersion(appName)
	fmt.Println(version)
}

func Test_GetApps(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FilePerm: nonRoot})
	installed, err := app.GetApps()
	fmt.Println(err)
	if err != nil {
		return
	}
	pprint.PrintJOSN(installed)
}

func Test_BackupApp(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FilePerm: nonRoot})
	installed, err := app.GetApps()
	fmt.Println(err)
	if err != nil {
		return
	}
	pprint.PrintJOSN(installed)
}

func Test_BackupApps(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FilePerm: nonRoot})
	installed, err := app.GetApps()
	fmt.Println(err)
	if err != nil {
		return
	}
	pprint.PrintJOSN(installed)
}
