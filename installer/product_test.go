package installer

import (
	"fmt"
	pprint "github.com/NubeIO/lib-rubix-installer/print"
	"testing"
)

func TestGetProduct(t *testing.T) {

	p, err := GetProduct()
	fmt.Println(err)
	pprint.PrintJOSN(p)

}
