package installer

import (
	"fmt"
	fileutils "github.com/NubeIO/lib-dirs/dirs"
	"testing"
)

func TestApp_MoveFile(t *testing.T) {
	var err error
	homeDir, err := fileutils.Dir()
	fmt.Println(homeDir, err)
	app := New(&App{DataDir: "/data", FilePerm: nonRoot, HostDownloadPath: fmt.Sprintf("%s/Downloads", homeDir)})

	err = app.MoveFile("/home/aidan/Downloads/nubeio-rubix-app-lora-serial-py-0.0.1-fefc55ac.amd64/nubeio-rubix-app-lora-serial-py", "/home/aidan/Downloads/nubeio-rubix-app-lora-serial-py-0.0.1-fefc55ac.amd64/app", true)
	fmt.Println(err)
	if err != nil {
		return
	}
}
