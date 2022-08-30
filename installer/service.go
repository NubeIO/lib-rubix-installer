package installer

import (
	"errors"
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/lib-systemctl-go/ctl"
	"github.com/NubeIO/lib-systemctl-go/systemctl"
	log "github.com/sirupsen/logrus"
	"time"
)

func (inst *App) SetServiceFileName(appName string) string {
	return fmt.Sprintf("nubeio-%s.service", appName)
}

func (inst *App) setServiceFileName(appName string) string {
	return fmt.Sprintf("nubeio-%s.service", appName)
}

func (inst *App) InstallService(app *Install) (*InstallResp, error) {
	var serviceName = inst.setServiceFileName(app.Name)
	if app.ServiceName != "" {
		serviceName = app.ServiceName
	}
	if serviceName == "" {
		return nil, errors.New("service name can not be empty, nubeio-flow-framework.service")
	}
	var serviceFilePath = app.Source
	if serviceFilePath == "" {
		return nil, errors.New("service file path can not be empty, /data/tmp/tmp_B8CB4D888176/nubeio-flow-framework.service")
	}
	found := fileutils.New().FileExists(serviceFilePath)
	if !found {
		return nil, errors.New(fmt.Sprintf("no service file found in path: %s", serviceFilePath))
	}
	found = inst.ConfirmAppDir(app.Name)
	if !found {
		return nil, errors.New(fmt.Sprintf("no app dir found for provided app: %s", app.Name))
	}
	found = inst.ConfirmAppInstallDir(app.Name)
	if !found {
		return nil, errors.New(fmt.Sprintf("no app install dir found for provided app: %s", app.Name))
	}
	return inst.installService(serviceName, serviceFilePath)
}

// InstallService a new linux service
//	- service: the service name (eg: nubeio-rubix-wires)
//	- path: the service file path and name (eg: "/tmp/rubix-bios.service")
func (inst *App) installService(service, tmpServiceFile string) (*InstallResp, error) {
	var err error
	ser := ctl.New(service, false, inst.DefaultTimeout)
	err = ser.TransferSystemdFile(tmpServiceFile)
	if err != nil {
		fmt.Println("full install error", err)
		return nil, err
	}
	return inst.systemCtlInstall(service)
}

type InstallResp struct {
	Install        string `json:"installed"`
	DaemonReload   string `json:"daemon_reload"`
	Enable         string `json:"enabled"`
	Restart        string `json:"restarted"`
	CheckIsRunning bool   `json:"check_is_running"`
}

// Install a new service
func (inst *App) systemCtlInstall(service string) (*InstallResp, error) {
	resp := &InstallResp{
		Install: "install ok",
	}
	systemCtl := systemctl.New(false, inst.DefaultTimeout)
	var ok = "action ok"
	// reload
	err := systemCtl.DaemonReload()
	if err != nil {
		log.Errorf("failed to DaemonReload%s: err: %s", service, err.Error())
		resp.DaemonReload = err.Error()
		return resp, err
	} else {
		resp.DaemonReload = ok
	}
	// enable
	err = systemCtl.Enable(service)
	if err != nil {
		log.Errorf("failed to enable%s: err: %s", service, err.Error())
		resp.Enable = err.Error()
		return resp, err
	} else {
		resp.Enable = ok
	}
	log.Infof("enable new service: %s", service)
	// start
	err = systemCtl.Restart(service)
	if err != nil {
		log.Errorf("failed to start %s: err: %s", service, err.Error())
		resp.Restart = err.Error()
		return resp, err
	} else {
		resp.Restart = ok
	}
	log.Infof("start new service: %s", service)
	time.Sleep(8 * time.Second)
	active, status, err := systemCtl.IsRunning(service)
	if err != nil {
		log.Errorf("service found or failed to check IsRunning: %s: %v", service, err)
		return nil, err
	} else {
		resp.CheckIsRunning = true
	}
	log.Infof("service: %s: isActive: %t status: %s", service, active, status)
	return resp, nil
}
