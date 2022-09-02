package installer

import (
	"errors"
	"fmt"
	"github.com/NubeIO/lib-systemctl-go/systemctl"
)

type SystemCtlBody struct {
	ServiceName  string   `json:"service_name"`  // "nubeio-flow-framework.service"
	Action       string   `json:"action"`        // start, stop, restart, enable, disable
	ServiceNames []string `json:"service_names"` // ["nubeio-rubix-edge.service", "nubeio-flow-framework.service"]
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

func (inst *App) SystemCtlAction(body *SystemCtlBody) (*SystemResponse, error) {
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
	case "restart":
		err = inst.SystemCtl.Restart(body.ServiceName)
	case "enable":
		err = inst.SystemCtl.Enable(body.ServiceName)
	case "disable":
		err = inst.SystemCtl.Disable(body.ServiceName)
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

func (inst *App) SystemCtlStatus(body *SystemCtlBody) (*systemctl.SystemState, error) {
	if body.ServiceName == "" {
		return nil, errors.New("service_name can not be empty")
	}
	resp, err := inst.SystemCtl.State(body.ServiceName)
	return &resp, err
}

func (inst *App) ServiceMassAction(body *SystemCtlBody) ([]MassSystemResponse, error) {
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

func (inst *App) ServiceMassStatus(body *SystemCtlBody) ([]systemctl.SystemState, error) {
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
