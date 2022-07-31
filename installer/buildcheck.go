package installer

import (
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

type BuildDetails struct {
	MatchedName    string `json:"matched_name"`
	MatchedVersion string `json:"matched_version"`
	MatchedArch    string `json:"matched_arch"`
	ZipName        string `json:"zip_name"`
}

func (inst *App) GetZipBuildDetails(zipName string) *BuildDetails {
	parts := strings.Split(zipName, "-")
	count := 0
	name := ""
	version := ""
	arch := ""
	for i, part := range parts { // match version
		p := strings.Split(part, ".")
		// if len is 3 eg, 0.0.1
		isNum := 0
		if len(p) == 3 || len(p) == 4 {
			// check if they are numbers
			for _, s := range p {
				if _, err := strconv.Atoi(s); err == nil {
					isNum++
				}
			}
			if isNum == 3 {
				count = i
				version = part
				version = strings.Trim(version, ".zip")
			}
			for _, s := range p { // match arch
				if s == "amd64" {
					arch = "amd64"
				}
				if s == "armv7" {
					arch = "armv7"
				}
			}
		}
	}
	name = strings.Join(parts[0:count], "-")
	return &BuildDetails{
		MatchedName:    name,
		MatchedVersion: version,
		MatchedArch:    arch,
		ZipName:        zipName,
	}
}

// GetBuildZipNameByArch // get a build by its arch
func (inst *App) GetBuildZipNameByArch(path, arch string) (*BuildDetails, error) {
	var out *BuildDetails
	details, err := getFileDetails(path)
	if err != nil {
		return out, err
	}
	for _, name := range details {
		app := inst.GetZipBuildDetails(name.Name)
		if app.MatchedArch == arch {
			out = app
		}
	}
	return out, nil
}

// GetBuildZipNames // get all the builds zips from a path
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

type fileDetails struct {
	Name      string `json:"name"`
	Extension string `json:"extension"`
	IsDir     bool   `json:"is_dir"`
}

func getFileDetails(dir string) ([]fileDetails, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []fileDetails
	var f fileDetails
	for _, file := range files {
		var extension = filepath.Ext(file.Name())
		f.Extension = extension
		f.Name = file.Name()
		f.IsDir = file.IsDir()
		out = append(out, f)
	}
	return out, nil
}
