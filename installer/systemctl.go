package installer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeIO/lib-systemctl-go/systemctl"
)

type SystemCtlBody struct {
	AppName      string   `json:"app_name"`      // "flow-framework"
	ServiceName  string   `json:"service_name"`  // "nubeio-flow-framework.service"
	Action       string   `json:"action"`        // start, stop, restart, enable, disable
	AppNames     []string `json:"app_names"`     // ["nubeio-rubix-edge.service", "nubeio-flow-framework.service"]
	ServiceNames []string `json:"service_names"` // ["rubix-edge", "flow-framework"]
}

type SystemResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type AppSystemState struct {
	ServiceName            string                  `json:"service_name,omitempty"`
	AppName                string                  `json:"app_name,omitempty"`
	State                  systemctl.UnitFileState `json:"state,omitempty"`        // enabled, disabled
	ActiveState            systemctl.ActiveState   `json:"active_state,omitempty"` // active, inactive
	SubState               systemctl.SubState      `json:"sub_state,omitempty"`    // running, dead
	ActiveEnterTimestamp   string                  `json:"active_enter_timestamp,omitempty"`
	InactiveEnterTimestamp string                  `json:"inactive_enter_timestamp,omitempty"`
	Restarts               string                  `json:"restarts,omitempty"` // number of restart
	IsInstalled            bool                    `json:"is_installed,omitempty"`
}

type MassSystemResponse struct {
	AppName     string `json:"app_name"`
	ServiceName string `json:"service_name"`
	Success     bool   `json:"success"`
	Message     string `json:"message"`
}

func (inst *App) SystemCtlAction(body *SystemCtlBody) (*SystemResponse, error) {
	resp := &SystemResponse{}
	if body.AppName == "" && body.ServiceName == "" {
		return nil, errors.New("app_name & service_name both can not be empty")
	}
	var serviceName string
	if body.ServiceName != "" {
		serviceName = body.ServiceName
	} else {
		serviceName = inst.GetServiceNameFromAppName(body.AppName)
	}
	var err error
	switch body.Action {
	case "start":
		err = inst.SystemCtl.Start(serviceName)
	case "stop":
		err = inst.SystemCtl.Stop(serviceName)
	case "restart":
		err = inst.SystemCtl.Restart(serviceName)
	case "enable":
		err = inst.SystemCtl.Enable(serviceName)
	case "disable":
		err = inst.SystemCtl.Disable(serviceName)
	default:
		err = errors.New("no valid action found try, start, stop, restart, enable or disable")
	}
	if err != nil {
		resp.Success = false
		resp.Message = err.Error()
	} else {
		resp.Success = true
		resp.Message = fmt.Sprintf("service: %s, action: %s is executed succefully", serviceName, body.Action)
	}
	return resp, nil
}

func (inst *App) SystemCtlStatus(body *SystemCtlBody) (*AppSystemState, error) {
	if body.AppName == "" && body.ServiceName == "" {
		return nil, errors.New("app_name and service_name can not be empty")
	}
	var serviceName string
	if body.ServiceName != "" {
		serviceName = body.ServiceName
	} else {
		serviceName = inst.GetServiceNameFromAppName(body.AppName)
	}
	systemState, err := inst.SystemCtl.State(serviceName)
	appSystemState := AppSystemState{}
	systemStateJson, _ := json.Marshal(systemState)
	_ = json.Unmarshal(systemStateJson, &appSystemState)
	appSystemState.AppName = body.AppName
	return &appSystemState, err
}

func (inst *App) ServiceMassAction(body *SystemCtlBody) ([]MassSystemResponse, error) {
	if len(body.AppNames) == 0 && len(body.ServiceNames) == 0 {
		return nil, errors.New("no app_names & service_names are provided")
	}
	var outputs []MassSystemResponse
	for _, appName := range body.AppNames {
		input := SystemCtlBody{AppName: appName, Action: body.Action}
		response, _ := inst.SystemCtlAction(&input)
		output := MassSystemResponse{AppName: appName, Success: response.Success, Message: response.Message}
		outputs = append(outputs, output)
	}
	for _, serviceName := range body.ServiceNames {
		input := SystemCtlBody{ServiceName: serviceName, Action: body.Action}
		response, _ := inst.SystemCtlAction(&input)
		output := MassSystemResponse{ServiceName: serviceName, Success: response.Success, Message: response.Message}
		outputs = append(outputs, output)
	}
	return outputs, nil
}

func (inst *App) ServiceMassStatus(body *SystemCtlBody) ([]AppSystemState, error) {
	if len(body.AppNames) == 0 && len(body.ServiceNames) == 0 {
		return nil, errors.New("no app_names & service_names are provided")
	}
	var outputs []AppSystemState
	for _, appName := range body.AppNames {
		input := SystemCtlBody{AppName: appName, Action: body.Action}
		response, _ := inst.SystemCtlStatus(&input)
		outputs = append(outputs, *response)
	}
	for _, serviceName := range body.ServiceNames {
		input := SystemCtlBody{ServiceName: serviceName, Action: body.Action}
		response, _ := inst.SystemCtlStatus(&input)
		outputs = append(outputs, *response)
	}
	return outputs, nil
}
