package installer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApp_GetServiceNameFromAppName(t *testing.T) {
	app := App{}
	serviceName := app.GetServiceNameFromAppName("rubix-wires")
	assert.Equal(t, serviceName, "nubeio-rubix-wires.service")
}

func TestApp_GetAppNameFromRepoName(t *testing.T) {
	app := App{}
	appName := app.GetAppNameFromRepoName("wires-builds")
	assert.Equal(t, appName, "rubix-wires")
}

func TestApp_GetRepoNameFromAppName(t *testing.T) {
	app := App{}
	repoName := app.GetRepoNameFromAppName("rubix-wires")
	assert.Equal(t, repoName, "wires-builds")
}

func TestApp_GetDataDirNameFromAppName(t *testing.T) {
	app := App{}
	repoName := app.GetDataDirNameFromAppName("rubix-point-server")
	assert.Equal(t, repoName, "point-server")
}
