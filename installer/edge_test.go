package installer

import (
	"fmt"
	pprint "github.com/NubeIO/lib-rubix-installer/print"
	"testing"
)

var appName = "flow-framework"
var serviceName = "nubeio-flow-framework"

func Test_ConfirmAppInstalled(t *testing.T) {
	//var err error
	//fmt.Println(err)
	//app := New(&App{DataDir: "/data", FilePerm: nonRoot})
	//installed, err := app.ConfirmAppInstalled(appName, serviceName)
	//fmt.Println(err)
	//pprint.PrintJOSN(installed)

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
	installed, err := app.ListApps()
	fmt.Println(err)
	if err != nil {
		return
	}
	pprint.PrintJOSN(installed)
}

func Test_ListNubeServicesFiles(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FilePerm: nonRoot})
	installed, err := app.ListNubeServiceFiles()
	fmt.Println(err)
	if err != nil {
		return
	}
	pprint.PrintJOSN(installed)
}
func Test_getNubeServiceFileName(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FilePerm: nonRoot})
	installed, err := app.GetNubeServiceFileName(appName)
	fmt.Println(err)
	if err != nil {
		return
	}
	pprint.PrintJOSN(installed)
}
func Test_ListNubeServices(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FilePerm: nonRoot})
	installed, err := app.ListNubeServices()
	fmt.Println(installed, err)
	if err != nil {
		return
	}
	pprint.PrintJOSN(installed)
}

func Test_ListAppsAndService(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FilePerm: nonRoot})
	installed, err := app.ListAppsAndService()
	fmt.Println(installed, err)
	if err != nil {
		return
	}
	pprint.PrintJOSN(installed)
}
