package installer

import (
	"fmt"
	fileutils "github.com/NubeIO/lib-dirs/dirs"
	"testing"
)

func Test_checks(t *testing.T) {
	var err error
	homeDir, err := fileutils.Dir()
	fmt.Println(homeDir, err)
	app := New(&App{DataDir: "/data", FilePerm: nonRoot, HostDownloadPath: fmt.Sprintf("%s/Downloads", homeDir)})

	appName := "rubix-wires"
	//appInstallName := "wires-builds"
	serviceName := "nubeio-rubix-wires"
	//appVersion := "v2.7.2"
	//appZip := "wires-builds-2.7.2.zip"
	isInstalled := app.ConfirmAppDir(appName)
	fmt.Println(isInstalled)
	installed := app.ConfirmAppInstalled(appName, serviceName)
	print(installed)

}
