package installer

import (
	"fmt"
	fileutils "github.com/NubeIO/lib-dirs/dirs"
	"testing"
)

func TestApp_uploadApp(t *testing.T) {
	//var err error
	homeDir, _ := fileutils.Dir()
	fmt.Println(homeDir)
	app := New(&App{DataDir: "/data", FilePerm: nonRoot, HostDownloadPath: fmt.Sprintf("%s/Downloads", homeDir)})

	appName := "rubix-wires"
	appInstallName := "wires-builds"
	serviceName := "nubeio-rubix-wires"
	appVersion := "v2.7.2"
	appZip := "wires-builds-2.7.2.zip"

	fmt.Println(app, appName, appInstallName, serviceName, appVersion, appZip)

	//err = app.uploadApp(appName, appInstallName, appVersion)
	//
	//fmt.Println(err)
}
