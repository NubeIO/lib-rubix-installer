package installer

import "fmt"

// This is for to support older deployments too...
// Otherwise we will use appName in new deployments
var appNameToServiceNameMap = map[string]string{
	"rubix-wires":                      "nubeio-rubix-wires.service",
	"rubix-plat":                       "nubeio-wires-plat.service",
	"rubix-point-server":               "nubeio-point-server.service",
	"nubeio-rubix-app-modbus-py":       "nubeio-rubix-app-modbus-py.service",
	"rubix-bacnet-server":              "nubeio-bacnet-server.service",
	"lora-raw":                         "nubeio-lora-raw.service",
	"rubix-bacnet-master":              "nubeio-rubix-bacnet-master.service",
	"lora-driver":                      "nubeio-rubix-app-lora-serial-py.service",
	"edge-28-driver":                   "nubeio-rubix-app-lora-serial-py.service",
	"flow-framework":                   "nubeio-flow-framework.service",
	"nubeio-rubix-app-rubix-broker-go": "nubeio-rubix-app-rubix-broker-go.service",
	"rubix-io-driver":                  "nubeio-rubix-app-pi_gpio-go.service",
	"rubix-assist":                     "nubeio-rubix-assist.service",
	"rubix-edge":                       "nubeio-rubix-edge.service",
	"rubix-user-management":            "nubeio-user-management.service",
	"rubix-data-push":                  "nubeio-data-push.service",
	"bacnet-server-driver":             "nubeio-bacnet-server-c.service",
}

var appNameToRepoNameMap = map[string]string{
	"rubix-wires":          "wires-builds",
	"rubix-plat":           "rubix-plat-build",
	"bacnet-server-driver": "bacnet-server-c",
	"rubix-io-driver":      "nubeio-rubix-app-pi-gpio-go",
	"lora-driver":          "nubeio-rubix-app-lora-serial-py",
	"edge-28-driver":       "bbb-rest-py",
}

var appNameToDataDirNameMap = map[string]string{
	"rubix-wires":                      "rubix-wires",
	"rubix-plat":                       "rubix-plat",
	"rubix-point-server":               "point-server",
	"nubeio-rubix-app-modbus-py":       "nubeio-rubix-app-modbus-py",
	"rubix-bacnet-server":              "bacnet-server",
	"lora-raw":                         "lora-raw",
	"rubix-bacnet-master":              "rubix-bacnet-master",
	"lora-driver":                      "nubeio-rubix-app-lora-serial-py",
	"edge-28-driver":                   "nubeio-rubix-app-bbb-rest-py",
	"flow-framework":                   "flow-framework",
	"nubeio-rubix-app-rubix-broker-go": "rubix-broker",
	"rubix-io-driver":                  "pi-gpio",
	"rubix-assist":                     "rubix-assist",
	"rubix-edge":                       "rubix-edge",
	"rubix-user-management":            "user-management",
	"rubix-data-push":                  "data-push",
	"bacnet-server-driver":             "bacnet-server-c",
}

func (inst *App) GetServiceNameFromAppName(appName string) string {
	if value, found := appNameToServiceNameMap[appName]; found {
		return value
	}
	return fmt.Sprintf("nubeio-%s.service", appName)
}

func (inst *App) GetAppNameFromRepoName(repoName string) string {
	for k := range appNameToRepoNameMap {
		if appNameToRepoNameMap[k] == repoName {
			return k
		}
	}
	return repoName
}

func (inst *App) GetRepoNameFromAppName(appName string) string {
	if value, found := appNameToRepoNameMap[appName]; found {
		return value
	}
	return appName
}

func (inst *App) GetDataDirNameFromAppName(appName string) string {
	if value, found := appNameToDataDirNameMap[appName]; found {
		return value
	}
	return appName
}
