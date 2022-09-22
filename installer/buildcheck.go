package installer

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type BuildDetails struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	Arch    string `json:"arch,omitempty"`
	ZipName string `json:"zip_name,omitempty"`
}

type FileDetails struct {
	Name      string `json:"name"`
	Extension string `json:"extension"`
	IsDir     bool   `json:"is_dir"`
}

func (inst *App) GetZipBuildDetails(zipName string) *BuildDetails {
	parts := strings.Split(zipName, "-")
	if len(parts) > 2 {
		version := parts[len(parts)-2]
		if !strings.Contains(version, "v") {
			version = fmt.Sprintf("v%s", version)
		}
		archContent := parts[len(parts)-1]
		archParts := strings.Split(archContent, ".")
		arch := ""
		if len(parts) > 1 {
			arch = archParts[1]
		}
		nameParts := parts[:len(parts)-2]
		name := strings.Join(nameParts, "-")
		return &BuildDetails{
			Name:    name,
			Version: version,
			Arch:    arch,
			ZipName: zipName,
		}
	}
	return &BuildDetails{
		ZipName: zipName,
	}
}

// GetBuildZipNameByArch get a build by its arch
func (inst *App) GetBuildZipNameByArch(path, arch string, doNotValidateArch bool) (*BuildDetails, error) {
	var out *BuildDetails
	details, err := getFileDetails(path)
	if err != nil {
		return out, err
	}
	for _, name := range details {
		app := inst.GetZipBuildDetails(name.Name)
		if app.Arch == arch {
			out = app
		}
		if doNotValidateArch { // most likely rubix-wires as it has no arch
			out = app
		}
	}
	return out, nil
}

// GetBuildZipNames get all the builds zips from a path
func (inst *App) GetBuildZipNames(path string) ([]BuildDetails, error) {
	var out []BuildDetails
	details, err := getFileDetails(path)
	if err != nil {
		return out, err
	}
	for _, name := range details {
		app := inst.GetZipBuildDetails(name.Name)
		out = append(out, *app)
	}
	return out, nil
}

func getFileDetails(dir string) ([]FileDetails, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []FileDetails
	var f FileDetails
	for _, file := range files {
		var extension = filepath.Ext(file.Name())
		f.Extension = extension
		f.Name = file.Name()
		f.IsDir = file.IsDir()
		out = append(out, f)
	}
	return out, nil
}
