package installer

import (
	"fmt"
	"testing"
)

func TestApp_uploadApp(t *testing.T) {
	//var err error
	///data/tmp/tmp_D2F8D40F77F8/nubeio-rubix-app-lora-serial-py-0.0.1-fefc55ac.armv7.zip dest:/data/rubix-service/apps/install/lora-driver/v0.0.1
	//err = app.uploadApp(appName, appInstallName, appVersion)
	//
	//fmt.Println(err)

	app := New(&App{DataDir: "/data", FilePerm: nonRoot})

	zip, err := app.unZip("/data/tmp/tmp_D2F8D40F77F8/nubeio-rubix-app-lora-serial-py-0.0.1-fefc55ac.armv7.zip", "/data/rubix-service/apps/install/lora-driver/v0.0.1")
	fmt.Println(err)
	fmt.Println(zip)
	if err != nil {
		return
	}

}
