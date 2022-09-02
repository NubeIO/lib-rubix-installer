package installer

import (
	"fmt"
	"testing"
)

func Test_UploadApp(t *testing.T) {
	app := New(&App{DataDir: "/data", FilePerm: filePerm})
	zip, err := app.unzip("/data/tmp/tmp_D2F8D40F77F8/nubeio-rubix-app-lora-serial-py-0.0.1-fefc55ac.armv7.zip", "/data/rubix-service/apps/install/lora-driver/v0.0.1")
	fmt.Println(err)
	fmt.Println(zip)
	if err != nil {
		return
	}
}
