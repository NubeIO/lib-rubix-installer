package installer

import (
	"fmt"
	fileutils "github.com/NubeIO/lib-dirs/dirs"
	"testing"
)

func Test_checkVersion(t *testing.T) {
	var err error
	homeDir, _ := fileutils.Dir()
	fmt.Println(homeDir)
	app := New(&App{DataDir: "/data", FilePerm: nonRoot, HostDownloadPath: fmt.Sprintf("%s/Downloads", homeDir)})

	//appName := "rubix-wires"
	appInstallName := "wires-builds"
	serviceName := "nubeio-rubix-wires"
	//appVersion := "v2.7.2"
	//appZip := "wires-builds-2.7.2.zip"
	version := app.GetAppVersion(appInstallName)

	fmt.Println(version)

	files, err := app.listFiles(fmt.Sprintf("%s/%s", "/data", appInstallName))
	fmt.Println(err)
	if err != nil {
		return
	}

	fmt.Println(files)
	fmt.Println(serviceName)

	//uninstall, err := app.UninstallService(appName, appInstallName, serviceName)
	//fmt.Println(err)
	//pprint.PrintJOSN(uninstall)

}
