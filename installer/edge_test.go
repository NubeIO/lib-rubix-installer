package installer

import (
	"fmt"
	fileutils "github.com/NubeIO/lib-dirs/dirs"
	pprint "github.com/NubeIO/lib-rubix-installer/print"
	"testing"
)

func Test_EdgeCheckInstallApps(t *testing.T) {
	var err error
	homeDir, err := fileutils.Dir()
	fmt.Println(homeDir, err)
	app := New(&App{DataDir: "/data", FilePerm: nonRoot})
	installed, err := app.DiscoverInstalled()
	fmt.Println(err)
	if err != nil {
		return
	}
	pprint.PrintJOSN(installed)

}
