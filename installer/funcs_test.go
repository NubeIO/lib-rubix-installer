package installer

import (
	"fmt"
	"testing"
)

func Test_unzip(t *testing.T) {
	err := unzip("/Users/raibnod/Downloads/bgis-users.zip", "/Users/raibnod/Downloads/test2")
	if err != nil {
		fmt.Println("Error:", err)
	}
}
