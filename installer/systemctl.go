package installer

import (
	"errors"
	"fmt"
	"github.com/NubeIO/lib-systemctl-go/systemctl"
)

type CtlBody struct {
	AppName      string   `json:"app_name"`
	ServiceName  string   `json:"service_name"`
	Action       string   `json:"action"`
	ServiceNames []string `json:"service_names"` // nubeio-flow-framework.service
	AppNames     []string `json:"app_names"`     // flow-framework
}

type SystemResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type MassSystemResponse struct {
	ServiceName string `json:"service_name"`
	Success     bool   `json:"success"`
	Message     string `json:"message"`
}

func (inst *App) SystemCtlAction(body *CtlBody) (*SystemResponse, error) {
	resp := &SystemResponse{}
	if body.ServiceName == "" {
		return nil, errors.New("service_name can not be empty")
	}
	var err error
	switch body.Action {
	case "start":
		err = inst.SystemCtl.Start(body.ServiceName)
	case "stop":
		err = inst.SystemCtl.Stop(body.ServiceName)
	case "enable":
		err = inst.SystemCtl.Enable(body.ServiceName)
	case "disable":
		err = inst.SystemCtl.Disable(body.ServiceName)
	case "restart":
		err = inst.SystemCtl.Restart(body.ServiceName)
	default:
		err = errors.New("no valid action found try, start, stop, restart, enable or disable")
	}
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
	} else {
		resp.Success = true
		resp.Message = fmt.Sprintf("service: %s, action: %s is executed succefully", body.ServiceName, body.Action)
	}
	return resp, nil
}

func (inst *App) SystemCtlStatus(body *CtlBody) (*systemctl.SystemState, error) {
	if body.ServiceName == "" {
		return nil, errors.New("service_name can not be empty")
	}
	resp, err := inst.SystemCtl.State(body.ServiceName)
	return &resp, err
}

func (inst *App) ServiceMassAction(body *CtlBody) ([]MassSystemResponse, error) {
	if len(body.ServiceNames) == 0 {
		return nil, errors.New("no service_names provided")
	}
	var outputs []MassSystemResponse
	for _, serviceName := range body.ServiceNames {
		body.ServiceName = serviceName
		response, _ := inst.SystemCtlAction(body)
		output := MassSystemResponse{ServiceName: serviceName, Success: response.Success, Message: response.Message}
		outputs = append(outputs, output)
	}
	return outputs, nil
}

func (inst *App) ServiceMassStatus(body *CtlBody) ([]systemctl.SystemState, error) {
	if len(body.ServiceNames) == 0 {
		return nil, errors.New("no service_names provided")
	}
	var outputs []systemctl.SystemState
	for _, serviceName := range body.ServiceNames {
		body.ServiceName = serviceName
		response, _ := inst.SystemCtlStatus(body)
		outputs = append(outputs, *response)
	}
	return outputs, nil
}
