package installer

import (
	"fmt"
	"github.com/NubeIO/lib-rubix-installer/pprint"
	"testing"
)

func Test_GetProduct(t *testing.T) {
	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FilePerm: filePerm})
	installed, err := app.GetProduct()
	fmt.Println(err)
	pprint.PrintJSON(installed)
}
