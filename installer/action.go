package installer

import "errors"

type ActionType int

//go:generate stringer -type=ActionType
const (
	start ActionType = iota
	stop
	status
	enable
	disable
	isRunning
	isInstalled
	isEnabled
	isActive
	isFailed
)

func CheckAction(s string) error {
	switch s {
	case start.String():
		return nil
	case stop.String():
		return nil
	case status.String():
		return nil
	case enable.String():
		return nil
	case disable.String():
		return nil
	case isRunning.String():
		return nil
	case isInstalled.String():
		return nil
	case isEnabled.String():
		return nil
	case isActive.String():
		return nil
	case isFailed.String():
		return nil
	}
	return errors.New("invalid action type, try start")

}
