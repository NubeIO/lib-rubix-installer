package installer

import (
	"fmt"
	pprint "github.com/NubeIO/lib-rubix-installer/print"
	"testing"
)

func TestGetProduct(t *testing.T) {

	var err error
	fmt.Println(err)
	app := New(&App{DataDir: "/data", FilePerm: nonRoot})
	installed, err := app.GetProduct()
	fmt.Println(err)
	pprint.PrintJOSN(installed)

}
