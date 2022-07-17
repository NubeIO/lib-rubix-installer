package installer

import (
	"errors"
	"github.com/NubeIO/lib-systemctl-go/systemctl"
)

func (inst *App) SystemCtlAction(serviceName string, action string, timeout int) (*SystemResponse, error) {
	actionResp := &SystemResponse{}
	switch action {
	case start.String():
		return inst.Start(serviceName, timeout)
	case stop.String():
		return inst.Stop(serviceName, timeout)
	case enable.String():
		return inst.Enable(serviceName, timeout)
	case disable.String():
		return inst.Disable(serviceName, timeout)
	}
	return actionResp, errors.New("no valid action found try, start, stop, enable or disable")
}

func (inst *App) SystemCtlStatus(serviceName, action string, timeout int) (*SystemResponseChecks, error) {
	actionResp := &SystemResponseChecks{}
	switch action {
	case isRunning.String():
		return inst.IsRunning(serviceName, timeout)
	case isInstalled.String():
		return inst.IsInstalled(serviceName, timeout)
	case isEnabled.String():
		return inst.IsEnabled(serviceName, timeout)
	case isActive.String():
		return inst.IsActive(serviceName, timeout)
	case isFailed.String():
		return inst.IsFailed(serviceName, timeout)
	}
	return actionResp, errors.New("no valid action found try, isRunning, isInstalled, isEnabled, isActive or isFailed")
}

type Mass struct {
	Action  string
	Timeout int `json:"timeout"`
}

type massResponse struct {
	AppName string `json:"app_name"`
	Action  string `json:"action"`
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
	Err     error  `json:"err"`
}

// ServiceMassAction mass start, stop, enable, disable, a service
func (inst *App) ServiceMassAction(serviceNames []string, action string, timeout int) ([]SystemResponse, error) {
	var out []SystemResponse
	for _, name := range serviceNames {
		ctlAction, _ := inst.SystemCtlAction(name, action, timeout)
		out = append(out, *ctlAction)
	}
	return out, nil
}

// ServiceMassCheck check if a service isRunning, isEnabled and so on
func (inst *App) ServiceMassCheck(serviceNames []string, action string, timeout int) ([]SystemResponseChecks, error) {
	var out []SystemResponseChecks
	for _, name := range serviceNames {
		ctlAction, _ := inst.SystemCtlStatus(name, action, timeout)
		out = append(out, *ctlAction)
	}
	return out, nil
}

type SystemResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

func (inst *App) Start(serviceName string, timeout int) (resp *SystemResponse, err error) {
	resp = &SystemResponse{}
	systemOpts.Timeout = timeout
	err = systemctl.Start(serviceName, systemOpts)
	if err != nil {
		resp.Message = "tried to start but failed"
		return resp, err
	}
	resp.Ok = true
	resp.Message = "start ok"
	return resp, nil
}

func (inst *App) Stop(serviceName string, timeout int) (resp *SystemResponse, err error) {
	resp = &SystemResponse{}
	systemOpts.Timeout = timeout
	err = systemctl.Stop(serviceName, systemOpts)
	if err != nil {
		resp.Message = "tried to stop but failed"
		return resp, err
	}
	resp.Ok = true
	resp.Message = "stop ok"
	return resp, nil
}

func (inst *App) Enable(serviceName string, timeout int) (resp *SystemResponse, err error) {
	resp = &SystemResponse{}
	systemOpts.Timeout = timeout
	err = systemctl.Enable(serviceName, systemOpts)
	if err != nil {
		resp.Message = "tried to enable but failed"
		return resp, err
	}
	resp.Ok = true
	resp.Message = "enable ok"
	return resp, nil
}

func (inst *App) Disable(serviceName string, timeout int) (resp *SystemResponse, err error) {
	resp = &SystemResponse{}
	systemOpts.Timeout = timeout
	err = systemctl.Disable(serviceName, systemOpts)
	if err != nil {
		resp.Message = "tried to disable but failed"
		return resp, err
	}
	resp.Ok = true
	resp.Message = "disable ok"
	return resp, nil
}

type SystemResponseChecks struct {
	Is      bool   `json:"is"`
	Message string `json:"message"`
}

func (inst *App) IsEnabled(serviceName string, timeout int) (resp *SystemResponseChecks, err error) {
	resp = &SystemResponseChecks{}
	systemOpts.Timeout = timeout
	out, err := systemctl.IsEnabled(serviceName, systemOpts)
	if err != nil || out == false {
		resp.Message = "is not enabled"
		return resp, err
	}
	resp.Is = out
	resp.Message = "is enabled"
	return
}

func (inst *App) IsFailed(serviceName string, timeout int) (resp *SystemResponseChecks, err error) {
	resp = &SystemResponseChecks{}
	systemOpts.Timeout = timeout
	out, err := systemctl.IsFailed(serviceName, systemOpts)
	if err != nil || out == true {
		resp.Message = "is failed"
		return resp, err
	}
	resp.Is = out
	resp.Message = "is not failed"
	return
}

func (inst *App) IsInstalled(serviceName string, timeout int) (resp *SystemResponseChecks, err error) {
	resp = &SystemResponseChecks{}
	systemOpts.Timeout = timeout
	out, err := systemctl.IsInstalled(serviceName, systemOpts)
	if err != nil || out == false {
		resp.Message = "is not installed"
		return resp, err
	}
	resp.Is = out
	resp.Message = "is installed"
	return
}

func (inst *App) IsActive(serviceName string, timeout int) (resp *SystemResponseChecks, err error) {
	resp = &SystemResponseChecks{}
	systemOpts.Timeout = timeout
	out, msg, err := systemctl.IsActive(serviceName, systemOpts)
	if err != nil || out == false {
		resp.Message = msg
		return resp, err
	}
	resp.Is = out
	resp.Message = msg
	return
}

func (inst *App) IsRunning(serviceName string, timeout int) (resp *SystemResponseChecks, err error) {
	resp = &SystemResponseChecks{}
	systemOpts.Timeout = timeout
	out, msg, err := systemctl.IsRunning(serviceName, systemOpts)
	if err != nil || out == false {
		resp.Message = msg
		return resp, err
	}
	resp.Is = out
	resp.Message = msg
	return
}

func (inst *App) Status(serviceName string, timeout int) (message string, err error) {
	systemOpts.Timeout = timeout
	return systemctl.Status(serviceName, systemOpts)
}

func (inst *App) ServiceStats(serviceName string, timeout int) (resp systemctl.SystemState, err error) {
	systemOpts.Timeout = timeout
	resp, err = systemctl.State(serviceName, systemOpts)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
