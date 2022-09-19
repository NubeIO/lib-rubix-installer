package installer

import (
	"errors"
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/lib-systemctl-go/systemd"
)

func (inst *App) InstallService(app *Install) (*systemd.InstallResponse, error) {
	if app.Source == "" {
		return nil, errors.New("service source file path can not be empty, for example: /data/tmp/tmp_1234/nubeio-flow-framework.service")
	}
	found := fileutils.FileExists(app.Source)
	if !found {
		return nil, errors.New(fmt.Sprintf("no service file found in path: %s", app.Source))
	}
	found = fileutils.DirExists(inst.GetAppDataPath(app.Name))
	if !found {
		return nil, errors.New(fmt.Sprintf("no app dir found for provided app: %s", app.Name))
	}
	found = fileutils.DirExists(inst.GetAppInstallPath(app.Name))
	if !found {
		return nil, errors.New(fmt.Sprintf("no app install dir found for provided app: %s", app.Name))
	}
	serviceName := inst.GetServiceNameFromAppName(app.Name)
	return inst.installService(serviceName, app.Source)
}

// InstallService a new linux service
//	- serviceName: the service_name (eg: nubeio-rubix-wires.service)
//	- tmpServiceFile: the service file path and name (eg: "/tmp/rubix-bios.service")
func (inst *App) installService(serviceName, tmpServiceFile string) (*systemd.InstallResponse, error) {
	systemdService := systemd.New(serviceName, false, inst.DefaultTimeout)
	err := systemdService.TransferSystemdFile(tmpServiceFile)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return systemdService.Install(), nil
}
