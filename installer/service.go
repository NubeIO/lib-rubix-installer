package installer

import (
	"errors"
	"fmt"
	fileutils "github.com/NubeIO/lib-dirs/dirs"
	"github.com/NubeIO/lib-systemctl-go/ctl"
	"github.com/NubeIO/lib-systemctl-go/systemctl"
	log "github.com/sirupsen/logrus"
	"time"
)

func (inst *App) InstallService(app *Install) (*InstallResp, error) {
	var serviceName = app.ServiceName
	var serviceFilePath = app.Source
	if serviceName == "" {
		return nil, errors.New("service name can not be empty")
	}
	if serviceFilePath == "" {
		return nil, errors.New("service file path can not be empty")
	}
	found := fileutils.New().FileExists(serviceFilePath)
	if !found {
		return nil, errors.New(fmt.Sprintf("no service file found in path:%s", serviceFilePath))
	}
	found = inst.ConfirmAppDir(app.Name)
	if !found {
		return nil, errors.New(fmt.Sprintf("no app dir found for provided app:%s", app.Name))
	}
	found = inst.ConfirmAppInstallDir(app.BuildName)
	if !found {
		return nil, errors.New(fmt.Sprintf("no app install dir found for provided app:%s", app.BuildName))
	}
	found = inst.ConfirmAppInstallDir(app.BuildName)
	if !found {
		return nil, errors.New(fmt.Sprintf("no app install dir found for provided app:%s", app.BuildName))
	}
	return inst.installService(serviceName, serviceFilePath)
}

// InstallService a new linux service
//	- service: the service name (eg: nubeio-rubix-wires)
//	- path: the service file path and name (eg: "/tmp/rubix-bios.service")
func (inst *App) installService(service, tmpServiceFile string) (*InstallResp, error) {
	var err error
	ser := ctl.New(service, tmpServiceFile)
	ser.InstallOpts = ctl.InstallOpts{
		Options: systemctl.Options{Timeout: inst.DefaultTimeout},
	}
	err = ser.TransferFile()
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

//Install a new service
func (inst *App) systemCtlInstall(service string) (*InstallResp, error) {
	resp := &InstallResp{
		Install: "install ok",
	}
	var ok = "action ok"
	//reload
	err := systemctl.DaemonReload(systemOpts)
	if err != nil {
		log.Errorf("failed to DaemonReload%s: err:%s", service, err.Error())
		resp.DaemonReload = err.Error()
		return resp, err
	} else {
		resp.DaemonReload = ok
	}
	//enable
	err = systemctl.Enable(service, systemOpts)
	if err != nil {
		log.Errorf("failed to enable%s: err:%s", service, err.Error())
		resp.Enable = err.Error()
		return resp, err
	} else {
		resp.Enable = ok
	}
	log.Infof("enable new service:%s", service)
	//start
	err = systemctl.Restart(service, systemOpts)
	if err != nil {
		log.Errorf("failed to start%s: err:%s", service, err.Error())
		resp.Restart = err.Error()
		return resp, err
	} else {
		resp.Restart = ok
	}
	log.Infof("start new service:%s", service)
	time.Sleep(8 * time.Second)
	active, status, err := systemctl.IsRunning(service, systemctl.Options{})
	if err != nil {
		log.Errorf("service found or failed to check IsRunning: %s: %v", service, err)
		return nil, err
	} else {
		resp.CheckIsRunning = true
	}
	log.Infof("service: %s: isActive: %t status: %s", service, active, status)
	return resp, nil
}
