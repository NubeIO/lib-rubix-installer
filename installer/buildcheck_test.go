package installer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApp_GetZipBuildDetails(t *testing.T) {
	zipName1 := "flow-framework-0.7.4-5d45159f.amd64.zip"
	zipName2 := "nubeio-rubix-app-lora-serial-py-1.0.0-db9b31f0.armv7.zip"
	inst := &App{}
	result1 := inst.GetZipBuildDetails(zipName1)
	result2 := inst.GetZipBuildDetails(zipName2)

	assert.Equal(t, "flow-framework", result1.Name)
	assert.Equal(t, "v0.7.4", result1.Version)
	assert.Equal(t, "amd64", result1.Arch)
	assert.Equal(t, "flow-framework-0.7.4-5d45159f.amd64.zip", result1.ZipName)

	assert.Equal(t, "nubeio-rubix-app-lora-serial-py", result2.Name)
	assert.Equal(t, "v1.0.0", result2.Version)
	assert.Equal(t, "armv7", result2.Arch)
	assert.Equal(t, "nubeio-rubix-app-lora-serial-py-1.0.0-db9b31f0.armv7.zip", result2.ZipName)
}
