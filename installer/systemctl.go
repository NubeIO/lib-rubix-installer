package installer

import (
	"errors"
	"github.com/NubeIO/lib-systemctl-go/systemctl"
)

type CtlBody struct {
	AppName      string   `json:"app_name"`
	Service      string   `json:"service"`
	Action       string   `json:"action"`
	Timeout      int      `json:"timeout"`
	ServiceNames []string `json:"service_names"` // nubeio-flow-framework.service
	AppNames     []string `json:"app_names"`     // flow-framework
}

func (inst *App) CtlAction(body *CtlBody) (*systemctl.SystemResponse, error) {
	if body.AppName != "" { // if user passes in the appName then get the serviceFile
		name, err := inst.GetNubeServiceFileName(body.AppName)
		if err != nil {
			return nil, err
		}
		body.Service = name
	}
	return inst.Ctl.CtlAction(body.Action, body.Service, body.Timeout)
}

func (inst *App) CtlStatus(body *CtlBody) (*systemctl.SystemState, error) {
	if body.AppName != "" {
		name, err := inst.GetNubeServiceFileName(body.AppName)
		if err != nil {
			return nil, err
		}
		body.Service = name
	}
	resp, err := inst.Ctl.ServiceState(body.Service, body.Timeout)
	return &resp, err
}

func (inst *App) ServiceMassAction(body *CtlBody) ([]systemctl.MassSystemResponse, error) {
	if len(body.AppNames) > 0 {
		for _, name := range body.AppNames { // if user passes in the appName then get the serviceFile
			serviceName, err := inst.GetNubeServiceFileName(name)
			if err != nil {
				return nil, err
			}
			body.ServiceNames = append(body.ServiceNames, serviceName)
		}
	}
	if len(body.ServiceNames) == 0 {
		return nil, errors.New("no services names provided")
	}
	return inst.Ctl.ServiceMassAction(body.ServiceNames, body.Action, body.Timeout)
}

func (inst *App) ServiceMassStatus(body *CtlBody) ([]systemctl.SystemState, error) {
	if len(body.AppNames) > 0 {
		for _, name := range body.AppNames {
			serviceName, err := inst.GetNubeServiceFileName(name)
			if err != nil {
				return nil, err
			}
			body.ServiceNames = append(body.ServiceNames, serviceName)
		}
	}
	if len(body.ServiceNames) == 0 {
		return nil, errors.New("no services names provided")
	}
	return inst.Ctl.ServiceStateMass(body.ServiceNames, body.Timeout)
}
