package installer

import (
	"errors"
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/lib-systemctl-go/systemd"
)

func (inst *App) InstallService(app *Install) (*systemd.InstallResponse, error) {
	if app.ServiceName == "" {
		return nil, errors.New("service_name can not be empty, for example: nubeio-flow-framework.service")
	}
	if app.Source == "" {
		return nil, errors.New("service source file path can not be empty, for example: /data/tmp/tmp_b8cb4d888176/nubeio-flow-framework.service")
	}
	found := fileutils.FileExists(app.Source)
	if !found {
		return nil, errors.New(fmt.Sprintf("no service file found in path: %s", app.Source))
	}
	found = inst.ConfirmAppDir(app.Name)
	if !found {
		return nil, errors.New(fmt.Sprintf("no app dir found for provided app: %s", app.Name))
	}
	found = inst.ConfirmAppInstallDir(app.Name)
	if !found {
		return nil, errors.New(fmt.Sprintf("no app install dir found for provided app: %s", app.Name))
	}
	return inst.installService(app.ServiceName, app.Source)
}

// InstallService a new linux service
//	- service: the service_name (eg: nubeio-rubix-wires)
//	- path: the service file path and name (eg: "/tmp/rubix-bios.service")
func (inst *App) installService(serviceName, tmpServiceFile string) (*systemd.InstallResponse, error) {
	systemdService := systemd.New(serviceName, false, inst.DefaultTimeout)
	err := systemdService.TransferSystemdFile(tmpServiceFile)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return systemdService.Install(), nil
}
